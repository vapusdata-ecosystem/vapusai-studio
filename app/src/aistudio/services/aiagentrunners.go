package services

import (
	"encoding/json"
	"strings"

	grpccodes "google.golang.org/grpc/codes"

	"github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	"github.com/vapusdata-oss/aistudio/core/globals"
	models "github.com/vapusdata-oss/aistudio/core/models"
	grpcops "github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func (s *VapusAIAgentThread) FileUploader(link *VapusAgentChainLink) error {
	var err error
	if link.userParams.Input != "" {
		err = s.RenderSteps(link)
		if err != nil {
			link.failedStep = "Input analyzation failed for agent " + link.vapusAgent.Name + " failed"
			s.logger.Error().Err(err).Msgf("error while rendering steps for agent %s", link.vapusAgent.Name)
			link.eolReason = mpb.EOSReasons_INVALID_PARAMS
			link.statusCode = grpccodes.InvalidArgument
			return err
		} else {
			s.streamServer.Send(&pb.AgentInvokeStreamResponse{
				StatusCode: int64(grpccodes.OK),
				Output: grpcops.BuildVapusStreamResponse(
					mpb.VapusStreamEvents_DATA,
					mpb.ContentFormats_JSON,
					s.SetStepResponse("Rendering", "Input analysed successfully for "+link.vapusAgent.Name),
					0,
					nil,
				),
			})
		}
	}
	validationErr := s.ValidateSteps(link)
	if validationErr != nil {
		return validationErr
	}

	fileBuildingErr := s.FileBuildingQueryStep(link)
	if fileBuildingErr != nil {
		return fileBuildingErr
	}
	s.streamServer.Send(&pb.AgentInvokeStreamResponse{
		StatusCode: int64(grpccodes.OK),
		Output: grpcops.BuildVapusStreamResponse(
			mpb.VapusStreamEvents_DATA,
			mpb.ContentFormats_JSON,
			s.SetStepResponse("File", "File uploading started"),
			0,
			nil,
		),
	})
	if len(link.dataFiles) == 0 {
		link.failedStep = "Uploading file"
		err = utils.ErrFileDataNotFound
		s.logger.Error().Err(err).Msg("error while uploading file")
		link.eolReason = mpb.EOSReasons_DATA_ERROR
		link.statusCode = grpccodes.Internal
	}
	for _, file := range link.dataFiles {
		reqbytes, err := json.Marshal(&models.FileManageOpts{
			Path:        file.Path,
			Data:        file.Data,
			FileName:    file.Name,
			ContentType: file.Format.String(),
		})
		if err != nil {
			s.logger.Err(err).Ctx(s.Ctx).Msg("error while marshalling request spec for file store plugin action")
			return err
		}
		_, err = s.NewPluginActionsAgent(s.Ctx, &pb.PluginActionRequest{
			PluginType: mpb.IntegrationPluginTypes_FILE_STORE.String(),
			Spec:       reqbytes,
		})

		if err != nil {
			link.failedStep = "Uploading file"
			s.logger.Error().Err(err).Msg("error while uploading file")
			link.eolReason = mpb.EOSReasons_SERVER_ERROR
			link.statusCode = grpccodes.Internal
			return err
		} else {
			s.logger.Info().Msg("file uploaded successfully")
			s.streamServer.Send(&pb.AgentInvokeStreamResponse{
				StatusCode: int64(grpccodes.OK),
				Output: grpcops.BuildVapusStreamResponse(
					mpb.VapusStreamEvents_DATA,
					mpb.ContentFormats_JSON,
					s.SetStepResponse("File", "File uploaded successfully"),
					0,
					nil,
				),
			})
		}
	}
	link.eolReason = mpb.EOSReasons_SUCCESSFULL

	return nil
}

