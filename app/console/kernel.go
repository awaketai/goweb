package console

import (
	"time"

	"github.com/awaketai/goweb/app/console/command/demo"
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/cobra"
	"github.com/awaketai/goweb/framework/command"
)

func RunCommand(container framework.Container) error {
	var rootCmd = &cobra.Command{
		Use:   "gw",
		Short: "gw command",
		Long:  "gw framework cli tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpCmd()
			return cmd.Help()
		},
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	rootCmd.SetContainer(container)
	// 绑定框架的命令
	command.AddKernelCommands(rootCmd)
	// 绑定业务命令
	AddAppCommand(rootCmd)
	return rootCmd.Execute()
}

func AddAppCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(demo.InitFoo())
	// 每秒执行一次
	// rootCmd.AddCronCommand("* * * * * *", demo.FooCommand)
	// 分布式调度任务，每个节点每5s调用一次Foo命令，抢占到了调度任务的节点将抢占锁持续挂载2s后释放
	rootCmd.AddDistributedCronCommand("foo_func_for_test", "*/5 * * * * *", demo.FooCommand, 2*time.Second)
}

func AddDeployCommand(rootCmd *cobra.Command){
	rootCmd.AddCommand()
}
