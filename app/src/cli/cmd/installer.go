package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	setupconfig "github.com/vapusdata-oss/aistudio/cli/internals/setup-config"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
	k8s "github.com/vapusdata-oss/aistudio/core/k8s"
	dmutils "github.com/vapusdata-oss/aistudio/core/utils"
)

var devMode, upgradeVapus bool
var configFile, vapusversion string
var err error

func NewInstallerCmd() *cobra.Command {
	installerCmd := &cobra.Command{
		Use:   pkg.InstallerOps,
		Short: "This command will allow you to install vapusdata application on kubernetes cluster.",
		Long:  `This command will allow you to install vapusdata application on kubernetes cluster.`,
		Run: func(cmd *cobra.Command, args []string) {
			installer := NewInstaller(configFile)
			err := installer.PreRunChecks()
			if err != nil {
				cobra.CheckErr(err)
			}
			err = installer.Run()
			if err != nil {
				cobra.CheckErr(err)
			}

		},
	}
	installerCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file containing the configuration of the vapusdata platform")
	installerCmd.PersistentFlags().StringVar(&vapusversion, "version", "", "Version of the vapusdata platform")
	installerCmd.PersistentFlags().BoolVar(&upgradeVapus, "upgrade", false, "Upgrade the vapusdata platform")
	return installerCmd
}

type VapusdataInstaller struct {
	CustomValues     string
	ConfigSpec       *setupconfig.VapusInstallerConfig
	namespace        string
	installationName string
	version          string
	packageUrl       string
}

func NewInstaller(f string) *VapusdataInstaller {
	fBytes, err := os.ReadFile(f)
	if err != nil {
		cobra.CheckErr(err)
	}
	config := &setupconfig.VapusInstallerConfig{}
	err = dmutils.GenericUnMarshaler(fBytes, config, dmutils.GetConfFileType(f))
	if err != nil {
		cobra.CheckErr(err)
	}
	return &VapusdataInstaller{CustomValues: f, ConfigSpec: config}
}

func (x *VapusdataInstaller) PreRunChecks() error {
	var err error
	vapusGlobals.logger.Info().Msg("Running pre-installation checks")
	vapusGlobals.logger.Info().Msg("Checking if kubectl is installed")
	err = x.checkKubectlInstalled()
	if err != nil {
		return err
	}
	vapusGlobals.logger.Info().Msg("Checking if helm is installed")
	err = x.checkHelmInstalled()
	if err != nil {
		return err
	}
	vapusGlobals.logger.Info().Msg("Pre-installation checks completed successfully")
	return nil
}

