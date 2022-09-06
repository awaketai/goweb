package command

import (
	"fmt"

	"github.com/awaketai/goweb/framework/cobra"
	"github.com/awaketai/goweb/framework/contract"
)

var DemoCommand = &cobra.Command{
	Use:   "demo",
	Short: "demo for framework",
	Run: func(cmd *cobra.Command, args []string) {
		container := cmd.GetContainer()
		appServicde := container.MustMake(contract.AppKey).(contract.App)
		fmt.Println("app base folder[by demo]:", appServicde.BaseFolder())
	},
}
