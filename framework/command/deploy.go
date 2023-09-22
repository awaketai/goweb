package command

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/cobra"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/provider/ssh"
	"github.com/awaketai/goweb/framework/util"
	"github.com/pkg/sftp"
)

func initDeployCommand() *cobra.Command {
	deployCommand.AddCommand(deployFrontendCommand)
	deployCommand.AddCommand(deployBackendCommand)
	deployCommand.AddCommand(deployAllCommand)
	deployCommand.AddCommand(deployRollbackCommand)
	return deployCommand
}

var deployCommand = &cobra.Command{
	Use:   "deploy",
	Short: "deploy project to server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var deployFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "部署前端",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		deployFolder, err := createDeployFolder(container)
		if err != nil {
			return err
		}
		// 编译
		if err := deployBuildFrontend(cmd, deployFolder); err != nil {
			return err
		}
		// 上传部署
		return deployUploadAction(container, deployFolder, "frontend")
	},
}

var deployBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "部署后端",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		deployFolder, err := createDeployFolder(container)
		if err != nil {
			return err
		}
		if err := deployBuildBackend(cmd, deployFolder); err != nil {
			return err
		}
		return deployUploadAction(container, deployFolder, "backend")
	},
}

var deployAllCommand = &cobra.Command{
	Use:   "all",
	Short: "deploy all",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		deployFolder, err := createDeployFolder(container)
		if err != nil {
			return err
		}
		// 编译
		if err := deployBuildFrontend(cmd, deployFolder); err != nil {
			return err
		}
		if err := deployBuildBackend(cmd, deployFolder); err != nil {
			return err
		}

		return deployUploadAction(container, deployFolder, "all")
	},
}

var deployRollbackCommand = &cobra.Command{
	Use:   "rollback",
	Short: "rollback",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		if len(args) < 2 {
			return errors.New("invalid param:./gw deploy rollback [version] [frontend/backend/all]")
		}
		version := args[0]
		end := args[1]
		appService := container.MustMake(contract.AppKey).(contract.App)
		deployFolder := filepath.Join(appService.DeployFolder(), version)

		return deployUploadAction(container, deployFolder, end)
	},
}

func deployBuildBackend(c *cobra.Command, deployFolder string) error {
	container := c.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	appService := container.MustMake(contract.AppKey).(contract.App)
	envService := container.MustMake(contract.EnvKey).(contract.Env)
	logService := container.MustMake(contract.LogKey).(contract.Log)
	env := envService.AppEnv()
	binFile := "gw"
	// compile
	path, err := exec.LookPath("go")
	if err != nil {
		return err
	}
	deployBinFile := filepath.Join(deployFolder, binFile)
	cmd := exec.Command(path, "build", "-o", deployBinFile, "./")
	cmd.Env = os.Environ()
	// set GOOS and GOARCH
	if configService.GetString("deploy.backend.goos") != "" {
		cmd.Env = append(cmd.Env, "GOOS="+configService.GetString("deploy.backend.goos"))
	}
	if configService.GetString("deploy.backend.goarch") != "" {
		cmd.Env = append(cmd.Env, "GOARCH="+configService.GetString("deploy.backend.goarch"))
	}
	// ececute command
	ctx := context.Background()
	out, err := cmd.CombinedOutput()
	fmt.Println("cmd:", string(out), err, "end:", cmd.Env)
	if err != nil {
		logService.Error(ctx, "go build err", map[string]any{
			"err": err,
			"out": string(out),
		})
		return err
	}
	logService.Info(ctx, "compile successed!", nil)
	// copy env file
	ok, err := util.Exists(filepath.Join(appService.BaseFolder(), ".env"))
	if err != nil {
		return err
	}
	if ok {
		err = util.CopyFile(filepath.Join(appService.BaseFolder(), ".env"), filepath.Join(deployFolder, ".env"))
		if err != nil {
			return err
		}
	}
	// copy config file
	deployConfigFolder := filepath.Join(deployFolder, "config", env)
	ok, err = util.Exists(deployConfigFolder)
	if err != nil {
		return err
	}
	if !ok {
		// create folder
		if err := os.MkdirAll(deployConfigFolder, os.ModePerm); err != nil {
			return err
		}
	}
	if err := util.CopyFolder(filepath.Join(appService.ConfigFolder(), env), deployConfigFolder); err != nil {
		return err
	}
	logService.Info(ctx, "build local ok", nil)
	return nil
}

