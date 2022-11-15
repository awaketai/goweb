package command

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/awaketai/goweb/framework/cobra"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/util"
	"github.com/jianfengye/collection"
)

// initCmdCommand
func initCmdCommand() *cobra.Command {
	cmdCommand.AddCommand(cmdListCommand)
	cmdCommand.AddCommand(cmdCreateCommand)
	return cmdCommand
}

// cmdCommand
var cmdCommand = &cobra.Command{
	Use:   "command",
	Short: "控制台命令相关",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

// cmdListCommand 列出所有控制台命令
var cmdListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有控制台命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmds := cmd.Root().Commands()
		ps := [][]string{}
		for _, cmd := range cmds {
			line := []string{cmd.Name(), cmd.Short}
			ps = append(ps, line)
		}
		util.PrettyPrint(ps)
		return nil
	},
}

var cmdCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个控制台命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		fmt.Println("开始创建控制台命令")
		var (
			name   string
			folder string
		)

		{
			prompt := &survey.Input{
				Message: "请输入控制台命令名称",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入文件夹名称(默认同控制台命令)",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}
		if folder == "" {
			folder = name
		}

		// check if file exists
		app := container.MustMake(contract.AppKey).(contract.App)
		pFolder := app.CommandFolder()
		subFolders, err := util.SubDir(pFolder)
		if err != nil {
			return err
		}
		subColl := collection.NewStrCollection(subFolders)
		if subColl.Contains(folder) {
			fmt.Println("目录名称已经存在")
			return nil
		}
		// create file
		if err := os.Mkdir(filepath.Join(pFolder, folder), 0700); err != nil {
			return err
		}
		// 创建title这个模板方法
		funcs := template.FuncMap{"title": strings.Title}
		{
			// 创建name.go
			file := filepath.Join(pFolder, folder, name+".go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}
			t := template.Must(template.New("cmd").Funcs(funcs).Parse(cmdTmpl))
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}

		fmt.Println("创建命令行工具成功，路径:", filepath.Join(pFolder, folder))
		fmt.Println("请记得开发完成后将命令行工具持载到 console/kernel.go")
		return nil
	},
}

var cmdTmpl string = `
package {{.}}

import (
	"fmt"
	"github.com/awaketai/goweb/framework/cobra"

)

var {{.|title}}Command = &cobra.Command{
	Use: "{{.}}",
	Short: "{{.}}",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("{{.}}")
		return nil
	},
}
`
