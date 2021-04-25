package invoice_microservice

import (
	"context"

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
	createdBy bool
}

type GetInvoiceListResponse struct {
	invoices []InvoiceResponseFormat `json:"invoices"`
}

type InvoiceResponseFormat struct {
	id           string  `json:"id"`
	amount       float32 `json:"amount"`
	state        string  `json:"state"`
	expDate      string  `json:"expDate"`
	withClientId string  `json:"withClientId"`
}

func MakeGetInvoiceListEndpoint(s InvoiceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetInvoiceListRequest)
		var invoicesRet []InvoiceResponseFormat
		invoices, err := s.GetInvoiceList(ctx, req.ClientID)

		for _, invoice := range invoices {
			// Si on veut les invoice créées et que l'utilisateur est le récepteur de l'invoice
			if req.createdBy && invoice.AccountReceiverId == req.ClientID {
				invoicesRet = append(invoicesRet, InvoiceResponseFormat{
					invoice.ID,
					float32(invoice.Amount),
					StateToString(invoice.State),
					invoice.ExpirationDate,
					invoice.AccountPayerId,
				})
			}
			// Si on veut les invoice reçues et que l'utilisateur et le payeur de l'invoice
			if !req.createdBy && invoice.AccountPayerId == req.ClientID {
				invoicesRet = append(invoicesRet, InvoiceResponseFormat{
					invoice.ID,
					float32(invoice.Amount),
					StateToString(invoice.State),
					invoice.ExpirationDate,
					invoice.AccountReceiverId,
				})
			}
		}

		return GetInvoiceListResponse{invoicesRet}, err
	}
}

type AddRequest struct {
	uid         string  // Id du client créant la facture
	emailClient string  // email du client payeur
	amount      float32 // montant de la facture
	expDate     string  // date d'expiration de la facture
}

type AddResponse struct {
	created bool `json:"created"`
}

func MakeAddEndpoint(s InvoiceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)

		id, err := s.GetIdFromMail(ctx, req.emailClient)

		if err != nil {
			return nil, err
		}

		i := Invoice{
			"",
			float64(req.amount),
			PENDING,
			req.expDate,
			id,
			req.uid,
		}

		_, err = s.Create(ctx, i)

		if err != nil {
			return AddResponse{true}, nil
		} else {
			return nil, err
		}
	}
}

type InvoicePaymentRequest struct {
	iid string
}

type InvoicePaymentResponse struct {
	paid bool `json:"paid"`
}

func MakeInvoicePaymentEndpoint(s InvoiceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(InvoicePaymentRequest)

		paid, err := s.PayInvoice(ctx, req.iid)

		return InvoicePaymentResponse{paid}, err
	}
}

type DeleteRequest struct {
	iid string
}

type DeleteResponse struct {
}

func MakeDeleteEndpoint(s InvoiceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)

		err := s.Delete(ctx, req.iid)

		return DeleteResponse{}, err
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
