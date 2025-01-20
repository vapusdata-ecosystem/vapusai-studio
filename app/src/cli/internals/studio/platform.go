package plclient

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/table"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
	"github.com/vapusdata-oss/aistudio/core/globals"
	gwcl "github.com/vapusdata-oss/aistudio/core/serviceops/httpcls"
	dmutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func (x *VapusPlatformClient) RetrieveLoginURL() (*pb.LoginHandlerResponse, error) {
	return x.GwClient.LoginHandler(&gwcl.HttpRequestGeneric{
		Method: globals.GET,
		Body:   []byte{},
	})
}

func (x *VapusPlatformClient) RetrieveAccessToken(code, host string) (string, string, error) {
	log.Println("Retrieving access token for code: ", code)
	log.Println("Retrieving access token for url: ", host)
	req := &pb.LoginCallBackRequest{Code: code, Host: host}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return "", "", err
	}
	result, err := x.GwClient.LoginCallback(&gwcl.HttpRequestGeneric{
		Method: globals.POST,
		Body:   reqBytes,
	})
	if err != nil {
		return "", "", err
	}
	// result, err := x.UserConn.LoginCallback(context.Background(), &pb.LoginCallBackRequest{Code: code, Host: host})
	// if err != nil {
	// 	return "", "", err
	// }
	return result.Token.GetAccessToken(), result.Token.GetIdToken(), nil
}

func (x *VapusPlatformClient) RetrievePlatformAccessToken(ctx context.Context, token, domain string) (string, error) {
	reqBytes, err := x.protojsonMarshal.Marshal(&pb.AccessTokenInterfaceRequest{Domain: domain, Utility: pb.AccessTokenAgentUtility_DOMAIN_LOGIN})
	if err != nil {
		return "", err
	}
	result, err := x.GwClient.AccessTokenInterface(&gwcl.HttpRequestGeneric{
		Token: token,
		Body:  reqBytes,
	})
	if err != nil {
		return "", err
	}
	return result.Token.AccessToken, nil
}

func (x *VapusPlatformClient) ListUser(ctx context.Context) error {
	result, err := x.GwClient.UserGetter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &pb.UserGetterRequest{})
	if err != nil {
		return err
	}
	pkg.LogTitles("Logged-in User Info: ", x.logger)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"UserId", "Name", "Added On", "Domain", "Roles"})
	for _, user := range result.Output.Users {
		tw.AppendRow(table.Row{user.UserId, user.DisplayName, user.InvitedOn, user.DomainRoles[0].DomainId, user.DomainRoles[0].Role})
	}
	tw.AppendSeparator()
	tw.Render()
	return nil
}

func (x *VapusPlatformClient) DescribeUser(ctx context.Context) error {
	result, err := x.GwClient.UserGetter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &pb.UserGetterRequest{
		Action: pb.UserGetterActions_GET_USER,
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("Logged-in User Info: ", x.logger)
	x.PrintDescribe(result.Output, "user")
	return nil
}

func (x *VapusPlatformClient) ListPlatformSpec() {
	fm := []interface{}{}
	for _, f := range mpb.ContentFormats_name {
		fm = append(fm, f)
	}
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Resource Name", "Spec Available", "Formats Available", "Generate command"})
	for _, spec := range mpb.RequestObjects_name {
		if spec == mpb.RequestObjects_INVALID_REQUEST_OBJECT.String() {
			continue
		}
		tw.AppendRow(table.Row{spec, true, pkg.NewListWritter(fm, list.StyleBulletSquare).Render(), pkg.APPNAME + " spec --name " + spec + " --generate-file=true --format yaml"})
		tw.AppendSeparator()
	}
	tw.Render()
}

func (x *VapusPlatformClient) GeneratePlatformSpec(token, specName, format string, withFakeData bool) error {
	log.Println("Generating spec for ", specName)
	log.Println("Generating spec in format ", x)
	val := SpecMap[mpb.RequestObjects(mpb.RequestObjects_value[specName])]
	specVal, err := dmutils.GenericMarshaler(val, format)
	if err != nil {
		return err
	}

	fileName := strings.ToLower(specName + "." + strings.ToLower(format))
	x.logger.Info().Msgf("Sample %v spec generated with file name - %v \n", specName, fileName)
	err = os.WriteFile(fileName, specVal, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (x *VapusPlatformClient) SearchVapusData(ctx context.Context) error {
	response, err := x.GwClient.VapusSearch(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &pb.VapusSearchRequest{
		Q:          x.ActionHandler.SearchQ,
		SearchType: pb.VapusSearchType(pb.VapusSearchType_value[strings.ToUpper(x.ActionHandler.Action)]),
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("Search Result:", x.logger)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Name", "Identifier", "Resource"})
	for _, item := range response.Results {
		tw.AppendRow(table.Row{item.Result.Name, item.Result.ResourceId, item.Result.Resource})
		tw.AppendSeparator()
	}
	tw.Render()
	return nil
}
