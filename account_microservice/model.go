package account_microservice

type Account struct {
	ClientID      string  `json:"client_id,omitempty"`
	Name          string  `json:"name,omitempty"`
	Surname       string  `json:"surname,omitempty"`
	PhoneNumber   string  `json:"phone,omitempty"`
	MailAdress    string  `json:"mail,omitempty"`
	AccountAmount float32 `json:"amount,omitempty"`
}
