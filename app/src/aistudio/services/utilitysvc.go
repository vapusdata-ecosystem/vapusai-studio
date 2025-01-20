package services

import (
	"context"
	"fmt"
	"strings"

	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	globals "github.com/vapusdata-oss/aistudio/core/globals"
	"github.com/vapusdata-oss/aistudio/core/models"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type BlobAgent struct {
	*models.VapusInterfaceAgentBase
	method        string
	uploadRequest *pb.UploadRequest
	uploadStream  pb.UtilityService_UploadStreamServer
	organization  *models.Organization
	uploadResult  *pb.UploadResponse
	*dmstores.DMStore
}

func (s *StudioServices) NewUtilityAgent(ctx context.Context, uploadRequest *pb.UploadRequest, uploadStream pb.UtilityService_UploadStreamServer) (*BlobAgent, error) {
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, encryption.ErrInvalidJWTClaims
	}
	organization, err := s.DMStore.GetOrganization(ctx, vapusStudioClaim[encryption.ClaimOrganizationKey], vapusStudioClaim)
	if err != nil {
		s.logger.Error().Err(err).Ctx(ctx).Msg("error while getting organization from datastore")
		return nil, dmerrors.DMError(utils.ErrOrganization404, err)
	}
	agent := &BlobAgent{
		uploadRequest: uploadRequest,
		uploadStream:  uploadStream,
		DMStore:       s.DMStore,
		organization:  organization,
		VapusInterfaceAgentBase: &models.VapusInterfaceAgentBase{
			CtxClaim: vapusStudioClaim,
			Ctx:      ctx,
			InitAt:   coreutils.GetEpochTime(),
		},
	}
	agent.SetAgentId()
	agent.Logger = pkgs.GetSubDMLogger(globals.AIPROMPTAGENT.String(), agent.AgentId)
	return agent, nil
}

func (v *BlobAgent) GetUploadedResult() *pb.UploadResponse {
	v.FinishAt = coreutils.GetEpochTime()
	v.FinalLog()
	return v.uploadResult
}

func (v *BlobAgent) Act() error {
	if v.uploadRequest != nil {
		v.uploadResult = &pb.UploadResponse{
			Output: []*pb.UploadResponse_ObjectUploadResult{},
		}
		return v.uploadFile()
	} else if v.uploadStream != nil {
		return v.uploadFileStream()
	}
	return nil
}

func (v *BlobAgent) uploadFile() error {
	if v.uploadRequest.Resource == "" {
		v.Logger.Error().Msg("resource is empty")
		return dmerrors.DMError(utils.ErrMissingUploadResourceName, nil)
	}
	if v.uploadRequest.ResourceId == "" {
		v.uploadRequest.ResourceId = v.CtxClaim[encryption.ClaimUserIdKey]
	}
	for _, fileData := range v.uploadRequest.Objects {
		if fileData.Name == "" && fileData.Data == nil {
			v.Logger.Error().Msg("file data is empty")
			v.uploadResult.Output = append(v.uploadResult.Output, &pb.UploadResponse_ObjectUploadResult{
				Error:  "file data is empty",
				Object: fileData,
			})
			continue
		}
		fType := coreutils.GetConfFileType(fileData.Name)
		fileData.Format = mpb.ContentFormats(mpb.ContentFormats_value[strings.ToUpper(fType)])
		if fileData.Name == "" {
			fileData.Name = coreutils.GetUUID() + "." + strings.ToLower(fileData.Format.String())
		} else {
			fTL := len(fileData.Format.String()) + 1
			if strings.ToLower(fileData.Name[len(fileData.Name)-fTL:]) != "."+strings.ToLower(fileData.Format.String()) {
				fileData.Name = fileData.Name + "." + strings.ToLower(fileData.Format.String())
			}
		}
		if fileData.Data == nil {
			v.Logger.Error().Msg("file data is empty")
			v.uploadResult.Output = append(v.uploadResult.Output, &pb.UploadResponse_ObjectUploadResult{
				Error:  "file data is empty",
				Object: fileData,
			})
			continue
		}

		keyPath := fmt.Sprintf("%s/%s/%s", v.uploadRequest.Resource, v.uploadRequest.ResourceId, fileData.Name)
		err := v.DMStore.BlobOps.UploadObject(v.Ctx, v.CtxClaim[encryption.ClaimOrganizationKey], keyPath, fileData.Data)
		fileData.Data = nil
		if err != nil {
			v.Logger.Error().Msg("error while uploading file to blob storage")
			v.uploadResult.Output = append(v.uploadResult.Output, &pb.UploadResponse_ObjectUploadResult{
				Error:  err.Error(),
				Object: fileData,
			})
		} else {
			v.uploadResult.Output = append(v.uploadResult.Output, &pb.UploadResponse_ObjectUploadResult{
				Object:       fileData,
				ResponsePath: keyPath,
				Fid:          coreutils.GetUUID(),
			})
		}
	}

	return nil
}

func (v *BlobAgent) uploadFileStream() error {
	// upload file stream
	return nil
}
