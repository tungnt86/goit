package model

type Product struct {
	BaseModel
	Name        string
	CategoryID  int64
	WarehouseID int64
}
