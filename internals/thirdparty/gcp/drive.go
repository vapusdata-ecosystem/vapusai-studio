package gcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/vapusdata-oss/aistudio/core/models"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

type GDriveOpts struct {
	Client *drive.Service
	logger zerolog.Logger
}

func NewDriveAgentOld(ctx context.Context, opts *GcpConfig, logger zerolog.Logger) (*GDriveOpts, error) {
	client, err := drive.NewService(ctx, option.WithCredentialsJSON(opts.ServiceAccountKey))
	if err != nil {
		logger.Err(err).Msgf("Error while creating credentials from json -- %v", err)
		return nil, err
	}
	return &GDriveOpts{
		Client: client,
		logger: logger,
	}, nil
}

func NewDriveAgent(ctx context.Context, opts *GcpConfig, userEmailAddr string, logger zerolog.Logger) (*GDriveOpts, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(string(opts.ServiceAccountKey))
	if err != nil {
		logger.Err(err).Msgf("Error while decoding the GCP KEY -- %v", err)
		return nil, err
	}
	creds, err := google.CredentialsFromJSON(ctx, decodedKey)
	if err != nil || creds == nil {
		logger.Err(err).Msgf("Error while creating credentials from json for GCP drive plugin-- %v", err)
		return nil, err
	}
	keyJson := map[string]interface{}{}
	err = json.Unmarshal(creds.JSON, &keyJson)
	if err != nil {
		logger.Err(err).Msgf("Error while unmarshalling the GCP KEY json -- %v", err)
		return nil, err
	}
	clEmail, ok := keyJson["client_email"].(string)
	if !ok {
		logger.Err(err).Msgf("Error while getting the client_email from the GCP KEY json -- %v", err)
		return nil, err
	}
	log.Println("gcp-SVC-KEY-EMAIL", string(decodedKey))
	log.Println("gcp-SVC-KEY-EMAIL", clEmail)
	log.Println("gcp-SVC-KEY-EMAIL", userEmailAddr)
	log.Println("gcp-SVC-KEY-TokenSource", creds.TokenSource)

	// Now impersonate the organization user
	tokenSource, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: clEmail,
		Subject:         userEmailAddr,
		Scopes:          []string{"https://www.googleapis.com/auth/drive"},
	}, option.WithCredentialsJSON(decodedKey))
	if err != nil {
		logger.Err(err).Msgf("Error while impersonating the user -- %v", err)
		return nil, err
	}
	client, err := drive.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		logger.Err(err).Msgf("Error while creating credentials from json -- %v", err)
		return nil, err
	}
	return &GDriveOpts{
		Client: client,
		logger: logger,
	}, nil
}

func (g *GDriveOpts) Upload(ctx context.Context, opts *models.FileManageOpts) error {
	var err error
	log.Println("gcp-DRIVE-UPLOAD", opts)
	tempDir := os.TempDir()
	tempFileName := opts.FileName
	tempFile := filepath.Join(tempDir, tempFileName)
	err = os.WriteFile(tempFile, opts.Data, 0644)
	if err != nil {
		g.logger.Err(err).Msgf("Error while writing the file -- %v", err)
		return err
	}
	defer os.Remove(tempFile)
	file, err := os.Open(tempFile)
	if err != nil {
		g.logger.Err(err).Msgf("Error while opening the file -- %v", err)
		return err
	}
	defer file.Close()
	fileMetadata := &drive.File{
		Name: filepath.Join(opts.Path, opts.FileName),
	}
	_, err = g.Client.Files.Create(fileMetadata).Media(file).Do()
	if err != nil {
		g.logger.Err(err).Msgf("Error while uploading the file -- %v", err)
		return err
	}
	return nil
}

func (g *GDriveOpts) ListFiles(ctx context.Context, opts *models.FileManageOpts) error {
	var err error
	_, err = g.Client.Files.List().Do()
	if err != nil {
		g.logger.Err(err).Msgf("Error while listing the files -- %v", err)
		return err
	}
	return nil
}

func (g *GDriveOpts) DeleteFiles(ctx context.Context, opts *models.FileManageOpts) error {
	var err error
	err = g.Client.Files.Delete(opts.FileName).Do()
	if err != nil {
		g.logger.Err(err).Msgf("Error while deleting the file -- %v", err)
		return err
	}
	return nil
}
