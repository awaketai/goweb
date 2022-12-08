package command

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/awaketai/goweb/framework/cobra"
	"github.com/awaketai/goweb/framework/util"
	"github.com/google/go-github/v39/github"
	"github.com/spf13/cast"
)

func initNewProjectCommand() *cobra.Command {
	return newProjectCommand
}

var newProjectCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "create a new project",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentPath, err := util.GetExecDir()
		if err != nil {
			return err
		}
		var (
			name    string
			folder  string
			mod     string
			version string
			release *github.RepositoryRelease
		)
		{
			prompt := &survey.Input{
				Message: "请输入目当名称",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
			folder = filepath.Join(currentPath, name)
			exists, _ := util.Exists(folder)
			if exists {
				isForce := false
				prompt2 := &survey.Confirm{
					Message: "目录" + folder + "已经存在？是否删除重新创建？",
					Default: false,
				}
				err := survey.AskOne(prompt2, &isForce)
				if err != nil {
					return err
				}
				if isForce {
					if err := os.RemoveAll(folder); err != nil {
						return err
					}
				} else {
					fmt.Println("目录已经存在，创建应用失败")
					return nil
				}
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入模块名称(go.mod中的module，默认为文件夹名称)",
			}
			err := survey.AskOne(prompt, &mod)
			if err != nil {
				return err
			}
			if mod == "" {
				mod = name
			}
		}
		{
			// get goweb framework version
			client := github.NewClient(nil)
			prompt := &survey.Input{
				Message: "请输入版本(参考 https://github.com/awaketai/goweb/releases)，默认为最新版",
			}
			err := survey.AskOne(prompt, &version)
			if err != nil {
				return err
			}
			if version != "" {
				release, _, err = client.Repositories.GetReleaseByTag(context.Background(), "awaketai", "goweb", version)
				if err != nil {
					return err
				}
				if release == nil {
					fmt.Println("版本不存在，创建应用失败")
					return nil
				}
			}
			if version == "" {
				release, _, err = client.Repositories.GetLatestRelease(context.Background(), "awaketai", "goweb")
				if err != nil {
					return err
				}
				version = release.GetTagName()
			}
			fmt.Println("========================")
			fmt.Println("create project start...")
			fmt.Println("project directory:", folder)
			fmt.Println("project name:", mod)
			fmt.Println("goweb framework version:", release.GetTagName())
			templateFolder := filepath.Join(currentPath, "template_goweb_"+version+"_"+cast.ToString(time.Now().Unix()))
			os.Mkdir(templateFolder, os.ModePerm)
			fmt.Println("创建临时目录：", templateFolder)
			// copy
			url := release.GetZipballURL()
			if url == "" {
				return fmt.Errorf("download zip file url blank")
			}
			fmt.Println("download url:", url)
			err = util.DownloadFile(filepath.Join(templateFolder, "template.zip"), url)
			if err != nil {
				os.RemoveAll(templateFolder)
				return err
			}
			fmt.Println("downlad zip file to template file success.")
			_, err = util.Unzip(filepath.Join(templateFolder, "template.zip"), templateFolder)
			if err != nil {
				os.RemoveAll(templateFolder)
				return err
			}
			fInfos, err := ioutil.ReadDir(templateFolder)
			if err != nil {
				os.RemoveAll(templateFolder)
				return err
			}
			for _, fInfo := range fInfos {
				if fInfo.IsDir() && strings.Contains(fInfo.Name(), "awaketai-goweb") {
					if err := os.Rename(filepath.Join(templateFolder, fInfo.Name()), folder); err != nil {
						os.RemoveAll(templateFolder)
						return err
					}
				}
			}
			fmt.Println("decommpress zip file success")
			if err := os.RemoveAll(templateFolder); err != nil {
				return err
			}
			fmt.Println("del template file success")
			err = os.RemoveAll(path.Join(folder, ".git"))
			if err != nil {
				return err
			}
			fmt.Println("del .git directory")

			err = os.RemoveAll(path.Join(folder, "framework"))
			if err != nil {
				return err
			}
			fmt.Println("del framework")
			filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				c, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				if path == filepath.Join(folder, "go.mod") {
					fmt.Println("update file:", path)
					c = bytes.ReplaceAll(c, []byte("module github.com/awaketai/goweb"), []byte("module "+mod))
					c = bytes.ReplaceAll(c, []byte("require ("), []byte("require (\n\tgithub.com/awaketai/goweb "+version))
					err = ioutil.WriteFile(path, c, 0644)
					if err != nil {
						return err
					}
					return nil
				}
				isContain := bytes.Contains(c, []byte("github.com/awaketai/goweb/app"))
				if isContain {
					fmt.Println("update file:", path)
					c = bytes.ReplaceAll(c, []byte("github.com/awaketai/goweb/app"), []byte(mod+"/app"))
					err = ioutil.WriteFile(path, c, 0644)
					if err != nil {
						return err
					}
				}
				return nil
			})
			fmt.Println("creat project sucess")
			fmt.Println("direcotry:", folder)
			fmt.Println("=================")
			return nil
		}
	},
}
