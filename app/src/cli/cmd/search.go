package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-oss/aistudio/cli/internals/studio"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

func NewSearchmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   pkg.SearchOpts,
		Short: "This command will allow you to perform search action on the resources provided based on actions provided.",
		Long:  `This command will allow you to perform search action on the resources provided based on actions provided.`,
		Run: func(cmd *cobra.Command, args []string) {
			// if len(args) < 1 {
			// 	cobra.CheckErr(fmt.Errorf("no resource provided for this command, please select resource from result of this command -> " + pkg.APPNAME + " " + pkg.ExplainOps))
			// }
			defer vapusGlobals.VapusStudioClient.Close()
			if la {
				vapusGlobals.VapusStudioClient.ListResourceActions("organizations")
				return
			}
			spinner := pkg.GetSpinner(36)
			spinner.Start()
			defer spinner.Stop()
			vapusGlobals.VapusStudioClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(currentAccessToken),
				Action:      action,
				File:        file,
				SearchQ:     searchQuery,
				Resource:    pkg.SearchOpts,
			}
			err := vapusGlobals.VapusStudioClient.HandleSearch()
			if err != nil {
				cobra.CheckErr(err)
			}
		},
	}
	return cmd
}
