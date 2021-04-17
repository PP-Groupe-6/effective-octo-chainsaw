package account_microservice

import (
	"context"
	"testing"
)

func TestAdd(t *testing.T) {
	// Instance du service utilisée pour les tests
	s := NewAccountService()

	// Compte mocké
	mockAccount := Account{
		"jsvkjjsbfvkjnsdfvknsdlfjkvnsdlkfnvuhdfovhj",
		"Valentin",
		"Roche",
		"+33678239003",
		"jsadjcb@jnskd.com",
		10000,
	}

	// Test avec un paramètre vide
	_, err := s.Add(context.TODO(), Account{})

	if err == nil {
		t.Errorf("Passed empty account param to add function, should have raised an error")
	}

	// Test avec un compte valide
	result, err := s.Add(context.TODO(), mockAccount)
	if err != nil {
		t.Errorf("Valid account, method should not fail")
	}

	if result != mockAccount {
		t.Errorf("Returned account is not the same as the param")
	}
}
