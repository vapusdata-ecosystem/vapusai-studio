package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-oss/aistudio/cli/internals/studio"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func NewUserCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   pkg.UserResource,
		Short: "This command is interface to interact with the platform for users resources.",
		Long:  `This command is interface to interact with the platform for users resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			vapusGlobals.VapusStudioClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(currentAccessToken),
				Action:      userActions(cmd.Parent().Use),
				Resource:    pkg.UserResource,
			}
			err := vapusGlobals.VapusStudioClient.HandleAction()
			if err != nil {
				cobra.CheckErr(err)
			}

			defer vapusGlobals.VapusStudioClient.Close()
		},
	}
	return cmd
}

func userActions(parentCmd string) string {
	switch parentCmd {
	case pkg.GetOps:
		return pb.UserGetterActions_LIST_USERS.String()
	case pkg.DescribeOps:
		return pb.UserGetterActions_GET_USER.String()
	default:
		return ""
	}
}

// func getuser() {
// 	err := vapusGlobals.VapusStudioClient.ListActions(pb.UserAgentOperations_GET_USER.String(), viper.GetString(currentAccessToken))
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func describeUser() {

// 	err := vapusGlobals.VapusStudioClient.DescribeActions(pb.UserAgentOperations_GET_USER.String(), viper.GetString(currentAccessToken), "")
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
