package model

type Model interface {
	GetID() int64
}

type BaseModel struct {
	ID int64
}

func (b *BaseModel) GetID() int64 {
	return b.ID
}