func (s *VapusAIAgentThread) EmailerAgent(link *VapusAgentChainLink) error {
	var err error

	if link.userParams.Input != "" {
		err = s.RenderSteps(link)
		if err != nil {
			link.failedStep = "Input analyzation failed for agent " + link.vapusAgent.Name + " failed"
			s.logger.Error().Err(err).Msgf("error while rendering steps for agent %s", link.vapusAgent.Name)
			link.eolReason = mpb.EOSReasons_INVALID_PARAMS
			link.statusCode = grpccodes.InvalidArgument
			return err
		} else {
			s.streamServer.Send(&pb.AgentInvokeStreamResponse{
				StatusCode: int64(grpccodes.OK),
				Output: grpcops.BuildVapusStreamResponse(
					mpb.VapusStreamEvents_DATA,
					mpb.ContentFormats_JSON,
					s.SetStepResponse("Rendering", "Input analysed successfully for"+link.vapusAgent.Name),
					0,
					nil,
				),
			})
		}
	}

	validationErr := s.ValidateSteps(link)
	if validationErr != nil {
		return validationErr
	}

	fileBuildingErr := s.FileBuildingQueryStep(link)
	if fileBuildingErr != nil {
		return fileBuildingErr
	}
	sub, ok := link.requestAgentMap[mpb.AgentStepEnum_AGENTST_EMAIL_SUBJECT.String()]
	if !ok {
		sub = ""
	} else {
		sub = "Context: " + sub
	}
	link.requestAgentMap[mpb.AgentStepEnum_AGENTST_EMAIL_SUBJECT.String()], err = s.GenerateContent(link,
		&pb.ChatRequest{
			ModelNodeId: link.userParams.ModelNodeId,
			Model:       link.userParams.ModelName,
			Temperature: 0.7,
			Contexts: []*mpb.Mapper{
				{
					Key:   "Dataset",
					Value: link.dataproductname,
				},
				{
					Key:   "User Input",
					Value: link.requestAgentMap[mpb.AgentStepEnum_AGENTST_EMAIL_SUBJECT.String()],
				},
			},
			Messages: []*pb.ChatMessageObject{
				{
					Role: pb.AIMessageRoles_USER,
					Content: link.linkSteps[mpb.AgentStepEnum_AGENTST_EMAIL_SUBJECT.String()].Prompt +
						"Input: Generate a generic email subject stating database name and report. " + "\n" + sub,
				}, {
					Role:    pb.AIMessageRoles_SYSTEM,
					Content: "You are a helpful assistant that generates email subject based on the dataset and user input.",
				},
			},
			PromptId: link.linkSteps[mpb.AgentStepEnum_AGENTST_EMAIL_SUBJECT.String()].PromptId,
			Mode:     pb.AIInterfaceMode_P2P,
		}, false)
	if err != nil {
		link.failedStep = mpb.AgentStepEnum_AGENTST_EMAIL_SUBJECT.String()
		s.logger.Error().Err(err).Msg("error while generating email subject")
		link.eolReason = mpb.EOSReasons_SERVER_ERROR
		link.statusCode = grpccodes.Internal
		s.logger.Error().Err(err).Msg("error while generating email subject")
		return err
	} else {
		s.streamServer.Send(&pb.AgentInvokeStreamResponse{
			StatusCode: int64(grpccodes.OK),
			Output: grpcops.BuildVapusStreamResponse(
				mpb.VapusStreamEvents_DATA,
				mpb.ContentFormats_JSON,
				s.SetStepResponse("Email", "Email subject generated successfully"),
				0,
				nil,
			),
		})
	}
	body, ok := link.requestAgentMap[mpb.AgentStepEnum_AGENTST_EMAIL_BODY.String()]
	if !ok {
		body = ""
	} else {
		body = "Context: " + body
	}
	link.requestAgentMap[mpb.AgentStepEnum_AGENTST_EMAIL_BODY.String()], err = s.GenerateContent(link,
		&pb.ChatRequest{
			ModelNodeId: link.userParams.ModelNodeId,
			Model:       link.userParams.ModelName,
			Temperature: 0.7,
			Contexts: []*mpb.Mapper{
				{
					Key:   "Fields in database",
					Value: strings.Join(link.dataFields, ", "),
				},
				{
					Key:   "User Input",
					Value: link.requestAgentMap[mpb.AgentStepEnum_AGENTST_EMAIL_BODY.String()],
				},
			},
			Messages: []*pb.ChatMessageObject{
				{
					Role: pb.AIMessageRoles_USER,
					Content: link.linkSteps[mpb.AgentStepEnum_AGENTST_EMAIL_BODY.String()].Prompt +
						"Input: Generate a generic email body with 40-50 words baseded on data fields and query as mentioned in context where sender name is " +
						s.CtxClaim[encryption.ClaimUserNameKey] +
						"If Recipient's is not mentioned then just say hello or greetings in body." +
						"\n" + body,
				}, {
					Role:    pb.AIMessageRoles_SYSTEM,
					Content: "You are a helpful assistant that generates email body based on data fields and query.",
				},
			},
			PromptId: link.linkSteps[mpb.AgentStepEnum_AGENTST_EMAIL_BODY.String()].PromptId,
			Mode:     pb.AIInterfaceMode_P2P,
		}, false)
	if err != nil {
		link.failedStep = mpb.AgentStepEnum_AGENTST_EMAIL_BODY.String()
		s.logger.Error().Err(err).Msg("error while generating email body")
		link.eolReason = mpb.EOSReasons_SERVER_ERROR
		link.statusCode = grpccodes.Internal
		s.logger.Error().Err(err).Msg("error while generating email body")
		return err
	} else {
		s.streamServer.Send(&pb.AgentInvokeStreamResponse{
			StatusCode: int64(grpccodes.OK),
			Output: grpcops.BuildVapusStreamResponse(
				mpb.VapusStreamEvents_DATA,
				mpb.ContentFormats_JSON,
				s.SetStepResponse("Email", "Email body generated successfully"),
				0,
				nil,
			),
		})
	}
	emailParams := &models.EmailOpts{
		SenderName:  s.CtxClaim[encryption.ClaimUserIdKey],
		Subject:     link.requestAgentMap[mpb.AgentStepEnum_AGENTST_EMAIL_SUBJECT.String()],
		Body:        link.requestAgentMap[mpb.AgentStepEnum_AGENTST_EMAIL_BODY.String()],
		To:          strings.Split(link.requestAgentMap[mpb.AgentStepEnum_AGENTST_EMAIL_RECIEVER.String()], ","),
		From:        s.CtxClaim[encryption.ClaimUserIdKey],
		Attachments: []*models.Attachment{},
	}
	fileName, ok := link.requestAgentMap[mpb.AgentStepEnum_AGENTST_FILENAME.String()]
	if !ok {
		fileName = s.AgentId + "." + strings.ToLower(link.requestAgentMap[mpb.AgentStepEnum_AGENTST_CONTENT_FORMAT.String()])
	} else {
		if !strings.Contains(fileName, "."+link.requestAgentMap[mpb.AgentStepEnum_AGENTST_CONTENT_FORMAT.String()]) {
			fileName = fileName + "." + link.requestAgentMap[mpb.AgentStepEnum_AGENTST_CONTENT_FORMAT.String()]
		}
	}
	for _, file := range link.dataFiles {
		emailParams.Attachments = append(emailParams.Attachments, &models.Attachment{
			Data:        file.Data,
			Format:      file.Format.String(),
			FileName:    file.Name,
			ContentType: globals.FileContentTypes[strings.ToLower(file.Format.String())],
		})
	}
	reqbytes, err := json.Marshal(emailParams)
	if err != nil {
		s.logger.Err(err).Ctx(s.Ctx).Msg("error while marshalling request spec for email plugin action")
		return err
	}
	_, err = s.NewPluginActionsAgent(s.Ctx, &pb.PluginActionRequest{
		PluginType: mpb.IntegrationPluginTypes_EMAIL.String(),
		Spec:       reqbytes,
	})
	if err != nil {
		link.failedStep = "Sending email"
		s.logger.Error().Err(err).Msg("error while sending email")
		link.eolReason = mpb.EOSReasons_SERVER_ERROR
		link.statusCode = grpccodes.Internal
		return err
	} else {
		s.logger.Info().Msg("email sent successfully")
		s.streamServer.Send(&pb.AgentInvokeStreamResponse{
			StatusCode: int64(grpccodes.OK),
			Output: grpcops.BuildVapusStreamResponse(
				mpb.VapusStreamEvents_DATA,
				mpb.ContentFormats_JSON,
				s.SetStepResponse("Email", "Email sent successfully to "+link.requestAgentMap[mpb.AgentStepEnum_AGENTST_EMAIL_RECIEVER.String()]),
				0,
				nil,
			),
		})
	}
	link.eolReason = mpb.EOSReasons_SUCCESSFULL
	return nil
}
