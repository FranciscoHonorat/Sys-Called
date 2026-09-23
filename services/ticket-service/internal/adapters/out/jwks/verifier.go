package jwks

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
)

var _ out.TokenVerifier = (*Verifier)(nil)

const (
	expectedIssuer     = "employees-service"
	expectedAudience   = "sys-called"
	MinRefreshInterval = 5 * time.Second
)

var (
	errUnknownKey      = errors.New("jwks: unknown signing key")
	errUnexpectedReply = errors.New("jwks: unexpected reply from the key set endpoint")
)

type claims struct {
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Option func(*Verifier)

func WithClock(now func() time.Time) Option {
	return func(v *Verifier) { v.now = now }
}

type Verifier struct {
	url    string
	client *http.Client
	parser *jwt.Parser
	now    func() time.Time

	refreshing sync.Mutex
	mu         sync.RWMutex
	keys       map[string]ed25519.PublicKey
	fetchedAt  time.Time
}

func NewVerifier(url string, client *http.Client, options ...Option) *Verifier {
	v := &Verifier{
		url:    url,
		client: client,
		now:    time.Now,
		parser: jwt.NewParser(
			jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}),
			jwt.WithIssuer(expectedIssuer),
			jwt.WithAudience(expectedAudience),
			jwt.WithExpirationRequired(),
		),
	}
	for _, option := range options {
		option(v)
	}
	return v
}

func (v *Verifier) Verify(ctx context.Context, token string) (actor.Actor, error) {
	parsed := &claims{}
	_, err := v.parser.ParseWithClaims(token, parsed, func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		return v.key(ctx, kid)
	})
	if err != nil {
		return actor.Actor{}, err
	}
	return actor.New(parsed.Subject, parsed.Name, parsed.Role)
}

func (v *Verifier) key(ctx context.Context, kid string) (ed25519.PublicKey, error) {
	if key, ok := v.cachedKey(kid); ok {
		return key, nil
	}

	v.refreshing.Lock()
	defer v.refreshing.Unlock()

	if key, ok := v.cachedKey(kid); ok {
		return key, nil
	}
	if !v.canRefresh() {
		return nil, errUnknownKey
	}
	if err := v.refresh(ctx); err != nil {
		return nil, err
	}
	if key, ok := v.cachedKey(kid); ok {
		return key, nil
	}
	return nil, errUnknownKey
}

func (v *Verifier) cachedKey(kid string) (ed25519.PublicKey, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	key, ok := v.keys[kid]
	return key, ok
}

func (v *Verifier) canRefresh() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.keys == nil || v.now().Sub(v.fetchedAt) >= MinRefreshInterval
}

func (v *Verifier) refresh(ctx context.Context) error {
	keys, err := v.fetchKeys(ctx)

	v.mu.Lock()
	defer v.mu.Unlock()
	v.fetchedAt = v.now()
	if err != nil {
		return err
	}
	v.keys = keys
	return nil
}

func (v *Verifier) fetchKeys(ctx context.Context) (map[string]ed25519.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errUnexpectedReply
	}

	var set struct {
		Keys []struct {
			ID string `json:"kid"`
			X  string `json:"x"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return nil, err
	}

	keys := make(map[string]ed25519.PublicKey, len(set.Keys))
	for _, k := range set.Keys {
		x, err := base64.RawURLEncoding.DecodeString(k.X)
		if err != nil {
			return nil, err
		}
		keys[k.ID] = ed25519.PublicKey(x)
	}
	return keys, nil
}
