package model

import (
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	Telephone string  `gorm:"varchar(11);not null;"`
	Goods     []Goods `gorm:"many2many:order_goods;"`
	Total     float32 `gorm:"not null;"`
	Number    string
}
