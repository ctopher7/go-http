package model

type HttpRes struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type HttpResLoan struct {
	Message string `json:"message,omitempty"`
	Data    []Loan `json:"data,omitempty"`
}
