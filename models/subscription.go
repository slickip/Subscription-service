package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Subscription struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	ServiceName string    `gorm:"type:text; not null"`
	Price       int       `gorm:"not null"`
	StartMonth  int       `gorm:"not null"`
	StartYear   int       `gorm:"not null"`
	EndMonth    *int
	EndYear     *int
	CreatedAt   time.Time
}
