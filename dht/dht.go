package dht

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/lbryio/lbry.go/errors"
	"github.com/lbryio/lbry.go/stopOnce"
	"github.com/lbryio/reflector.go/dht/bits"
	"github.com/spf13/cast"

	log "github.com/sirupsen/logrus"
)

func init() {
	//log.SetFormatter(&log.TextFormatter{ForceColors: true})
	//log.SetLevel(log.DebugLevel)
}

const (
	Network = "udp4"

	// TODO: all these constants should be defaults, and should be used to set values in the standard Config. then the code should use values in the config
	// TODO: alternatively, have a global Config for constants. at least that way tests can modify the values
	alpha           = 3             // this is the constant alpha in the spec
	bucketSize      = 8             // this is the constant k in the spec
	nodeIDLength    = bits.NumBytes // bytes. this is the constant B in the spec
	nodeIDBits      = bits.NumBits  // number of bits in node ID
	messageIDLength = 20            // bytes.

	udpRetry            = 3
	udpTimeout          = 5 * time.Second
	udpMaxMessageLength = 1024 // bytes. I think our longest message is ~676 bytes, so I rounded up

	maxPeerFails = 3 // after this many failures, a peer is considered bad and will be removed from the routing table
	//tExpire     = 60 * time.Minute // the time after which a key/value pair expires; this is a time-to-live (TTL) from the original publication date
	tReannounce = 50 * time.Minute // the time after which the original publisher must republish a key/value pair
	tRefresh    = 1 * time.Hour    // the time after which an otherwise unaccessed bucket must be refreshed
	//tReplicate   = 1 * time.Hour    // the interval between Kademlia replication events, when a node is required to publish its entire database
	//tNodeRefresh = 15 * time.Minute // the time after which a good node becomes questionable if it has not messaged us

	compactNodeInfoLength = nodeIDLength + 6 // nodeID + 4 for IP + 2 for port

	tokenSecretRotationInterval = 5 * time.Minute // how often the token-generating secret is rotated
)

// Config represents the configure of dht.
type Config struct {
	// this node's address. format is `ip:port`
	Address string
	// the seed nodes through which we can join in dht network
	SeedNodes []string
	// the hex-encoded node id for this node. if string is empty, a random id will be generated
	NodeID string
	// print the state of the dht every X time
	PrintState time.Duration
}

// NewStandardConfig returns a Config pointer with default values.
func NewStandardConfig() *Config {
	return &Config{
		Address: "0.0.0.0:4444",
		SeedNodes: []string{
			"lbrynet1.lbry.io:4444",
			"lbrynet2.lbry.io:4444",
			"lbrynet3.lbry.io:4444",
		},
	}
}

// DHT represents a DHT node.
type DHT struct {
	// config
	conf *Config
	// local contact
	contact Contact
	// node
	node *Node
	// stopper to shut down DHT
	stop *stopOnce.Stopper
	// channel is closed when DHT joins network
	joined chan struct{}
	// lock for announced list
	lock *sync.RWMutex
	// list of bitmaps that need to be reannounced periodically
	announced map[bits.Bitmap]bool
}

// New returns a DHT pointer. If config is nil, then config will be set to the default config.
func New(config *Config) *DHT {
	if config == nil {
		config = NewStandardConfig()
	}

	d := &DHT{
		conf:      config,
		stop:      stopOnce.New(),
		joined:    make(chan struct{}),
		lock:      &sync.RWMutex{},
		announced: make(map[bits.Bitmap]bool),
	}
	return d
}

func (dht *DHT) connect(conn UDPConn) error {
	contact, err := getContact(dht.conf.NodeID, dht.conf.Address)
	if err != nil {
		return err
	}

	dht.contact = contact
	dht.node = NewNode(contact.ID)

	err = dht.node.Connect(conn)
	if err != nil {
		return err
	}
	return nil
}

// Start starts the dht
func (dht *DHT) Start() error {
	listener, err := net.ListenPacket(Network, dht.conf.Address)
	if err != nil {
		return errors.Err(err)
	}
	conn := listener.(*net.UDPConn)

	err = dht.connect(conn)
	if err != nil {
		return err
	}

	dht.join()
	log.Debugf("[%s] DHT ready on %s (%d nodes found during join)",
		dht.node.id.HexShort(), dht.contact.Addr().String(), dht.node.rt.Count())

	go dht.startReannouncer()

	return nil
}

// join makes current node join the dht network.
func (dht *DHT) join() {
	defer close(dht.joined) // if anyone's waiting for join to finish, they'll know its done

	log.Debugf("[%s] joining network", dht.node.id.HexShort())

	// ping nodes, which gets their real node IDs and adds them to the routing table
	atLeastOneNodeResponded := false
	for _, addr := range dht.conf.SeedNodes {
		err := dht.Ping(addr)
		if err != nil {
			log.Error(errors.Prefix(fmt.Sprintf("[%s] join", dht.node.id.HexShort()), err))
		} else {
			atLeastOneNodeResponded = true
		}
	}

	if !atLeastOneNodeResponded {
		log.Errorf("[%s] join: no nodes responded to initial ping", dht.node.id.HexShort())
		return
	}

	// now call iterativeFind on yourself
	_, _, err := FindContacts(dht.node, dht.node.id, false, dht.stop.Ch())
	if err != nil {
		log.Errorf("[%s] join: %s", dht.node.id.HexShort(), err.Error())
	}

	// TODO: after joining, refresh all buckets further away than our closest neighbor
	// http://xlattice.sourceforge.net/components/protocol/kademlia/specs.html#join
}

