package node

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"crypto/ecdsa"

	libp2p "github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"

	"github.com/proullon/bitcoin/blockchain"
)

const (
	protocolID         = "/bitcon/0.0.1"
	DiscoveryNamespace = "bitcoin"
)

const (
	TxTopic    = "transactions"
	BlockTopic = "blocks"
)

type Status string

const (
	Ready      Status = "ready"
	Failed     Status = "failed"
	Retrieving Status = "retrieving blocks"
	Waiting    Status = "waiting for transaction"
	Building   Status = "building block"
	Stopping   Status = "stopping"
)

type discoveryNotifee struct{}

func (n *discoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	fmt.Println("found peer", peerInfo.String())
}

type Node struct {
	_status Status
	m       sync.RWMutex

	privateKey     *ecdsa.PrivateKey
	bc             *blockchain.Blockchain
	tentativeBlock *blockchain.Block

	self    peer.ID
	host    host.Host
	txTopic *pubsub.Topic
	bTopic  *pubsub.Topic
	txSub   *pubsub.Subscription
	bSub    *pubsub.Subscription
}

// Run initialise the node network with provided private key and disk working directory
func Run() error {

	node, err := NewNode()
	if err != nil {
		return err
	}

	node.Start()

	// Exit properly on SIGTERM and SIGINT
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh

	fmt.Printf("\nShutting down\n")
	node.Stop()
	return nil
}

func NewNode() (*Node, error) {
	var err error
	ctx := context.Background()

	node := &Node{_status: Failed}

	node.host, err = libp2p.New(
		libp2p.Ping(false),
	)
	if err != nil {
		return nil, err
	}

	fmt.Println(node.host.Addrs())

	// PubSub
	ps, err := pubsub.NewGossipSub(ctx, node.host)
	if err != nil {
		return nil, fmt.Errorf("cannot start pubsub service: %s", err)
	}
	// join topics for transactions and blocks
	node.txTopic, err = ps.Join(TxTopic)
	if err != nil {
		return nil, fmt.Errorf("cannot join transaction topic: %s", err)
	}
	node.txSub, err = node.txTopic.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("cannot subscribe to transaction topic: %s", err)
	}
	node.bTopic, err = ps.Join(BlockTopic)
	if err != nil {
		return nil, fmt.Errorf("cannot join transaction topic: %s", err)
	}
	node.bSub, err = node.bTopic.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("cannot subscribe to transaction topic: %s", err)
	}

	// Peer discovery
	mdnsService := mdns.NewMdnsService(node.host, DiscoveryNamespace, &discoveryNotifee{})
	err = mdnsService.Start()
	if err != nil {
		return nil, fmt.Errorf("cannot start mdns service: %s", err)
	}

	node.self = node.host.ID()
	node._status = Ready
	return node, nil
}

func (n *Node) Start() {
	n.m.Lock()
	n._status = Retrieving
	n.m.Unlock()

	go n.txReader()
	go n.blockReader()
}

func (n *Node) Stop() {
	n.m.Lock()
	n._status = Stopping
	n.m.Unlock()

	n.host.Close()
}

func (n *Node) status() Status {
	n.m.RLock()
	s := n._status
	n.m.RUnlock()
	return s
}

func (n *Node) incomingTransaction(tx *blockchain.Transaction) {
}

func (n *Node) shareTransaction(tx *blockchain.Transaction) {
}

func (n *Node) incomingBlock(b *blockchain.Block) {
}

func (n *Node) startBuildingBlock() {
}

func (n *Node) stopBuildingBlock() {
}

func (n *Node) shareBlock() error {
	data, err := n.tentativeBlock.Encode()
	if err != nil {
		return err
	}

	err = n.bTopic.Publish(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) txReader() {
	ctx := context.Background()

	for n.status() != Stopping {
		msg, err := n.txSub.Next(ctx)
		if err != nil {
			fmt.Printf("cannot read from transaction sub: %s\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if msg.ReceivedFrom != n.self {
			continue
		}
		tx, err := blockchain.DecodeTransaction(msg.Data)
		if err != nil {
			fmt.Printf("cannot decode transaction: %s\n", err)
			continue
		}

		n.incomingTransaction(tx)
	}
}

func (n *Node) blockReader() {
	ctx := context.Background()

	for n.status() != Stopping {
		msg, err := n.bSub.Next(ctx)
		if err != nil {
			fmt.Printf("cannot read from bloc sub: %s\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if msg.ReceivedFrom != n.self {
			continue
		}
		b, err := blockchain.DecodeBlock(msg.Data)
		if err != nil {
			fmt.Printf("cannot decode block: %s\n", err)
			continue
		}

		n.incomingBlock(b)
	}
}
