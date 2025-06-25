package crypto

import (
	"api/models"
	"api/repos/operatorsRepo"
	"api/repos/permissionsRepo"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	cry "github.com/zarkones/xena-crypto"
)

var (
	ErrMissingAuthHeader       = errors.New("missing Authorization header")
	ErrInvalidJWTSegments      = errors.New("invalid size of jwt segments")
	ErrNilInsecurePayload      = errors.New("insecure parsed token payload appears as nil")
	ErrInvalidInsecureUsername = errors.New("insecure parsed token username is invalid")
	ErrInvalidSigningAlg       = errors.New("invalid signing algorithm")
	ErrInvalidSig              = errors.New("invalid token's signature")
	ErrInvalidClaims           = errors.New("claims casting did not go well")
)

func VerifyToken(rawToken string) (username string, permissions []models.Permission, err error) {
	if rawToken == "" {
		return "", nil, ErrMissingAuthHeader
	}

	segments := strings.Split(rawToken, ".")
	if len(segments) != 3 {
		return "", nil, ErrInvalidJWTSegments
	}

	insecurePayloadJsonStr, err := base64.RawStdEncoding.DecodeString(segments[1])
	if err != nil {
		return "", nil, err
	}

	var insecurePayload map[string]string

	if err := json.Unmarshal(insecurePayloadJsonStr, &insecurePayload); err != nil {
		return "", nil, err
	}

	if insecurePayload == nil {
		return "", nil, ErrNilInsecurePayload
	}

	insecureUsername := fmt.Sprint(insecurePayload["u"])
	if len(insecureUsername) == 0 || len(insecureUsername) > 32 {
		return "", nil, ErrInvalidInsecureUsername
	}

	operator, err := operatorsRepo.Get(insecureUsername)
	if err != nil {
		return "", nil, err
	}

	token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidSigningAlg
		}

		publicKeyPem, err := hex.DecodeString(operator.PublicKeyHex)
		if err != nil {
			return nil, err
		}

		publicKey, err := cry.ImportPubKeyPEM(publicKeyPem)
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	})
	if err != nil {
		return "", nil, err
	}

	if !token.Valid {
		return "", nil, ErrInvalidSig
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, ErrInvalidClaims
	}

	verifiedUsername := fmt.Sprint(claims["u"])

	permissions, err = permissionsRepo.GetMultipleByUsername(verifiedUsername)

	return verifiedUsername, permissions, err
}
