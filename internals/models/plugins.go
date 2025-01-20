package models

import (
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type Plugin struct {
	VapusBase     `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	PluginType    string               `json:"pluginType,omitempty" yaml:"pluginType"`
	PluginService string               `json:"pluginService,omitempty" yaml:"pluginService"`
	Name          string               `bun:"unique,no_null" json:"name,omitempty" yaml:"name"`
	NetworkParams *PluginNetworkParams `bun:"type:jsonb" json:"networkParams,omitempty" yaml:"networkParams"`
	DynamicParams []*Mapper            `bun:"type:jsonb[]" json:"dynamicParams,omitempty" yaml:"dynamicParams"`
	Editable      bool                 `json:"editable,omitempty" yaml:"editable" default:"true"`
	Status        string               `json:"status,omitempty" yaml:"status"`
}

func (dm *Plugin) ConvertToPb() *mpb.Plugin {
	if dm == nil {
		return nil
	}
	obj := &mpb.Plugin{
		PluginType:    mpb.IntegrationPluginTypes(mpb.IntegrationPluginTypes_value[dm.PluginType]),
		PluginService: mpb.IntegrationPlugins(mpb.IntegrationPlugins_value[dm.PluginService]),
		Name:          dm.Name,
		NetworkParams: dm.NetworkParams.ConvertToPb(),
		DynamicParams: MapperSliceToPb(dm.DynamicParams),
		Org:           dm.Organization,
		Editable:      dm.Editable,
		Status:        dm.Status,
		PluginId:      dm.VapusID,
		ResourceBase:  dm.ConvertToPbBase(),
	}
	obj.Scope = dm.Scope
	return obj
}

func (dm *Plugin) ConvertFromPb(pb *mpb.Plugin) *Plugin {
	if pb == nil {
		return nil
	}
	dm.PluginType = pb.PluginType.String()
	dm.Name = pb.Name
	dm.NetworkParams = (&PluginNetworkParams{}).ConvertFromPb(pb.NetworkParams)
	dm.DynamicParams = MapperSliceFromPb(pb.DynamicParams)
	dm.Scope = pb.Scope
	dm.Editable = pb.Editable
	dm.Status = pb.Status
	dm.PluginService = pb.PluginService.String()
	return dm
}

func (dm *Plugin) PreSaveCreate(authzClaim map[string]string) {
	if dm == nil {
		return
	}
	dm.PreSaveVapusBase(authzClaim)
	dm.Status = mpb.CommonStatus_ACTIVE.String()
}

func (dn *Plugin) PreSaveUpdate(userId string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = userId
	dn.UpdatedAt = utils.GetEpochTime()
}

type PluginNetworkParams struct {
	URL                 string                  `json:"url,omitempty" yaml:"url"`
	Port                string                  `json:"port,omitempty" yaml:"port"`
	Version             string                  `json:"version,omitempty" yaml:"version"`
	Credentials         *GenericCredentialModel `json:"credentials,omitempty" yaml:"credentials"`
	Name                string                  `json:"name,omitempty" yaml:"name"`
	SecretName          string                  `json:"secretName,omitempty" yaml:"secretName"`
	IsAlreadyInSecretBS bool                    `json:"isAlreadyInSecretBS,omitempty" yaml:"isAlreadyInSecretBS"`
}

func (dm *PluginNetworkParams) ConvertToPb() *mpb.PluginNetworkParams {
	if dm == nil {
		return nil
	}
	return &mpb.PluginNetworkParams{
		Url:                 dm.URL,
		Port:                dm.Port,
		Version:             dm.Version,
		Credentials:         dm.Credentials.ConvertToPb(),
		Name:                dm.Name,
		SecretName:          dm.SecretName,
		IsAlreadyInSecretBs: dm.IsAlreadyInSecretBS,
	}
}

func (dm *PluginNetworkParams) ConvertFromPb(pb *mpb.PluginNetworkParams) *PluginNetworkParams {
	if pb == nil {
		return nil
	}
	dm.URL = pb.Url
	dm.Port = pb.Port
	dm.Version = pb.Version
	dm.Credentials = (&GenericCredentialModel{}).ConvertFromPb(pb.Credentials)
	dm.Name = pb.Name
	dm.SecretName = pb.SecretName
	dm.IsAlreadyInSecretBS = pb.IsAlreadyInSecretBs
	return dm
}
