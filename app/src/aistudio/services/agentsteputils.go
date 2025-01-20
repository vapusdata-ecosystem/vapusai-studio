package services

import (
	"fmt"
	"log"
	"strings"

	grpccodes "google.golang.org/grpc/codes"
	"gopkg.in/yaml.v3"

	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-oss/vapusdata/aistudio/pkgs"
	utils "github.com/vapusdata-oss/vapusdata/aistudio/utils"
	encryption "github.com/vapusdata-oss/vapusdata/core/encryption"
	models "github.com/vapusdata-oss/vapusdata/core/models"
	grpcops "github.com/vapusdata-oss/vapusdata/core/serviceops/grpcops"
	coreutils "github.com/vapusdata-oss/vapusdata/core/utils"
)

func (s *VapusAIAgentThread) ValidateSteps(link *VapusAgentChainLink) error {
	var err error
	log.Println("steps validation started", link.requestAgentMap)
	for _, step := range link.vapusAgent.Steps {
		log.Println("step id: ", step.Id)
		if step.Required {
			val, ok := link.requestAgentMap[step.Id]
			nameL := strings.ReplaceAll(strings.ToLower(step.Id), "_", " ")
			nameL = strings.ReplaceAll(nameL, "AGENTST", "")
			if !ok {
				s.logger.Error().Msg("error while validating steps, step not found")
				err = fmt.Errorf("step %s is required", nameL)
			}
			if len(val) == 0 {
				s.logger.Error().Msg("error while validating steps, step value not found")
				err = fmt.Errorf("step %s is required", nameL)
			}
		}
		link.linkSteps[step.Id] = step
	}
	if err != nil {
		link.failedStep = "Validation failed: " + err.Error() + ". Try again with more information in the input"
		s.logger.Error().Err(err).Msg("error while validating steps")
		link.eolReason = mpb.EOSReasons_INVALID_PARAMS
		link.statusCode = grpccodes.InvalidArgument
	} else {
		s.logger.Info().Msg("steps validated")
		s.streamServer.Send(&pb.AgentInvokeStreamResponse{
			StatusCode: int64(grpccodes.OK),
			Output: grpcops.BuildVapusStreamResponse(
				mpb.VapusStreamEvents_DATA,
				mpb.ContentFormats_JSON,
				s.SetStepResponse("Validation", "Agent steps validation passed"),
				0,
				nil,
			),
		})
	}
	return err
}

func (s *VapusAIAgentThread) RenderSteps(link *VapusAgentChainLink) error {
	type localTemplate struct {
		values []*models.Steps `yaml:"values"`
	}
	link.userParams.Input = link.userParams.Input + "\n" + `Add Sender name: ` + s.CtxClaim[encryption.ClaimUserNameKey] + ` in the email body with thanks and regards.`
	// link.userParams.Input = link.userParams.Input + "\n" + pkgs.AIAgentDataQueryToolCallContext
	stepsVal, err := s.GenerateContent(link,
		&pb.ChatRequest{
			ModelNodeId: link.userParams.ModelNodeId,
			Model:       link.userParams.ModelName,
			Tools: func() []*mpb.ToolCall {
				var toolCalls []*mpb.ToolCall
				f := models.GetFunctionCallFromString(link.vapusAgent.Settings.ToolCallSchema)
				toolCalls = append(toolCalls, &mpb.ToolCall{
					Type: strings.ToLower(mpb.AIToolCallType_FUNCTION.String()),
					FunctionSchema: &mpb.FunctionCall{
						Name:           f.Name,
						Arguments:      f.GetStringParamSchema(),
						Description:    f.Description,
						RequiredFields: f.RequiredFields,
					},
				})
				return toolCalls
			}(),
			Messages: []*pb.ChatMessageObject{
				{
					Role:    pb.AIMessageRoles_USER,
					Content: link.userParams.Input,
				}, {
					Role:    pb.AIMessageRoles_SYSTEM,
					Content: pkgs.AIAgentSystemMessage,
				},
			},
			Temperature: 0.7,
			Contexts:    []*mpb.Mapper{},
			Mode:        pb.AIInterfaceMode_P2P,
		}, true)
	result := map[string]interface{}{}
	log.Println("stepsVal:+++++++++++++++++++++ ", stepsVal)
	err = yaml.Unmarshal([]byte(stepsVal), &result)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while unmarshalling steps from yaml")
		return utils.ErrInvalidAgentChainParams
	}
	for key, value := range result {
		log.Println("key: ", key, "value: ", value)
		link.requestAgentMap[key] = value.(string)
	}

	return nil
}

