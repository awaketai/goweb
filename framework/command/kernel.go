package command

import "github.com/awaketai/goweb/framework/cobra"

func AddKernelCommands(root *cobra.Command) {
	root.AddCommand(DemoCommand)
	root.AddCommand(initAppCommand())
}
