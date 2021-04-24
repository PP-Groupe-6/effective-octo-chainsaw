package invoice_microservice

type InvoiceState struct {
	ID   string `json:"state_id,omitempty"`
	Name string `json:"state_name,omitempty"`
}

type Invoice struct {
	ID                string  `json:"invoice_id,omitempty"`
	Amount            float64 `json:"invoice_amount,omitempty"`
	State             int     `json:"invoice_state,omitempty"`
	ExpirationDate    string  `json:"invoice_expiration_date,omitempty"`
	AccountPayerId    string  `json:"invoice_payer_id,omitempty"`
	AccountReceiverId string  `json:"invoice_receveiver_id,omitempty"`
}
