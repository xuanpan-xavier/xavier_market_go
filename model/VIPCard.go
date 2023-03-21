package model

type VipCard struct {
	Telephone string `gorm:"varchar(11);not null;unique;primary_key"`
	Point     uint   `json:"Point" gorm:"type:int"`
	StartTime Time   `json:"StartTime"gorm:"type:timestamp"`
	EndTime   Time   `json:"EndTime"gorm:"type:timestamp"`
}
