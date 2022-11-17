package command

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/awaketai/goweb/framework/cobra"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/util"
	"github.com/go-git/go-git/v5"
	"github.com/jianfengye/collection"
)

var ginContribAddr = "https://github.com/gin-contrib/%s.git"

func initMiddlewareCommand() *cobra.Command {
	middlewareCommand.AddCommand(middlewareAllCommand)
	middlewareCommand.AddCommand(middlewareCreateCommand)
	middlewareCommand.AddCommand(middlewareMigrateCommand)
	return middlewareCommand
}

// middlewareCommand
var middlewareCommand = &cobra.Command{
	Use:   "middleware",
	Short: "中间件相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

// middlewareAllCommand show all middleware
var middlewareAllCommand = &cobra.Command{
	Use:   "list",
	Short: "显示所有中间件",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		middlewarePath := path.Join(appService.BaseFolder(), "http", "middleware")
		// read file
		files, err := ioutil.ReadDir(middlewarePath)
		if err != nil {
			return err
		}
		for _, f := range files {
			if f.IsDir() {
				fmt.Println(f.Name())
			}
		}
		return nil
	},
}

// middlewareCreateCommand create mieeleware command
var middlewareCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个中间件",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		fmt.Println("创建一个中间件")
		var (
			name   string
			folder string
		)
		{
			prompt := &survey.Input{
				Message: "请输入中间件名称",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入中间件所在日录名称(默认和名称相同)",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}
		if folder == "" {
			folder = name
		}
		app := container.MustMake(contract.AppKey).(contract.App)
		pFolder := app.MiddlewareFolder()
		subFolders, err := util.SubDir(pFolder)
		if err != nil {
			return err
		}
		subColl := collection.NewStrCollection(subFolders)
		if subColl.Contains(folder) {
			fmt.Println("目录已经存在")
			return nil
		}
		// create file
		if err := os.Mkdir(filepath.Join(pFolder, folder), 0700); err != nil {
			return err
		}

		funcs := template.FuncMap{"title": strings.Title}
		{
			file := filepath.Join(pFolder, folder, "middleware.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}
			t := template.Must(template.New("middleware").Funcs(funcs).Parse(middlewareTmpl))
			if err := t.Execute(f, name); err != nil {
				return err
			}

		}
		fmt.Println("中间件创建成功：目录：", filepath.Join(pFolder, folder))
		return nil
	},
}

var middlewareMigrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "迁移gin-contrib中间件：迁移地址：" + ginContribAddr,
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		fmt.Println("迁移一个中间件")
		// 1.获取参数
		var repo string
		{
			prompt := &survey.Input{
				Message: "请输入中间件名称",
			}
			err := survey.AskOne(prompt, &repo)
			if err != nil {
				return err
			}
			// 2.下载git到指定目录中
			appService := container.MustMake(contract.AppKey).(contract.App)
			middlewarePath := appService.MiddlewareFolder()
			url := fmt.Sprintf(ginContribAddr, repo)
			fmt.Println("download middleware from gin-contrib:")
			fmt.Println(url)
			_, err = git.PlainClone(path.Join(middlewarePath, repo), false, &git.CloneOptions{
				URL:      url,
				Progress: os.Stdout,
			})
			if err != nil {
				return err
			}
			// 3.删除go.mod go.sum .git
			repoFolder := path.Join(middlewarePath, repo)
			fmt.Println("remove " + path.Join(repoFolder, "go.mod"))
			os.Remove(path.Join(repoFolder, "go.mod"))
			fmt.Println("remove " + path.Join(repoFolder, "go.sum"))
			os.Remove(path.Join(repoFolder, "go.sum"))
			fmt.Println("remove " + path.Join(repoFolder, ".git"))
			os.Remove(path.Join(repoFolder, ".git"))

			// 4.替换关键词
			filepath.Walk(repoFolder, func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				if filepath.Ext(path) != ".go" {
					return nil
				}

				c, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				isContain := bytes.Contains(c, []byte("github.com/gin-gonic/gin"))
				if isContain {
					fmt.Println("更新文件：" + path)
					c = bytes.ReplaceAll(c, []byte("github.com/gin-gonic/gin"), []byte("github.com/awaketai/goweb/framework/gin"))
					err = ioutil.WriteFile(path, c, 0644)
					if err != nil {
						return err
					}
				}
				// replace refer
				//  "github.com/gin-contrib/cors"
				// "github.com/awaketai/goweb/app/http/middleware/cors"
				isContain = bytes.Contains(c, []byte("github.com/gin-contrib/"+repo))
				if isContain {
					fmt.Println("更新文件：" + path)
					c = bytes.ReplaceAll(c, []byte("github.com/gin-contrib/"+repo), []byte("github.com/awaketai/goweb/app/http/middleware/"+repo))
					err = ioutil.WriteFile(path, c, 0644)
					if err != nil {
						return err
					}
				}
				return nil
			})

		}
		return nil
	},
}

var middlewareTmpl string = `
package {{.}}

import (
	"github.com/awaketai/goweb/framework/gin"
)

func {{.|title}}Middleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
	}
}

`
