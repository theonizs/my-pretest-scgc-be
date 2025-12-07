package model

type HeatData struct {
	ID          int     `json:"id"`
	Temp float64 `json:"temp" binding:"required,min=0,max=100"`
	Time        string  `json:"time" binding:"required"`
	CreatedAt   string  `json:"created_at"`
}

type HeatDataResponse struct {
	Data       []HeatData `json:"data"`
	Prediction *float64   `json:"prediction"`
}

// Standard API Response
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}