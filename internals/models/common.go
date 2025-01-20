package models

import (
	"slices"

	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type SupportedPackageTypes struct {
	TypeId           string `json:"typeId"`
	PackageType      string `json:"packageType"`
	PackageExtension string `json:"packageExtension"`
}

type VapusBase struct {
	CreatedAt    int64    `json:"createdAt" bun:"created_at"`
	CreatedBy    string   `json:"createdBy" bun:"created_by"`
	DeletedAt    int64    `json:"deletedAt" bun:"deleted_at,nullzero"`
	DeletedBy    string   `json:"deletedBy" bun:"deleted_by"`
	UpdatedAt    int64    `json:"updatedAt" bun:"updated_at,nullzero"`
	UpdatedBy    string   `json:"updatedBy" bun:"updated_by"`
	OwnerAccount string   `json:"ownerAccount" bun:"owner_account"`
	ID           int64    `json:"id" bun:"id,pk,autoincrement"`
	VapusID      string   `json:"vId" bun:"vapus_id,unique,notnull"`
	LastAuditID  string   `json:"lastAuditId" bun:"last_audit_id"`
	ErrorLogs    []string `json:"errorLogs" bun:"error_logs"`
	// Labels       []string `json:"labels" bun:"labels,array"`
	Organization string   `json:"organization" bun:"organization"`
	Status       string   `json:"status" bun:"status"`
	Editors      []string `json:"editors" bun:"editors,array"`
	Scope        string   `json:"scope" bun:"scope"`
}

func (dm *VapusBase) PreSaveVapusBase(authzClaim map[string]string) {
	if dm.CreatedBy == utils.EMPTYSTR {
		dm.CreatedBy = authzClaim[encryption.ClaimUserIdKey]
	}
	if dm.CreatedAt == 0 {
		dm.CreatedAt = utils.GetEpochTime()
	}
	if dm.OwnerAccount == utils.EMPTYSTR {
		dm.OwnerAccount = authzClaim[encryption.ClaimAccountKey]
	}
	if dm.VapusID == utils.EMPTYSTR {
		dm.VapusID = utils.GetUUID()
	}
	dm.Organization = authzClaim[encryption.ClaimOrganizationKey]

	if dm.Editors == nil {
		dm.Editors = []string{authzClaim[encryption.ClaimUserIdKey]}
	} else if !slices.Contains(dm.Editors, authzClaim[encryption.ClaimUserIdKey]) {
		dm.Editors = append(dm.Editors, authzClaim[encryption.ClaimUserIdKey])
	}
}

func (dn *VapusBase) PreDeleteVapusBase(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.DeletedBy = authzClaim[encryption.ClaimUserIdKey]
	dn.DeletedAt = utils.GetEpochTime()
	dn.Status = mpb.CommonStatus_DELETED.String()
}

func (dn *VapusBase) ConvertToPbBase() *mpb.VapusBase {
	if dn == nil {
		return &mpb.VapusBase{}
	}
	return &mpb.VapusBase{
		CreatedAt: dn.CreatedAt,
		CreatedBy: dn.CreatedBy,
		DeletedAt: dn.DeletedAt,
		DeletedBy: dn.DeletedBy,
		UpdatedAt: dn.UpdatedAt,
		UpdatedBy: dn.UpdatedBy,
	}
}

type JWTParams struct {
	Name                string `json:"name"`
	PublicJwtKey        string `json:"publicJwtKey"`
	PrivateJwtKey       string `json:"privateJwtKey"`
	VId                 string `json:"vId"`
	SigningAlgorithm    string `json:"signingAlgorithm"`
	IsAlreadyInSecretBs bool   `json:"isAlreadyInSecretBs"`
	Status              string `json:"status"`
	GenerateInStudio    bool   `json:"generateInStudio"`
}

func (j *JWTParams) Reset() {
	j = nil
}

func (j *JWTParams) GetName() string {
	if j != nil && j.Name != "" {
		return j.Name
	}
	return ""
}

func (j *JWTParams) ConvertToPb() *mpb.JWTParams {
	if j != nil {
		return &mpb.JWTParams{
			Name:                j.Name,
			PublicJwtKey:        j.PublicJwtKey,
			PrivateJwtKey:       j.PrivateJwtKey,
			VId:                 j.VId,
			SigningAlgorithm:    mpb.EncryptionAlgo(mpb.EncryptionAlgo_value[j.SigningAlgorithm]),
			IsAlreadyInSecretBs: j.IsAlreadyInSecretBs,
			Status:              j.Status,
			GenerateInStudio:    j.GenerateInStudio,
		}
	}
	return nil
}

func (j *JWTParams) ConvertFromPb(pb *mpb.JWTParams) *JWTParams {
	if pb == nil {
		return nil
	}
	return &JWTParams{
		Name:                pb.GetName(),
		PublicJwtKey:        pb.GetPublicJwtKey(),
		PrivateJwtKey:       pb.GetPrivateJwtKey(),
		VId:                 pb.GetVId(),
		SigningAlgorithm:    pb.GetSigningAlgorithm().String(),
		IsAlreadyInSecretBs: pb.GetIsAlreadyInSecretBs(),
		Status:              pb.GetStatus(),
		GenerateInStudio:    pb.GetGenerateInStudio(),
	}
}

type BackendStorages struct {
	BesType       string         `json:"besType"`
	BesOnboarding string         `json:"besOnboarding"`
	BesService    string         `json:"besService"`
	NetParams     *NetworkParams `json:"netParams"`
	Status        string         `json:"status"`
	BesEngine     string         `json:"besEngine"`
}

func (b *BackendStorages) ConvertToPb() *mpb.BackendStorages {
	if b != nil {
		return &mpb.BackendStorages{
			BesType:       mpb.DataSourceType(mpb.DataSourceType_value[b.BesType]),
			BesOnboarding: mpb.BackendStorageOnboarding(mpb.BackendStorageOnboarding_value[b.BesOnboarding]),
			BesService:    mpb.StorageService(mpb.StorageService_value[b.BesService]),
			NetParams:     b.NetParams.ConvertToPb(),
			Status:        b.Status,
			BesEngine:     mpb.StorageEngine(mpb.StorageEngine_value[b.BesEngine]),
		}
	}
	return nil
}

func (b *BackendStorages) ConvertFromPb(pb *mpb.BackendStorages) *BackendStorages {
	if pb == nil {
		return nil
	}
	return &BackendStorages{
		BesType:       pb.GetBesType().String(),
		BesOnboarding: pb.GetBesOnboarding().String(),
		BesService:    pb.GetBesService().String(),
		NetParams:     (&NetworkParams{}).ConvertFromPb(pb.GetNetParams()),
		Status:        pb.GetStatus(),
		BesEngine:     pb.GetBesEngine().String(),
	}
}

type AuthnOIDC struct {
	Callback            string `json:"callback"`
	ClientId            string `json:"clientId"`
	ClientSecret        string `json:"clientSecret"`
	VId                 string `json:"vId"`
	IsAlreadyInSecretBs bool   `json:"isAlreadyInSecretBs"`
	Status              string `json:"status"`
}

func (a *AuthnOIDC) ConvertToPb() *mpb.AuthnOIDC {
	if a != nil {
		return &mpb.AuthnOIDC{
			Callback:            a.Callback,
			ClientId:            a.ClientId,
			ClientSecret:        a.ClientSecret,
			VId:                 a.VId,
			IsAlreadyInSecretBs: a.IsAlreadyInSecretBs,
			Status:              a.Status,
		}
	}
	return nil
}

func (a *AuthnOIDC) ConvertFromPb(pb *mpb.AuthnOIDC) *AuthnOIDC {
	if pb == nil {
		return nil
	}
	return &AuthnOIDC{
		Callback:            pb.GetCallback(),
		ClientId:            pb.GetClientId(),
		ClientSecret:        pb.GetClientSecret(),
		VId:                 pb.GetVId(),
		IsAlreadyInSecretBs: pb.GetIsAlreadyInSecretBs(),
		Status:              pb.GetStatus(),
	}
}

type DigestVal struct {
	Algo   string `json:"algo"`
	Digest string `json:"digest"`
}

func (d *DigestVal) GetAlgo() string {
	if d != nil && d.Algo != "" {
		return d.Algo
	}
	return ""
}

func (d *DigestVal) GetDigest() string {
	if d != nil && d.Digest != "" {
		return d.Digest
	}
	return ""
}

func (d *DigestVal) ConvertToPb() *mpb.DigestVal {
	if d != nil {
		return &mpb.DigestVal{
			Algo:   mpb.HashAlgos(mpb.HashAlgos_value[d.Algo]),
			Digest: d.GetDigest(),
		}
	}
	return nil
}

func (d *DigestVal) ConvertFromPb(pb *mpb.DigestVal) *DigestVal {
	if pb == nil {
		return nil
	}
	return &DigestVal{
		Algo:   pb.GetAlgo().String(),
		Digest: pb.GetDigest(),
	}
}

type Mapper struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (m *Mapper) GetKey() string {
	if m != nil && m.Key != "" {
		return m.Key
	}
	return ""
}

func (m *Mapper) GetValue() string {
	if m != nil && m.Value != "" {
		return m.Value
	}
	return ""
}

func (m *Mapper) ConvertToPb() *mpb.Mapper {
	if m != nil {
		return &mpb.Mapper{
			Key:   m.GetKey(),
			Value: m.GetValue(),
		}
	}
	return nil
}

func (m *Mapper) ConvertFromPb(pb *mpb.Mapper) *Mapper {
	if pb == nil {
		return nil
	}
	return &Mapper{
		Key:   pb.GetKey(),
		Value: pb.GetValue(),
	}
}

func MapperSliceToPb(mappers []*Mapper) []*mpb.Mapper {
	if mappers == nil {
		return nil
	}
	var list []*mpb.Mapper
	for _, m := range mappers {
		list = append(list, m.ConvertToPb())
	}
	return list
}

func MapperSliceFromPb(pbs []*mpb.Mapper) []*Mapper {
	if pbs == nil {
		return nil
	}
	var list []*Mapper
	for _, pb := range pbs {
		list = append(list, (&Mapper{}).ConvertFromPb(pb))
	}
	return list
}

type BaseIdentifier struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

func (b *BaseIdentifier) GetName() string {
	if b != nil && b.Name != "" {
		return b.Name
	}
	return ""
}

func (b *BaseIdentifier) GetType() string {
	if b != nil && b.Type != "" {
		return b.Type
	}
	return ""
}

func (b *BaseIdentifier) GetIdentifier() string {
	if b != nil && b.Identifier != "" {
		return b.Identifier
	}
	return ""
}

func (b *BaseIdentifier) ConvertToPb() *mpb.BaseIdentifier {
	if b != nil {
		return &mpb.BaseIdentifier{
			Name:       b.GetName(),
			Type:       b.GetType(),
			Identifier: b.GetIdentifier(),
		}
	}
	return nil
}

func (b *BaseIdentifier) ConvertFromPb(pb *mpb.BaseIdentifier) *BaseIdentifier {
	if pb == nil {
		return nil
	}
	return &BaseIdentifier{
		Name:       pb.GetName(),
		Type:       pb.GetType(),
		Identifier: pb.GetIdentifier(),
	}
}

type SyncSchedule struct {
	Frequency string `json:"frequency"`
	Value     int64  `json:"value"`
	Limit     int64  `json:"limit"`
}

func (s *SyncSchedule) GetFrequency() string {
	if s != nil && s.Frequency != "" {
		return s.Frequency
	}
	return ""
}

func (s *SyncSchedule) GetValue() int64 {
	if s != nil {
		return s.Value
	}
	return 0
}

func (s *SyncSchedule) GetLimit() int64 {
	if s != nil {
		return s.Limit
	}
	return 0
}

func (s *SyncSchedule) ConvertToPb() *mpb.SyncSchedule {
	if s != nil {
		return &mpb.SyncSchedule{
			Frequency: mpb.Frequency(mpb.Frequency_value[s.Frequency]),
			Value:     s.GetValue(),
			Limit:     s.GetLimit(),
		}
	}
	return nil
}

func (s *SyncSchedule) ConvertFromPb(pb *mpb.SyncSchedule) *SyncSchedule {
	if pb == nil {
		return nil
	}
	return &SyncSchedule{
		Frequency: pb.GetFrequency().String(),
		Value:     pb.GetValue(),
		Limit:     pb.GetLimit(),
	}
}

func GetMapperPbList(mappers []*Mapper) []*mpb.Mapper {
	if mappers == nil {
		return nil
	}
	var list []*mpb.Mapper
	for _, m := range mappers {
		list = append(list, m.ConvertToPb())
	}
	return list
}

func GetMapperObjList(pbs []*mpb.Mapper) []*Mapper {
	if pbs == nil {
		return nil
	}
	var list []*Mapper
	for _, pb := range pbs {
		list = append(list, (&Mapper{}).ConvertFromPb(pb))
	}
	return list
}

func SetVapusId(name string) string {
	return utils.SlugifyBase(name)
}
