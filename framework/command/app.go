package command

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/awaketai/goweb/framework/cobra"
	"github.com/awaketai/goweb/framework/contract"
)

var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，其包含业务启动，关闭，重启，查询等功能",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 打印帮助功
		cmd.Help()
		return nil

	},
}

func initAppCommand() *cobra.Command {
	appCommand.AddCommand(appStartCommand)
	return appCommand

}

var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个web服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		core := kernelService.HttpEngine()

		server := &http.Server{
			Handler: core,
			Addr:    ":8080",
		}
		go func() {
			server.ListenAndServe()
		}()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-quit
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		return nil
	},
}
