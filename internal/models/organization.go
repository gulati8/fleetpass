package models

import "time"

type Organization struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Slug      string    `json:"slug" gorm:"type:varchar(255);uniqueIndex;not null"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Organization) TableName() string {
	return "organizations"
}

type CreateOrganizationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type UpdateOrganizationRequest struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	IsActive bool   `json:"is_active"`
}
