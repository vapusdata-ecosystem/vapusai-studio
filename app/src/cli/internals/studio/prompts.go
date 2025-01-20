package plclient

import (
	"context"
	"log"

	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
	dpb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func (x *VapusPlatformClient) DataProductPrompt(ctx context.Context) error {
	response, err := x.DomainConn.DataProductServer(ctx, &dpb.DataProductServerRequest{
		DataProductId: x.ActionHandler.Params[pkg.DataproductKey],
		Query:         x.ActionHandler.Params[pkg.SearchqueryKey],
	})
	if err != nil {
		return err
	}
	log.Println(response.Output)
	return nil
}
