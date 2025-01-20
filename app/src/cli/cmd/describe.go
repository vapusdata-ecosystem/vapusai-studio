package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

func NewDescribeCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   pkg.DescribeOps,
		Short: "This command will allow you to perform describe actions based on the resources provided.",
		Long:  `This command will allow you to perform describe actions based on the resources provided.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cobra.CheckErr(fmt.Errorf("no resource provided for this command, please select resource from result of this command -> " + pkg.APPNAME + " " + pkg.ExplainOps))
			}
		},
	}
	cmd.AddCommand(NewUserCmd(), NewDomainCmd(), NewAIModelNodeCmd())
	return cmd
}
