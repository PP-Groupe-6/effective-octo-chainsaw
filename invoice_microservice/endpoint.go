package invoice_microservice

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

const (
	PENDING = 0
	PAID    = 1
	EXPIRED = 2
)

type InvoiceEndpoints struct {
	GetInvoiceListEndpoint  endpoint.Endpoint
	AddEndpoint             endpoint.Endpoint
	DeleteEndpoint          endpoint.Endpoint
	InvoicePaiementEndpoint endpoint.Endpoint
}

func MakeInvoiceEndpoints(s InvoiceService) InvoiceEndpoints {
	return InvoiceEndpoints{
		GetInvoiceListEndpoint:  MakeGetInvoiceListEndpoint(s),
		AddEndpoint:             MakeAddEndpoint(s),
		DeleteEndpoint:          MakeDeleteEndpoint(s),
		InvoicePaiementEndpoint: MakeInvoicePaymentEndpoint(s),
	}
}

// Si created by est à true on retourne les invoices créées par le client si il est à false on retourne celles reçues par le client
type GetInvoiceListRequest struct {
	ClientID  string
	CreatedBy bool
}

type GetInvoiceListResponse struct {
	Invoices []InvoiceResponseFormat `json:"invoices"`
}

type InvoiceResponseFormat struct {
	Name      string `json:"name"`
	Mail      string `json:"mail"`
	Phone     string `json:"phone"`
	Amount    string `json:"amount"`
	State     string `json:"state"`
	ExpDate   string `json:"expDate"`
	InvoiceID string `json:"InvoiceID"`
}

func MakeGetInvoiceListEndpoint(s InvoiceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetInvoiceListRequest)
		InvoicesRet := []InvoiceResponseFormat{}
		invoices, err := s.GetInvoiceList(ctx, req.ClientID)
		for _, Invoice := range invoices {
			// Si on veut les invoice créées et que l'utilisateur est le récepteur de l'invoice
			if req.CreatedBy && Invoice.AccountReceiverId == req.ClientID {
				otherAccount, err := s.GetAccountInformation(ctx, Invoice.AccountPayerId)

				if err != nil {
					return GetInvoiceListResponse{InvoicesRet}, err
				}

				InvoicesRet = append(InvoicesRet, InvoiceResponseFormat{
					otherAccount.Name + " " + otherAccount.Surname,
					otherAccount.Mail,
					otherAccount.Phone,
					fmt.Sprint(Invoice.Amount),
					StateToString(Invoice.State),
					Invoice.ExpirationDate,
					Invoice.ID,
				})
			}
			// Si on veut les invoice reçues et que l'utilisateur et le payeur de l'invoice
			if !req.CreatedBy && Invoice.AccountPayerId == req.ClientID {
				otherAccount, err := s.GetAccountInformation(ctx, Invoice.AccountReceiverId)

				if err != nil {
					return GetInvoiceListResponse{InvoicesRet}, err
				}

				InvoicesRet = append(InvoicesRet, InvoiceResponseFormat{
					otherAccount.Name + " " + otherAccount.Surname,
					otherAccount.Mail,
					otherAccount.Phone,
					fmt.Sprint(Invoice.Amount),
					StateToString(Invoice.State),
					Invoice.ExpirationDate,
					Invoice.ID,
				})
			}
		}
		return GetInvoiceListResponse{InvoicesRet}, err
	}
}

type AddRequest struct {
	Uid         string  // Id du client créant la facture
	EmailClient string  // email du client payeur
	Amount      float32 // montant de la facture
	ExpDate     string  // date d'expiration de la facture
}

type AddResponse struct {
	Created bool `json:"created"`
}

func MakeAddEndpoint(s InvoiceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)

		id, err := s.GetIdFromMail(ctx, req.EmailClient)

		if err != nil {
			return nil, err
		}

		i := Invoice{
			"",
			float64(req.Amount),
			PENDING,
			req.ExpDate,
			id,
			req.Uid,
		}

		_, err = s.Create(ctx, i)

		if err == nil {
			return AddResponse{true}, nil
		} else {
			return nil, err
		}
	}
}

type InvoicePaymentRequest struct {
	Iid string
}

type InvoicePaymentResponse struct {
	Paid bool `json:"paid"`
}

func MakeInvoicePaymentEndpoint(s InvoiceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(InvoicePaymentRequest)

		paid, err := s.PayInvoice(ctx, req.Iid)

		return InvoicePaymentResponse{paid}, err
	}
}

type DeleteRequest struct {
	Iid string
}

type DeleteResponse struct {
	Deleted bool `json:"deleted"`
}

func MakeDeleteEndpoint(s InvoiceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)

		err := s.Delete(ctx, req.Iid)

		if err != nil {
			return DeleteResponse{false}, err
		} else {
			return DeleteResponse{true}, nil
		}
	}
}

func StateToString(stateID int) string {
	switch stateID {
	case PENDING:
		return "Pending"
	case PAID:
		return "Paid"
	case EXPIRED:
		return "Expired"
	}
	return ""
}
