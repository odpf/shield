package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type TokenJWKSHandler struct {
	keySet jwk.Set
}

func NewTokenJWKSHandler(RSAKeyPath string) (*TokenJWKSHandler, error) {
	// if no config provided for path, ignore jwks
	if len(RSAKeyPath) == 0 {
		return &TokenJWKSHandler{
			keySet: jwk.NewSet(),
		}, nil
	}

	privateSet, err := jwk.ReadFile(RSAKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read rsa key path: %w", err)
	}

	// convert private to public
	publicKeySet := jwk.NewSet()
	for iter := privateSet.Keys(context.Background()); iter.Next(context.Background()); {
		pair := iter.Pair()
		key := pair.Value.(jwk.Key)

		pubKey, err := key.PublicKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate public key from private rsa: %w", err)
		}
		publicKeySet.AddKey(pubKey)
	}
	return &TokenJWKSHandler{
		keySet: publicKeySet,
	}, nil
}

// ServeHTTP at "/jwks.json" with rsa public keys endpoint
// generate keys via shield cli "shield server gen rsa"
func (t TokenJWKSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(t.keySet); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
