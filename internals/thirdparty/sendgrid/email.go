package sendgrid

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
	"github.com/vapusdata-oss/aistudio/core/models"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type Sendgrid interface {
	SendEmail(ctx context.Context, opts *models.EmailOpts, agentId string) error
	SendRawEmail(ctx context.Context, opts *models.EmailOpts, agentId string) error
}

type SendgridClient struct {
	client *sendgrid.Client
	config *SendgridConfig
	logger zerolog.Logger
	Params []*models.Mapper
}

type SendgridConfig struct {
	APIKey string
	Host   string
	Params []*models.Mapper
}

func New(ctx context.Context, opts *SendgridConfig, logger zerolog.Logger) (Sendgrid, error) {
	cl := sendgrid.NewSendClient(opts.APIKey)
	if opts.Host != "" {
		cl.Request.BaseURL = opts.Host
	} else {
		opts.Host = "https://api.sendgrid.com"
	}
	return &SendgridClient{
		client: cl,
		config: opts,
		logger: dmlogger.GetSubDMLogger(logger, "emailer", "sendgrid"),
	}, nil
}

func (s *SendgridClient) GetParams(pType string) string {
	for _, param := range s.config.Params {
		log.Println(param.Key, "====================+++++++++++++????", param.Value, "====================+++++++++++++>>")
		if param.Key == pType {
			log.Println(param.Value, "====================+++++++++++++")
			return param.Value
		}
	}
	return ""
}

func (s *SendgridClient) SendEmail(ctx context.Context, opts *models.EmailOpts, agentId string) error {
	var err error
	message := mail.NewSingleEmail(mail.NewEmail(opts.SenderName, opts.From),
		opts.Subject,
		mail.NewEmail("", opts.To[0]),
		opts.Body,
		opts.HtmlTemplateBody)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err = client.Send(message)
	if err != nil {
		s.logger.Error().Err(err).Msgf("Error sending email for agent %s", agentId)
	} else {
		s.logger.Info().Msgf("Email sent for agent %s", agentId)
	}
	return err
}

func (s *SendgridClient) SendRawEmail(ctx context.Context, opts *models.EmailOpts, agentId string) error {
	var err error
	log.Println(s.GetParams(mpb.EmailSettings_SENDER_NAME.String()), s.GetParams(mpb.EmailSettings_SENDER_EMAIL.String()), "====================+++++++++++++")
	message := mail.NewV3Mail()
	message.SetFrom(mail.NewEmail(s.GetParams(mpb.EmailSettings_SENDER_NAME.String()), s.GetParams(mpb.EmailSettings_SENDER_EMAIL.String())))
	message.Subject = opts.Subject
	personalization := mail.NewPersonalization()
	ttos := func() []*mail.Email {
		var tos []*mail.Email
		for _, to := range opts.To {
			tos = append(tos, mail.NewEmail("", to))
		}
		return tos
	}
	personalization.AddTos(ttos()...)
	personalization.Subject = opts.Subject
	if opts.HtmlTemplateBody != "" {
		content := mail.NewContent("text/html", opts.HtmlTemplateBody)
		message.AddContent(content)
	} else {
		content := mail.NewContent("text/plain", opts.Body)
		message.AddContent(content)
	}
	message.AddPersonalizations(personalization)
	for _, attach := range opts.Attachments {
		attachmentData := base64.StdEncoding.EncodeToString(attach.Data)
		att := mail.NewAttachment()
		att.SetContent(attachmentData)
		att.Type = attach.ContentType
		att.Filename = attach.FileName
		att.SetDisposition("attachment")
		message.AddAttachment(att)
	}
	request := sendgrid.GetRequest(s.config.APIKey, "/v3/mail/send", s.config.Host)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(message)
	response, err := sendgrid.API(request)
	if response != nil {
		log.Println(response.StatusCode, response.Body, response.Headers, "====================+++++++++++++")
	}
	if err != nil {
		s.logger.Error().Err(err).Msgf("Error sending email for agent %s", agentId)
	} else {
		s.logger.Info().Msgf("Email sent for agent %s with message id %s", agentId, response.Body)
	}
	return err
}
