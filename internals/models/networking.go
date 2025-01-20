package models

import mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"

type StoreParams struct {
	Creds                 *NetworkParams         `validate:"required" json:"dataSourceCreds" yaml:"dataSourceCreds" toml:"dataSourceCreds"`
	DataSourceEngine      string                 `validate:"required" json:"dataSourceEngine" yaml:"dataSourceEngine" toml:"dataSourceEngine"`
	DataSourceType        string                 `validate:"required" json:"dataSourceType" yaml:"dataSourceType" toml:"dataSourceType"`
	DataSourceService     string                 `json:"dataSourceService" yaml:"dataSourceService" toml:"dataSourceService"`
	Params                map[string]interface{} `json:"params" yaml:"params" toml:"params"`
	DataSourceSvcProvider string                 `json:"dataSourceSvcProvider" yaml:"dataSourceSvcProvider" toml:"dataSourceSvcProvider"`
	DataStorelist         []string               `json:"dataStorelist" yaml:"dataStorelist" toml:"dataStorelist"`
	DataStorePrefixList   []string               `json:"dataStorePrefixList" yaml:"dataStorePrefixList" toml:"dataStorePrefixList"`
	Version               string                 `json:"version" yaml:"version" toml:"version"`
}

func (x *StoreParams) GetSourceCredentials() (bool, *GenericCredentialModel) {
	if x == nil {
		return false, nil
	}
	if x.Creds == nil && x.Creds.DsCreds != nil && x.Creds.DsCreds.GenericCredentialModel != nil {
		return false, x.Creds.DsCreds.GenericCredentialModel
	}
	return false, nil
}

type NetworkParams struct {
	Address          string              `json:"address,omitempty" yaml:"address"`
	Port             int32               `json:"port,omitempty" yaml:"port"`
	Databases        []string            `json:"databases,omitempty" yaml:"databases"`
	DsCreds          *NetworkCredentials `json:"dsCreds,omitempty" yaml:"dsCreds"`
	DatabasePrefixes []string            `json:"databasePrefixes,omitempty" yaml:"databasePrefixes"`
	Version          string              `json:"version,omitempty" yaml:"version"`
}

func (j *NetworkParams) ConvertToPb() *mpb.NetworkParams {
	if j != nil {
		obj := &mpb.NetworkParams{
			Address:          j.Address,
			Port:             j.Port,
			Databases:        j.Databases,
			DatabasePrefixes: j.DatabasePrefixes,
			Version:          j.Version,
			DsCreds:          j.DsCreds.ConvertToPb(),
		}
		return obj
	}
	return nil
}

func (j *NetworkParams) ConvertFromPb(pb *mpb.NetworkParams) *NetworkParams {
	if pb == nil {
		return nil
	}
	return &NetworkParams{
		Address:          pb.GetAddress(),
		Port:             pb.GetPort(),
		Databases:        pb.GetDatabases(),
		DatabasePrefixes: pb.GetDatabasePrefixes(),
		Version:          pb.GetVersion(),
		DsCreds:          (&NetworkCredentials{}).ConvertFromPb(pb.DsCreds),
	}
}

type NetworkCredentials struct {
	Name                    string `json:"name,omitempty" yaml:"name"`
	IsAlreadyInSecretBs     bool   `json:"isAlreadyInSecretBS,omitempty" yaml:"isAlreadyInSecretBS"`
	*GenericCredentialModel `json:"credentials,omitempty" yaml:"credentials"`
	Priority                int32  `json:"priority,omitempty" yaml:"priority"`
	AccessScope             string `json:"accessScope,omitempty" yaml:"accessScope"`
	DB                      string `json:"db,omitempty" yaml:"db"`
	SecretName              string `json:"secretName,omitempty" yaml:"secretName"`
}

func (j *NetworkCredentials) ConvertToPb() *mpb.NetworkCredentials {
	if j != nil {
		return &mpb.NetworkCredentials{
			Name:                j.Name,
			IsAlreadyInSecretBs: j.IsAlreadyInSecretBs,
			Credentials:         j.GenericCredentialModel.ConvertToPb(),
			Priority:            j.Priority,
			Db:                  j.DB,
			SecretName:          j.SecretName,
		}
	}
	return nil
}

func (j *NetworkCredentials) ConvertFromPb(pb *mpb.NetworkCredentials) *NetworkCredentials {
	if pb == nil {
		return nil
	}
	return &NetworkCredentials{
		Name:                   pb.GetName(),
		IsAlreadyInSecretBs:    pb.GetIsAlreadyInSecretBs(),
		GenericCredentialModel: (&GenericCredentialModel{}).ConvertFromPb(pb.GetCredentials()),
		Priority:               pb.GetPriority(),
		DB:                     pb.GetDb(),
		SecretName:             pb.GetSecretName(),
	}
}

type TlsConfig struct {
	TlsType        string `json:"tlsType" yaml:"tlsType"`
	CaCertFile     string `json:"caCertFile" yaml:"caCertFile"`
	ServerKeyFile  string `json:"serverKeyFile" yaml:"serverKeyFile"`
	ServerCertFile string `json:"serverCertFile" yaml:"serverCertFile"`
}

func (j *TlsConfig) ConvertToPb() *mpb.TlsConfig {
	if j != nil {
		return &mpb.TlsConfig{
			TlsType:        mpb.TlsType(mpb.TlsType_value[j.TlsType]),
			CaCertFile:     j.CaCertFile,
			ServerKeyFile:  j.ServerKeyFile,
			ServerCertFile: j.ServerCertFile,
		}
	}
	return nil
}

func (j *TlsConfig) ConvertFromPb(pb *mpb.TlsConfig) *TlsConfig {
	if pb == nil {
		return nil
	}
	return &TlsConfig{
		TlsType:        mpb.TlsType_name[int32(pb.GetTlsType())],
		CaCertFile:     pb.GetCaCertFile(),
		ServerKeyFile:  pb.GetServerKeyFile(),
		ServerCertFile: pb.GetServerCertFile(),
	}
}
