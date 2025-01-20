package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	setup "github.com/vapusdata-oss/aistudio/cli/internals/setup-config"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
	"github.com/vapusdata-oss/aistudio/core/authn"
	encrytion "github.com/vapusdata-oss/aistudio/core/encryption"
	models "github.com/vapusdata-oss/aistudio/core/models"
	dmutils "github.com/vapusdata-oss/aistudio/core/utils"
)

func NewInstallerSpecGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pkg.InstallerSecretSpecGenOps,
		Short: "This command is interface to generate installer spec files for the platform installation.",
		Long:  `This command is interface to generate installer spec files for the platform installation.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := generateInstallerSpecTemplate()
			if err != nil {
				cobra.CheckErr(err)
			}
		},
	}
	return cmd
}

func generateInstallerSpecTemplate() error {
	secM := &models.DataSourceCredsParams{
		DataSourceCreds: &models.DataSourceSecrets{
			GenericCredentialModel: &models.GenericCredentialModel{
				AwsCreds:   &models.AWSCreds{},
				GcpCreds:   &models.GCPCreds{},
				AzureCreds: &models.AzureCreds{},
			},
			DB:   "",
			URL:  "",
			Port: 0,
		},
	}
	specVal, err := dmutils.GenericMarshaler(&setup.VapusSecretInstallerConfig{
		SecretStore:       secM,
		BackendDataStore:  secM,
		ArtifactStore:     secM,
		BackendCacheStore: secM,
		AuthnSecrets: &authn.AuthnSecrets{
			OIDCSecrets: &authn.OIDCSecrets{},
		},
		JWTAuthnSecrets: &encrytion.JWTAuthn{},
	}, "YAML")
	if err != nil {
		return err
	}
	log.Println("Generated installer spec: ", string(specVal))
	fileName := strings.ToLower("vapusdata-secrets.yaml")
	vapusGlobals.logger.Info().Msgf("Sample installer %v spec generated with file name - %v \n", specName, fileName)
	err = os.WriteFile(fileName, specVal, 0644)
	if err != nil {
		return err
	}
	return nil
}
