package yaakcli

import (
	"crypto/subtle"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
)

type OAuthRedirectHandler struct {
	State        string
	CodeVerifier string
	OAuthConfig  *oauth2.Config
}

func (h *OAuthRedirectHandler) ExchangeCode(r *http.Request) (string, error) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	if subtle.ConstantTimeCompare([]byte(h.State), []byte(state)) == 0 {
		return "", fmt.Errorf("invalid state")
	}

	if code == "" {
		return "", fmt.Errorf("missing code")
	}

	token, err := h.OAuthConfig.Exchange(
		r.Context(),
		code,
		oauth2.VerifierOption(h.CodeVerifier),
	)
	if err != nil {
		return "", fmt.Errorf("could not exchange code for token: %w", err)
	}

	return token.AccessToken, nil
}
