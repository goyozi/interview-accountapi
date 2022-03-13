package interviewaccountapi

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var client = NewAccountsClient("http://localhost:8080")

var country string = "GB"
var classification = "Personal"
var version = int64(0)
var account = AccountData{
	ID:             "0568dec3-e8af-433b-906e-d9e616b30843",
	OrganisationID: "c50b6049-ef55-4c5b-8923-403758fc7170",
	Type:           "accounts",
	Attributes: &AccountAttributes{
		Country:               &country,
		BaseCurrency:          "GBP",
		BankID:                "400302",
		BankIDCode:            "GBDSC",
		AccountNumber:         "10000004",
		Iban:                  "GB28NWBK40030212764204",
		Bic:                   "NWBKGB42",
		Name:                  []string{"Samantha Holder"},
		AccountClassification: &classification},
	Version: &version}

func TestCreateFetchDelete(t *testing.T) {
	res, err := client.Create(&account)

	assert.Nil(t, err)
	assert.Equal(t, account, *res)

	res, err = client.Fetch(account.ID)

	assert.Nil(t, err)
	assert.Equal(t, account, *res)

	err = client.Delete(&account)
	assert.Nil(t, err)
}

func TestCreateWrongRequest(t *testing.T) {
	_, err := client.Create(&AccountData{})

	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "validation failure list")
	}
}

func TestFetchInvalidId(t *testing.T) {
	_, err := client.Fetch("invalid")

	if assert.NotNil(t, err) {
		assert.Equal(t, "id is not a valid uuid", err.Error())
	}
}

func TestFetchNonExistentAccount(t *testing.T) {
	uuid := uuid.NewString()

	_, err := client.Fetch(uuid)

	if assert.NotNil(t, err) {
		assert.Equal(t, "record "+uuid+" does not exist", err.Error())
	}
}

func TestDeleteInvalidId(t *testing.T) {
	invalidAccount := account
	invalidAccount.ID = "invalid"

	err := client.Delete(&invalidAccount)

	if assert.NotNil(t, err) {
		assert.Equal(t, "id is not a valid uuid", err.Error())
	}
}

func TestDeleteNonExistentAccount(t *testing.T) {
	err := client.Delete(&account)

	if assert.NotNil(t, err) {
		assert.Equal(t, "record 0568dec3-e8af-433b-906e-d9e616b30843 does not exist", err.Error())
	}
}
