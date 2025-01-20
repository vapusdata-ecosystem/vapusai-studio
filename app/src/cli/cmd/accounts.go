package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-oss/aistudio/cli/internals/studio"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

func NewAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pkg.AccountResource,
		Short: "This command is interface to interact with the platform for account resources.",
		Long:  `This command is interface to interact with the platform for account resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			if la {
				vapusGlobals.VapusStudioClient.ListResourceActions("account")
				return
			}
			resAct := getAccountAction(cmd.Parent().Use, action)
			// spinner := pkg.GetSpinner(36)
			// spinner.Start()
			vapusGlobals.VapusStudioClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(currentAccessToken),
				Action:      resAct,
				File:        file,
				Resource:    pkg.AccountResource,
			}
			log.Println("Action: ", resAct)
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

func getAccountAction(parentCmd string, action string) string {
	switch parentCmd {
	case pkg.ActOps:
		return action
	default:
		return ""
	}
}

// func accountActions(parentCmd string) {
// 	switch parentCmd {
// 	case pkg.GetOps:
// 		getAccount()
// 	case pkg.DescribeOps:
// 		describeAccount()
// 	default:
// 		cobra.CheckErr("Invalid action")
// 	}
// }

// func getAccount() {
// 	err := vapusGlobals.VapusStudioClient.ListActions(pb.AccountAgentActions_LIST_ACCOUNT.String(), viper.GetString(currentAccessToken))
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func describeAccount() {
// 	err := vapusGlobals.VapusStudioClient.DescribeActions(pb.AccountAgentActions_LIST_ACCOUNT.String(), viper.GetString(currentAccessToken), "")
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
