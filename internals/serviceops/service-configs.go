package svcops

import (
	filepath "path/filepath"
)

type StudioBootConfig struct {
	StudioOwners  []string `yaml:"platformOwners" validate:"required"`
	StudioAccount struct {
		Name    string `yaml:"name" validate:"required"`
		Creator string `yaml:"creator" validate:"required"`
	} `yaml:"platformAccount" validate:"required"`
	StudioAccountOrganization struct {
		Name string `yaml:"name" `
	} `yaml:"platformAccountOrganization"`
	DataMarketplace struct {
		Name    string `yaml:"name" validate:"required"`
		Creator string `yaml:"creator" validate:"required"`
	} `yaml:"datamarketplace" validate:"required"`
}

type VapusAISvcConfig struct {
	Path                 string
	VapusBESecretStorage struct {
		FilePath string `yaml:"filePath"`
		Secret   string `yaml:"secret"`
	} `yaml:"vapusBESecretStorage"`
	VapusBEDbStorage struct {
		FilePath string `yaml:"filePath"`
		Secret   string `yaml:"secret"`
	} `yaml:"vapusBEDbStorage"`
	VapusBECacheStorage struct {
		FilePath string `yaml:"filePath"`
		Secret   string `yaml:"secret"`
	} `yaml:"vapusBECacheStorage"`
	VapusFileStorage struct {
		FilePath string `yaml:"filePath"`
		Secret   string `yaml:"secret"`
	} `yaml:"vapusFileStorage"`
	NetworkConfigFile string `yaml:"networkConfigFile"`
	JWTAuthnSecrets   struct {
		FilePath string `yaml:"filePath"`
		Secret   string `yaml:"secret"`
	} `yaml:"JWTAuthnSecrets"`
	AuthnSecrets struct {
		FilePath string `yaml:"filePath"`
		Secret   string `yaml:"secret"`
	} `yaml:"authnSecrets"`
	AuthnMethod       string            `yaml:"authnMethod"`
	StudioBaseAccount *StudioBootConfig `yaml:"platformBaseAccount"`
	SelfSignup        bool              `yaml:"selfSignup"`
	ServerCerts       struct {
		Mtls           bool   `yaml:"mtls"`
		PlainTls       bool   `yaml:"plainTls"`
		Insecure       bool   `yaml:"insecure"`
		CaCertFile     string `yaml:"caCertFile"`
		ServerCertFile string `yaml:"serverCertFile"`
		ServerKeyFile  string `yaml:"serverKeyFile"`
		ClientCertFile string `yaml:"serverCertFile"`
		ClientKeyFile  string `yaml:"serverKeyFile"`
	} `yaml:"serverCerts"`
}

func (sc *VapusAISvcConfig) GetFileStorePath() string {
	return sc.VapusFileStorage.Secret
}

func (sc *VapusAISvcConfig) GetSecretStoragePath() string {
	return filepath.Join(sc.Path, sc.VapusBESecretStorage.FilePath)
}

func (sc *VapusAISvcConfig) GetDBStoragePath() string {
	return sc.VapusBEDbStorage.Secret
}

func (sc *VapusAISvcConfig) GetCachStoragePath() string {
	return sc.VapusBECacheStorage.Secret
}

func (sc *VapusAISvcConfig) GetJwtAuthSecretPath() string {
	return sc.JWTAuthnSecrets.Secret
}

func (sc *VapusAISvcConfig) GetMtlsCerts() (string, string, string) {
	return sc.ServerCerts.CaCertFile,
		sc.ServerCerts.ServerCertFile,
		sc.ServerCerts.ServerKeyFile
}

func (sc *VapusAISvcConfig) GetPlainTlsCerts() (string, string) {
	return sc.ServerCerts.ServerCertFile,
		sc.ServerCerts.ServerKeyFile
}

func (sc *VapusAISvcConfig) GetCaCert() string {
	return sc.ServerCerts.CaCertFile
}

func (sc *VapusAISvcConfig) GetClientMtlsCerts() (string, string, string) {
	return sc.ServerCerts.CaCertFile,
		sc.ServerCerts.ClientCertFile,
		sc.ServerCerts.ClientKeyFile
}

func (sc *VapusAISvcConfig) GetClientPlainTlsCerts() (string, string) {
	return sc.ServerCerts.ClientCertFile,
		sc.ServerCerts.ClientCertFile
}
