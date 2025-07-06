package crypto

import (
	"encoding/hex"
	"errors"
	"testing"

	cry "github.com/zarkones/xena-crypto"
)

func TestOperatorAccountCreation(t *testing.T) {
	username := "bob"

	operator, phrase, hexEncodedPrivateKey, err := CreateAdminOperator(username)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(phrase) != SECURE_AMOUNT_OF_WORDS {
		t.Log("recovery wordphrase has invalid length")
	}

	if operator.Username != username {
		t.Log("username is screwed up")
		t.FailNow()
	}

	pemEncodedPrivateKey, err := hex.DecodeString(hexEncodedPrivateKey)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	privateKey, err := cry.ImportPrivKeyPEM(pemEncodedPrivateKey)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	secret, err := cry.GenSecret()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	encryptedData, err := cry.EncryptRSAOAEPEncodeHex(privateKey.PublicKey, string(secret))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	decryptedSecret, err := cry.DecryptRSAOAEPDecodeHex(*privateKey, encryptedData)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if string(secret) != decryptedSecret {
		t.Log("secret which was encrypted and then decrypted appears malformed")
		t.Log(string(secret))
		t.Log(decryptedSecret)
		t.FailNow()
	}
}

func TestOperatorAccountCreationFailsOnInvalidUsername(t *testing.T) {
	if _, _, _, err := CreateAdminOperator(""); !errors.Is(err, ErrUsername) {
		t.Log("error is not the one which was expected")
	}

	if _, _, _, err := CreateAdminOperator("ssssfgdddfsdfsdfsdfsdfsdfsdfsdfghgd"); !errors.Is(err, ErrUsername) {
		t.Log("error is not the one which was expected")
	}
}
