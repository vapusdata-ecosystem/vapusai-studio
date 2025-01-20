package plclient

import (
	"context"
	"log"
	"strings"

	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
	gwcl "github.com/vapusdata-oss/aistudio/core/serviceops/httpcls"
	dmutils "github.com/vapusdata-oss/aistudio/core/utils"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func (x *VapusPlatformClient) HandleAccountAct(ctx context.Context) error {
	var fileBytes []byte
	var err error
	if x.ActionHandler.File != "" {
		x.protoyamlUnMarshal.Path = x.ActionHandler.File
		fileBytes, err = dmutils.ReadFile(x.ActionHandler.File)
		if err != nil {
			return err
		}
		x.inputFormat = strings.ToUpper(dmutils.GetConfFileType(x.ActionHandler.File))
	}
	requestSpec := &pb.AccountManagerRequest{}
	err = x.protoyamlUnMarshal.Unmarshal(fileBytes, requestSpec)
	if err != nil {
		return err
	}
	fBytes, err := x.protojsonMarshal.Marshal(requestSpec)
	if err != nil {
		return err
	}
	return x.AccountManagerClient(ctx, fBytes)
}

func (x *VapusPlatformClient) ListAccount(ctx context.Context) error {
	log.Println("Listing account")
	result, err := x.GwClient.AccountGetter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	})
	if err != nil {
		return err
	}
	obj, err := x.protoyamlMarshal.Marshal(result.Output)
	if err != nil {
		return err
	}
	return pkg.ParseAndBuildYamlTable(obj)
	// maps, err := dmutils.StructToMap(result.Output)
	// pkg.LogTitles("Account Info of current loggedin instance: ", x.logger)
	// tw := pkg.NewTableWritter()
	// // tw.AppendHeader(table.Row{"Account name", "Account Id", "Data Store", "Secret Store", "Artifact Store", "Authn Method", "Authz"})
	// tw.AppendHeader(table.Row{"Heading", "Value"})
	// for k, v := range maps {
	// 	tw.AppendRow(table.Row{k, v})
	// }
	// // tw.AppendRow(table.Row{result.Output.Name, result.Output.AccountId, result.Output.BackendDataStorage.BesService, result.Output.BackendSecretStorage.BesService,
	// // 	result.Output.ArtifactStorage.BesService, result.Output.AuthnMethod, result.Output.DmAccessJwtKeys.SigningAlgorithm})
	// tw.AppendSeparator()
	// tw.Render()
	// return nil
}

func (x *VapusPlatformClient) AccountManagerClient(ctx context.Context, requestSpec []byte) error {
	log.Println(string(requestSpec))
	response, err := x.GwClient.AccountManager(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
		Body:  requestSpec,
	})
	// response, err := x.PlConn.AccountManager(ctx, requestSpec)
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Account Configured Successfully with ID:", x.logger, response.Output.AccountId)
	return nil
}
