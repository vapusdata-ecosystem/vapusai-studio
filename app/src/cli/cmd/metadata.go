package cmd

// import (
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
// 	dpb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
// 	plclient "github.com/vapusdata-oss/aistudio/cli/internals/studio"
// 	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
// )

// func NewMetadataCmd() *cobra.Command {
// 	var dsId string
// 	cmd := &cobra.Command{
// 		Use:   pkg.MetaDataResource,
// 		Short: "This command is interface to interact with the platform for metadata resources.",
// 		Long:  `This command is interface to interact with the platform for metadata resources.`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			defer vapusGlobals.VapusStudioClient.Close()
// 			if la {
// 				vapusGlobals.VapusStudioClient.ListResourceActions("metadata")
// 				return
// 			}
// 			resAct := getMetadataAction(cmd.Parent().Use, action)
// 			spinner := pkg.GetSpinner(36)
// 			spinner.Prefix = "Agent is running"
// 			spinner.Start()
// 			vapusGlobals.VapusStudioClient.ActionHandler = plclient.ActionHandlerOpts{
// 				ParentCmd:   cmd.Parent().Use,
// 				Args:        args,
// 				AccessToken: viper.GetString(currentAccessToken),
// 				Action:      resAct,
// 				File:        file,
// 				Params:      map[string]string{pkg.DatasourceKey: dsId},
// 			}

// 			err := vapusGlobals.VapusStudioClient.HandleAction()
// 			if err != nil {
// 				spinner.Stop()
// 				cobra.CheckErr(err)
// 			}

// 			spinner.Stop()

// 		},
// 	}
// 	cmd.PersistentFlags().StringVar(&dsId, "datasource", "", "Data product Id to perform the action ons")
// 	return cmd
// }

// func getMetadataAction(parentCmd string, action string) string {
// 	switch parentCmd {
// 	case pkg.GetOps:
// 		return dpb.DataSourceAgentActions_LIST_DATASOURCE.String()
// 	case pkg.DescribeOps:
// 		return dpb.DataSourceAgentActions_DESCRIBE_DATASOURCE.String()
// 	case pkg.ActOps:
// 		return action
// 	default:
// 		return pkg.ErrInvalidAction.Error()
// 	}
// }

// // func getDataSource() {
// // 	err := vapusGlobals.VapusStudioClient.ListActions(pb.DataSourceAgentActions_LIST_DATASOURCE.String(), viper.GetString(currentAccessToken))
// // 	if err != nil {
// // 		cobra.CheckErr(err)
// // 	}
// // }

// // func describeDataSource(args []string) {
// // 	if len(args) < 1 {
// // 		cobra.CheckErr("Invalid number of arguments, please provide the dataSource ID")
// // 	}
// // 	err := vapusGlobals.VapusStudioClient.DescribeActions(pb.DataSourceAgentActions_DESCRIBE_DATASOURCE.String(), viper.GetString(currentAccessToken), args[0])
// // 	if err != nil {
// // 		cobra.CheckErr(err)
// // 	}
// // }
