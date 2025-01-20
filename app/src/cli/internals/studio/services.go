package plclient

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
)

func (x *VapusPlatformClient) ListResourceActions(resource string) {
	xx := text.FormatUpper.Apply("Actions for Resource: ")
	xx = text.Underline.Sprintf(xx)
	x.logger.Info().Msgf("\n%v", xx)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Actions", "Version", "Commands"})
	for _, action := range x.ResourceActionMap[resource] {
		if strings.Contains(action.(string), "INVALID") {
			continue
		}
		tw.AppendRow(table.Row{action, "v1alpha1", pkg.APPNAME + " act " + resource + " --file <Input File> "})
		tw.AppendSeparator()
	}
	tw.Render()
}

func (x *VapusPlatformClient) HandleAction() error {
	if x.ActionHandler.ParentCmd == "" {
		return errors.New("invalid operations")
	}
	switch x.ActionHandler.ParentCmd {
	case pkg.GetOps:
		return x.HandleGet()
	case pkg.DescribeOps:
		return x.HandleDescription()
	case pkg.ActOps:
		return x.HandleAct()
	case pkg.SearchOpts:
		return x.HandleSearch()
	case pkg.GetPrompt:
		return x.HandlePrompt()
	default:
		return pkg.ErrInvalidAction
	}
}

func (x *VapusPlatformClient) HandleAct() error {
	ctx := context.Background()
	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
	switch x.ActionHandler.Resource {
	case pkg.AccountResource:
		return x.HandleAccountAct(newCtx)
	case pkg.DataProductResource:
		return x.HandleDataProductAct(newCtx)
	case pkg.DataWorkerResource:
		return x.HandleDataworkerAct(newCtx)
	case pkg.MarketplaceResource:
		return x.HandleDataMarketplaceAct(newCtx)
	case pkg.DataSourceResource:
		return x.HandleDataSourceAct(newCtx)
	case pkg.DomainResource:
		return x.HandleDomainAct(newCtx)
	case pkg.DataWorkerDeploymentResource:
		return x.HandleDataWorkerDeploymentAct(newCtx)
	case pkg.AIModelNodeResource:
		return x.HandleAIStudioAct(newCtx)
	case pkg.DataProductDeploymentResource:
		return x.HandleVDCDeploymentAct(newCtx)
	default:
		return pkg.ErrInvalidResource
	}
}

func (x *VapusPlatformClient) HandleGet() error {
	log.Println("Getting resource..............................")
	ctx := context.Background()
	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
	switch x.ActionHandler.Resource {
	case pkg.AccountResource:
		return x.ListAccount(newCtx)
	case pkg.UserResource:
		return x.ListUser(newCtx)
	case pkg.DomainResource:
		return x.ListDomains(newCtx)
	case pkg.DataSourceResource:
		return x.ListDataSources(newCtx)
	case pkg.MarketplaceResource:
		return x.ListDataMarketplace(newCtx)
	case pkg.DataWorkerResource:
		return x.ListDataWorkers(newCtx)
	case pkg.DataProductResource:
		return x.ListDataProducts(newCtx)
	case pkg.DataProductDeploymentResource:
		return x.ListVDCOrchestrators(newCtx)
	case pkg.AIModelNodeResource:
		return x.getAIModelNodes(newCtx)
	case pkg.DataWorkerDeploymentResource:
		return x.listDataworkerDeploments(newCtx)
	default:
		return pkg.ErrInvalidResource
	}
}

func (x *VapusPlatformClient) HandleDescription() error {
	var iden string = ""
	if x.ActionHandler.Resource != pkg.DomainResource {
		if len(x.ActionHandler.Args) < 1 {
			return pkg.ErrNoArgs
		}
		iden = pkg.GetDescId(x.ActionHandler.Args)
	}
	ctx := context.Background()
	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
	log.Println("Describing with resource identifier: ", iden)
	switch x.ActionHandler.Resource {
	case pkg.UserResource:
		return x.DescribeUser(newCtx)
	case pkg.DomainResource:
		return x.DescribeDomains(newCtx)
	case pkg.DataSourceResource:
		return x.DescribeDataSources(newCtx, iden)
	case pkg.MarketplaceResource:
		return x.DescribeDataMarketplace(newCtx, iden)
	case pkg.DataWorkerResource:
		return x.DescribeDataWorkers(newCtx, iden)
	case pkg.DataProductResource:
		return x.DescribeDataProducts(newCtx, iden)
	case pkg.DataProductDeploymentResource:
		return x.DescribeVDCOrchestrator(newCtx, iden)
	case pkg.AIModelNodeResource:
		return x.describeAIModelNode(newCtx, iden)
	case pkg.DataWorkerDeploymentResource:
		return x.describeDataworkerDeploment(newCtx, iden)
	default:
		return pkg.ErrInvalidAction
	}
}

// func (x *VapusPlatformClient) HandleConfigure() error {
// 	var fileBytes []byte
// 	var err error
// 	log.Println("File Path: ", x.ActionHandler.File)
// 	if x.ActionHandler.File != "" {
// 		x.protoyamlUnMarshal.Path = x.ActionHandler.File
// 		fileBytes, err = dmutils.ReadFile(x.ActionHandler.File)
// 		if err != nil || len(fileBytes) == 0 {
// 			return err
// 		}
// 		x.inputFormat = strings.ToUpper(dmutils.GetConfFileType(x.ActionHandler.File))
// 	} else {
// 		return pkg.ErrNoFileInput
// 	}

// 	ctx := context.Background()
// 	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
// 	switch x.ActionHandler.Resource {
// 	// case pkg.DataProductResource:
// 	// 	return x.ConfigureDataProduct(newCtx, fileBytes)
// 	// case pkg.DataProductDeploymentResource:
// 	// 	return x.DeployDataProduct(newCtx, fileBytes)
// 	// case pkg.DataWorkerDeploymentResource:
// 	// 	return x.DeployDataWorkers(newCtx, fileBytes)
// 	// case pkg.AIModelNodeResource:
// 	// 	return x.ConfigureAIModelNode(newCtx, fileBytes)
// 	// case pkg.DataSourceResource:
// 	// 	return x.ConfigureDataSource(newCtx, fileBytes)
// 	// case pkg.DomainResource:
// 	// 	return x.ConfigureDomain(newCtx, fileBytes)
// 	// case pkg.DataWorkerResource:
// 	// 	return x.ConfigureDataWorker(newCtx, fileBytes)
// 	default:
// 		return pkg.ErrInvalidAction
// 	}
// }

func (x *VapusPlatformClient) HandleSearch() error {
	ctx := context.Background()
	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
	switch x.ActionHandler.Resource {
	case pkg.SearchResource:
		return x.SearchVapusData(newCtx)
	default:
		return pkg.ErrInvalidAction
	}
}

func (x *VapusPlatformClient) HandlePrompt() error {
	ctx := context.Background()
	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
	return x.DataProductPrompt(newCtx)
}
