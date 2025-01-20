package plclient

import (
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	aipb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

var SpecMap = map[mpb.RequestObjects]interface{}{
	mpb.RequestObjects_VAPUS_DATAPRODUCTS:               dataproductManagerRequest,
	mpb.RequestObjects_VAPUS_DATASOURCES:                datasourceManagerRequest,
	mpb.RequestObjects_VAPUS_DOMAINS:                    domainManagerRequest,
	mpb.RequestObjects_VAPUS_DATAWORKERS:                dataworkerManagerRequest,
	mpb.RequestObjects_VAPUS_AIMODEL_NODES:              ainodeConfiguratorRequest,
	mpb.RequestObjects_VAPUS_DATA_CONTAINER_DEPLOYMENTS: vdcDeploymentManagerRequest,
	mpb.RequestObjects_VAPUS_DATA_WORKER_DEPLOYMENTS:    dataworkerDeploymentManagerRequest,
	mpb.RequestObjects_VAPUS_DATAMARKETPLACE:            datamarketplaceManagerRequest,
	mpb.RequestObjects_VAPUS_ACCOUNT:                    accountManagerRequest,
	mpb.RequestObjects_VAPUS_AIPROMPTS:                  aiPromptManagerRequest,
}

var accountManagerRequest *pb.AccountManagerRequest = &pb.AccountManagerRequest{
	Spec: &mpb.Account{
		AiAttributes: &mpb.AIAttributes{},
	},
}

var datamarketplaceManagerRequest *pb.DataMarketplaceManagerRequest = &pb.DataMarketplaceManagerRequest{
	Spec: &mpb.DataMarketplace{
		Attributes: &mpb.DataMarketplaceAttributes{},
	},
}

var vdcDeploymentManagerRequest *pb.VDCOrchestratorManagerRequest = &pb.VDCOrchestratorManagerRequest{
	Spec: &mpb.VDCOrchestrator{
		DeploymentSpec: &mpb.VDCK8SDeploymentSpec{},
	},
}

var dataworkerDeploymentManagerRequest *pb.DataWorkerDeploymentManagerRequest = &pb.DataWorkerDeploymentManagerRequest{
	Spec: &mpb.DataWorkerDeployment{
		WorkerDeploymentSpec: &mpb.WorkerK8SDeploymentSpec{},
	},
}

var ainodeConfiguratorRequest *aipb.AIModelNodeConfiguratorRequest = &aipb.AIModelNodeConfiguratorRequest{
	Spec: []*mpb.AIModelNode{
		{
			Attributes: &mpb.AIModelNodeAttributes{
				NetworkParams: &mpb.AIModelNodeNetworkParams{
					Credentials: &mpb.GenericCredentialObj{},
				},
			},
		},
	},
}

var aiPromptManagerRequest *aipb.PromptManagerRequest = &aipb.PromptManagerRequest{
	Spec: []*mpb.AIModelPrompt{
		{
			Prompt: &mpb.Prompt{
				Sample: &mpb.Sample{},
			},
		},
	},
}

var dataworkerManagerRequest *pb.DataWorkerManagerRequest = &pb.DataWorkerManagerRequest{
	Spec: &mpb.VapusDataWorker{
		DataWorkerType: mpb.DataWorkerType_ETL.String(),
		WorkerEngine: &mpb.WorkerEngine{
			Extracter: []*mpb.ExtractionRule{{
				DataRule: &mpb.ExtractionDataRule{
					RawQuery:         &mpb.RawQueryOpts{},
					CustomQueryParam: &mpb.ExtractedDataCustomQuery{},
				},
			},
			},
			Loader: []*mpb.LoadingRule{{
				DataRule: &mpb.LoadingDataRule{},
			},
			},
			Transformers: &mpb.Transformer{
				Globals: &mpb.GlobalTransformer{
					Vars:       []*mpb.GlobalVarTransformer{},
					Classified: []*mpb.ClassifiedTransformer{},
				},
				Steps: []*mpb.TransformerSteps{
					{
						Column: &mpb.ColumnTransformer{
							Add: []*mpb.ColumnProperties{
								{}},
							Rename: []*mpb.ColumnProperties{
								{}},
							Drop: []*mpb.ColumnProperties{
								{}},
							Update: []*mpb.ColumnProperties{
								{}},
						},
						Row: &mpb.RowTransformer{
							Drop: []*mpb.RowProperties{
								{},
							},
							Update: []*mpb.RowProperties{
								{},
							},
						},
					},
				},
			},
		},
	},
}

var dataproductManagerRequest *pb.DataProductManagerRequest = &pb.DataProductManagerRequest{
	Spec: &mpb.DataProduct{
		Labels: []*mpb.Mapper{{}},
		AccessPolicies: &mpb.AccessControlPolicy{
			UserRules: []*mpb.DataProductUsers{{}},
			RoleRules: []*mpb.DataProductUserRole{{}},
		},
		Contract: &mpb.DataProductContract{
			Readiness: &mpb.Readiness{},
			IoPorts: &mpb.DataProductIO{
				InputPorts: []*mpb.ProductDataSource{{
					Datarules: &mpb.ProductDataRules{},
				}},
				ProductOutputPorts: []*mpb.ProductOutputPorts{{
					Params: []*mpb.Mapper{{}},
				}},
				QueryPrompts: []*mpb.QueryPrompts{{}},
			},
			Governance: &mpb.DataGovernance{
				AttributePolicies: []*mpb.AttributePolicies{
					{
						Filters: []*mpb.AttributeFilter{{}},
					},
				},
				TransformerPolicies: []*mpb.DataProductTransformer{
					{
						Transformers: &mpb.Transformer{
							Globals: &mpb.GlobalTransformer{
								Vars:       []*mpb.GlobalVarTransformer{},
								Classified: []*mpb.ClassifiedTransformer{},
							},
						},
					},
				},
			},
		},
	},
}

var datasourceManagerRequest *pb.DataSourceManagerRequest = &pb.DataSourceManagerRequest{
	Spec: &mpb.DataSource{
		NetParams: &mpb.DataSourceNetParams{
			DsCreds: []*mpb.DataSourceCreds{
				{
					Credentials: &mpb.GenericCredentialObj{
						AwsCreds:   &mpb.AWSCreds{},
						GcpCreds:   &mpb.GCPCreds{},
						AzureCreds: &mpb.AzureCreds{},
					},
				},
			},
		},
		Attributes:    &mpb.DataSourceAttributes{},
		Tags:          []*mpb.Mapper{{}},
		SharingParams: &mpb.DataSourceSharingParams{},

		MetaSyncSchedule: &mpb.SyncSchedule{},
	},
}

var dataSourceCreds *mpb.GenericCredentialObj = &mpb.GenericCredentialObj{
	AwsCreds:   &mpb.AWSCreds{},
	GcpCreds:   &mpb.GCPCreds{},
	AzureCreds: &mpb.AzureCreds{},
}

var domainManagerRequest *pb.DomainManagerRequest = &pb.DomainManagerRequest{
	Spec: &mpb.Domain{
		SecretPasscode: &mpb.CredentialSalt{},
		Attributes: &mpb.DomainAttributes{
			AuthnJwtParams: &mpb.JWTParams{},
		},
		BackendSecretStorage: &mpb.BackendStorages{},
		ArtifactStorage:      &mpb.BackendStorages{},
		DomainArtifacts: []*mpb.DomainArtifacts{
			{
				Artifacts: []*mpb.PlatformArtifact{},
			},
		},
		DataProductInfraPlatform: []*mpb.K8SInfraParams{
			{
				Credentials: &mpb.GenericCredentialObj{
					AwsCreds:   &mpb.AWSCreds{},
					GcpCreds:   &mpb.GCPCreds{},
					AzureCreds: &mpb.AzureCreds{},
				},
			},
		},
	},
}
