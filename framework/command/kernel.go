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
	// create middleware command
	root.AddCommand(initMiddlewareCommand())
	root.AddCommand(initNewProjectCommand())
	// create swagger
	root.AddCommand(initSwaggerCommand())
	// build project
	root.AddCommand(initBuildCommand())
	// deploy command
	root.AddCommand(initDeployCommand())
}
