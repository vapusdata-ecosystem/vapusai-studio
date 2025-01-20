package emailer

import (
	"context"
	"encoding/base64"

	"github.com/rs/zerolog"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"github.com/vapusdata-oss/aistudio/core/models"
	aws "github.com/vapusdata-oss/aistudio/core/thirdparty/aws"
	gcp "github.com/vapusdata-oss/aistudio/core/thirdparty/gcp"
	sendgrid "github.com/vapusdata-oss/aistudio/core/thirdparty/sendgrid"
	nabhikutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type Emailer interface {
	SendEmail(ctx context.Context, opts *models.EmailOpts, agentId string) error
	SendRawEmail(ctx context.Context, opts *models.EmailOpts, agentId string) error
}

type EmailerClient struct {
	logger zerolog.Logger
	Emailer
}

func New(ctx context.Context,
	pType string,
	netOps *models.PluginNetworkParams,
	ops []*models.Mapper,
	logger zerolog.Logger) (Emailer, error) {
	if netOps == nil || netOps.Credentials == nil {
		logger.Error().Msg("Invalid emailer credentials")
		return nil, dmerrors.DMError(nabhikutils.ErrInvalidEmailerCredentials, nabhikutils.ErrEmailerConn)
	}
	emailer := &EmailerClient{
		logger: logger,
	}
	switch pType {
	case mpb.IntegrationPlugins_AMAZON_SES.String():
		client, err := aws.NewAwsSesClient(ctx, &aws.AWSConfig{
			Region:          netOps.Credentials.AwsCreds.Region,
			AccessKeyId:     netOps.Credentials.AwsCreds.AccessKeyId,
			SecretAccessKey: netOps.Credentials.AwsCreds.SecretAccessKey,
		}, logger)
		if err != nil {
			logger.Err(err).Msg("Error creating aws ses client")
			return nil, err
		}
		emailer.Emailer = client
	case mpb.IntegrationPlugins_GMAIL.String():
		decodeData, err := base64.StdEncoding.DecodeString(netOps.Credentials.GcpCreds.ServiceAccountKey)
		if err != nil {
			logger.Err(err).Msg("Error decoding gcp service account key")
			return nil, err
		}
		client, err := gcp.NewGmailAgent(ctx, &gcp.GcpConfig{
			ServiceAccountKey: []byte(decodeData),
			ProjectID:         netOps.Credentials.GcpCreds.ProjectId,
			Region:            netOps.Credentials.GcpCreds.Region,
		}, logger)
		if err != nil {
			logger.Err(err).Msg("Error creating gcp secret manager client")
			return nil, err
		}
		emailer.Emailer = client
	case mpb.IntegrationPlugins_SENDGRID.String():
		client, err := sendgrid.New(ctx, &sendgrid.SendgridConfig{
			APIKey: netOps.Credentials.GetApiToken(),
			Host:   netOps.URL,
			Params: ops,
		}, logger)
		if err != nil {
			logger.Err(err).Msg("Error creating sendgrid client")
			return nil, err
		}
		emailer.Emailer = client
	default:
		logger.Error().Msg("Invalid emailer engine")
		return nil, dmerrors.DMError(nabhikutils.ErrInvalidEmailerEngine, nabhikutils.ErrEmailerConn)
	}
	return emailer, nil
}
