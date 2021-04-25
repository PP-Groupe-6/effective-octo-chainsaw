package invoice_microservice

import (
	"context"
	"errors"
	"fmt"
)

type InvoiceService interface {
	Create(ctx context.Context, invoice Invoice) (Invoice, error)
	Read(ctx context.Context, id string) (Invoice, error)
	Update(ctx context.Context, id string, invoice Invoice) (Invoice, error)
	Delete(ctx context.Context, id string) error
	GetInvoiceList(ctx context.Context, id string) ([]*Invoice, error)
}

var (
	ErrNotAnId         = errors.New("not an ID")
	ErrNotFound        = errors.New("invoice not found")
	ErrNoTransfer      = errors.New("invoice field is empty")
	ErrNoUpdate        = errors.New("could not update invoice")
	ErrNoDb            = errors.New("could not access database")
	ErrAlreadyExist    = errors.New("invoice id already exists")
	ErrNoInsert        = errors.New("insert did not go through")
	ErrInconsistentIDs = errors.New("could not access database")
)

type invoiceService struct {
	DbInfos dbConnexionInfo
}

func NewInvoiceService(dbinfos dbConnexionInfo) InvoiceService {
	return &invoiceService{
		DbInfos: dbinfos,
	}
}

func (s *invoiceService) GetInvoiceList(ctx context.Context, id string) ([]*Invoice, error) {
	db := GetDbConnexion(s.DbInfos)

	invoices := make([]*Invoice, 0)
	rows, err := db.Queryx("SELECT * FROM invoice WHERE account_invoice_payer_id=$1 OR account_invoice_receiver_id=$1", id)

	for rows.Next() {
		var i Invoice
		if err := rows.StructScan(&i); err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return invoices, err
}

func (s *invoiceService) Create(ctx context.Context, invoice Invoice) (Invoice, error) {
	if (invoice == Invoice{}) {
		return Invoice{}, ErrNoTransfer
	}

	if testID, _ := s.Read(ctx, invoice.ID); (testID != Invoice{}) {
		return Invoice{}, ErrAlreadyExist
	}

	db := GetDbConnexion(s.DbInfos)
	tx := db.MustBegin()
	res := tx.MustExec("INSERT INTO invoice VALUES('" + invoice.ID + "','" + fmt.Sprint(invoice.Amount) + "','" + fmt.Sprint(invoice.State) + "','" + invoice.ExpirationDate + "','" + invoice.AccountPayerId + "','" + invoice.AccountReceiverId + "')")
	tx.Commit()
	db.Close()

	if nRows, err := res.RowsAffected(); nRows != 1 || err != nil {
		if err != nil {
			return Invoice{}, err
		}
		return Invoice{}, ErrNoInsert
	}

	return s.Read(ctx, invoice.ID)
}

func (s *invoiceService) Read(ctx context.Context, id string) (Invoice, error) {
	db := GetDbConnexion(s.DbInfos)

	res := Invoice{}
	err := db.Get(&res, "SELECT * FROM invoice WHERE invoice_id=$1", id)

	if err != nil {
		return Invoice{}, err
	}

	return res, nil
}

func (s *invoiceService) Update(ctx context.Context, id string, invoice Invoice) (Invoice, error) {
	if (invoice == Invoice{}) {
		return Invoice{}, ErrNoTransfer
	}

	if testID, _ := s.Read(ctx, id); (testID == Invoice{}) {
		return Invoice{}, ErrNotFound
	}

	db := GetDbConnexion(s.DbInfos)
	tx := db.MustBegin()
	res := tx.MustExec("UPDATE invoice SET invoice_amount = '"+fmt.Sprint(invoice.Amount)+"', invoice_state ='"+fmt.Sprint(invoice.State)+"', invoice_expiration_date = '"+invoice.ExpirationDate+"', account_invoice_payer_id = '"+invoice.AccountPayerId+"', account_invoice_receiver_id = '"+invoice.AccountReceiverId+"' WHERE invoice_id=$1", id)
	tx.Commit()
	db.Close()

	if nRows, err := res.RowsAffected(); nRows != 1 || err != nil {
		if err != nil {
			return Invoice{}, err
		}
		return Invoice{}, ErrNoInsert
	}

	return s.Read(ctx, invoice.ID)
}

func (s *invoiceService) Delete(ctx context.Context, id string) error {
	if testID, _ := s.Read(ctx, id); (testID == Invoice{}) {
		return ErrNotFound
	}
	db := GetDbConnexion(s.DbInfos)
	tx := db.MustBegin()
	res := tx.MustExec("DELETE FROM invoice WHERE invoice_id=$1", id)

	if nRows, err := res.RowsAffected(); nRows != 1 || err != nil {
		if err != nil {
			return err
		}
	}
	tx.Commit()
	db.Close()

	return nil
}
