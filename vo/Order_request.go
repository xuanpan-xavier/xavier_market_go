package vo

type CreateOrderRequest struct {
	Telephone string   `json:"Telephone" binding:"required,max=11,min=11"`
	GoodsID   []string `json:"goodsID" binding:"required"`
}
type UpdateOrderRequest struct {
	GoodsID []string `json:"goodsID" binding:"required"`
}
type ReadOrderRequest struct {
	Telephone string `json:"Telephone" binding:"required,max=11,min=11"`
}
