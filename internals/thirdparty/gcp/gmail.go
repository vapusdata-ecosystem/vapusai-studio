package gcp

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/rs/zerolog"
	"github.com/vapusdata-oss/aistudio/core/models"
	"google.golang.org/api/gmail/v1"
	gapi "google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type Gmail interface {
	SendEmail(ctx context.Context, opts *models.EmailOpts, agentId string) error
	SendRawEmail(ctx context.Context, opts *models.EmailOpts, agentId string) error
}

type GmailClient struct {
	Client *gmail.Service
	logger zerolog.Logger
}

func NewGmailAgent(ctx context.Context, opts *GcpConfig, logger zerolog.Logger) (Gmail, error) {
	client, err := gmail.NewService(ctx, option.WithCredentialsJSON(opts.ServiceAccountKey))
	if err != nil {
		logger.Err(err).Msgf("Error while creating credentials from json for GMAIL -- %v", err)
		return nil, err
	}
	return &GmailClient{
		Client: client,
		logger: logger,
	}, nil
}

func (x *GmailClient) SendEmail(ctx context.Context, opts *models.EmailOpts, agentId string) error {
	// Create the email message
	msg := &gmail.Message{
		ThreadId: "test_thread", // optional, depending on your use case
		Payload: &gmail.MessagePart{
			Body: &gmail.MessagePartBody{},
		},
	}
	bo := ""
	if opts.HtmlTemplateBody != "" {
		bo = opts.HtmlTemplateBody
	} else {
		bo = opts.Body
	}
	// Build the email content
	email := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"From: " + opts.From + "\r\n" +
		"To: " + strings.Join(opts.To, ",") + "\r\n" +
		"Subject: " + opts.Subject + "\r\n\r\n" +
		bo)

	// Encode the email content as base64
	encodedMessage := base64.URLEncoding.EncodeToString(email)

	// Set the raw email content
	msg.Raw = encodedMessage

	// Send the email
	_, err := x.Client.Users.Messages.Send("me", msg).Do()
	return err
}

func (x *GmailClient) SendRawEmail(ctx context.Context, opts *models.EmailOpts, agentId string) error {
	// Create the email message
	msg := &gmail.Message{
		ThreadId: "test_thread", // optional, depending on your use case
		Payload: &gmail.MessagePart{
			Body: &gmail.MessagePartBody{},
		},
	}
	bo := ""
	if opts.HtmlTemplateBody != "" {
		bo = opts.HtmlTemplateBody
	} else {
		bo = opts.Body
	}
	// Build the email content
	email := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"From: " + opts.From + "\r\n" +
		"To: " + strings.Join(opts.To, ",") + "\r\n" +
		"Subject: " + opts.Subject + "\r\n\r\n" +
		bo)

	// Encode the email content as base64
	encodedMessage := base64.StdEncoding.EncodeToString(email)

	// Set the raw email content
	msg.Raw = encodedMessage

	// Send the email
	att := opts.Attachments[0]
	decodedAttachment, err := base64.StdEncoding.DecodeString(string(att.Data))
	if err != nil {
		x.logger.Err(err).Msgf("Error while decoding attachment -- %v", err)
		return err
	}

	_, err = x.Client.Users.Messages.Send("me", msg).Media(strings.NewReader(string(decodedAttachment)),
		gapi.ContentType(att.ContentType),
	).Do()
	return err
}
