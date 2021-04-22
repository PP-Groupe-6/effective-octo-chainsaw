package account_microservice

import (
	"context"
	"testing"
)

type TestData struct {
	s            AccountService
	mockAccount  Account
	otherAccount Account
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

	// autre compte mocké
	otherAccount := Account{
		"jsvkjjsbfvkjnsdfvknsdlfjkvnsdlkfnvuhdfovhj",
		"Valentin",
		"Roche",
		"+33678239003",
		"jsadjcb@jnskd.com",
		100000,
	}

	return TestData{
		s,
		mockAccount,
		otherAccount,
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

func TestGetAccountByID(t *testing.T) {
	testData := NewTestData()

	// Test avec un paramètre vide
	_, err := testData.s.GetAccountByID(context.TODO(), "")

	if err == nil {
		t.Errorf("Passed empty account param to getByID function, should have raised an error")
	}

	// Test avec un id valide
	result, err := testData.s.GetAccountByID(context.TODO(), testData.mockAccount.ClientID)
	if err != nil {
		t.Errorf("Valid account, method should not fail")
	}

	if result != testData.mockAccount {
		t.Errorf("Returned account is not the same as the on specified")
	}
}

func TestUpdate(t *testing.T) {
	testData := NewTestData()

	// Test avec le param id vide
	errEmptyID := testData.s.Update(context.TODO(), "", testData.mockAccount)

	if errEmptyID == nil {
		t.Errorf("Passed empty id to Update function, should have raised an error")
	}

	// Test avec le param account vide
	errEmptyAccount := testData.s.Update(context.TODO(), testData.mockAccount.ClientID, Account{})

	if errEmptyAccount == nil {
		t.Errorf("Passed empty account to Update function, should have raised an error")
	}

	// Test avec des id inconsistants (l'id en param est différent de celui du param)
	errInconsistentIDs := testData.s.Update(context.TODO(), "lmao", testData.mockAccount)

	if errInconsistentIDs == nil {
		t.Errorf("Passed inconsistent IDs to Update function, should have raised an error")
	}

	if errInconsistentIDs != ErrInconsistentIDs {
		t.Errorf("Function should have raised an inconsistent id error and instead raised another one")
	}

	// Test avec un fonctionnement valide
	err := testData.s.Update(context.TODO(), testData.mockAccount.ClientID, testData.otherAccount)
	if err != nil {
		t.Errorf("Valid account, method should not fail")
	}

	dbResult, err := testData.s.GetAccountByID(context.TODO(), testData.mockAccount.ClientID)

	if err != nil {
		t.Errorf("Error during fetch")
	}

	if dbResult != testData.otherAccount && dbResult == testData.mockAccount {
		t.Errorf("Update did not go through")
	}

	if dbResult != testData.otherAccount && dbResult != testData.mockAccount {
		t.Errorf("Fetched result is not the test account we wanted")
	}
}

func TestDelete(t *testing.T) {
	testData := NewTestData()

	// Cas avec un ID vide
	errEmptyID := testData.s.Delete(context.TODO(), "")
	if errEmptyID == nil {
		t.Errorf("Passed empty ID as param, should have failed")
	}

	// Cas avec un ID invalide
	errInvalidID := testData.s.Delete(context.TODO(), "sjdhfbviujas")
	if errInvalidID == nil {
		t.Errorf("Passed wrong ID, should raise an error")
	}

	// Cas avec un ID valide
	err := testData.s.Delete(context.TODO(), testData.mockAccount.ClientID)
	if err != nil {
		t.Errorf("Passed existing ID, should not raise an error")
	}
}

func TestGetAmountForID(t *testing.T) {
	testData := NewTestData()

	// Cas avec un ID vide
	_, errEmptyID := testData.s.GetAmountForID(context.TODO(), "")
	if errEmptyID == nil {
		t.Errorf("Passed empty ID as param, should have failed")
	}

	// Cas avec un ID invalide
	_, errInvalidID := testData.s.GetAmountForID(context.TODO(), "sjdhfbviujas")
	if errInvalidID == nil {
		t.Errorf("Passed wrong ID, should raise an error")
	}

	// Cas avec un ID valide
	_, err := testData.s.GetAmountForID(context.TODO(), testData.mockAccount.ClientID)
	if err != nil {
		t.Errorf("Passed existing ID, should not raise an error")
	}

}
