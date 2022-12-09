package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/awaketai/goweb/framework/cobra"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/util"
	"github.com/swaggo/swag/gen"
)

func initSwaggerCommand() *cobra.Command {
	swaggerCommand.AddCommand(swaggerGenCommand)
	return swaggerCommand
}

var swaggerCommand = &cobra.Command{
	Use:   "swagger",
	Short: "swagger command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var swaggerGenCommand = &cobra.Command{
	Use:   "gen",
	Short: "生成对应的swagger文件,包括：swagger.yaml,doc.go",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		fmt.Println("app:", appService.AppFolder())
		outputDir := filepath.Join(appService.AppFolder(), "http", "swagger")
		fmt.Println("file:", outputDir)
		if exists, _ := util.Exists(outputDir); !exists {
			err := os.MkdirAll(outputDir, 0644)
			if err != nil {
				return err
			}
		}
		// 生成swagger.json
		httpFolder := filepath.Join(appService.AppFolder(), "http")
		conf := &gen.Config{
			// 遍历需要查询注释的目录
			SearchDir: httpFolder,
			// 不包含哪些文件
			Excludes: "",
			// 生成文件的输出目录
			OutputDir: outputDir,
			// 整个swagger接口的说明文档注释
			MainAPIFile: "swagger.go",
			// 名字的显示策略，比如首字母大写等
			PropNamingStrategy: "",
			// 是否需要解析vendor目录
			ParseVendor: false,
			// 是否需要解析外部依赖库的包
			ParseDependency: false,
			// 是否要查找markdown文件，这个markdown文件能用来为tag增加说明格式
			MarkdownFilesDir: httpFolder,
			// 是否应该在docs.go中生成时间戳
			GeneratedTime: false,
		}
		err := gen.New().Build(conf)
		if err != nil {
			return fmt.Errorf("swag gen err:%v", err)
		}
		return nil
	},
}
