package transfer_microservice

type Transfer struct {
	ID                string  `json:"transfer_id,omitempty"`
	Type              string  `json:"transfer_type,omitempty"`
	Amount            float64 `json:"transfer_amount,omitempty"`
	AccountPayerId    string  `json:"transfer_payer_id,omitempty"`
	AccountReceiverId string  `json:"transfer_receiver_id,omitempty"`
	ReceiverQuestion  string  `json:"receiver_question,omitempty"`
	ReceiverAnswer    string  `json:"receiver_answer,omitempty"`
	ScheduledDate     string  `json:"scheduled_transfer_date,omitempty"`
	ExecutedDate      string  `json:"executed_transfer_date,omitempty"`
}
