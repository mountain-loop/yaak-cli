package yaakcli

import (
	"errors"
	"github.com/zalando/go-keyring"
)

const keyringUser = "yaak"

var keyringService = prodStagingDevStr("app.yaak.cli.Token", "app.yaak.cli.staging.Token", "app.yaak.cli.dev.Token")

func getAuthToken() (bool, string, error) {
	token, err := keyring.Get(keyringService, keyringUser)
	if errors.Is(err, keyring.ErrNotFound) {
		return false, "", nil
	} else if err != nil {
		return false, "", err
	}

	return true, token, nil
}

func storeAuthToken(token string) error {
	return keyring.Set(keyringService, keyringUser, token)
}

func deleteAuthToken() error {
	err := keyring.Delete(keyringService, keyringUser)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	return err
}
