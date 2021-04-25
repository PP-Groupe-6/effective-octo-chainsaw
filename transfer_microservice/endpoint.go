package transfer_microservice

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type TransferEndpoints struct {
	GetTransferListEndpoint    endpoint.Endpoint
	GetWaitingTransferEndpoint endpoint.Endpoint
	CreateEndpoint             endpoint.Endpoint
	PostTransferStatusEndpoint endpoint.Endpoint
}

func MakeTransferEndpoints(s TransferService) TransferEndpoints {
	return TransferEndpoints{
		GetTransferListEndpoint:    MakeGetTransferListEndpoint(s),
		GetWaitingTransferEndpoint: MakeGetWaitingTransferEndpoint(s),
		CreateEndpoint:             MakeCreateEndpoint(s),
		PostTransferStatusEndpoint: MakePostTransferStatusEndpoint(s),
	}
}

type GetTransferListRequest struct {
	ClientID string
}

type FormatedTransfer struct {
	Type     string  `json:"ype"`
	Role     string  `json:"role"`
	FullName string  `json:"name"`
	Amount   float64 `json:"transactionAmount"`
	Date     string  `json:"transactionDate"`
}

type GetTransferListResponse struct {
	Transfers []FormatedTransfer `json:"transfers"`
}

func MakeGetTransferListEndpoint(s TransferService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetTransferListRequest)
		transfers, err := s.GetTransferList(ctx, req.ClientID)
		if err != nil {
			return nil, err
		}
		accountInfo, err := s.GetAccountInformation(ctx, req.ClientID)
		if err != nil {
			return nil, err
		}
		formatedName := accountInfo.Surname + " " + accountInfo.Name
		response := make([]FormatedTransfer, 0)
		for _, transfer := range transfers {
			response = append(response, FormatedTransfer{
				Type:     "transfer",
				Amount:   transfer.Amount,
				FullName: formatedName,
				Date:     transfer.ExecutionDate,
			})

			if transfer.AccountPayerId == req.ClientID {
				response[len(response)-1].Role = "payer"
			} else if transfer.AccountReceiverId == req.ClientID {
				response[len(response)-1].Role = "receiver"
			}
		}
		return GetTransferListResponse{response}, err

	}
}

type GetWaitingTransferRequest struct {
	ClientID string
}

type FormatedWaitingTransfer struct {
	ID               string  `json:"transferId"`
	Mail             string  `json:"mailAdressTransferPayer"`
	Amount           float64 `json:"transferAmount"`
	ExecutionDate    string  `json:"executionTransferDate"`
	ReceiverQuestion string  `json:"receiverQuestion"`
	ReceiverAnswer   string  `json:"receiverAnswer"`
}

type GetWaitingTransferListResponse struct {
	Transfers []FormatedWaitingTransfer `json:"transfers"`
}

func MakeGetWaitingTransferEndpoint(s TransferService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetWaitingTransferRequest)
		transfers, err := s.GetTransferList(ctx, req.ClientID)
		if err != nil {
			return nil, err
		}
		accountInfo, err := s.GetAccountInformation(ctx, req.ClientID)
		if err != nil {
			return nil, err
		}
		response := make([]FormatedWaitingTransfer, 0)

		for _, transfer := range transfers {
			response = append(response, FormatedWaitingTransfer{
				ID:               transfer.ID,
				Mail:             accountInfo.Mail,
				Amount:           transfer.Amount,
				ExecutionDate:    transfer.ExecutionDate,
				ReceiverQuestion: transfer.ReceiverQuestion,
				ReceiverAnswer:   transfer.ReceiverAnswer,
			})
		}
		return GetWaitingTransferListResponse{response}, err

	}
}

type CreateRequest struct {
	emailAdressTransferPayer    string
	emailAdressTransferReceiver string
	transferAmount              float64
	transferType                string
	receiverQuestion            string
	receiverAnswer              string
	executionTransferDate       string
}

type CreateResponse struct {
	Added bool `json:"added"`
}

func MakeCreateEndpoint(s TransferService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		idPayer, err := s.GetIdFromMail(ctx, req.emailAdressTransferPayer)
		if err != nil {
			return nil, err
		}
		idReceiver, err := s.GetIdFromMail(ctx, req.emailAdressTransferReceiver)
		if err != nil {
			return nil, err
		}
		var date string
		if req.transferType == "instant" {
			date = time.Now().Format("2006-01-02")
		} else if req.transferType == "scheduled" {
			date = req.executionTransferDate
		}
		toAdd := Transfer{
			ID:                "",
			Type:              req.transferType,
			State:             0,
			Amount:            req.transferAmount,
			AccountPayerId:    idPayer,
			AccountReceiverId: idReceiver,
			ReceiverQuestion:  req.receiverQuestion,
			ReceiverAnswer:    req.receiverAnswer,
			ExecutionDate:     date,
		}

		transfer, err := s.Create(ctx, toAdd)
		if (err == nil && transfer != Transfer{}) {
			return true, nil
		} else {
			return false, err
		}
	}
}

type PostTransferStatusRequest struct {
	ID string `json:"transfer_id"`
}

type PostTransferStatusResponse struct {
	Done bool `json:"done"`
}

func MakePostTransferStatusEndpoint(s TransferService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PostTransferStatusRequest)
		res, err := s.PostTransferStatus(ctx, req.ID)

		if err == nil && res {
			return true, nil
		} else {
			return false, err
		}

	}
}