func (s *VapusAIAgentThread) GetFileName(link *VapusAgentChainLink) string {
	fileName := link.requestAgentMap[mpb.AgentStepEnum_AGENTST_FILENAME.String()]
	if fileName == "" {
		fileName = s.AgentId + "." + link.requestAgentMap[mpb.AgentStepEnum_AGENTST_CONTENT_FORMAT.String()]
	}
	if !strings.Contains(fileName, "."+link.requestAgentMap[mpb.AgentStepEnum_AGENTST_CONTENT_FORMAT.String()]) {
		fileName = fileName + "." + link.requestAgentMap[mpb.AgentStepEnum_AGENTST_CONTENT_FORMAT.String()]
	}
	return fileName
}

func (s *VapusAIAgentThread) GenerateFileBytes(link *VapusAgentChainLink) error {
	var err error
	format, ok := link.requestAgentMap[mpb.AgentStepEnum_AGENTST_CONTENT_FORMAT.String()]
	if !ok {
		return utils.InvalidContentFormatParam
	}
	for _, field := range link.resultData {
		for k := range field {
			link.dataFields = append(link.dataFields, k)
		}
		break
	}
	fileName := s.GetFileName(link)
	var fBytes []byte
	switch format {
	case mpb.ContentFormats_JSON.String():
		fBytes, err = coreutils.GenericMarshaler(link.resultData, format)
	case mpb.ContentFormats_YAML.String():
		fBytes, err = coreutils.GenericMarshaler(link.resultData, format)
	case mpb.ContentFormats_CSV.String():
		fBytes, err = coreutils.MapArrayCSVMarshaler(link.resultData)
	default:
		link.requestAgentMap[mpb.AgentStepEnum_AGENTST_CONTENT_FORMAT.String()] = mpb.ContentFormats_CSV.String()
		fBytes, err = coreutils.MapArrayCSVMarshaler(link.resultData)
	}
	if err != nil {
		s.logger.Error().Err(err).Msg("error while generating file content for format " + format)
		return err
	}
	var path string = ""
	val, ok := link.requestAgentMap[mpb.AgentStepEnum_AGENTST_FILEPATH.String()]
	if !ok {
		path = ""
	} else {
		path = val
	}
	link.dataFiles = append(link.dataFiles, &mpb.FileData{
		Name:   fileName,
		Data:   fBytes,
		Path:   path,
		Format: mpb.ContentFormats(mpb.ContentFormats_value[link.requestAgentMap[mpb.AgentStepEnum_AGENTST_CONTENT_FORMAT.String()]]),
	})
	return nil
}

func (s *VapusAIAgentThread) GenerateContent(link *VapusAgentChainLink,
	request *pb.ChatRequest, isToolCall bool) (string, error) {
	genAgent, err := s.NewAIModelAgent(s.Ctx, request, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while generating content")
		return "", err
	}
	err = genAgent.Act()
	if err != nil {
		s.logger.Error().Err(err).Msg("error while generating content")
		return "", err
	}
	res := genAgent.GetResult()
	if len(res.Choices) == 0 {
		return "", fmt.Errorf("error while generating content")
	}
	result := ""
	for _, message := range res.Choices {
		if isToolCall {
			if len(message.Messages.ToolCalls) > 0 {
				for _, toolCall := range message.Messages.ToolCalls {
					if strings.ToLower(toolCall.Type) == strings.ToLower(mpb.AIToolCallType_FUNCTION.String()) {
						return toolCall.FunctionSchema.GetArguments(), nil
					}
				}
			}
		} else {
			result = message.Messages.Content
		}
	}
	return result, err
}

