package svcops

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	datasvc "github.com/vapusdata-oss/aistudio/core/dataservices"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
	models "github.com/vapusdata-oss/aistudio/core/models"
)

type DataStoreOptions struct {
	creds *models.StoreParams
}

type dataStoreOpts func(*DataStoreOptions)

func WithDataStoreCreds(r *models.StoreParams) dataStoreOpts {
	return func(s *DataStoreOptions) {
		s.creds = r
	}
}

func NewDataStore(ctx context.Context, log zerolog.Logger, opts ...dataStoreOpts) (*datasvc.DataStoreClient, error) {
	s := &DataStoreOptions{}
	for _, opt := range opts {
		opt(s)
	}
	return datasvc.New(ctx, datasvc.WithStoreParams(s.creds), datasvc.WithLogger(dmlogger.GetSubDMLogger(log, "vapusloader Package", "DataStoreClient")))
}

func ListResourceWithGovernance(ctxClaim map[string]string) string {
	if ctxClaim != nil {
		return fmt.Sprintf(`status = 'ACTIVE' AND deleted_at IS NULL AND 
		(
	(scope='PLATFORM_SCOPE') OR (scope='DOMAIN_SCOPE' AND organization='%s') OR (scope='USER_SCOPE' AND organization='%s' AND created_by='%s')
		)
		 ORDER BY created_at DESC`, ctxClaim[encryption.ClaimOrganizationKey], ctxClaim[encryption.ClaimOrganizationKey], ctxClaim[encryption.ClaimUserIdKey])
	}
	return ""
}
