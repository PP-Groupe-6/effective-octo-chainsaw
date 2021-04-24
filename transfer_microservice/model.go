package transfer_microservice

type Transfer struct {
	ID                string  `json:"transfer_id,omitempty" db:"transfer_id"`
	Type              string  `json:"transfer_type,omitempty" db:"transfer_type"`
	Amount            float64 `json:"transfer_amount,omitempty" db:"transfer_amount"`
	AccountPayerId    string  `json:"transfer_payer_id,omitempty" db:"account_transfert_payer_id"`
	AccountReceiverId string  `json:"transfer_receiver_id,omitempty" db:"account_transfert_receiver_id"`
	ReceiverQuestion  string  `json:"receiver_question,omitempty" db:"receiver_question"`
	ReceiverAnswer    string  `json:"receiver_answer,omitempty" db:"receiver_answer"`
	ScheduledDate     string  `json:"scheduled_transfer_date,omitempty" db:"scheduled_transfert_date"`
	ExecutedDate      string  `json:"executed_transfer_date,omitempty" db:"executed_transfer_date"`
}
