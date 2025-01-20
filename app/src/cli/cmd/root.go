package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"

	plclient "github.com/vapusdata-oss/aistudio/cli/internals/studio"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

var (
	defaultCfgFolder   = ".vapusdata"
	defaultCfgFileType = "toml"
	defaultCfgFileName = "config"
	version            = "0.0.1"
)

func NewRootCmd() *cobra.Command {
	rootCmd = &cobra.Command{
		Use:     pkg.APPNAME,
		Version: version,
		Short:   "vapusctl is a cli tool tht provides an interface to interact with different services offered by VapusData Platform",
		Long:    `This cli tool will allow you to perform different operation on VapusData platform under the current context. `,
		Run: func(cmd *cobra.Command, args []string) {
			vapusGlobals.logger.Info().Msg("Welcome to VapusData CLI")
			// cmd.Help()
		},
	}
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		setGlobalPersist()
		_, exist := ignoreConnMap[cmd.Name()]
		if !exist {
			initCurrentContextInstance()
		}
		return nil
	}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vapusdata/config.toml)")
	rootCmd.PersistentFlags().BoolVar(&debugLogFlag, "debug", true, "Enable debug/verbose logging mode")
	rootCmd.PersistentFlags().StringVarP(&action, "action", "a", "", "Action for the platform that should be executed on current resource with params in a file")
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "File containing the parameters for the action")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "File containing the output of the action")
	rootCmd.PersistentFlags().StringVarP(&searchQuery, "query", "q", "", "Search query for the resource")
	rootCmd.PersistentFlags().BoolVarP(&la, "listactions", "l", false, "List down all the actions that can be performed on the current resource")

	_ = viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	_ = viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	_ = viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("debug", debugLogFlag)
	rootCmd.AddCommand(NewClearCmd(), NewConfigCmd(), NewExplainOps(), NewAuthCmd(), NewGetCmd(), NewConnectCmd(),
		NewDescribeCmd(), NewRequestSpecCmd(), NewActCmd(), NewSearchmd(), NewPromptCmd(), NewConfigureCmd(), NewOperatorCmd())
	return rootCmd
}

func defaultConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return homeDir + "/" + defaultCfgFolder

}

func initConfig() {
	logger.Debug().Msgf("Initiating config... at %v", cfgFile)
	// Read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory to get/create config file in there.
		cfgFile = filepath.Join(defaultConfigDir(), defaultCfgFileName+"."+defaultCfgFileType)

		// Search config in home directory with name ".vapusdata" (without extension).
		viper.AddConfigPath(defaultConfigDir())
		viper.SetConfigName(defaultCfgFileName)
		viper.SetConfigType(defaultCfgFileType)

		// Write the file only and gracefully handles if file already exists
		cfgErr := viper.SafeWriteConfig()
		if cfgErr != nil {
			existCfgError := viper.ConfigFileAlreadyExistsError("")
			if !errors.As(cfgErr, &existCfgError) {
				cobra.CheckErr(cfgErr)
			}
		}
		cobra.CheckErr(viper.ReadInConfig())
		viper.SetConfigFile(cfgFile)
	}
	// Viper read the config file on desired path
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func setGlobalPersist() {
	pkg.InitLogger(debugLogFlag)
	vapusGlobals = &GlobalsPersists{
		logger:        pkg.DmLogger,
		cfgFile:       cfgFile,
		debugLogFlag:  debugLogFlag,
		cfgDir:        filepath.Dir(cfgFile),
		ctx:           context.Background(),
		AgentsActions: plclient.AgentGoals,
	}
}

func initCurrentContextInstance() {
	// Initialize the platform client
	currentContext := viper.GetString(currentContextKey)
	if currentContext == "" {
		return
	}
	currentContextParams := viper.GetStringMapString(currentContextParamsKey)
	if currentContextParams == nil {
		vapusGlobals.logger.Error().Msg("No context found, please add a context")
	} else {
		if vapusGlobals.VapusStudioClient == nil {
			client, err := plclient.NewPlatFormClient(currentContextParams, vapusGlobals.logger)
			if err != nil {
				vapusGlobals.logger.Error().Msgf("Error connecting to current context platform instance: %v", err)
				cobra.CheckErr(pkg.ErrVapusDataPlatformNotConnected)
				return
			}
			vapusGlobals.VapusStudioClient = client
			vapusGlobals.CurrentContext = currentContextParams["name"]
			// vapusGlobals.ResourceActionMap = vapusGlobals.AgentsActions
		}
	}
}
