package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-oss/aistudio/cli/internals/studio"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

func NewPromptCmd() *cobra.Command {
	var dataproduct, query string
	cmd := &cobra.Command{
		Use:   pkg.GetPrompt,
		Short: "This command will provide a prompt to query the vapusdata products.",
		Long:  `This command will provide a prompt to query the vapusdata products.`,
		Run: func(cmd *cobra.Command, args []string) {
			defer vapusGlobals.VapusStudioClient.Close()
			// spinner := pkg.GetSpinner(36)
			// spinner.Prefix = fmt.Sprintf("Querying Data Product %v with the given query", dataproduct)
			// spinner.Start()
			vapusGlobals.VapusStudioClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Use,
				Args:        args,
				AccessToken: viper.GetString(currentProductAccessToken),
				Action:      pkg.GetPrompt,
				Params:      map[string]string{pkg.DataproductKey: dataproduct, pkg.SearchqueryKey: query},
			}
			err := vapusGlobals.VapusStudioClient.HandleAction()
			if err != nil {
				// spinner.Stop()
				cobra.CheckErr(err)
			}
			// spinner.Stop()
		},
	}
	cmd.PersistentFlags().StringVar(&dataproduct, "dataproduct", "", "Data product to query")
	cmd.PersistentFlags().StringVar(&query, "query", "", "Query to search for the data product")
	cmd.AddCommand()
	return cmd
}
