package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	plclient "github.com/vapusdata-oss/aistudio/cli/internals/studio"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
	dpb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func NewDomainCmd() *cobra.Command {
	var dmId string
	cmd := &cobra.Command{
		Use:   pkg.DomainResource,
		Short: "This command is interface to interact with the platform for organization resources.",
		Long:  `This command is interface to interact with the platform for organization resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			defer vapusGlobals.VapusStudioClient.Close()
			if la {
				vapusGlobals.VapusStudioClient.ListResourceActions("organizations")
				return
			}
			spinner := pkg.GetSpinner(39)
			spinner.Start()
			defer spinner.Stop()
			vapusGlobals.VapusStudioClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(currentAccessToken),
				Action:      getDomainAction(cmd.Parent().Use, action),
				File:        file,
				Params:      map[string]string{pkg.DomainKey: dmId},
				Resource:    pkg.DomainResource,
			}
			err := vapusGlobals.VapusStudioClient.HandleAction()
			if err != nil {
				cobra.CheckErr(err)
			}
		},
	}
	cmd.PersistentFlags().StringVar(&dmId, "organization", "", "Data product Id to perform the action on")
	return cmd
}

func getDomainAction(parentCmd string, action string) string {
	switch parentCmd {
	case pkg.GetOps:
		return dpb.Organization.String()
	case pkg.DescribeOps:
		return dpb.DomainAgentActions_DESCRIBE_DOMAIN.String()
	case pkg.ActOps:
		return action
	default:
		return pkg.ErrInvalidAction.Error()
	}
}

// func getDomainAction(parentCmd string, action string) string {
// 	switch parentCmd {
// 	case pkg.GetOps:
// 		return dpb.DomainAgentActions_LIST_DOMAINS.String()
// 	case pkg.DescribeOps:
// 		return dpb.DomainAgentActions_LIST_DOMAINS.String()
// 	case pkg.ActOps:
// 		return action
// 	case pkg.ConfigureOps:

// 	default:
// 		return pkg.ErrInvalidAction.Error()
// 	}
// }

// func (x DomainHandler) getDomain() {
// 	err := vapusGlobals.VapusStudioClient.ListActions(dpb.DomainAgentActions_LIST_DOMAINS.String(), x.accessToken)
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func (x DomainHandler) describeDomain() {
// 	if len(x.args) < 1 {
// 		cobra.CheckErr("Invalid number of arguments, please provide the organization ID")
// 	}
// 	err := vapusGlobals.VapusStudioClient.DescribeActions(dpb.DomainAgentActions_LIST_DOMAINS.String(), x.accessToken, x.args[0])
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func (x DomainHandler) act() {
// 	if x.action == "" {
// 		cobra.CheckErr("No action provided")
// 	}
// 	if x.file == "" {
// 		cobra.CheckErr("No input provided")
// 	}
// 	err := vapusGlobals.VapusStudioClient.PerformAct(x.action, x.accessToken, x.file)
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
