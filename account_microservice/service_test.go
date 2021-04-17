package account_microservice

import (
	"context"
	"testing"
)

type TestData struct {
	s           AccountService
	mockAccount Account
}

func NewTestData() TestData {
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

	return TestData{
		s,
		mockAccount,
	}
}

func TestAdd(t *testing.T) {
	testData := NewTestData()

	// Test avec un paramètre vide
	_, err := testData.s.Add(context.TODO(), Account{})

	if err == nil {
		t.Errorf("Passed empty account param to add function, should have raised an error")
	}

	// Test avec un compte valide
	result, err := testData.s.Add(context.TODO(), testData.mockAccount)
	if err != nil {
		t.Errorf("Valid account, method should not fail")
	}

	if result != testData.mockAccount {
		t.Errorf("Returned account is not the same as the param")
	}
}
