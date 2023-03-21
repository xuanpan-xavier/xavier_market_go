package model

import (
	"github.com/jinzhu/gorm"
)

type Goods struct {
	gorm.Model
	Name     string  `json:"Name" gorm:"type:varchar(50);not null"`
	Category string  `json:"Category" gorm:"type:varchar(50);not null"`
	Price    float32 `json:"Price" gorm:"not null"`
	//Stock    uint      `json:"Stock" gorm:"not null"`
	//Img      string    `json:"Img"`
}

//func (goods *Goods) BeforeCreate(scope *gorm.Scope) error {
//	return scope.SetColumn("ID", uuid.NewV4())
//}
