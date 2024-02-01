package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hamiddarani/anakonda/cmd"
	"github.com/spf13/cobra"
)

func main() {
	const description = "Anakonda Application"

	root := &cobra.Command{Short: description}

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	root.AddCommand(
		cmd.Server{}.Command(trap),
		cmd.Controller{}.Command(trap),
		cmd.Worker{}.Command(trap),
	)

	if err := root.Execute(); err != nil {
		log.Fatalf("failed to execute root command\n%v", err)
	}
}
