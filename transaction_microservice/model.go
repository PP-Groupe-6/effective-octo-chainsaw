package transaction_microservice

type Transfert struct {
	TransferId                string  `json:"transfer_id,omitempty"`
	TransferType              string  `json:"transfer_type,omitempty"`
	TransferAmount            float64 `json:"transfer_amount,omitempty"`
	AccountTransferPayerId    string  `json:"transfer_payer_id,omitempty"`
	AccountTransferReceiverId string  `json:"transfer_receiver_id,omitempty"`
	ReceiverQuestion          string  `json:"receiver_question,omitempty"`
	ReceiverAnswer            string  `json:"receiver_answer,omitempty"`
	ScheduledTransferDate     string  `json:"scheduled_transfer_date,omitempty"`
	ExecutedTransferDate      string  `json:"executed_transfer_date,omitempty"`
}

type Invoice struct {
	InvoiceId                string       `json:"invoice_id,omitempty"`
	InvoiceAmount            float64      `json:"invoice_amount,omitempty"`
	InvoiceState             InvoiceState `json:"invoice_state,omitempty"`
	InvoiceExpirationDate    string       `json:"invoice_expiration_date,omitempty"`
	AccountInvoicePayerId    string       `json:"invoice_payer_id,omitempty"`
	AccountInvoiceReceiverId string       `json:"invoice_receveiver_id,omitempty"`
}

type InvoiceState struct {
	StateId   string `json:"state_id,omitempty"`
	StateName string `json:"state_name,omitempty"`
}
