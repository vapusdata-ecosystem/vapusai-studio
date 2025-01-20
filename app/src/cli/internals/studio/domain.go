package plclient

import (
	"context"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
	gwcl "github.com/vapusdata-oss/aistudio/core/serviceops/httpcls"
	dmutils "github.com/vapusdata-oss/aistudio/core/utils"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func (x *VapusPlatformClient) HandleDomainAct(ctx context.Context) error {
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
	requestSpec := &pb.DomainManagerRequest{}
	err = x.protoyamlUnMarshal.Unmarshal(fileBytes, requestSpec)
	if err != nil {
		return err
	}
	fBytes, err := x.protojsonMarshal.Marshal(requestSpec)
	if err != nil {
		return err
	}
	switch requestSpec.Actions.String() {
	case pb.DomainAgentActions_CONFIGURE_DOMAIN.String():
		return x.ConfigureDomain(ctx, fBytes)
	case pb.DomainAgentActions_UPGRADE_DOMAIN_ARTIFACTS.String():
		return x.UpgradeDomainArtifacts(ctx, fBytes)
	case pb.DomainAgentActions_PATCH_DOMAIN.String():
		return x.PatchDomain(ctx, fBytes)
	default:
		return pkg.ErrInvalidAction
	}
}

func (x *VapusPlatformClient) ConfigureDomain(ctx context.Context, requestSpec []byte) error {
	response, err := x.GwClient.DomainManager(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
		Body:  requestSpec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Domain Configured Successfully with ID:", x.logger, response.Output.Domains[0].DomainId)
	return nil
}

func (x *VapusPlatformClient) UpgradeDomainArtifacts(ctx context.Context, requestSpec []byte) error {
	response, err := x.GwClient.DomainManager(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
		Body:  requestSpec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Artifacts successfully upgraded for domain :", x.logger, response.Output.Domains[0].DomainId)
	return nil
}

func (x *VapusPlatformClient) PatchDomain(ctx context.Context, requestSpec []byte) error {
	response, err := x.GwClient.DomainManager(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
		Body:  requestSpec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Domain deployment infra is successfully upgraded for domain :", x.logger, response.Output.Domains[0].DomainId)
	return nil
}

func (x *VapusPlatformClient) ListDomains(ctx context.Context) error {
	// reqBytes, err := json.Marshal(&pb.DomainGetterRequest{
	// 	Actions: pb.DomainAgentActions_LIST_DOMAINS,
	// })
	result, err := x.GwClient.DomainGetter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &pb.DomainGetterRequest{})
	if err != nil {
		return err
	}
	// result, err := vapushttpcl.DataProductManager(x.ActionHandler.AccessToken, reqBytes)
	// if err != nil {
	// 	return err
	// }
	pkg.LogTitles("List of Domain Info: ", x.logger)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Name", "Id", "Datasource Count", "Type", "Total users", "status"})
	for _, dm := range result.Output.Domains {
		tw.AppendRow(table.Row{dm.Name, dm.DomainId, len(dm.DataSources), dm.DomainType.String(), len(dm.GetUsers()), dm.Status})
		tw.AppendSeparator()
	}
	tw.Render()
	return nil
}

func (x *VapusPlatformClient) DescribeDomains(ctx context.Context) error {
	result, err := x.GwClient.DomainGetter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &pb.DomainGetterRequest{})
	if err != nil {
		return err
	}
	pkg.LogTitles("Logged-In user's Domain Info: ", x.logger)
	x.PrintDescribe(result.Output.Domains[0], "domain")
	return nil
}
