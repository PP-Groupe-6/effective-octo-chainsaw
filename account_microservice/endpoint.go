package account_microservice

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

type AccountEndpoints struct {
	GetAmountEndpoint          endpoint.Endpoint
	GetUserInformationEndpoint endpoint.Endpoint
	AddEndpoint                endpoint.Endpoint
	UpdateEndpoint             endpoint.Endpoint
	DeleteEndpoint             endpoint.Endpoint
}

func MakeAccountEndpoints(s AccountService) AccountEndpoints {
	return AccountEndpoints{
		GetAmountEndpoint:          MakeGetAmountEndpoint(s),
		GetUserInformationEndpoint: MakeGetUserInformationEndpoint(s),
		AddEndpoint:                MakeAddEndpoint(s),
	}
}

type GetAmountRequest struct {
	ClientID string
}

type GetAmountResponse struct {
	AccountAmount float64 `json:"amount"`
}

func MakeGetAmountEndpoint(s AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAmountRequest)
		amount, err := s.GetAmountForID(ctx, req.ClientID)
		return GetAmountResponse{amount}, err
	}
}

type GetUserInformationRequest struct {
	ClientID string
}

type GetUserInformationResponse struct {
	FullName    string `json:"fullName"`
	MailAdress  string `json:"mailAdress"`
	PhoneNumber string `json:"phoneNumber"`
}

func MakeGetUserInformationEndpoint(s AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserInformationRequest)
		account, err := s.GetAccountByID(ctx, req.ClientID)
		formated_name := account.Surname + " " + account.Name
		return GetUserInformationResponse{formated_name, account.MailAdress, account.PhoneNumber}, err
	}
}

type AddRequest struct {
	ClientID    string
	FullName    string
	PhoneNumber int
	MailAdress  string
}

type AddResponse struct {
	Added bool `json:"added"`
}

func MakeAddEndpoint(s AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)
		sep_name := strings.Split(req.FullName, " ")
		phone := strconv.Itoa(req.PhoneNumber)
		to_add := Account{
			req.ClientID,
			sep_name[0],
			sep_name[1],
			phone,
			req.MailAdress,
			0,
		}
		account, err := s.Add(ctx, to_add)
		if (err == nil && account != Account{}) {
			return true, nil
		} else {
			return false, err
		}
	}
}
