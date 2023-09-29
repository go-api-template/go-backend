package cmd

import (
	"github.com/go-api-template/go-backend/modules/server"
	"github.com/spf13/cobra"
)

type runOptions struct {
	Server *server.Server
}

func newRunCmd() *cobra.Command {

	// Command options
	o := runOptions{
		Server: &server.Server{},
	}

	// Command
	runCmd := &cobra.Command{
		Use:     "serve",
		GroupID: "server",
		Short:   "Run the server",
		Aliases: []string{"r", "run", "serve", "s"},
		Args:    cobra.NoArgs,
		RunE:    o.runCmd,
	}

	// Silence usage when an error occurs
	runCmd.SilenceUsage = true

	return runCmd
}

func (o *runOptions) runCmd(_ *cobra.Command, _ []string) error {
	return o.Server.Run()
}
