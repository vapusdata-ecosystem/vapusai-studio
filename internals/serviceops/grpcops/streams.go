package grpcops

import (
	"log"

	utils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	"google.golang.org/protobuf/types/known/structpb"
)

func BuildVapusStreamResponse(event mpb.VapusStreamEvents,
	contentType mpb.ContentFormats,
	content any,
	reason mpb.EOSReasons,
	resErr error) *mpb.VapusStreamResponse {
	var contents = ""
	data := &mpb.VapusContentObject{
		ContentType: contentType,
		Final:       &mpb.VapusEOL{},
	}
	switch contentType {
	case mpb.ContentFormats_JSON:
		contents, err := utils.GenericMarshaler(content, contentType.String())
		if err != nil {
			data.Content = "error while processing content"
		} else {
			data.Content = string(contents)
		}
	case mpb.ContentFormats_YAML:
		contents, err := utils.GenericMarshaler(content, contentType.String())
		if err != nil {
			data.Content = "error while processing content"
		} else {
			data.Content = string(contents)
		}
	case mpb.ContentFormats_DATASET:
		vv, ok := content.(*structpb.Struct)
		if !ok {
			data.Content = "error while processing content"
		} else {
			data.Dataset = vv
		}
	case mpb.ContentFormats_PLAIN_TEXT:
		data.Content = content.(string)
	default:
		contents, err := utils.GenericMarshaler(content, contentType.String())
		if err != nil {
			data.Content = "error while processing content"
		} else {
			data.Content = string(contents)
		}
	}

	if event == mpb.VapusStreamEvents_DATA ||
		event == mpb.VapusStreamEvents_DATASET_START ||
		event == mpb.VapusStreamEvents_START ||
		event == mpb.VapusStreamEvents_DATASET_END ||
		event == mpb.VapusStreamEvents_STATE ||
		event == mpb.VapusStreamEvents_FILE_DATA ||
		event == mpb.VapusStreamEvents_RESPONSE_ID {
		data.Final = nil
	} else {
		data.Final.Metadata = func() string {
			if resErr == nil {
				log.Println(string(contents))
				return string(contents)
			}
			return resErr.Error()
		}()
		data.Final.Reason = reason
	}
	obj := &mpb.VapusStreamResponse{
		EventAt: utils.GetEpochTime(),
		Event:   event,
		Data:    data,
	}
	if event == mpb.VapusStreamEvents_FILE_DATA {
		fd, ok := content.(*mpb.FileData)
		if !ok {
			event = mpb.VapusStreamEvents_DATA
		} else {
			obj.Data.Content = ""
			obj.Files = fd
		}
	}
	return obj
}
