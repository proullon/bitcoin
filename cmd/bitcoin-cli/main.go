package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"os"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/spf13/cobra"

	"github.com/proullon/bitcoin/node"
)

func generateKey(cmd *cobra.Command, args []string) error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("unexpected error generating private key: %s", err)
	}

	key, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("cannot encode public key to ASN1 DER: %s", err)
	}

	fmt.Printf("key: %x  (len=%d)", key, len(key))
	return nil
}

type Notifee struct {
	l sync.Cond
}

func (n *Notifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	fmt.Printf("Sending tx\n")
	n.l.Broadcast()
}

// generateTx connects to bitcoin network and publish a transaction
func generateTx(cmd *cobra.Command, args []string) error {

	host, err := libp2p.New()
	if err != nil {
		return err
	}

	m := new(sync.Mutex)
	notifee := &Notifee{
		l: *sync.NewCond(m),
	}

	notifee.l.L.Lock()
	err = mdns.NewMdnsService(host, node.DiscoveryNamespace, notifee).Start()
	if err != nil {
		return nil
	}

	notifee.l.Wait()
	return nil
}

func main() {

	rootCmd := &cobra.Command{
		Use:   `bitcoin-cli`,
		Short: `CLI utility to interact with bitcoin-node network`,
		Long:  `CLI utility to interact with bitcoin-node network`,
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   `version`,
		Short: `Software version number`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("This is a naive implementation, do not use it.\n")
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   `generate-key`,
		Short: `Generate an ECDSA key to be used on bitcoin network`,
		RunE:  generateKey,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   `generate-tx`,
		Short: `Generate a transaction request on bitcoin network`,
		RunE:  generateTx,
	})

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
