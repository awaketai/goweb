package command

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/awaketai/goweb/framework/cobra"
)

func initBuildCommand() *cobra.Command {
	buildCommand.AddCommand(buildSelfCommand)
	buildCommand.AddCommand(buildBackendCommand)
	buildCommand.AddCommand(buildFrontendCommand)
	buildCommand.AddCommand(buildAllCommand)
	return buildCommand
}

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "execute command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var buildSelfCommand = &cobra.Command{
	Use:   "self",
	Short: "compile go web",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := exec.LookPath("go")
		if err != nil {
			log.Fatalln("not found go program,install first?")
		}
		c := exec.Command(path, "build", "-o", "hade", "./")
		out, err := c.CombinedOutput()
		if err != nil {
			fmt.Println("go build err:", err)
			return err
		}
		fmt.Println("compiled success", string(out))
		return nil
	},
}

var buildBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "compile goweb use go ",
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildSelfCommand.RunE(cmd, args)
	},
}

var buildFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "use npm compile frontend",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := exec.LookPath("npm")
		if err != nil {
			log.Fatalln("not found npm program,install first?")
		}
		c := exec.Command(path, "run", "build")
		out, err := c.CombinedOutput()
		if err != nil {
			fmt.Println("npm build err:", err)
			return err
		}
		fmt.Println("compile frontend success:", string(out))
		return nil
	},
}

var buildAllCommand = &cobra.Command{
	Use:   "all",
	Short: "compile frontend and backend together",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := buildFrontendCommand.RunE(cmd, args)
		if err != nil {
			return err
		}
		err = buildBackendCommand.RunE(cmd, args)
		if err != nil {
			return err
		}

		return nil
	},
}
