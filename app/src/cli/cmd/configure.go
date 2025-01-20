package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

func NewConfigureCmd() *cobra.Command {

	getCmd := &cobra.Command{
		Use:   pkg.ConfigureOps,
		Short: "This command will allow you to perform configure opertions based on the resources provided.",
		Long:  `This command will allow you to perform configure operations based on the resources provided.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cobra.CheckErr(fmt.Errorf("no resource provided for this command, please select resource from result of this command -> " + pkg.APPNAME + " " + pkg.ExplainOps))
			}
		},
	}
	getCmd.AddCommand(NewAccountCmd(), NewUserCmd(), NewDomainCmd(), NewAIModelNodeCmd())
	return getCmd
}
