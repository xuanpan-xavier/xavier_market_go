package vo

type UpdateVouchersRequest struct {
	ID string `json:"ID"`
	//IsUsed int    `json:"IsUsed"`
}
type ReadVouchersRequest struct {
	Telephone string `json:"Telephone"`
}
