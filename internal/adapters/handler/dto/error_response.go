package dto

type ErrorDTO struct {
	Message string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}