func (s *VapusAIAgentThread) FileBuildingQueryStep(link *VapusAgentChainLink) error {
	var nError error
	if len(link.resultData) == 0 {
		if len(link.userParams.Steps) == 0 {
			return utils.ErrFileDataNotFound
		}
		var fType string = ""
		fType, _ = link.requestAgentMap[mpb.AgentStepEnum_AGENTST_CONTENT_FORMAT.String()]
		fPath, ok := link.requestAgentMap[mpb.AgentStepEnum_AGENTST_FILEPATH.String()]
		if !ok {
			fPath = ""
		}
		fileObj := &mpb.FileData{
			Name:   s.GetFileName(link),
			Format: mpb.ContentFormats(mpb.ContentFormats_value[fType]),
			Path:   fPath,
		}
		for _, step := range link.userParams.Steps {
			if step.StepId == mpb.AgentStepEnum_AGENTST_DATASET.String() {
				fileObj.Data = step.Data
				detFType := coreutils.DetectFileTypeFromContent(step.Data)
				if strings.ToLower(fType) != strings.ToLower(detFType) {
					// link.failedStep = mpb.AgentStepEnum_AGENTST_CONTENT_FORMAT.String()
					// link.eolReason = mpb.EOSReasons_DATA_ERROR
					// link.statusCode = grpccodes.InvalidArgument
					s.streamServer.Send(&pb.AgentInvokeStreamResponse{
						StatusCode: int64(grpccodes.OK),
						Output: grpcops.BuildVapusStreamResponse(
							mpb.VapusStreamEvents_DATA,
							mpb.ContentFormats_JSON,
							s.SetStepResponse("File", fmt.Sprintf("file type mismatch with conent format, current format: %s. Converting it into expected format: %s", detFType, fType)),
							0,
							nil,
						),
					})
					result := []map[string]interface{}{}
					err := coreutils.GenericUnMarshaler(step.Data, &result, detFType)
					if err != nil {
						s.logger.Error().Err(err).Msg("error while unmarshalling file data")
						link.eolReason = mpb.EOSReasons_DATA_ERROR
						link.statusCode = grpccodes.InvalidArgument
						return utils.ErrInvalidDataFormat
					}
					link.resultData = result
					err = s.GenerateFileBytes(link)
					if err != nil {
						s.logger.Error().Err(err).Msg("error while generating file bytes")
						link.eolReason = mpb.EOSReasons_DATA_ERROR
						link.statusCode = grpccodes.Internal
						return err
					} else {
						return nil
					}
					// return fmt.Sprintf("file type mismatch with conent format, file type: %s, content format: %s", detFType, fType)
				} else {
					link.dataFiles = append(link.dataFiles, fileObj)
				}
			}
		}

		return nil
	}
	nError = s.GenerateFileBytes(link)
	if nError != nil {
		s.logger.Error().Err(nError).Msg("error while generating file bytes")
		link.eolReason = mpb.EOSReasons_DATA_ERROR
		link.statusCode = grpccodes.Internal
	} else {
		link.failedStep = mpb.AgentStepEnum_AGENTST_CONTENT_FORMAT.String()
		s.streamServer.Send(&pb.AgentInvokeStreamResponse{
			StatusCode: int64(grpccodes.OK),
			Output: grpcops.BuildVapusStreamResponse(
				mpb.VapusStreamEvents_DATA,
				mpb.ContentFormats_JSON,
				s.SetStepResponse("File", "File prepared"),
				0,
				nil,
			),
		})
	}
	return nError
}
