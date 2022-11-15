package command

import "github.com/awaketai/goweb/framework/cobra"

func AddKernelCommands(root *cobra.Command) {
	root.AddCommand(DemoCommand)
	root.AddCommand(initAppCommand())
	root.AddCommand(initCronCommand())
	root.AddCommand(initEnvCommand())
	// create provider command
	root.AddCommand(initProviderCommand())
	// create cmd command
	root.AddCommand(initCmdCommand())
}
