package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

func NewActCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   pkg.ActOps,
		Short: "This command will allow you to perform actions the resources provided based on actions provided.",
		Long:  `This command will allow you to perform actions the resources provided based on actions provided.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cobra.CheckErr(fmt.Errorf("no resource provided for this command, please select resource from result of this command -> " + pkg.APPNAME + " " + pkg.ExplainOps))
			}
		},
	}
	cmd.AddCommand(NewAccountCmd(), NewUserCmd(), NewDomainCmd(), NewAIModelNodeCmd())
	return cmd
}
