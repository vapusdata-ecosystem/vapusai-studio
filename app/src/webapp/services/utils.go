package services

import (
	"log"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/core/globals"
	"github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	"github.com/vapusdata-oss/aistudio/webapp/models"
	pkgs "github.com/vapusdata-oss/aistudio/webapp/pkgs"
	routes "github.com/vapusdata-oss/aistudio/webapp/routes"
	utils "github.com/vapusdata-oss/aistudio/webapp/utils"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

var OrganizationMap map[string]string

func (x *WebappService) BaseGroupClientCalls(c echo.Context, obj *models.GlobalContexts) error {
	errCh := make(chan error, 3)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		currentOrganization, err := x.grpcClients.GetCurrentOrganization(c)

		if err != nil {
			errCh <- err
			return
		}
		obj.CurrentOrganization = currentOrganization
		errCh <- nil
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := x.grpcClients.GetAccountInfo(c)
		if err != nil {
			errCh <- err
			return
		}
		obj.Account = result.Output
		errCh <- nil
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		res, OrganizationMap, err := x.grpcClients.GetUserInfo(c, "")

		if err != nil {
			errCh <- err
			return
		}
		obj.UserInfo = res
		obj.OrganizationMap = OrganizationMap
		errCh <- nil
	}()

	wg.Wait()
	close(errCh)
	if err := <-errCh; err != nil {
		// If any goroutine returns an error, log it and return the error
		x.logger.Err(err).Msg("error while getting data from server")
		return err
	}
	obj.CurrentUrl = c.Request().URL.EscapedPath()
	log.Println("obj.Account", obj.OrganizationMap)
	// obj.OrganizationMap = x.LoadOrganizationMap(c)
	return nil
}

// TODO - Return error if account is nil or any other error
func (x *WebappService) getDataMarketplaceSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {
	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		CurrentSideBar: currentpage,
		LoginUrl:       routes.Login,
		AccessTokenKey: pkgs.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getHomeSectionGlobals(c echo.Context) (*models.GlobalContexts, error) {
	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		CurrentNav:     routes.HomeNav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: pkgs.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getDatamanagerSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {

	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		CurrentSideBar: currentpage,
		LoginUrl:       routes.Login,
		AccessTokenKey: pkgs.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// func (x *WebappService) getExploreSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {

// 	obj := &models.GlobalContexts{
// 		NavMenuMap:     routes.NavMenuList,
// 		SidebarMap:     routes.ExploreSideList,
// 		CurrentSideBar: currentpage,
// 		CurrentNav:     routes.ExploreNav.String(),
// 		LoginUrl:       routes.Login,
// 		AccessTokenKey: pkgs.ACCESS_TOKEN,
// 	}
// 	err := x.BaseGroupClientCalls(c, obj)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj, nil
// }

func (x *WebappService) getAiStudioSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {

	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		SidebarMap:     routes.ManageAINavSideList,
		CurrentSideBar: currentpage,
		CurrentNav:     routes.ManageAINav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: pkgs.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getStudioSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {

	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		CurrentSideBar: currentpage,
		CurrentNav:     routes.VapusStudioNav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: pkgs.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getSettingsSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {
	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		CurrentSideBar: currentpage,
		CurrentNav:     routes.SettingsNav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: pkgs.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getDeveloperSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {
	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		CurrentSideBar: currentpage,
		CurrentNav:     routes.DevelopersNav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: pkgs.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) responseHandler(c echo.Context, status int, template string, data map[string]interface{}) error {
	return c.Render(status, template, data)
}

func GetUrlParams(c echo.Context, resourceName string) (string, error) {
	_, ok := globals.UrlResouceIdMap[resourceName]
	if !ok {
		return "", utils.ErrInvalidURLResourceName
	}
	return c.Param(resourceName), nil
}

func GetProtoYamlString(obj protoreflect.ProtoMessage) string {
	bytess, err := grpcops.ProtoYamlMarshal(obj)
	if err != nil {
		return ""
	}
	return string(bytess)
}

func GetUserCurrentOrganizationRole(user *mpb.User, currOrganization string) []string {
	for _, Organization := range user.OrganizationRoles {
		if Organization.OrganizationId == currOrganization {
			return Organization.Role
		}
	}
	return []string{}
}
