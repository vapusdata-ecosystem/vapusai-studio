package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-oss/aistudio/cli/internals/studio"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

func NewAIModelNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pkg.AIModelNodeResource,
		Short: "This command is interface to interact with the ai studio resources.",
		Long:  `This command is interface to interact with the ai studio resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			if la {
				vapusGlobals.VapusStudioClient.ListResourceActions(pkg.AIModelNodeResource)
				return
			}
			resAct := action
			// spinner := pkg.GetSpinner(36)
			// spinner.Start()
			vapusGlobals.VapusStudioClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(currentAccessToken),
				Action:      resAct,
				File:        file,
				Resource:    pkg.AIModelNodeResource,
			}
			err := vapusGlobals.VapusStudioClient.HandleAction()
			// spinner.Stop()
			if err != nil {
				cobra.CheckErr(err)
			}

			defer vapusGlobals.VapusStudioClient.Close()

		},
	}
	return cmd
}
