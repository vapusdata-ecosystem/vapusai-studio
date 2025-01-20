package models

import mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"

type GenericCredentialModel struct {
	Username     string      `json:"username,omitempty" yaml:"username"`
	Password     string      `json:"password,omitempty" yaml:"password"`
	ApiToken     string      `json:"apiToken,omitempty" yaml:"apiToken"`
	ApiTokenType string      `json:"apiTokenType,omitempty" yaml:"apiTokenType"`
	AwsCreds     *AWSCreds   `json:"awsCreds,omitempty" yaml:"awsCreds"`
	GcpCreds     *GCPCreds   `json:"gcpCreds,omitempty" yaml:"gcpCreds"`
	AzureCreds   *AzureCreds `json:"azureCreds,omitempty" yaml:"azureCreds"`
	ClientId     string      `json:"clientId,omitempty" yaml:"clientId"`
	ClientSecret string      `json:"clientSecret,omitempty" yaml:"clientSecret"`
}

func (m *GenericCredentialModel) Reset() {
	m = nil
}

func (m *GenericCredentialModel) ConvertToPb() *mpb.GenericCredentialObj {
	if m != nil {
		return &mpb.GenericCredentialObj{
			Username:     m.Username,
			Password:     m.Password,
			ApiToken:     m.ApiToken,
			ApiTokenType: mpb.ApiTokenType(mpb.ApiTokenType_value[m.ApiTokenType]),
			AwsCreds:     m.AwsCreds.ConvertToPb(),
			GcpCreds:     m.GcpCreds.ConvertToPb(),
			AzureCreds:   m.AzureCreds.ConvertToPb(),
			ClientId:     m.ClientId,
			ClientSecret: m.ClientSecret,
		}
	}
	return nil
}

func (m *GenericCredentialModel) ConvertFromPb(pb *mpb.GenericCredentialObj) *GenericCredentialModel {
	if pb != nil {
		return &GenericCredentialModel{
			Username:     pb.GetUsername(),
			Password:     pb.GetPassword(),
			ApiToken:     pb.GetApiToken(),
			ApiTokenType: pb.GetApiTokenType().String(),
			AwsCreds:     (&AWSCreds{}).ConvertFromPb(pb.GetAwsCreds()),
			GcpCreds:     (&GCPCreds{}).ConvertFromPb(pb.GetGcpCreds()),
			AzureCreds:   (&AzureCreds{}).ConvertFromPb(pb.GetAzureCreds()),
			ClientId:     pb.GetClientId(),
			ClientSecret: pb.GetClientSecret(),
		}
	}
	return nil
}

func (m *GenericCredentialModel) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *GenericCredentialModel) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *GenericCredentialModel) GetApiToken() string {
	if m != nil {
		return m.ApiToken
	}
	return ""
}

func (m *GenericCredentialModel) GetApiTokenType() string {
	if m != nil {
		return m.ApiTokenType
	}
	return ""
}

func (m *GenericCredentialModel) GetAwsCreds() *AWSCreds {
	if m != nil {
		return m.AwsCreds
	}
	return nil
}

func (m *GenericCredentialModel) GetGcpCreds() *GCPCreds {
	if m != nil {
		return m.GcpCreds
	}
	return nil
}

func (m *GenericCredentialModel) GetAzureCreds() *AzureCreds {
	if m != nil {
		return m.AzureCreds
	}
	return nil
}

type AWSCreds struct {
	AccessKeyId     string `json:"accessKeyId,omitempty" yaml:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey,omitempty" yaml:"secretAccessKey"`
	Region          string `json:"region,omitempty" yaml:"region"`
	SessionToken    string `json:"sessionToken,omitempty" yaml:"sessionToken"`
	RoleArn         string `json:"roleArn,omitempty" yaml:"roleArn"`
}

func (m *AWSCreds) ConvertToPb() *mpb.AWSCreds {
	if m != nil {
		return &mpb.AWSCreds{
			AccessKeyId:     m.AccessKeyId,
			SecretAccessKey: m.SecretAccessKey,
			Region:          m.Region,
			SessionToken:    m.SessionToken,
			RoleArn:         m.RoleArn,
		}
	}
	return nil
}

func (m *AWSCreds) ConvertFromPb(pb *mpb.AWSCreds) *AWSCreds {
	if pb != nil {
		return &AWSCreds{
			AccessKeyId:     pb.GetAccessKeyId(),
			SecretAccessKey: pb.GetSecretAccessKey(),
			Region:          pb.GetRegion(),
			SessionToken:    pb.GetSessionToken(),
			RoleArn:         pb.GetRoleArn(),
		}
	}
	return nil
}

type GCPCreds struct {
	ServiceAccountKey string `json:"serviceAccountKey,omitempty" yaml:"serviceAccountKey"`
	Base64Encoded     bool   `json:"base64Encoded,omitempty" yaml:"base64Encoded"`
	Region            string `json:"region,omitempty" yaml:"region"`
	ProjectId         string `json:"projectId,omitempty" yaml:"projectId"`
	Zone              string `json:"zone,omitempty" yaml:"zone"`
}

func (m *GCPCreds) ConvertToPb() *mpb.GCPCreds {
	if m != nil {
		return &mpb.GCPCreds{
			ServiceAccountKey: m.ServiceAccountKey,
			Base64Encoded:     m.Base64Encoded,
			Region:            m.Region,
			ProjectId:         m.ProjectId,
			Zone:              m.Zone,
		}
	}
	return nil
}

func (m *GCPCreds) ConvertFromPb(pb *mpb.GCPCreds) *GCPCreds {
	if pb != nil {
		return &GCPCreds{
			ServiceAccountKey: pb.GetServiceAccountKey(),
			Base64Encoded:     pb.GetBase64Encoded(),
			Region:            pb.GetRegion(),
			ProjectId:         pb.GetProjectId(),
			Zone:              pb.GetZone(),
		}
	}
	return nil
}

type AzureCreds struct {
	TenantId     string `json:"tenantId,omitempty" yaml:"tenantId"`
	ClientId     string `json:"clientId,omitempty" yaml:"clientId"`
	ClientSecret string `json:"clientSecret,omitempty" yaml:"clientSecret"`
}

func (m *AzureCreds) ConvertToPb() *mpb.AzureCreds {
	if m != nil {
		return &mpb.AzureCreds{
			TenantId:     m.TenantId,
			ClientId:     m.ClientId,
			ClientSecret: m.ClientSecret,
		}
	}
	return nil
}

func (m *AzureCreds) ConvertFromPb(pb *mpb.AzureCreds) *AzureCreds {
	if pb != nil {
		return &AzureCreds{
			TenantId:     pb.GetTenantId(),
			ClientId:     pb.GetClientId(),
			ClientSecret: pb.GetClientSecret(),
		}
	}
	return nil
}