func (x *VapusdataInstaller) writeFinalSpec() error {
	fBytes, err := dmutils.GenericMarshaler(x.ConfigSpec, dmutils.GetConfFileType(x.CustomValues))
	if err != nil {
		return err
	}

	err = os.WriteFile(x.CustomValues, fBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (x *VapusdataInstaller) Run() error {
	kPath, err := pkg.KubeConfigPath.Run()
	if err != nil {
		return err
	}
	_, k8sConfig, err := k8s.GetLocalK8SClientSet(kPath)
	if err != nil {
		return err
	}
	vapusGlobals.logger.Info().Msgf("Connected to k8s cluster: %s", k8sConfig.Host)
	// Install vapusdata with below Steps
	err = x.setupEnv()
	if err != nil {
		return err
	}
	// Install Vault
	if x.ConfigSpec.App.Dev {
		vapusGlobals.logger.Info().Msg("Vapusdata running in development mode")
		err = x.installHashiCorpVault(x.namespace)
		if err != nil {
			return err
		}
	}
	err = x.writeFinalSpec()
	if err != nil {
		return err
	}
	// var lintoutput []byte
	// logger.Info().Msg("Linting the chart")
	// lintoutput, err = exec.Command("helm", "lint", chartPackagePath).CombinedOutput()
	// if err != nil {
	// 	logger.Info().Msgf("Helmer lint output: %v", string(lintoutput))
	// 	logger.Fatal().Msgf("Helmer lint failed: %v", err)
	// }
	command := "install"
	if upgradeVapus {
		command = "upgrade"
	}
	vapusGlobals.logger.Info().Msgf("Installing vapusdata platform with version %s using helm chart %s", x.version, x.packageUrl)
	cmd := exec.Command("helm",
		command,
		x.installationName,
		"--create-namespace",
		pkg.HemlInstallerPackage,
		"--namespace", x.namespace,
		"--version", x.version,
		"--set", "app.namespace="+x.namespace,
		"--set", "app.dev="+dmutils.BoolToString(x.ConfigSpec.App.Dev),
		"-f", x.CustomValues)
	log.Println(strings.Join(cmd.Args, " "))
	log.Println(cmd.String())
	if err := cmd.Run(); err != nil {
		vapusGlobals.logger.Err(err).Msg("Error while installing vapusdata")
		return err
	}
	log.Println("Vapusdata operation completed successfully")
	return nil
}

func (x *VapusdataInstaller) checkHelmInstalled() error {
	cmd := exec.Command("helm", "version")
	if err := cmd.Run(); err != nil {
		return err
	}
	vapusGlobals.logger.Info().Msg("Helm is installed")
	return nil
}

func (x *VapusdataInstaller) checkKubectlInstalled() error {
	cmd := exec.Command("kubectl", "version")
	if err := cmd.Run(); err != nil {
		return err
	}
	vapusGlobals.logger.Info().Msg("Kubectl is installed")
	return nil
}

func (x *VapusdataInstaller) setupEnv() error {
	var err error
	x.namespace, err = pkg.Namespace.Run()
	if err != nil {
		return err
	}
	vapusGlobals.logger.Info().Msgf("Vapusdata installation started in namespace : %s", x.namespace)
	devModeStr, err := pkg.VapusDevMode.Run()
	if err != nil {
		return err
	}
	x.installationName, err = pkg.VapusInstallationName.Run()
	if err != nil {
		return err
	}
	if version == "" {
		x.version, err = pkg.VapusPlatformVersion.Run()
		if err != nil {
			return err
		}
	} else {
		x.version = vapusversion
	}
	x.ConfigSpec.App.Dev = dmutils.StringToBool(devModeStr)
	return nil
}

func (x *VapusdataInstaller) installHashiCorpVault(namespace string) error {
	var helmRepoName = "hashicorp"
	var helmRepoURL = "https://helm.releases.hashicorp.com"
	chartName := "vault"
	vapusGlobals.logger.Info().Msg("Adding Hashicorp Helm repo")
	cmd := exec.Command("helm", "repo", "add", helmRepoName, helmRepoURL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	vapusGlobals.logger.Info().Msgf("Helm repo addition output: %s", string(output))

	// Update Helm repos
	cmd = exec.Command("helm", "repo", "update")
	if err := cmd.Run(); err != nil {
		return err
	}
	vapusGlobals.logger.Info().Msg("Helm repos updated successfully, proceeding with Vault chart installation")
	// Install Vault Helm chart
	args := []string{"install", chartName, helmRepoName + "/" + chartName, "--namespace", namespace, "--create-namespace", "--set", "server.standalone.enabled=true"}
	cmd = exec.Command("helm", args...)
	output, err = cmd.CombinedOutput()
	vapusGlobals.logger.Info().Msgf("Vault chart installation output: %s", string(output))
	if err != nil {
		return err
	}
	vapusGlobals.logger.Info().Msg("Waiting for Vault to be ready...")
	time.Sleep(12 * time.Second)
	vapusGlobals.logger.Info().Msg("Vault chart installed successfully.")
	args = []string{"-n", namespace, "exec", "vault-0", "--", "vault", "operator", "init"}
	// args = []string{"get", "pods"}
	log.Println(strings.Join(args, " "))
	cmd = exec.Command("kubectl", args...)
	output, err = cmd.CombinedOutput()
	vapusGlobals.logger.Info().Msgf("Vault operator init output:\n %s", string(output))
	if err != nil {
		vapusGlobals.logger.Err(err).Msg("Error while initializing vault operator")
		return err
	}
	vapusGlobals.logger.Info().Msg("Vault operator initialized successfully, Copy the unseal keys and root token as these are required for vapusdata installation. They are for one time use display only.")
	return nil
}
