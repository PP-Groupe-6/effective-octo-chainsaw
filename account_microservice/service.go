package account_microservice

import (
	"context"
	"errors"
	"fmt"
)

// Déclaration de l'interface exposant les différentes méthodes du service
type AccountService interface {
	GetAccountByID(ctx context.Context, id string) (Account, error)
	Add(ctx context.Context, account Account) (Account, error)
	Update(ctx context.Context, id string, account Account) error
	Delete(ctx context.Context, id string) error
	GetAmountForID(ctx context.Context, id string) (float64, error)
}

// Structure représentant l'instance du service
type accountService struct {
	DbInfos dbConnexionInfo
}

// Déclaration des différentes erreurs de la couche service
var (
	ErrNotAnId           = errors.New("not an ID")
	ErrNoAccount         = errors.New("account param is empty")
	ErrNotFound          = errors.New("account not found not found")
	ErrNoDb              = errors.New("could not access database")
	ErrInconsistentIDs   = errors.New("inconsistent IDs during account update")
	ErrNoInsert          = errors.New("insert did not go through")
	ErrAlreadyExistingID = errors.New("ID already exists in db")
)

// Fonction permettant d'instancer le service
func NewAccountService(dbinfos dbConnexionInfo) AccountService {
	s := &accountService{
		dbinfos,
	}

	return s
}

func (s *accountService) GetAccountByID(ctx context.Context, id string) (Account, error) {
	if id == "" {
		return Account{}, ErrNotAnId
	}

	db := GetDbConnexion(s.DbInfos)

	res := Account{}
	err := db.Get(&res, "SELECT * FROM account WHERE client_id=$1", id)

	if err != nil {
		return Account{}, err
	}

	return res, nil
}

func (s *accountService) Add(ctx context.Context, account Account) (Account, error) {
	if (account == Account{}) {
		return Account{}, ErrNoAccount
	}

	if testID, _ := s.GetAccountByID(ctx, account.ClientID); (testID != Account{}) {
		return Account{}, ErrAlreadyExistingID
	}

	db := GetDbConnexion(s.DbInfos)

	tx := db.MustBegin()
	res := tx.MustExec("INSERT INTO account VALUES ('" + account.ClientID + "','" + account.Name + "','" + account.Surname + "','" + account.PhoneNumber + "','" + account.MailAdress + "'," + fmt.Sprint(account.AccountAmount) + ")")
	tx.Commit()
	db.Close()

	if nRows, err := res.RowsAffected(); nRows != 1 || err != nil {
		if err != nil {
			return Account{}, err
		}
		return Account{}, ErrNoInsert
	}

	insertedAccount, _ := s.GetAccountByID(ctx, account.ClientID)

	return insertedAccount, nil
}

func (s *accountService) Update(ctx context.Context, id string, account Account) error {
	if id != account.ClientID {
		return ErrInconsistentIDs
	}

	if id == "" {
		return ErrNotAnId
	}

	empty := Account{}

	if account == empty {
		return ErrNoAccount
	}

	db := GetDbConnexion(s.DbInfos)

	tx := db.MustBegin()
	tx.MustExec("UPDATE account SET name='" + account.Name + "',surname='" + account.Surname + "',phone_number='" + account.PhoneNumber + "',mail_adress='" + account.MailAdress + "',account_amount=" + fmt.Sprint(account.AccountAmount) + " WHERE client_id = '" + account.ClientID + "'")
	tx.Commit()
	db.Close()

	return nil
}

func (s *accountService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrNotAnId
	}

	if testID, _ := s.GetAccountByID(ctx, id); (testID == Account{}) {
		return ErrNotAnId
	}

	db := GetDbConnexion(s.DbInfos)

	tx := db.MustBegin()
	tx.MustExec("DELETE FROM account WHERE client_id = '" + id + "'")
	tx.Commit()
	db.Close()

	return nil
}

func (s *accountService) GetAmountForID(ctx context.Context, id string) (float64, error) {
	if id == "" {
		return 0, ErrNotAnId
	}

	if testID, _ := s.GetAccountByID(ctx, id); (testID == Account{}) {
		return 0, ErrNotAnId
	}
	acc, err := s.GetAccountByID(ctx, id)

	if err != nil {
		return 0, err
	}

	return acc.AccountAmount, nil
}
