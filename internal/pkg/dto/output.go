package dto

type (
	ItemsOutputDTO struct {
		Items interface{} `json:"items"`
		Count int64       `json:"count"`
	}
)
