package yaakcli

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
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

	ctx := context.WithValue(r.Context(), oauth2.HTTPClient, &YaakHttpClient)

	token, err := h.OAuthConfig.Exchange(ctx, code, oauth2.VerifierOption(h.CodeVerifier))
	if err != nil {
		return "", fmt.Errorf("could not exchange code for token: %w", err)
	}

	return token.AccessToken, nil
}
