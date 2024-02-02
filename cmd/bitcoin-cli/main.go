package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
