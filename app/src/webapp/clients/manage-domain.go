package clients

import (
	"github.com/labstack/echo/v4"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	dpb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func (s *GrpcClient) GetCurrentOrganization(eCtx echo.Context) (*mpb.Organization, error) {
	result, err := s.SvcConn.OrganizationManager(s.SetAuthCtx(eCtx), &dpb.OrganizationManagerRequest{
		Actions: dpb.OrganizationAgentActions_DESCRIBE_ORG,
	})
	if err != nil || len(result.Output.GetOrgs()) == 0 {
		s.logger.Err(err).Msg("error while getting Organization current logged in info")
		return nil, err
	}
	return result.Output.GetOrgs()[0], nil
}

func (s *GrpcClient) GetOrganizations(eCtx echo.Context) []*mpb.Organization {
	result, err := s.SvcConn.OrganizationGetter(s.SetAuthCtx(eCtx), &dpb.OrganizationGetterRequest{})
	if err != nil {
		s.logger.Err(err).Msg("error while getting list of Organizations")
		return []*mpb.Organization{}
	}
	return result.Output.GetOrgs()
}
