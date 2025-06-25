package crypto

import (
	"api/db"
	"api/repos/operatorsRepo"
	"crypto/rsa"
	"encoding/hex"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	cry "github.com/zarkones/xena-crypto"
)

func TestJWTVerification(t *testing.T) {
	if err := db.Init(""); err != nil {
		t.Log(err)
		t.FailNow()
	}

	username := "bob"

	operator, _, hexEncodedPrivateKey, err := CreateAdminOperator(username)
	if err != nil {
		t.Log(err)
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

	if err := operatorsRepo.Insert(&operator); err != nil {
		t.Log(err)
		t.FailNow()
	}

	authToken, err := generateJwt(username, privateKey)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	verifiedUsername, _, err := VerifyToken(authToken)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if verifiedUsername != username {
		t.Log("unexpected parsed username")
		t.FailNow()
	}
}

func generateJwt(username string, privateKey *rsa.PrivateKey) (token string, err error) {
	t := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"u": username,
	})

	return t.SignedString(privateKey)
}
