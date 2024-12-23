package model

type Hunt struct {
	BaseModel
	Name string `gorm:"type:text;"`
}

type CreateHuntRequest struct {
	Name string `json:"name" binding:"required,min:3,max:100"`
}

type UpdateHuntRequest struct {
	Name string `json:"name" binding:"required,min:3,max:100"`
}

type HuntResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
