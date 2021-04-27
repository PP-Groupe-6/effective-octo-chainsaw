package account_microservice

import (
	"context"
	"fmt"
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
	AccountAmount string `json:"amount"`
}

func MakeGetAmountEndpoint(s AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAmountRequest)
		amount, err := s.GetAmountForID(ctx, req.ClientID)

		return GetAmountResponse{fmt.Sprint(amount)}, err
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
		formatedName := account.Surname + " " + account.Name
		return GetUserInformationResponse{formatedName, account.MailAdress, account.PhoneNumber}, err
	}
}

type AddRequest struct {
	ClientID    string
	FullName    string
	PhoneNumber string
	MailAdress  string
}

type AddResponse struct {
	Account Account `json:"client"`
}

func MakeAddEndpoint(s AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)
		sepName := strings.Split(req.FullName, " ")
		toAdd := Account{
			req.ClientID,
			sepName[0],
			sepName[1],
			req.PhoneNumber,
			req.MailAdress,
			0,
		}
		account, err := s.Add(ctx, toAdd)
		formatedName := account.Surname + " " + account.Name
		if (err == nil && account != Account{}) {
			return GetUserInformationResponse{formatedName, account.MailAdress, account.PhoneNumber}, nil
		} else {
			return GetUserInformationResponse{formatedName, account.MailAdress, account.PhoneNumber}, err
		}
	}
}
