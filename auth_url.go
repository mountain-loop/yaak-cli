package yaakcli

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/oauth2"
)

type AuthURL struct {
	URL          string
	State        string
	CodeVerifier string
}

func (u *AuthURL) String() string {
	return u.URL
}

func AuthorizationURL(config *oauth2.Config) (*AuthURL, error) {
	codeVerifier, verifierErr := randomBytesInHex(32) // 64-character string here
	if verifierErr != nil {
		return nil, fmt.Errorf("could not create a code verifier: %w", verifierErr)
	}
	sha2 := sha256.New()
	_, _ = io.WriteString(sha2, codeVerifier)
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))

	state, stateErr := randomBytesInHex(24)
	if stateErr != nil {
		return nil, fmt.Errorf("could not generate random state: %w", stateErr)
	}

	authUrl := config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	return &AuthURL{
		URL:          authUrl,
		State:        state,
		CodeVerifier: codeVerifier,
	}, nil
}

func randomBytesInHex(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf("could not generate %d random bytes: %w", count, err)
	}

	return hex.EncodeToString(buf), nil
}
