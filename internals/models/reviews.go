package models

import (
	encrytion "github.com/vapusdata-oss/aistudio/core/encryption"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
)

type UpDownVote struct {
	VapusBase  `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Resource   string `bun:"resource" json:"resource,omitempty" yaml:"resource,omitempty" toml:"resource,omitempty"`
	ResourceID string `bun:"resource_id" json:"resource_id,omitempty" yaml:"resource_id,omitempty" toml:"resource_id,omitempty"`
	UpVote     int    `bun:"up_vote" json:"up_vote,omitempty" yaml:"up_vote,omitempty" toml:"up_vote,omitempty"`
	DownVote   int    `bun:"down_vote" json:"down_vote,omitempty" yaml:"down_vote,omitempty" toml:"down_vote,omitempty"`
}

func (u *UpDownVote) PreSaveCreate(authzClaim map[string]string) {
	if u.CreatedBy == "" {
		u.CreatedBy = authzClaim[encrytion.ClaimUserIdKey]
	}
	if u.CreatedAt == 0 {
		u.CreatedAt = utils.GetEpochTime()
	}
	if u.OwnerAccount == utils.EMPTYSTR {
		u.OwnerAccount = authzClaim[encrytion.ClaimAccountKey]
	}
	if u.VapusID == utils.EMPTYSTR {
		u.VapusID = utils.GetUUID()
	}
}

func (u *UpDownVote) PreSaveUpdate(userId string) {
	if u == nil {
		return
	}
	u.UpdatedBy = userId
	u.UpdatedAt = utils.GetEpochTime()
}

type StarReview struct {
	VapusBase  `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Resource   string  `bun:"resource" json:"resource,omitempty" yaml:"resource,omitempty" toml:"resource,omitempty"`
	ResourceID string  `bun:"resource_id" json:"resource_id,omitempty" yaml:"resource_id,omitempty" toml:"resource_id,omitempty"`
	Star       float32 `bun:"star" json:"star,omitempty" yaml:"star,omitempty" toml:"star,omitempty"`
}

func (u *StarReview) PreSaveCreate(authzClaim map[string]string) {
	if u.CreatedBy == "" {
		u.CreatedBy = authzClaim[encrytion.ClaimUserIdKey]
	}
	if u.CreatedAt == 0 {
		u.CreatedAt = utils.GetEpochTime()
	}
	if u.OwnerAccount == utils.EMPTYSTR {
		u.OwnerAccount = authzClaim[encrytion.ClaimAccountKey]
	}
	if u.VapusID == utils.EMPTYSTR {
		u.VapusID = utils.GetUUID()
	}
}

func (u *StarReview) PreSaveUpdate(userId string) {
	if u == nil {
		return
	}
	u.UpdatedBy = userId
	u.UpdatedAt = utils.GetEpochTime()
}
