package crypto

import (
	"api/models"
	"api/repos/permissionsRepo"
	"encoding/hex"
	"errors"
	"strings"

	cry "github.com/zarkones/xena-crypto"
)

const (
	SECURE_AMOUNT_OF_WORDS = 26
)

var (
	ErrUsername = errors.New("username is empty or too long, must be equal or less to 32 characters")
)

// CreateAdminOperator creates a really powerful administrative account with all permissions assigned.
// Recommendation is to have only one such account, and make second account for yourself which is less
// powerful, so that you can use it as a daily driver.
func CreateAdminOperator(username string) (operator models.Operator, recoveryWordphrase []string, hexEncodedPrivateKey string, err error) {
	if len(username) == 0 || len(username) > 32 {
		return operator, nil, "", ErrUsername
	}

	privateKey, err := cry.GenPrivKey()
	if err != nil {
		return operator, nil, "", err
	}

	serializedPrivateKey, err := cry.PrivKeyToPEM(privateKey)
	if err != nil {
		return operator, nil, "", err
	}

	serializedPublicKey, err := cry.PubKeyToPEM(&privateKey.PublicKey)
	if err != nil {
		return operator, nil, "", err
	}

	wordPhrase := []string{}
	for i := 0; i < SECURE_AMOUNT_OF_WORDS; i++ {
		word, err := RandomPopularWord()
		if err != nil {
			return operator, nil, "", err
		}
		wordPhrase = append(wordPhrase, word)
	}

	encryptedPrivateKey, err := cry.EncryptAESEncodeHex([]byte(serializedPrivateKey), []byte(NormalizeWordPhrase(wordPhrase)))
	if err != nil {
		return operator, nil, "", err
	}

	allowedPermissions := []models.PermissionKey{
		models.PERMISSION_AGENTS_LIST,
		models.PERMISSION_AGENTS_STATS,
		models.PERMISSION_AGENTS_LIST_MESSAGES,
		models.PERMISSION_AGENTS_INSERT_MESSAGE,
		models.PERMISSION_CHAT_LIST_CHANNELS,
		models.PERMISSION_CHAT_LIST_CHANNEL_MESSAGES,
		models.PERMISSION_CHAT_INSERT_CHANNEL,
		models.PERMISSION_CHAT_INSERT_CHANNEL_MESSAGE,
		models.PERMISSION_CHAT_DELETE_CHANNEL,
		models.PERMISSION_CHAT_DELETE_CHANNEL_MESSAGE,
		models.PERMISSION_OPERATORS_INSERT,
		models.PERMISSION_OPERATORS_DELETE,
		models.PERMISSION_OPERATORS_LIST,
		models.PERMISSION_INSERT,
		models.PERMISSION_DELETE,
		models.PERMISSION_LIST,
	}

	for _, key := range allowedPermissions {
		if err := permissionsRepo.Insert(&models.Permission{
			Username: username,
			Key:      key,
		}); err != nil {
			return operator, nil, "", err
		}
	}

	return models.Operator{
			Username:               username,
			PublicKeyHex:           hex.EncodeToString([]byte(serializedPublicKey)),
			EncryptedPrivateKeyHex: encryptedPrivateKey,
		},
		wordPhrase,
		hex.EncodeToString([]byte(serializedPrivateKey)),
		nil
}

func NormalizeWordPhrase(phrase []string) (normalized string) {
	if phrase == nil || len(phrase) == 0 {
		return ""
	}
	return strings.Join(phrase, "")
}
