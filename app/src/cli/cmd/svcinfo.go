package cmd

import (
	"github.com/spf13/cobra"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

func NewSvcInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pkg.SvcInfoResource,
		Short: "This command is interface to list the services information.",
		Long:  `This command is interface to list the services information.`,
		Run: func(cmd *cobra.Command, args []string) {
			// spinner := pkg.GetSpinner(36)
			// spinner.Start()
			// ctx, cancel := context.WithCancel(context.Background())
			// err := vapusGlobals.VapusStudioClient.GetSvcInfo(ctx)
			// // spinner.Stop()
			// defer vapusGlobals.VapusStudioClient.Close()
			// defer cancel()
			// if err != nil {
			// 	cobra.CheckErr(err)
			// }

		},
	}
	return cmd
}
