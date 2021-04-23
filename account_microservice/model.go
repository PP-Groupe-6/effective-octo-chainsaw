package account_microservice

type Account struct {
	ClientID      string  `json:"client_id,omitempty" db:"client_id"`
	Name          string  `json:"name,omitempty" db:"name"`
	Surname       string  `json:"surname,omitempty" db:"surname"`
	PhoneNumber   string  `json:"phone,omitempty" db:"phone_number"`
	MailAdress    string  `json:"mail,omitempty" db:"mail_adress"`
	AccountAmount float64 `json:"amount,omitempty" db:"account_amount"`
}
