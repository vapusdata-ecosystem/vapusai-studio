package models

import (
	"log"
	"slices"
	"strings"
	"time"

	guuid "github.com/google/uuid"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type Users struct {
	VapusBase         `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	DisplayName       string                  `bun:"display_name" json:"displayName,omitempty" yaml:"displayName"`
	UserId            string                  `bun:"user_id,notnull,unique" json:"userId,omitempty" yaml:"userId"`
	Email             string                  `bun:"email,notnull,unique" json:"email,omitempty" yaml:"email"`
	OrganizationRoles []*UserOrganizationRole `bun:"organization_roles,type:jsonb[]" json:"organizationRoles,omitempty" yaml:"organizationRoles"`
	InviteId          string                  `bun:"invite_id" json:"inviteId,omitempty" yaml:"inviteId"`
	StudioRoles       []string                `bun:"platform_roles,array" json:"platformRole,omitempty" yaml:"platformRole"`
	InvitedOn         int64                   `bun:"invited_on" json:"invitedOn,omitempty" yaml:"invitedOn"`
	InviteExpiresOn   int64                   `bun:"invite_expires_on" json:"inviteExpiresOn,omitempty" yaml:"inviteExpiresOn"`
	FirstName         string                  `bun:"first_name" json:"firstName,omitempty" yaml:"firstName"`
	LastName          string                  `bun:"last_name" json:"lastName,omitempty" yaml:"lastName"`
	InvitedType       string                  `bun:"invited_type" json:"invitedType,omitempty" yaml:"invitedType"`
	StudioPolicies    []string                `bun:"platform_policies,array" json:"platformPolicies,omitempty" yaml:"platformPolicies"`
	Profile           *UserProfile            `bun:"profile,type:jsonb" json:"profile,omitempty" yaml:"profile"`
}

func (dm *Users) ValidateJwtClaim(claimCtx map[string]string) bool {
	if dm == nil {
		return false
	}
	if dm.UserId != claimCtx[encryption.ClaimUserIdKey] {
		log.Println(dm.UserId, "--------------1")
		log.Println(claimCtx[encryption.ClaimUserIdKey], "--------------2")
		return false
	}
	// if dm.OwnerAccount != claimCtx[encryption.ClaimAccountKey] {
	// 	log.Println(dm.OwnerAccount, "--------------3")
	// 	log.Println(claimCtx[encryption.ClaimAccountKey], "--------------4")
	// 	return false
	// }
	for _, dm := range dm.GetOrganizationRoles() {
		if claimCtx[encryption.ClaimOrganizationKey] == dm.OrganizationId {
			claimRoles := strings.Split(claimCtx[encryption.ClaimOrganizationRolesKey], "|")
			log.Println(claimRoles, claimCtx[encryption.ClaimOrganizationRolesKey])
			if len(claimRoles) == 0 {
				return false
			} else if len(claimRoles) == 1 {
				return slices.Contains(dm.RoleArns, claimCtx[encryption.ClaimOrganizationRolesKey])
			} else {
				valid := false
				for _, r := range claimRoles {
					if !slices.Contains(dm.RoleArns, r) {
						valid = false
						break
					} else {
						valid = true
					}
				}
				return valid
			}
		}
	}
	return false
}

func (dm *Users) IsValidUserByOrganization(organization string) bool {
	if dm == nil {
		return false
	}
	for _, dm := range dm.GetOrganizationRoles() {
		if organization == dm.OrganizationId {
			return true
		}
	}
	return false
}

func (dm *Users) GetRoleByOrganization(organization string) []string {
	if dm == nil {
		return []string{}
	}
	for _, dm := range dm.GetOrganizationRoles() {
		if organization == dm.OrganizationId {
			return dm.RoleArns
		}
	}
	return []string{}
}

func (dm *Users) SetAccountId(accountId string) {
	if dm != nil {
		dm.OwnerAccount = accountId
	}
}

func (dm *Users) GetDisplayName() string {
	if dm.DisplayName == utils.EMPTYSTR {
		return dm.FirstName + dm.LastName
	}
	return dm.DisplayName
}
func (u *Users) GetUserId() string {
	if u != nil {
		return u.UserId
	}
	return ""
}

func (u *Users) GetEmail() string {
	if u != nil {
		return u.Email
	}
	return ""
}

func (u *Users) GetOrganizationRoles() []*UserOrganizationRole {
	if u != nil {
		return u.OrganizationRoles
	}
	return nil
}

func (u *Users) GetProfile() *UserProfile {
	if u != nil {
		return u.Profile
	}
	return nil
}

func (u *Users) GetInviteId() string {
	if u != nil {
		return u.InviteId
	}
	return ""
}

func (u *Users) GetStudioRole() []string {
	if u != nil {
		return u.StudioRoles
	}
	return []string{}
}

func (u *Users) GetStatus() string {
	if u != nil {
		return u.Status
	}
	return ""
}

func (u *Users) GetInvitedOn() int64 {
	if u != nil {
		return u.InvitedOn
	}
	return 0
}

func (u *Users) GetInviteExpiresOn() int64 {
	if u != nil {
		return u.InviteExpiresOn
	}
	return 0
}

func (u *Users) GetFirstName() string {
	if u != nil {
		return u.FirstName
	}
	return ""
}

func (u *Users) GetLastName() string {
	if u != nil {
		return u.LastName
	}
	return ""
}

func (u *Users) GetInvitedType() string {
	if u != nil {
		return u.InvitedType
	}
	return ""
}

func (u *Users) GetStudioPolicies() []string {
	if u != nil {
		return u.StudioPolicies
	}
	return nil
}

func (m *Users) ConvertToPb(organization string) *mpb.User {
	if m != nil {
		return &mpb.User{
			DisplayName: m.DisplayName,
			UserId:      m.UserId,
			Email:       m.Email,
			OrganizationRoles: func(s []*UserOrganizationRole) (pbs []*mpb.UserOrganizationRole) {
				for _, v := range s {
					// if v.OrganizationId == organization {
					pbs = append(pbs, v.ConvertToPb())
					// }
				}
				return

			}(m.OrganizationRoles),
			InviteId: m.InviteId,
			StudioRoles: func() (roles []mpb.StudioUserRoles) {
				for _, r := range m.StudioRoles {
					roles = append(roles, mpb.StudioUserRoles(mpb.StudioUserRoles_value[r]))
				}
				return roles
			}(),
			Status:          m.Status,
			InvitedOn:       m.InvitedOn,
			InviteExpiresOn: m.InviteExpiresOn,
			FirstName:       m.FirstName,
			LastName:        m.LastName,
			InvitedType:     mpb.UserInviteType(mpb.UserInviteType_value[m.InvitedType]),
			StudioPolicies:  m.StudioPolicies,
			Profile:         m.Profile.ConvertToPb(),
		}
	}
	return nil
}

func (m *Users) ConvertFromPb(pb *mpb.User) *Users {
	if pb != nil {
		return &Users{
			DisplayName: pb.GetDisplayName(),
			UserId:      pb.GetUserId(),
			Email:       pb.GetEmail(),
			OrganizationRoles: func(s []*mpb.UserOrganizationRole) (pbs []*UserOrganizationRole) {
				for _, v := range s {
					pbs = append(pbs, (&UserOrganizationRole{}).ConvertFromPb(v))
				}
				return pbs

			}(pb.GetOrganizationRoles()),
			InviteId: pb.GetInviteId(),
			StudioRoles: func() (roles []string) {
				for _, r := range pb.GetStudioRoles() {
					roles = append(roles, r.String())
				}
				return roles
			}(),
			InvitedOn:       pb.GetInvitedOn(),
			InviteExpiresOn: pb.GetInviteExpiresOn(),
			FirstName:       pb.GetFirstName(),
			LastName:        pb.GetLastName(),
			InvitedType:     pb.GetInvitedType().String(),
			StudioPolicies:  pb.GetStudioPolicies(),
			Profile:         (&UserProfile{}).ConvertFromPb(pb.GetProfile()),
		}
	}
	return nil
}

type UserOrganizationRole struct {
	OrganizationId string   `json:"organizationId" yaml:"organizationId"`
	RoleArns       []string `json:"roleArns" yaml:"roleArns"`
	InvitedOn      int64    `json:"invitedOn" yaml:"invitedOn"`
	ValidTill      int64    `json:"validTill" yaml:"validTill"`
	Policies       []string `json:"policies" yaml:"policies"`
	IsDefault      bool     `json:"isDefault" yaml:"isDefault"`
}

func (m *UserOrganizationRole) ConvertToPb() *mpb.UserOrganizationRole {
	if m != nil {
		return &mpb.UserOrganizationRole{
			OrganizationId: m.OrganizationId,
			Role:           m.RoleArns,
			InvitedOn:      m.InvitedOn,
			ValidTill:      m.ValidTill,
			Policies:       m.Policies,
			IsDefault:      m.IsDefault,
		}
	}
	return nil
}

func (m *UserOrganizationRole) ConvertFromPb(pb *mpb.UserOrganizationRole) *UserOrganizationRole {
	if pb != nil {
		return &UserOrganizationRole{
			OrganizationId: pb.GetOrganizationId(),
			RoleArns:       pb.GetRole(),
			InvitedOn:      pb.GetInvitedOn(),
			ValidTill:      pb.GetValidTill(),
			Policies:       pb.GetPolicies(),
			IsDefault:      pb.GetIsDefault(),
		}
	}
	return nil
}

type UserProfile struct {
	Addresses   []*Address `json:"addresses"`
	Avatar      string     `json:"avatar"`
	Description string     `json:"description"`
	Gender      string     `json:"gender"`
	DateOfBirth string     `json:"dateOfBirth"`
}

func (a *UserProfile) GetAddress() []*Address {
	if a == nil {
		return nil
	}
	return a.Addresses
}

func (a *UserProfile) GetAvatar() string {
	if a == nil {
		return ""
	}
	return a.Avatar
}

func (a *UserProfile) GetDescription() string {
	if a == nil {
		return ""
	}
	return a.Description
}

func (x *UserProfile) ConvertFromPb(a *mpb.UserProfile) *UserProfile {
	if x == nil {
		return nil
	}
	obj := &UserProfile{
		Avatar:      a.GetAvatar(),
		Description: a.GetDescription(),
		Addresses: func() []*Address {
			var address []*Address
			for _, d := range a.GetAddresses() {
				address = append(address, (&Address{}).ConvertFromPb(d))
			}
			return address
		}(),
	}
	return obj
}

func (a *UserProfile) ConvertToPb() *mpb.UserProfile {
	if a == nil {
		return nil
	}
	obj := &mpb.UserProfile{
		Avatar:      a.Avatar,
		Description: a.Description,
		Addresses: func() []*mpb.Address {
			var address []*mpb.Address
			for _, d := range a.Addresses {
				address = append(address, d.ConvertToPb())
			}
			return address
		}(),
	}
	return obj
}

type Team struct {
	VapusBase    `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name         string   `bun:"name" json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	Users        []*Users `bun:"users,type:jsonb" json:"users,omitempty" yaml:"users,omitempty" toml:"users,omitempty"`
	Description  string   `bun:"description" json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	TeamId       string   `bun:"team_id" json:"teamId,omitempty" yaml:"teamId,omitempty" toml:"teamId,omitempty"`
	Organization string   `bun:"organization" json:"organization,omitempty" yaml:"organization,omitempty" toml:"organization,omitempty"`
}

func (m *Team) SetAccountId(accountId string) {
	if m != nil {
		m.OwnerAccount = accountId
	}
}

func (m *Team) ConvertToPb() *mpb.Team {
	if m == nil {
		return nil
	}
	return &mpb.Team{
		Name: m.Name,
		Users: func(s []*Users) (pb []*mpb.User) {
			for _, v := range s {
				pb = append(pb, v.ConvertToPb(m.Organization))
			}
			return
		}(m.Users),
		Description: m.Description,
		TeamId:      m.TeamId,
	}
}

func (m *Team) ConvertFromPb(pb *mpb.Team) *Team {
	if pb == nil {
		return nil
	}
	return &Team{
		Name: pb.GetName(),
		Users: func(s []*mpb.User) (pb []*Users) {
			for _, v := range s {
				pb = append(pb, (&Users{}).ConvertFromPb(v))
			}
			return
		}(pb.GetUsers()),
		Description: pb.GetDescription(),
		TeamId:      pb.GetTeamId(),
	}
}

type JwtLog struct {
	VapusBase    `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	JwtId        string `bun:"jwt_id" json:"jwtId,omitempty" yaml:"jwtId,omitempty" toml:"jwtId,omitempty"`
	UserId       string `bun:"user_id" json:"userId,omitempty" yaml:"userId,omitempty" toml:"userId,omitempty"`
	Organization string `bun:"organization" json:"organization,omitempty" yaml:"organization,omitempty" toml:"organization,omitempty"`
	Scope        string `bun:"scope" json:"scope,omitempty" yaml:"scope,omitempty" toml:"scope,omitempty"`
	DataProduct  string `bun:"data_product" json:"dataProduct,omitempty" yaml:"dataProduct,omitempty" toml:"dataProduct,omitempty"`
}

func (m *JwtLog) SetAccountId(accountId string) {
	if m != nil {
		m.OwnerAccount = accountId
	}
}

type RefreshTokenLog struct {
	VapusBase    `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	JwtId        string `bun:"jwt_id" json:"jwtId,omitempty" yaml:"jwtId,omitempty" toml:"jwtId,omitempty"`
	TokenHash    string `bun:"token_hash" json:"tokenHash,omitempty" yaml:"tokenHash,omitempty" toml:"tokenHash,omitempty"`
	UserId       string `bun:"user_id" json:"userId,omitempty" yaml:"userId,omitempty" toml:"userId,omitempty"`
	Organization string `bun:"organization" json:"organization,omitempty" yaml:"organization,omitempty" toml:"organization,omitempty"`
	Scope        string `bun:"scope" json:"scope,omitempty" yaml:"scope,omitempty" toml:"scope,omitempty"`
	ValidTill    int64  `bun:"valid_till" json:"validTill,omitempty" yaml:"validTill,omitempty" toml:"validTill,omitempty"`
}

func (m *RefreshTokenLog) SetAccountId(accountId string) {
	if m != nil {
		m.OwnerAccount = accountId
	}
}

func (dm *Users) SetUserId() {
	if dm != nil {
		if dm.UserId == utils.EMPTYSTR {
			dm.UserId = dm.Email
		}
		if dm.VapusID == utils.EMPTYSTR {
			dm.VapusID = guuid.New().String()
		}
	}
}

func (dm *Users) PreSaveCreate(authzClaim map[string]string) {
	if dm.CreatedBy == utils.EMPTYSTR {
		dm.CreatedBy = authzClaim[encryption.ClaimUserIdKey]
	}
	if dm.CreatedAt == 0 {
		dm.CreatedAt = utils.GetEpochTime()
	}
	if dm.OwnerAccount == utils.EMPTYSTR {
		dm.OwnerAccount = authzClaim[encryption.ClaimAccountKey]
	}
}

func (dm *Users) PreSaveInvite(authzClaim map[string]string, duration time.Duration) {
	if dm.CreatedBy == utils.EMPTYSTR {
		dm.CreatedBy = authzClaim[encryption.ClaimUserIdKey]
	}
	if dm.CreatedAt == 0 {
		dm.CreatedAt = utils.GetEpochTime()
	}
	if dm.InviteId == utils.EMPTYSTR {
		dm.InviteId = guuid.New().String()
	}
	dm.InviteExpiresOn = time.Now().Add(duration).Unix()
	dm.InvitedOn = utils.GetEpochTime()
	dm.OwnerAccount = authzClaim[encryption.ClaimAccountKey]
}

func (dm *Users) PreSaveUpdate(userId string) {
	dm.UpdatedBy = userId
	dm.UpdatedAt = utils.GetEpochTime()
}

func (dm *Users) GetOrganizationRole(organizationId string) []*UserOrganizationRole {
	if organizationId == utils.EMPTYSTR {
		return dm.OrganizationRoles
	}
	if dm == nil || dm.OrganizationRoles == nil {
		return []*UserOrganizationRole{}
	}
	for _, val := range dm.OrganizationRoles {
		if organizationId == val.OrganizationId {
			return []*UserOrganizationRole{val}
		}
	}
	return []*UserOrganizationRole{}
}

func (dm *Users) GetDefaultOrganization() string {
	if dm == nil || dm.OrganizationRoles == nil {
		return ""
	}
	for _, val := range dm.OrganizationRoles {
		if val.IsDefault {
			return val.OrganizationId
		}
	}
	if len(dm.OrganizationRoles) > 0 {
		return dm.OrganizationRoles[0].OrganizationId
	}
	return ""
}

func (dm *Users) SetDefaultOrganization(organization string) {
	if dm == nil {
		return
	}
	for _, val := range dm.OrganizationRoles {
		if val.OrganizationId == organization {
			val.IsDefault = true
		} else {
			val.IsDefault = false
		}
	}
}

func (dm *Team) SetGroupId() {
	if dm.TeamId == utils.EMPTYSTR {
		dm.TeamId = guuid.New().String()
	}
}

func (dm *Team) PreSaveCreate(userId string) {
	if dm.CreatedBy == utils.EMPTYSTR {
		dm.CreatedBy = userId
	}
	if dm.CreatedAt == 0 {
		dm.CreatedAt = utils.GetEpochTime()
	}
}

func (dm *Team) PreSaveUpdate(userId string) {
	dm.UpdatedBy = userId
	dm.UpdatedAt = utils.GetEpochTime()
}
