package command

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/cobra"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/util"
	"github.com/erikdubbelboer/gspt"
	"github.com/sevlyar/go-daemon"
)

var (
	// serve default listening address
	appAddress = ""
	// daemon run
	appDaemon = false
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
	appStartCommand.Flags().BoolVarP(&appDaemon, "daemon", "d", false, "start app daemon")
	appStartCommand.Flags().StringVar(&appAddress, "address", "", "app start address,default:8080")

	appCommand.AddCommand(appStartCommand)
	appCommand.AddCommand(appRestartCommand)
	appCommand.AddCommand(appStopCommand)
	appCommand.AddCommand(appStateCommand)
	return appCommand

}

var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "start web serve",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		core := kernelService.HttpEngine()
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		if appAddress == "" {
			envService := container.MustMake(contract.EnvKey).(contract.Env)
			if envService.Get("ADDRESS") != "" {
				appAddress = envService.Get("ADDRESS")
			} else {

				if configService.IsExists("app.address") {
					appAddress = configService.GetString("app.address")
				} else {
					appAddress = ":8080"
				}
			}
		}
		if !strings.HasPrefix(appAddress, ":") {
			appAddress = ":" + appAddress
		}
		fmt.Println("listening:", appAddress)
		// ns -> s: t*1000000000
		server := &http.Server{
			Handler:           core,
			Addr:              appAddress,
			ReadHeaderTimeout: time.Duration(configService.GetInt("http.server.read_header_timeout")) * 1000000000,
			ReadTimeout:       time.Duration(configService.GetInt("http.server.read_timeout")) * 1000000000,
			IdleTimeout:       time.Duration(configService.GetInt("http.server.idle_timeout")) * 1000000000,
		}

		appService := container.MustMake(contract.AppKey).(contract.App)
		pidFolder := appService.RuntimeFolder()
		exists, _ := util.Exists(pidFolder)
		if !exists {
			if err := os.MkdirAll(pidFolder, os.ModePerm); err != nil {
				return err
			}
		}
		serverPidFile := filepath.Join(pidFolder, "app.pid")
		logFolder := appService.LogFolder()
		exists, _ = util.Exists(logFolder)
		if !exists {
			if err := os.MkdirAll(logFolder, os.ModePerm); err != nil {
				return err
			}
		}
		serveLogFile := filepath.Join(logFolder, "app.log")
		currentFolder, err := util.GetExecDir()
		if err != nil {
			return err
		}
		// daemon pattern
		if appDaemon {
			cntxt := &daemon.Context{
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				LogFileName: serveLogFile,
				LogFilePerm: 0640,
				WorkDir:     currentFolder,
				Umask:       027,
				// child process arg: ./gw app start --daemon=true
				Args: []string{"", "app", "start", "--daemon=true"},
			}
			fmt.Println("process...")
			// start child process
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}
			// parent process
			if d != nil {
				fmt.Println("start app success,pid:", d.Pid)
				fmt.Println("log file:", serveLogFile)
				return nil
			}
			defer cntxt.Release()
			// subprocess execute real app start
			fmt.Println("daemon started")
			gspt.SetProcTitle("gw app")
			if err := startAppServe(server, container); err != nil {
				fmt.Println("start app err:", err)
			}
			return nil
		}
		// not daemon
		content := strconv.Itoa(os.Getegid())
		fmt.Println("not daemon pattern:[PID]", content)
		err = ioutil.WriteFile(serverPidFile, []byte(content), 0644)
		if err != nil {
			return err
		}
		gspt.SetProcTitle("gw app")
		if err := startAppServe(server, container); err != nil {
			fmt.Println("start arr err:", err)
		}
		return nil
	},
}

var appRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "restart app",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
		content, err := os.ReadFile(serverPidFile)

		if err != nil {

			return fmt.Errorf("read file err:%w", err)
		}
		if len(content) != 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			exists, _ := util.CheckProcessExists(pid)
			fmt.Println("old pid:", pid)
			if exists {
				// kill process
				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
					return err
				}
				closeWait := 5
				configService := container.MustMake(contract.ConfigKey).(contract.Config)
				if configService.IsExists("app.close_wait") {
					closeWait = configService.GetInt("app.close_wait")
				}
				// if process killed
				for i := 0; i < closeWait*2; i++ {
					if exists, _ = util.CheckProcessExists(pid); !exists {
						break
					}
					time.Sleep(1 * time.Second)
				}
				// killed process failed
				if exists, _ = util.CheckProcessExists(pid); exists {
					fmt.Println("kill process failed:", pid)
					return fmt.Errorf("kill process failed")
				}
				if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
					return err
				}
				fmt.Println("kill process success:", pid)
			}
		}
		appDaemon = true
		return appStartCommand.RunE(cmd, args)
	},
}

var appStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "strop app",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		if len(content) != 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			// SIGTERM
			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}
			if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("stop app:", pid)
		}
		return nil
	},
}

var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "get app pid",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		if len(content) != 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if exists, _ := util.CheckProcessExists(pid); exists {
				fmt.Println("app start pid:", pid)
				return nil
			}
			fmt.Println("app not start")
		}
		return nil
	},
}

func startAppServe(server *http.Server, containter framework.Container) error {
	go func() {
		server.ListenAndServe()
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTTIN, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	closeWait := 5
	configService := containter.MustMake(contract.ConfigKey).(contract.Config)
	if configService.IsExists("app.close_wait") {
		closeWait = configService.GetInt("app.close_wait")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(closeWait)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
