package transfer_microservice

import (
	"context"
	"errors"
	"fmt"
)

type TransferService interface {
	Create(ctx context.Context, transfer Transfer) (Transfer, error)
	Read(ctx context.Context, id string) (Transfer, error)
	Update(ctx context.Context, id string, transfer Transfer) (Transfer, error)
	Delete(ctx context.Context, id string) error
}

var (
	ErrNotAnId         = errors.New("not an ID")
	ErrNotFound        = errors.New("transaction not found")
	ErrNoTransfer      = errors.New("transfer field is empty")
	ErrNoUpdate        = errors.New("could not update transfer")
	ErrNoDb            = errors.New("could not access database")
	ErrAlreadyExist    = errors.New("transfer id already exists")
	ErrNoInsert        = errors.New("insert did not go through")
	ErrInconsistentIDs = errors.New("could not access database")
)

type transferService struct {
	DbInfos dbConnexionInfo
}

func NewTransactionService(dbinfos dbConnexionInfo) TransferService {
	return &transferService{
		DbInfos: dbinfos,
	}
}

func (s *transferService) Create(ctx context.Context, transfer Transfer) (Transfer, error) {
	if (transfer == Transfer{}) {
		return Transfer{}, ErrNoTransfer
	}

	if testID, _ := s.Read(ctx, transfer.ID); (testID != Transfer{}) {
		return Transfer{}, ErrAlreadyExist
	}

	db := GetDbConnexion(s.DbInfos)
	tx := db.MustBegin()
	res := tx.MustExec("INSERT INTO transfer VALUES ('" + transfer.ID + "','" + transfer.Type + "'," + fmt.Sprint(transfer.Amount) + ",'" + transfer.AccountPayerId + "','" + transfer.AccountReceiverId + "','" + transfer.ReceiverQuestion + "','" + transfer.ReceiverAnswer + "','" + transfer.ScheduledDate + "','" + transfer.ExecutedDate + "')")
	tx.Commit()
	db.Close()

	if nRows, err := res.RowsAffected(); nRows != 1 || err != nil {
		if err != nil {
			return Transfer{}, err
		}
		return Transfer{}, ErrNoInsert
	}

	return s.Read(ctx, transfer.ID)

}

func (s *transferService) Read(ctx context.Context, id string) (Transfer, error) {
	db := GetDbConnexion(s.DbInfos)

	res := Transfer{}
	err := db.Get(&res, "SELECT * FROM transfer WHERE transfer_id=$1", id)

	if err != nil {
		return Transfer{}, err
	}

	return res, nil
}

func (s *transferService) Update(ctx context.Context, id string, transfer Transfer) (Transfer, error) {
	if (transfer == Transfer{}) {
		return Transfer{}, ErrNoTransfer
	}

	if testID, _ := s.Read(ctx, id); (testID == Transfer{}) {
		return Transfer{}, ErrNotFound
	}

	db := GetDbConnexion(s.DbInfos)
	tx := db.MustBegin()
	res := tx.MustExec("UPDATE transfer SET transfer_type = '" + transfer.Type + "', transfer_amount =" + fmt.Sprint(transfer.Amount) + ", account_transfer_payer_id = '" + transfer.AccountPayerId + "', account_transfer_receiver_id = '" + transfer.AccountReceiverId + "', receiver_question = '" + transfer.ReceiverQuestion + "', receiver_answer = '" + transfer.ReceiverAnswer + "', scheduled_transfer_date = '" + transfer.ScheduledDate + "', executed_transfer_date = '" + transfer.ExecutedDate + "' WHERE transfer_id=$1")
	tx.Commit()
	db.Close()

	if nRows, err := res.RowsAffected(); nRows != 1 || err != nil {
		if err != nil {
			return Transfer{}, err
		}
		return Transfer{}, ErrNoInsert
	}

	return s.Read(ctx, transfer.ID)
}

func (s *transferService) Delete(ctx context.Context, id string) error {

	if testID, _ := s.Read(ctx, id); (testID == Transfer{}) {
		return ErrNotFound
	}
	db := GetDbConnexion(s.DbInfos)
	tx := db.MustBegin()
	res := tx.MustExec("DELETE FROM transfer WHERE transfer_id=$1", id)

	if nRows, err := res.RowsAffected(); nRows != 1 || err != nil {
		if err != nil {
			return err
		}
	}
	tx.Commit()
	db.Close()

	return nil
}
