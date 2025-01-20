package models

import (
	fmt "fmt"

	guuid "github.com/google/uuid"
	encrytion "github.com/vapusdata-oss/aistudio/core/encryption"
	"github.com/vapusdata-oss/aistudio/core/globals"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type Organization struct {
	VapusBase            `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name                 string           `bun:"name,notnull,unique" json:"name,omitempty" yaml:"name"`
	DisplayName          string           `bun:"display_name" json:"displayName,omitempty" yaml:"displayName"`
	Users                []string         `bun:"users,array" json:"users,omitempty" yaml:"users"`
	SecretPasscode       string           `bun:"salt_val" json:"saltVal,omitempty" yaml:"saltVal"`
	AuthnJwtParams       *JWTParams       `bun:"authn_jwt_params,type:jsonb" json:"authnJwtParams,omitempty" yaml:"authnJwtParams"`
	DataSources          []string         `bun:"data_sources,array" json:"dataSources,omitempty" yaml:"dataSources"`
	OrganizationType     string           `bun:"organization_type" json:"organizationType,omitempty" yaml:"organizationType"`
	CatalogIndex         string           `bun:"catalog_index" json:"catalogIndex,omitempty" yaml:"catalogIndex"`
	BackendSecretStorage *BackendStorages `bun:"backend_secret_storage,type:jsonb" json:"backendSecretStorage,omitempty" yaml:"backendSecretStorage"`
	ArtifactStorage      *BackendStorages `bun:"artifact_storage,type:jsonb" json:"artifactStorage,omitempty" yaml:"artifactStorage"`
}

func (m *Organization) SetAccountId(accountId string) {
	if m != nil {
		m.OwnerAccount = accountId
	}
}

func (d *Organization) GetName() string {
	return d.Name
}

func (d *Organization) GetDisplayName() string {
	return d.DisplayName
}

func (d *Organization) GetOrganizationId() string {
	return d.VapusID
}

func (d *Organization) GetUsers() []string {
	return d.Users
}

func (d *Organization) GetSecretPasscode() string {
	return d.SecretPasscode
}

func (d *Organization) GetStatus() string {
	return d.Status
}

func (d *Organization) GetDataSources() []string {
	return d.DataSources
}

func (d *Organization) GetOrganizationType() string {
	return d.OrganizationType
}

func (d *Organization) GetCatalogIndex() string {
	return d.CatalogIndex
}

func (d *Organization) GetBackendSecretStorage() *BackendStorages {
	return d.BackendSecretStorage
}

func (d *Organization) GetArtifactStorage() *BackendStorages {
	return d.ArtifactStorage
}

func (d *Organization) GetAuthnJwtParams() *JWTParams {
	return d.AuthnJwtParams
}

func (dmn *Organization) ConvertToPb() *mpb.Organization {
	if dmn != nil {
		obj := &mpb.Organization{
			Name:                 dmn.Name,
			DisplayName:          dmn.DisplayName,
			OrgId:                dmn.VapusID,
			Users:                dmn.Users,
			SecretPasscode:       &mpb.CredentialSalt{SaltVal: dmn.SecretPasscode},
			Status:               dmn.Status,
			BackendSecretStorage: dmn.BackendSecretStorage.ConvertToPb(),
			ArtifactStorage:      dmn.ArtifactStorage.ConvertToPb(),
			ResourceBase:         dmn.ConvertToPbBase(),
			OrgType:              mpb.OrganizationTypes(mpb.OrganizationTypes_value[dmn.OrganizationType]),
		}
		return obj
	}
	return nil
}

func (dmn *Organization) ConvertFromPb(pb *mpb.Organization) *Organization {
	if pb == nil {
		return nil
	}
	obj := &Organization{
		Name:                 pb.GetName(),
		DisplayName:          pb.GetDisplayName(),
		Users:                pb.GetUsers(),
		SecretPasscode:       pb.GetSecretPasscode().GetSaltVal(),
		BackendSecretStorage: (&BackendStorages{}).ConvertFromPb(pb.GetBackendSecretStorage()),
		ArtifactStorage:      (&BackendStorages{}).ConvertFromPb(pb.GetArtifactStorage()),
		OrganizationType:     pb.GetOrgType().String(),
	}
	return obj
}

func (dmn *Organization) SetOrganizationId() {
	if dmn == nil {
		return
	}
	if dmn.VapusID == "" {
		dmn.VapusID = fmt.Sprintf(globals.DOMAINIDTEMPLATE, guuid.New())
	}
}

func (dmn *Organization) PreSaveCreate(authzClaim map[string]string) {
	if dmn == nil {
		return
	}
	if dmn.CreatedBy == utils.EMPTYSTR {
		dmn.CreatedBy = authzClaim[encrytion.ClaimUserIdKey]
	}
	if dmn.CreatedAt == 0 {
		dmn.CreatedAt = utils.GetEpochTime()
	}
	if dmn.DisplayName == "" {
		dmn.DisplayName = dmn.Name
	}
	dmn.OwnerAccount = authzClaim[encrytion.ClaimAccountKey]
}

func (dn *Organization) PreSaveUpdate(userId string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = userId
	dn.UpdatedAt = utils.GetEpochTime()
}

func (dn *Organization) Delete(userId string) {
	if dn == nil {
		return
	}
	dn.DeletedBy = userId
	dn.DeletedAt = utils.GetEpochTime()
}
