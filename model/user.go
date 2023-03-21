package model

type User struct {
	Telephone string `gorm:"varchar(11);not null;unique;primary_key"`
	Name      string `gorm:"type:varchar(15);not null"`
	Password  string `gorm:"size:255;not null"`
	//StartTime Time    `json:"StartTime"gorm:"type:timestamp"`
	//EndTime   Time    `json:"EndTime"gorm:"type:timestamp"`
	VipCard VipCard `gorm:"foreignKey:Telephone"`
}
