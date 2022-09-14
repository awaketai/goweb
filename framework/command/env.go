package command

import (
	"fmt"

	"github.com/awaketai/goweb/framework/cobra"
	"github.com/awaketai/goweb/framework/contract"
)

var envCommand = &cobra.Command{
	Use:   "env",
	Short: "获取当前的App环境",
	Run: func(cmd *cobra.Command, args []string) {
		container := cmd.GetContainer()
		service := container.MustMake(contract.EnvKey).(contract.Env)
		fmt.Printf("env:%s\n", service.AppEnv())
	},
}

func initEnvCommand() *cobra.Command {
	// env child command

	return envCommand
}
