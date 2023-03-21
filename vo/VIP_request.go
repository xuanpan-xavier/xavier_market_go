package vo

import "xmarket_gin/model"

type CreateVIPRequest struct {
	Name      string `json:"Name" binding:"required,max=10,min=3"`
	Telephone string `json:"Telephone" binding:"required,max=11,min=11"`
	Password  string `json:"Password" binding:"required,max=15,min=6"`
	//StartTime model.Time `json:"StartTime"`
	//EndTIme   model.Time `json:"EndTIme"`
}
type CVIPRequest struct {
	Name      string `json:"Name" binding:"required,max=10,min=3"`
	Telephone string `json:"Telephone" binding:"required,max=11,min=11"`
	//Password  string `json:"Password" binding:"required,max=15,min=6"`
	//StartTime model.Time `json:"StartTime"`
	//EndTIme   model.Time `json:"EndTIme"`
}
type UpdateVIPRequest struct {
	//Name      string `json:"Name" binding:"required,max=10,min=3"`
	//Telephone string `json:"Telephone" binding:"required,max=11,min=11"`
	//Password  string `json:"Password" binding:"required,max=15,min=6"`
	StartTime model.Time `json:"StartTime"`
	EndTime   model.Time `json:"EndTIme"`
}
type User2VIPRequest struct {
	IsVIP string `json:"IsVIP"`
	Year  string `json:"Year"`
}
type UpdatePointRequest struct {
	Minus     uint   `json:"Minus"`
	Telephone string `json:"Telephone" binding:"required,max=11,min=11"`
}
