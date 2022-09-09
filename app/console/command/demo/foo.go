package demo

import (
	"log"

	"github.com/awaketai/goweb/framework/cobra"
)

func InitFoo() *cobra.Command {
	FooCommand.AddCommand(Foo1Command)
	return FooCommand
}

var FooCommand = &cobra.Command{
	Use:     "foo",
	Short:   "foo short",
	Long:    "foo long",
	Aliases: []string{"fo", "f"},
	Example: "foo命令的例子",
	RunE: func(cmd *cobra.Command, args []string) error {
		// container := cmd.GetContainer()
		// log.Println(container)
		log.Println("execute foo command")
		return nil
	},
}

var Foo1Command = &cobra.Command{
	Use:     "foo1",
	Short:   "foo1 short description",
	Long:    "foo1 long description",
	Aliases: []string{"fo1", "f1"},
	Example: "foo1命令的例子",
	Run: func(cmd *cobra.Command, args []string) {
		container := cmd.GetContainer()
		log.Println(container)

	},
}