// deployUploadAction 上传部署文件夹，并执行前置后置shell
func deployUploadAction(container framework.Container, deployFolder, end string) error {
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	sshService := container.MustMake(contract.SSHKey).(contract.SSH)
	logService := container.MustMake(contract.LogKey).(contract.Log)
	// 遍历所有服务器
	nodes := configService.GetStringSlice("deploy.connections")
	if len(nodes) == 0 {
		return fmt.Errorf("not found any server to delploy the project")
	}
	remoteFolder := configService.GetString("deploy.remote_folder")
	if remoteFolder == "" {
		return fmt.Errorf("remote folder not assigned")
	}
	preActions := make([]string, 0, 1)
	postActions := make([]string, 0, 1)
	// 可部署前端和后端代码
	if end == "frontend" || end == "all" {
		preActions = append(preActions, configService.GetStringSlice("deploy.frontend.pre_action")...)
		postActions = append(postActions, configService.GetStringSlice("deploy.frontend.post_action")...)
	}
	if end == "backend" || end == "all" {
		preActions = append(preActions, configService.GetStringSlice("deploy.backend.pre_action")...)
		postActions = append(postActions, configService.GetStringSlice("deploy.backend.post_action")...)
	}
	// for every server,execute deploy action
	for _, node := range nodes {
		sshClient, err := sshService.GetClient(ssh.WithConfigPath(node))
		if err != nil {
			return fmt.Errorf("execute deploy upload action err when get ssh client:%w", err)
		}
		client, err := sftp.NewClient(sshClient)
		if err != nil {
			return fmt.Errorf("execute deploy upload action err when get sftp client:%w", err)
		}
		// execute all pre command
		for _, action := range preActions {
			session, err := sshClient.NewSession()
			if err != nil {
				return err
			}
			logService.Info(context.Background(), "execute pre action start", map[string]any{
				"cmd":        action,
				"connection": node,
			})
			// execute command,and waitting result
			bts, err := session.CombinedOutput(action)
			if err != nil {
				session.Close()
				return err
			}
			session.Close()
			// execute successed
			logService.Info(context.Background(), "execute pre action result", map[string]any{
				"cmd":        action,
				"connection": node,
				"out":        strings.ReplaceAll(string(bts), "\n", ""),
			})
		}

		// upload executable file to sever
		if err := uploadFolderToSFTP(container, deployFolder, remoteFolder, client); err != nil {
			return err
		}
		logService.Info(context.Background(), "upload success", nil)
		// execute post command
		for _, action := range postActions {
			session, err := sshClient.NewSession()
			if err != nil {
				return err
			}

			logService.Info(context.Background(), "execute post action start", map[string]any{
				"cmd":        action,
				"connection": node,
			})
			bts, err := session.CombinedOutput(action)
			if err != nil {
				session.Close()
				return err
			}
			session.Close()
			logService.Info(context.Background(), "execute post action result", map[string]any{
				"cmd":        action,
				"connection": node,
				"out":        strings.ReplaceAll(string(bts), "\n", ""),
			})
		}
	}

	return nil
}

func uploadFolderToSFTP(container framework.Container, localFolder, remoteFolder string, client *sftp.Client) error {
	logService := container.MustMake(contract.LogKey).(contract.Log)
	return filepath.Walk(localFolder, func(path string, info fs.FileInfo, err error) error {
		// 获取除了folder前缀的后续就文件名称
		relPath := strings.Replace(path, localFolder, "", 1)
		fmt.Println("rel path:", relPath)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			logService.Info(context.Background(), "mkdir:"+filepath.Join(remoteFolder, relPath), nil)
			// create folder
			return client.MkdirAll(filepath.Join(remoteFolder, relPath))
		}
		// open local file
		rf, err := os.Open(filepath.Join(localFolder, relPath))
		if err != nil {
			return err
		}
		defer rf.Close()
		// inspect file size
		rfStat, err := rf.Stat()
		if err != nil {
			return err
		}
		f, err := client.Create(filepath.Join(remoteFolder, relPath))
		if err != nil {
			return err
		}
		defer f.Close()
		logStr := "upload:" + filepath.Join(localFolder, relPath) + "to remote:" + filepath.Join(remoteFolder, relPath)
		// 大于2M的文件显示进度
		if rfStat.Size() > 2*1024*1024 {
			logService.Info(context.Background(), logStr+" start", nil)
			go showUploadProgress(localFolder, remoteFolder, relPath, client, logService, rfStat)
		}
		// 将本地文件并发读取到远端
		if _, err := f.ReadFromWithConcurrency(rf, 10); err != nil {
			return fmt.Errorf("read local file to remote err:%w", err)
		}
		logService.Info(context.Background(), logStr+" finished", nil)

		return nil
	})
}

// showUploadProgress 显示文件上传进度
func showUploadProgress(localFolder, remoteFolder, relPath string, client *sftp.Client, logService contract.Log, rfStat os.FileInfo) {
	ticker := time.NewTicker(2 * time.Second)
	remoteFile := filepath.Join(remoteFolder, relPath)
	for range ticker.C {
		remoteFileInfo, err := client.Stat(remoteFile)
		if err != nil {
			logService.Error(context.Background(), "stat err", map[string]any{
				"err":         err,
				"remote_file": remoteFile,
			})
			continue
		}
		remoteSize := remoteFileInfo.Size()
		// upload end
		if remoteSize >= rfStat.Size() {
			break
		}
		percent := int(remoteSize * 100 / rfStat.Size())
		logService.Info(context.Background(), "upload "+filepath.Join(localFolder, relPath)+
			"to remote:"+filepath.Join(remoteFolder, relPath)+fmt.Sprintf("%v%% %v/%v", percent, remoteSize, rfStat.Size()), nil)
	}
}

func deployBuildFrontend(c *cobra.Command, deployFolder string) error {
	container := c.GetContainer()
	appService := container.MustMake(contract.AppKey).(contract.App)
	// 编译前端
	if err := buildFrontendCommand.RunE(c, []string{}); err != nil {
		return err
	}
	frontendFolder := filepath.Join(deployFolder, "dist")
	if err := os.Mkdir(frontendFolder, os.ModePerm); err != nil {
		return err
	}
	buildFolder := filepath.Join(appService.BaseFolder(), "dist")
	if err := util.CopyFile(buildFolder, frontendFolder); err != nil {
		return err
	}

	return nil
}

func createDeployFolder(contaner framework.Container) (string, error) {
	appService := contaner.MustMake(contract.AppKey).(contract.App)
	deployFolder := appService.DeployFolder()
	deployVersion := time.Now().Format("20060102150405")
	versionFolder := filepath.Join(deployFolder, deployVersion)
	exists, err := util.Exists(versionFolder)
	if err != nil {
		return "", err
	}
	if !exists {
		err = os.MkdirAll(versionFolder, os.ModePerm)
		return versionFolder, err
	}

	return versionFolder, nil
}
