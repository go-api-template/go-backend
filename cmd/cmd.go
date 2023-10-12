package cmd

import (
	"github.com/go-api-template/go-backend/modules/config"
	"github.com/gookit/color"
	"github.com/spf13/cobra"

	_ "github.com/go-api-template/go-backend/modules/logger/main"
)

// Options stores the options which are used in the command
type options struct {
	RunCmd *cobra.Command
}

// This is the main command of the application
// It is used to create sub commands
func newCmd() *cobra.Command {

	// Command options
	o := options{
		RunCmd: newRunCmd(),
	}

	// Create the main command
	mainCmd := &cobra.Command{
		Use:     config.Config.App.Name,
		Version: config.Config.App.Version,
		RunE:    o.mainCmd,
	}

	// Sub commands
	mainCmd.AddGroup(&cobra.Group{ID: "server", Title: color.HiGreen.Sprint("Server:")})
	mainCmd.AddCommand(o.RunCmd)

	return mainCmd
}

// If no command is provided then run the server
func (o *options) mainCmd(_ *cobra.Command, args []string) error {
	return o.RunCmd.RunE(o.RunCmd, args)
}

func Execute() (err error) {

	// Construct the root command
	mainCmd := newCmd()

	// Execute the command
	if err = mainCmd.Execute(); err != nil {
		return
	}

	return nil
}
