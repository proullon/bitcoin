package node

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"crypto/ecdsa"

	libp2p "github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"

	"github.com/proullon/bitcoin/blockchain"
)

const (
	protocolID         = "/bitcon/0.0.1"
	discoveryNamespace = "bitcoin"
)

const (
	TxTopic    = "transactions"
	BlockTopic = "blocks"
)

type discoveryNotifee struct{}

func (n *discoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	fmt.Println("found peer", peerInfo.String())
}

// Run initialise the node network with provided private key and disk working directory
func Run() error {
	ctx := context.Background()

	host, err := libp2p.New(
		libp2p.Ping(false),
	)
	if err != nil {
		return err
	}
	defer host.Close()

	fmt.Println(host.Addrs())

	// PubSub
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return fmt.Errorf("cannot start pubsub service: %s", err)
	}
	// join topics for transactions and blocks
	txTopic, err := ps.Join(TxTopic)
	if err != nil {
		return fmt.Errorf("cannot join transaction topic: %s", err)
	}
	txSub, err := txTopic.Subscribe()
	if err != nil {
		return fmt.Errorf("cannot subscribe to transaction topic: %s", err)
	}
	bTopic, err := ps.Join(BlockTopic)
	if err != nil {
		return fmt.Errorf("cannot join transaction topic: %s", err)
	}
	bSub, err := bTopic.Subscribe()
	if err != nil {
		return fmt.Errorf("cannot subscribe to transaction topic: %s", err)
	}

	_ = txSub
	_ = bSub

	// Peer discovery
	mdnsService := mdns.NewMdnsService(host, discoveryNamespace, &discoveryNotifee{})
	err = mdnsService.Start()
	if err != nil {
		return fmt.Errorf("cannot start mdns service: %s", err)
	}

	// Exit properly on SIGTERM and SIGINT
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh

	fmt.Printf("\nShutting down\n")
	return nil
}

type Node struct {
	privateKey     *ecdsa.PrivateKey
	bc             *blockchain.Blockchain
	tentativeBlock *blockchain.Block

	txTopic *pubsub.Topic
	bTopic  *pubsub.Topic
	txSub   *pubsub.Subscription
	bSub    *pubsub.Subscription
}

func (n *Node) Start() {
}

func (n *Node) Stop() {
}

func (n *Node) run() {
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
