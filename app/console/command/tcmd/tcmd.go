
package tcmd

import (
	"fmt"
	"github.com/awaketai/goweb/framework/cobra"

)

var TcmdCommand = &cobra.Command{
	Use: "tcmd",
	Short: "tcmd",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("tcmd")
		return nil
	},
}
