package svcops

import (
	"context"

	"github.com/rs/zerolog"
	datasvc "github.com/vapusdata-oss/aistudio/core/dataservices"
	models "github.com/vapusdata-oss/aistudio/core/models"
)

type BlobStore struct {
	Path  string
	creds *models.StoreParams
	datasvc.BlobStore
}

func NewBlobStoreClient(ctx context.Context, creds *models.StoreParams, log zerolog.Logger) (*BlobStore, error) {
	s := &BlobStore{}
	client, err := datasvc.NewBlobStoreClient(ctx, creds, log)
	if err != nil {
		return nil, err
	}
	s.creds = creds
	s.BlobStore = client
	return s, nil
}

func (s *BlobStore) GetCreds() *models.NetworkParams {
	return s.GetCreds()
}
