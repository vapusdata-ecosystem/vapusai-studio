package cmd

import (
	"os"

	cobra "github.com/spf13/cobra"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

func NewClearCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pkg.ClearOps,
		Version: version,
		Short:   "This command will clear the config file of vapus cli that involves the context of the different vapusdata platform instance instances.",
		Long:    `This command will clear the config file of vapus cli that involves the context of the different vapusdata platform instance instances.`,
		Run: func(cmd *cobra.Command, args []string) {
			clearConfigDir(vapusGlobals.cfgDir)
		},
	}
	return cmd
}

func clearConfigDir(cfgDir string) {
	vapusGlobals.logger.Debug().Msgf("clearing config from dir %v", cfgDir)
	files, err := os.ReadDir(cfgDir)
	if err != nil {
		cobra.CheckErr(err)
	}
	for _, file := range files {
		if file.IsDir() {
			vapusGlobals.logger.Debug().Msgf("clearing config from sub dir %v", file.Name())
			clearConfigDir(file.Name())
		}
		vapusGlobals.logger.Debug().Msgf("clearing config file %v", file.Name())
		err := os.Remove(cfgDir + "/" + file.Name())
		if err != nil {
			cobra.CheckErr(err)
		}
	}
}
