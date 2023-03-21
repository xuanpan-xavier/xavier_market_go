package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Voucher struct {
	ID        uuid.UUID `json:"ID" gorm:"type:char(36);primary_key"`
	Telephone string    `gorm:"varchar(11);not null;"`
	Point     uint      `json:"Point" gorm:"type:int"`
	IsUsed    int       `json:"IsUsed" gorm:"type:int"`
}

func (voucher *Voucher) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
}
