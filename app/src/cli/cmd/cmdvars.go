package cmd

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	plclient "github.com/vapusdata-oss/aistudio/cli/internals/studio"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

type StudioInstanceClients struct {
	*plclient.VapusStudioClient
}

var (

	// RootCmd is the root command for vapusdt
	rootCmd                                                       *cobra.Command
	cfgFile                                                       string
	debugLogFlag                                                  bool
	vapusGlobals                                                  *GlobalsPersists
	logger                                                        zerolog.Logger
	currentIdToken, currentAccessToken, currentProductAccessToken string = "currentIdToken", "currentAccessToken", "currentProductAccessToken"
	GlobalVar                                                     string
	action, file, searchQuery, outputFile                         string
	la                                                            bool
)

var ignoreConnMap = map[string]bool{
	pkg.ConfigResource: true,
	pkg.ClearOps:       true,
	pkg.ExplainOps:     true,
	pkg.OperatorOps:    true,
	pkg.InstallerOps:   true,
}

type GlobalsPersists struct {
	CurrentContext    string
	logger            zerolog.Logger
	cfgFile           string
	cfgDir            string
	debugLogFlag      bool
	AgentsActions     map[string][]interface{}
	AgentsUtilities   map[string][]interface{}
	AgentsReflexes    map[string][]interface{}
	AgentInterfaceMap map[string]string

	*plclient.VapusStudioClient
	currentIdToken, currentAccessToken string
	ctx                                context.Context
}
