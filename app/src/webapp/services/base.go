package services

import (
	"github.com/rs/zerolog"
	"github.com/vapusdata-oss/aistudio/webapp/clients"
	"github.com/vapusdata-oss/aistudio/webapp/pkgs"
)

type WebappService struct {
	logger      zerolog.Logger
	grpcClients *clients.GrpcClient
	templateDir string
}

var WebappServiceManager *WebappService

func NewWebappSvc() *WebappService {
	svc := &WebappService{}
	svc.logger = pkgs.GetSubDMLogger("webapp", "services")
	svc.templateDir = "datamarketplace"
	clients.InitGrpcClient()
	svc.grpcClients = clients.GrpcClientManager
	return svc
}

func InitWebappSvc() {
	if WebappServiceManager == nil {
		WebappServiceManager = NewWebappSvc()
	}
}
