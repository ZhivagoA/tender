package models

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationType string

const (
	IE  OrganizationType = "IE"
	LLC OrganizationType = "LLC"
	JSC OrganizationType = "JSC"
)

type Organization struct {
	ID        uuid.UUID        `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string           `json:"name" gorm:"size:50;not null"`
	Type      OrganizationType `json:"type" gorm:"type:organization_type;not null"`
	CreatedAt time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username  string    `json:"username" gorm:"size:50;unique;not null"`
	FirstName string    `json:"first_name" gorm:"size:50"`
	LastName  string    `json:"last_name" gorm:"size:50"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Tender struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title             string    `json:"title" gorm:"size:100;not null"`
	Description       string    `json:"description" gorm:"type:text"`
	OrganizationID    uuid.UUID `json:"organization_id" gorm:"not null"`
	Status            string    `json:"status" gorm:"size:50;not null"`
	Version           int       `json:"version" gorm:"not null;default:1"`
	ResponsibleUserID uuid.UUID `json:"responsible_user_id" gorm:"not null"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type OrganizationResponsible struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"not null"`
	UserID         uuid.UUID `json:"user_id" gorm:"not null"`
}

type Feedback struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BidID     uuid.UUID `json:"bid_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Comment   string    `json:"comment" gorm:"type:text;not null"`
	Rating    int       `json:"rating" gorm:"type:int;not null;check:rating >= 1 and rating <= 5"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Feedback) TableName() string {
	return "feedback"
}

type Bid struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Amount        float64   `json:"amount" gorm:"type:numeric(10,2);not null"`
	TenderID      uuid.UUID `json:"tender_id" gorm:"not null"`
	UserID        uuid.UUID `json:"user_id" gorm:"not null"`
	Status        string    `json:"status" gorm:"size:50;not null"`
	Version       int       `json:"version" gorm:"not null;default:1"`
	ApprovalCount int       `json:"approval_count" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Bid) TableName() string {
	return "bid"
}

type TenderVersion struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenderID          uuid.UUID `json:"tender_id" gorm:"not null"`
	Title             string    `json:"title" gorm:"size:100;not null"`
	Description       string    `json:"description" gorm:"type:text"`
	OrganizationID    uuid.UUID `json:"organization_id" gorm:"not null"`
	ResponsibleUserID uuid.UUID `json:"responsible_user_id" gorm:"type:uuid;"`
	Status            string    `json:"status" gorm:"size:50;not null"`
	Version           int       `json:"version" gorm:"not null"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type BidVersion struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BidID         uuid.UUID `json:"bid_id" gorm:"not null"`
	Amount        float64   `json:"amount" gorm:"type:numeric(10,2);not null"`
	TenderID      uuid.UUID `json:"tender_id" gorm:"not null"`
	UserID        uuid.UUID `json:"user_id" gorm:"not null"`
	Status        string    `json:"status" gorm:"size:50;not null"`
	Version       int       `json:"version" gorm:"not null"`
	ApprovalCount int       `json:"approval_count" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}
