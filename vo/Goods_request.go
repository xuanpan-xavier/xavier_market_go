package vo

type CreateGoodsRequest struct {
	//ID       uuid.UUID `json:"goods_id" gorm:"type:char(36);primary_key"`
	Name     string `json:"Name" binding:"required,max=30"`
	Category string `json:"Category" binding:"required,max=10"`
	Price    string `json:"Price" binding:"required,min=0"`
	//Stock    string `json:"Stock" binding:"required,min=0"`
	//ImgData  string  `json:"ImgData" binding:"required"`
}

type UpdateGoodsRequest struct {
	//ID       uuid.UUID `json:"goods_id" gorm:"type:char(36);primary_key"`
	//Name     string  `json:"goods_name" binding:"required,max=30"`
	//Category string  `json:"Category" binding:"required,max=10"`
	Price string `json:"Price" binding:"required,min=0"`
	//Stock string `json:"Stock" binding:"required,min=0"`
	//Img   string `json:"Img" binding:"required"`
}
