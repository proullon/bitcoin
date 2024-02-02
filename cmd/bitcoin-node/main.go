package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/proullon/bitcoin/node"
)

func run(cmd *cobra.Command, args []string) error {
	return node.Run()
}

func main() {

	rootCmd := &cobra.Command{
		Use:   `bitcoin-node`,
		Short: `bitcoin node`,
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   `version`,
		Short: `Software version number`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("This is a naive implementation, do not use it.\n")
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   `run`,
		Short: `run the bitcoin node`,
		RunE:  run,
	})

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