// WaitUntilJoined blocks until the node joins the network.
func (dht *DHT) WaitUntilJoined() {
	if dht.joined == nil {
		panic("dht not initialized")
	}
	<-dht.joined
}

// Shutdown shuts down the dht
func (dht *DHT) Shutdown() {
	log.Debugf("[%s] DHT shutting down", dht.node.id.HexShort())
	dht.stop.StopAndWait()
	dht.node.Shutdown()
	log.Debugf("[%s] DHT stopped", dht.node.id.HexShort())
}

// Ping pings a given address, creates a temporary contact for sending a message, and returns an error if communication
// fails.
func (dht *DHT) Ping(addr string) error {
	raddr, err := net.ResolveUDPAddr(Network, addr)
	if err != nil {
		return err
	}

	tmpNode := Contact{ID: bits.Rand(), IP: raddr.IP, Port: raddr.Port}
	res := dht.node.Send(tmpNode, Request{Method: pingMethod})
	if res == nil {
		return errors.Err("no response from node %s", addr)
	}

	return nil
}

// Get returns the list of nodes that have the blob for the given hash
func (dht *DHT) Get(hash bits.Bitmap) ([]Contact, error) {
	contacts, found, err := FindContacts(dht.node, hash, true, dht.stop.Ch())
	if err != nil {
		return nil, err
	}

	if found {
		return contacts, nil
	}
	return nil, nil
}

// Add adds the hash to the list of hashes this node has
func (dht *DHT) Add(hash bits.Bitmap) error {
	// TODO: calling Add several times quickly could cause it to be announced multiple times before dht.announced[hash] is set to true
	dht.lock.RLock()
	exists := dht.announced[hash]
	dht.lock.RUnlock()
	if exists {
		return nil
	}
	return dht.announce(hash)
}

// Announce announces to the DHT that this node has the blob for the given hash
func (dht *DHT) announce(hash bits.Bitmap) error {
	contacts, _, err := FindContacts(dht.node, hash, false, dht.stop.Ch())
	if err != nil {
		return err
	}

	// if we found less than K contacts, or current node is closer than farthest contact
	if len(contacts) < bucketSize || dht.node.id.Xor(hash).Less(contacts[bucketSize-1].ID.Xor(hash)) {
		// pop last contact, and self-store instead
		contacts[bucketSize-1] = dht.contact
	}

	wg := &sync.WaitGroup{}
	for _, c := range contacts {
		wg.Add(1)
		go func(c Contact) {
			dht.storeOnNode(hash, c)
			wg.Done()
		}(c)
	}

	wg.Wait()

	dht.lock.Lock()
	dht.announced[hash] = true
	dht.lock.Unlock()

	return nil
}

func (dht *DHT) startReannouncer() {
	tick := time.NewTicker(tReannounce)
	for {
		select {
		case <-dht.stop.Ch():
			return
		case <-tick.C:
			dht.lock.RLock()
			for h := range dht.announced {
				dht.stop.Add(1)
				go func(bm bits.Bitmap) {
					defer dht.stop.Done()
					err := dht.announce(bm)
					if err != nil {
						log.Error("error re-announcing bitmap - ", err)
					}
				}(h)
			}
			dht.lock.RUnlock()
		}
	}
}

func (dht *DHT) storeOnNode(hash bits.Bitmap, c Contact) {
	// self-store
	if dht.contact.Equals(c) {
		dht.node.Store(hash, c)
		return
	}

	resCh, cancel := dht.node.SendCancelable(c, Request{
		Method: findValueMethod,
		Arg:    &hash,
	})

	var res *Response

	select {
	case res = <-resCh:
	case <-dht.stop.Ch():
		cancel()
		return
	}

	if res == nil {
		return // request timed out
	}

	resCh, cancel = dht.node.SendCancelable(c, Request{
		Method: storeMethod,
		StoreArgs: &storeArgs{
			BlobHash: hash,
			Value: storeArgsValue{
				Token:  res.Token,
				LbryID: dht.contact.ID,
				Port:   dht.contact.Port,
			},
		},
	})

	go func() {
		select {
		case <-resCh:
		case <-dht.stop.Ch():
			cancel()
		}
	}()
}

// PrintState prints the current state of the DHT including address, nr outstanding transactions, stored hashes as well
// as current bucket information.
func (dht *DHT) PrintState() {
	log.Printf("DHT node %s at %s", dht.contact.String(), time.Now().Format(time.RFC822Z))
	log.Printf("Outstanding transactions: %d", dht.node.CountActiveTransactions())
	log.Printf("Stored hashes: %d", dht.node.store.CountStoredHashes())
	log.Printf("Buckets:")
	for _, line := range strings.Split(dht.node.rt.BucketInfo(), "\n") {
		log.Println(line)
	}
}

func (dht DHT) ID() bits.Bitmap {
	return dht.contact.ID
}

func getContact(nodeID, addr string) (Contact, error) {
	var c Contact
	if nodeID == "" {
		c.ID = bits.Rand()
	} else {
		c.ID = bits.FromHexP(nodeID)
	}

	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		return c, errors.Err(err)
	} else if ip == "" {
		return c, errors.Err("address does not contain an IP")
	} else if port == "" {
		return c, errors.Err("address does not contain a port")
	}

	c.IP = net.ParseIP(ip)
	if c.IP == nil {
		return c, errors.Err("invalid ip")
	}

	c.Port, err = cast.ToIntE(port)
	if err != nil {
		return c, errors.Err(err)
	}

	return c, nil
}
