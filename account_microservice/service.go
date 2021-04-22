package account_microservice

import (
	"context"
	"errors"
)

// Déclaration de l'interface exposant les différentes méthodes du service
type AccountService interface {
	GetAccountByID(ctx context.Context, id string) (Account, error)
	Add(ctx context.Context, account Account) (Account, error)
	Update(ctx context.Context, id string, account Account) error
	Delete(ctx context.Context, id string) error
	GetAmountForID(ctx context.Context, id string) (float32, error)
}

// Structure représentant l'instance du service
type accountService struct {
}

// Déclaration des différentes erreurs de la couche service
var (
	ErrNotAnId         = errors.New("not an ID")
	ErrNotFound        = errors.New("post not found")
	ErrNoDb            = errors.New("could not access database")
	ErrInconsistentIDs = errors.New("could not access database")
)

// Fonction permettant d'instancer le service
func NewAccountService() AccountService {
	s := &accountService{}

	return s
}

func (s *accountService) GetAccountByID(ctx context.Context, id string) (Account, error) {
	return Account{}, nil
}

func (s *accountService) Add(ctx context.Context, account Account) (Account, error) {
	return Account{}, nil
}

func (s *accountService) Update(ctx context.Context, id string, account Account) error {
	return nil
}

func (s *accountService) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *accountService) GetAmountForID(ctx context.Context, id string) (float32, error) {
	return 0, nil
}
