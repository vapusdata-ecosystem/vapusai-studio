package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dmCtxCmd represents the dmCtx command
func NewContextCurrentCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   currentContextCmd,
		Short: "This command will use the context of the current VapusData platform instance.",
		Long:  `This command will use the context of the current VapusData platform instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			getCurrentContext()
		},
	}

	return cmd
}

func getCurrentContext() {
	currentContextParams = viper.GetStringMapString(currentContextParamsKey)
	vapusGlobals.logger.Info().Msgf("Current context - %v", currentContextParams["name"])
}
