package svcops

import (
	"github.com/vapusdata-oss/aistudio/core/authn"
	encrytion "github.com/vapusdata-oss/aistudio/core/encryption"
	models "github.com/vapusdata-oss/aistudio/core/models"
)

type StudioInstallerConfig struct {
	App struct {
		Name         string `yaml:"name"`
		Namespace    string `yaml:"namespace"`
		Organization string `yaml:"organization"`
		Address      string `yaml:"address"`
		Dev          bool   `yaml:"dev"`
	} `yaml:"app"`
	AccountBootstrap struct {
		StudioOwners  []string `yaml:"platformOwners"`
		StudioAccount struct {
			Name    string `yaml:"name"`
			Creator string `yaml:"creator"`
		} `yaml:"platformAccount"`
		StudioAccountOrganization struct {
			Name string `yaml:"name"`
		} `yaml:"platformAccountOrganization"`
		Datamarketplace struct {
			Name    string `yaml:"name"`
			Creator string `yaml:"creator"`
		} `yaml:"datamarketplace"`
	} `yaml:"accountBootstrap"`
	Secrets    *StudioSecretsMap `yaml:"secrets"`
	Postgresql struct {
		FullnameOverride string `yaml:"fullnameOverride"`
		Auth             struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
		} `yaml:"auth"`
	} `yaml:"postgresql"`
	Vault   *Vault `yaml:"vault"`
	TLSCert struct {
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	} `yaml:"tlsCert"`
	CreateDatabase bool                  `yaml:"createDatabase"`
	SecretStore    *models.NetworkParams `yaml:"secretStore"`
	DevSecretStore *models.NetworkParams `yaml:"devSecretStore"`
}

type Vault struct {
	FullnameOverride string `yaml:"fullnameOverride"`
	Server           struct {
		Standalone struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"standalone"`
	} `yaml:"server"`
}

type StudioSecretInstallerConfig struct {
	SecretStore       *models.NetworkParams `yaml:"secretStore"`
	DevSecretStore    *models.NetworkParams `yaml:"devSecretStore"`
	BackendDataStore  *models.NetworkParams `yaml:"backendDataStore"`
	BackendCacheStore *models.NetworkParams `yaml:"backendCacheStore"`
	FileStore         *models.NetworkParams `yaml:"fileStore"`
	JWTAuthnSecrets   *encrytion.JWTAuthn   `yaml:"JWTAuthnSecrets"`
	AuthnSecrets      *authn.AuthnSecrets   `yaml:"authnSecrets"`
}

type StudioSecretsMap struct {
	BackendSecretStore struct {
		Secret string `yaml:"secret"`
	} `yaml:"backendSecretStore"`
	BackendDataStore struct {
		Secret string `yaml:"secret"`
	} `yaml:"backendDataStore"`
	BackendCacheStore struct {
		Secret string `yaml:"secret"`
	} `yaml:"backendCacheStore"`
	FileStore struct {
		Secret string `yaml:"secret"`
	} `yaml:"fileStore"`
	JWTAuthnSecrets struct {
		Secret string `yaml:"secret"`
	} `yaml:"JWTAuthnSecrets"`
	AuthnSecrets struct {
		Secret string `yaml:"secret"`
	} `yaml:"authnSecrets"`
}
