package plugins

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vapusdata-oss/aistudio/core/emailer"
	filemanager "github.com/vapusdata-oss/aistudio/core/filemanager"
	"github.com/vapusdata-oss/aistudio/core/models"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type PluginsBase struct {
	Emailer     emailer.Emailer
	FileManager filemanager.FileManager
}

type VapusPlugins struct {
	StudioPlugins       *PluginsBase
	OrganizationPlugins map[string]*PluginsBase
	UserPlugins         map[string]*PluginsBase
	logger              zerolog.Logger
	PluginsNotConencted []string
}

var PluginTypeScopeMap = map[string][]string{
	mpb.IntegrationPluginTypes_EMAIL.String(): {
		mpb.ResourceScope_ACCOUNT_SCOPE.String(),
	},
	mpb.IntegrationPluginTypes_FILE_STORE.String(): {
		mpb.ResourceScope_ACCOUNT_SCOPE.String(),
		mpb.ResourceScope_USER_SCOPE.String(),
	},
}

func NewVapusPlugins(ctx context.Context, pluginList []*models.Plugin, selectedPlugins []string, logger zerolog.Logger) *VapusPlugins {
	obj := &VapusPlugins{
		StudioPlugins:       &PluginsBase{},
		OrganizationPlugins: make(map[string]*PluginsBase),
		logger:              logger,
	}
	logger.Info().Msg("Creating new vapus plugins pool")
	for _, plugin := range pluginList {
		// if len(selectedPlugins) > 0 && !slices.Contains(selectedPlugins, plugin.Name) {
		switch plugin.PluginType {
		case mpb.IntegrationPluginTypes_EMAIL.String():
			pl, err := emailer.New(ctx, plugin.PluginService, plugin.NetworkParams, plugin.DynamicParams, logger)
			if err != nil {
				logger.Err(err).Msg("Error creating email plugin")
				obj.PluginsNotConencted = append(obj.PluginsNotConencted, plugin.Name)
				continue
			}
			obj.StudioPlugins.Emailer = pl
			// case mpb.IntegrationPluginTypes_FILE_STORE.String():
			// pl, err := filemanager.NewFileManagerClient()
		}
		// }
	}
	return obj
}
