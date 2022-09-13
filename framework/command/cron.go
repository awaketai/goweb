package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/awaketai/goweb/framework/cobra"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/util"
	"github.com/erikdubbelboer/gspt"
	"github.com/sevlyar/go-daemon"
)

var cronDaemon = false

func initCronCommand() *cobra.Command {
	cronStartCommand.Flags().BoolVarP(&cronDaemon, "daemon", "d", false, "start serve daemon")
	cronCommand.AddCommand(cronStartCommand)
	cronCommand.AddCommand(cronRestartCommand)
	cronCommand.AddCommand(cronListCommand)
	cronCommand.AddCommand(cronStateCommand)
	cronCommand.AddCommand(cronStopCommand)
	return cronCommand
}

var cronCommand = &cobra.Command{
	Use:   "cron",
	Short: "定时任务相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var cronListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有的定时任务",
	RunE: func(cmd *cobra.Command, args []string) error {
		cronSpecs := cmd.Root().CronSpecs
		ps := [][]string{}
		for _, spec := range cronSpecs {
			line := []string{spec.Type, spec.Spec, spec.Cmd.Use, spec.Cmd.Short, spec.ServiceName}
			ps = append(ps, line)
		}

		util.PrettyPrint(ps)
		return nil
	},
}

var cronStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.Root().GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		pidFolder := appService.RuntimeFolder()
		serverPidFile := filepath.Join(pidFolder, "cron.pid")
		logFolder := appService.LogFolder()
		serverLogFile := filepath.Join(logFolder, "cron.log")
		currentFolder := appService.BaseFolder()

		if cronDaemon {
			daemonCtx := &daemon.Context{
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				LogFileName: serverLogFile,
				LogFilePerm: 0640,
				WorkDir:     currentFolder,
				Umask:       027,
				// 子进程参数 ./gw cron start --daemon=true
				Args: []string{"", "cron", "start", "--daemon=true"},
			}
			// 启动子进程，d不为空表示当前是父进程，d为空表示当前是子进程
			d, err := daemonCtx.Reborn()
			if err != nil {
				return err
			}

			if d != nil {
				// 父进程不做处理
				fmt.Println("corn serve started,pid:", d.Pid)
				fmt.Println("log file:", serverLogFile)
				return nil
			}
			// 子进程执行cron.Run
			defer daemonCtx.Release()
			fmt.Println("daemon started")
			gspt.SetProcTitle("gw crno")
			cmd.Root().Cron.Run()
			return nil
		}
		fmt.Println("start cron job[no daemon]")
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID] ", content)
		err := ioutil.WriteFile(serverPidFile, []byte(content), 0664)
		if err != nil {
			return err
		}

		gspt.SetProcTitle("gw cron")
		cmd.Root().Cron.Run()
		return nil
	},
}

var cronRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重启cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.Root().GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		if len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			processExists, err := util.CheckProcessExists(pid)
			if err != nil {
				return err
			}
			if processExists {
				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
					return err
				}
				// check process
				for i := 0; i < 10; i++ {
					processExists, err = util.CheckProcessExists(pid)
					if err != nil {
						return err
					}
					if !processExists {
						break
					}
					time.Sleep(1 * time.Second)

				}
				fmt.Println("kill process:", strconv.Itoa(pid))
			}
		}
		cronDaemon = true
		return cronStartCommand.RunE(cmd, args)

	},
}

var cronStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		service := container.MustMake(contract.AppKey).(contract.App)
		serverPidFile := filepath.Join(service.RuntimeFolder(), "cron.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}
			if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("stop pid:", pid)
		}
		return nil
	},
}

var cronStateCommand = &cobra.Command{
	Use:   "state",
	Short: "corn常驻进程状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		service := container.MustMake(contract.AppKey).(contract.App)
		serverPidFile := filepath.Join(service.RuntimeFolder(), "cron.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			exists, err := util.CheckProcessExists(pid)
			if err != nil {
				return err
			}
			if exists {
				fmt.Println("cron server started,pid:", pid)
				return nil
			}
		}
		fmt.Println("no cron server start")
		return nil
	},
}
