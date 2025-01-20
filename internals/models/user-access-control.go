package models

import (
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type VapusRoles struct {
	VapusBase   `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name        string `bun:"name" yaml:"name"`
	Action      string `bun:"action" yaml:"action"`
	Arn         string `bun:"arn" yaml:"arn"`
	Description string `bun:"description" yaml:"description"`
}

func (x *VapusRoles) ConvertToPb() *mpb.VapusRoles {
	return &mpb.VapusRoles{
		Name:        x.Name,
		Action:      x.Action,
		Arn:         x.Arn,
		Scope:       mpb.VapusRoleScope(mpb.VapusRoleScope_value[x.Scope]),
		Description: x.Description,
	}
}

func (x *VapusRoles) ConvertFromPb(pb *mpb.VapusRoles) *VapusRoles {
	obj := &VapusRoles{
		Name:        pb.Name,
		Action:      pb.Action,
		Arn:         pb.Arn,
		Description: pb.Description,
	}
	obj.Scope = pb.Scope.String()
	return obj

}

type UserRoleMap struct {
	VapusBase      `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	RoleArns       []string `bun:"role,array" yaml:"role"`
	UserID         string   `bun:"user_id" yaml:"userId"`
	OrganizationID string   `bun:"organization_id" yaml:"organizationId"`
	ValidTill      int64    `bun:"valid_till" yaml:"validTill"`
	IsDefault      bool     `bun:"is_default" yaml:"isDefault"`
}

type VapusResourceArn struct {
	VapusBase    `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty" json:"base,omitempty"`
	ResourceName string            `bun:"resource_name" yaml:"resourceName" json:"resourceName"`
	ResourceId   string            `bun:"resource_id" yaml:"resourceId" json:"resourceId"`
	ResourceARN  string            `bun:"resource_arn,unique" yaml:"resourceARN" json:"resourceARN"`
	BlockRules   []*ResourceAclMap `bun:"block_rules,type:jsonb[]" yaml:"blockRules" json:"blockRules"`
	AllowedRules []*ResourceAclMap `bun:"allowed_rules,type:jsonb[]" yaml:"allowedRules" json:"allowedRules"`
	BlockedUsers []string          `bun:"blocked_users,array" yaml:"blockedUsers" json:"blockedUsers"`
}

func (x *VapusResourceArn) PreSaveCreate(authzClaim map[string]string) {
	if x == nil {
		return
	}
	x.PreSaveVapusBase(authzClaim)
}

func (x *VapusResourceArn) PreSaveUpdate(userId string) {
	x.UpdatedBy = userId
	x.UpdatedAt = utils.GetEpochTime()
}

func (x *VapusResourceArn) ConvertToPb() *mpb.VapusResourceArn {
	var BlockedUsers []*mpb.ResourceAclMap
	for _, v := range x.BlockRules {
		BlockedUsers = append(BlockedUsers, v.ConvertToPb())
	}
	var AllowedUsers []*mpb.ResourceAclMap
	for _, v := range x.AllowedRules {
		AllowedUsers = append(AllowedUsers, v.ConvertToPb())
	}
	return &mpb.VapusResourceArn{
		ResourceName: x.ResourceName,
		ResourceId:   x.ResourceId,
		ResourceArn:  x.ResourceARN,
		AllowedRules: AllowedUsers,
		BlockedRules: BlockedUsers,
		BlockedUsers: x.BlockedUsers,
	}
}

type ResourceAclMap struct {
	Users        []string `bun:"users,array" yaml:"users"`
	Organization string   `bun:"organization" yaml:"organization"`
}

func (x *ResourceAclMap) ConvertToPb() *mpb.ResourceAclMap {
	return &mpb.ResourceAclMap{
		Users:        x.Users,
		Organization: x.Organization,
	}
}

func (x *ResourceAclMap) ConvertFromPb(pb *mpb.ResourceAclMap) *ResourceAclMap {
	return &ResourceAclMap{
		Users:        pb.Users,
		Organization: pb.Organization,
	}
}

func (x *VapusResourceArn) ConvertFromPb(pb *mpb.VapusResourceArn) *VapusResourceArn {
	var ar []*ResourceAclMap
	for _, v := range pb.AllowedRules {
		ar = append(ar, &ResourceAclMap{
			Users:        v.Users,
			Organization: v.Organization,
		})
	}
	var br []*ResourceAclMap
	for _, v := range pb.BlockedRules {
		br = append(br, &ResourceAclMap{
			Users:        v.Users,
			Organization: v.Organization,
		})
	}
	return &VapusResourceArn{
		ResourceName: pb.ResourceName,
		ResourceId:   pb.ResourceId,
		ResourceARN:  pb.ResourceArn,
		BlockRules:   br,
		AllowedRules: ar,
	}
}
