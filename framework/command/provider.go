package command

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/cobra"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/util"
	"github.com/jianfengye/collection"
)

// provierCommand 二级命令
var provierCommand = &cobra.Command{
	Use:   "provider",
	Short: "服务提供相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

// initProviderCommand 初始化provier相关服务
func initProviderCommand() *cobra.Command {
	provierCommand.AddCommand(providerCreateCommand)
	provierCommand.AddCommand(providerListCommand)
	return provierCommand
}

// providerListCommand 列出容器内所有服务
var providerListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出容器内所有服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		webContainer := container.(*framework.WebContainer)
		list := webContainer.NameList()
		for _, line := range list {
			println(line)
		}
		return nil
	},
}

var providerCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		fmt.Println("创建一个服务")
		var (
			name   string
			folder string
		)
		{
			prompt := &survey.Input{
				Message: "请输入服务名称：",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入服务所在目录名称：",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}
		// check service if not exists
		providers := container.(*framework.WebContainer).NameList()
		providerColl := collection.NewStrCollection(providers)
		if providerColl.Contains(name) {
			fmt.Println("服务已经存在")
			return nil
		}
		if folder == "" {
			folder = name
		}
		app := container.MustMake(contract.AppKey).(contract.App)
		pFolder := app.ProviderFolder()
		subFolders, err := util.SubDir(pFolder)
		if err != nil {
			return err
		}
		subColl := collection.NewStrCollection(subFolders)
		if subColl.Contains(folder) {
			fmt.Println("目录名称已经存在")
			return nil
		}

		// create assign file /gw/app/provider/{user}
		if err := os.Mkdir(filepath.Join(pFolder, folder), 0700); err != nil {
			return err
		}
		// 创建title这个模版方法
		funcs := template.FuncMap{"title": strings.Title}
		{
			// 创建contract.go
			file := filepath.Join(pFolder, folder, "contract.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}
			// 使用contractTmp模版来初始化template，并且让这个模看到支持title法，即支持{{.|title}}
			t := template.Must(template.New("contract").Funcs(funcs).Parse(contractTmp))
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}
		{
			// 创建provider.go
			file := filepath.Join(pFolder, folder, "provider.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}
			t := template.Must(template.New("provider").Funcs(funcs).Parse(providerTmp))
			if err := t.Execute(f, name); err != nil {
				return nil
			}
		}

		{
			// 创建service.go
			file := filepath.Join(pFolder, folder, "service.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}
			t := template.Must(template.New("service").Funcs(funcs).Parse(serviceTmp))
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}
		fmt.Println("创建服成功，文件夹地址：", filepath.Join(pFolder, folder))
		fmt.Println("不要忘记挂载新创的服务")
		return nil
	},
}

var contractTmp string = `
package {{.}}

const {{.|title}}Key = "web:{{.}}"

type Service interface {
	// define the function
	Foo() string
}
`

var providerTmp string = `
package {{.}}

import (
	"github.com/awaketai/goweb/framework"

)

type {{.|title}}Provider struct {
	framework.ServiceProvider

	c framework.Container
}

func (sp *{{.|title}}Provider) Name() string {
	return {{.|title}}Key
}

func (sp *{{.|title}}Provider) Register(c framework.Container) framework.NewInstance {
	return New{{.|title}}Service
}

func (sp *{{.|title}}Provider) IsDefer() bool {
	return false
}

func (sp *{{.|title}}Provider) Params(c framework.Container) []any{
	return []any{c}
}

func (sp *{{.|title}}Provider) Boot(c framework.Container) error {
	return nil
}
`

var serviceTmp string = `
package {{.}}

import "github.com/awaketai/goweb/framework"

type {{.|title}}Service struct {
	container framework.Container
}

func New{{.|title}}Service(params ...any) (any,error) {
	container := params[0].(framework.Container)
	return &{{.|title}}Service{container: container},nil
}

func (s *{{.|title}}Service) Foo() string {
	return ""
}
`
