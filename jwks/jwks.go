package jwks

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWKS (JSON Web Key Set) struct to hold cached issuer.
type JWKS struct {
	url   string
	mu    sync.RWMutex
	keys  map[string]*rsa.PublicKey
	exp   time.Time
	ttl   time.Duration
	httpc *http.Client
}

// NewJWKSCache creates a new cached instance.
func NewJWKSCache(url string, ttl time.Duration) *JWKS {
	return &JWKS{
		url:   url,
		ttl:   ttl,
		keys:  map[string]*rsa.PublicKey{},
		httpc: &http.Client{Timeout: 5 * time.Second},
	}
}

// get returns keys if cached,
// otherwise requests them from issuer and returns them.
func (c *JWKS) get(kid string) (*rsa.PublicKey, error) {
	c.mu.RLock()
	pk, ok := c.keys[kid]
	expired := time.Now().After(c.exp)
	c.mu.RUnlock()
	if ok && !expired {
		return pk, nil
	}

	resp, err := c.httpc.Get(c.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type jwks struct {
		Kty string `json:"kty"`
		Kid string `json:"kid"`
		N   string `json:"n"`
		E   string `json:"e"`
	}

	type doc struct {
		Keys []jwks `json:"keys"`
	}

	var p doc
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}

	tmp := make(map[string]*rsa.PublicKey, len(p.Keys))
	for _, k := range p.Keys {
		if k.Kty != "RSA" || k.N == "" || k.E == "" {
			continue
		}
		nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
		if err != nil {
			continue
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
		if err != nil || len(eBytes) == 0 {
			continue
		}
		var eInt int
		if len(eBytes) < 4 {
			eInt = 0
			for _, b := range eBytes {
				eInt = (eInt << 8) | int(b)
			}
		} else {
			eInt = int(binary.BigEndian.Uint32(eBytes))
		}
		tmp[k.Kid] = &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: eInt,
		}
	}

	c.mu.Lock()
	c.keys = tmp
	c.exp = time.Now().Add(c.ttl)
	pk, ok = c.keys[kid]
	c.mu.Unlock()
	if !ok {
		return nil, errors.New("kid not found")
	}
	return pk, nil
}

// Keyfunc is used to validate each key in the token.
func (c *JWKS) Keyfunc(token *jwt.Token) (any, error) {
	kid, _ := token.Header["kid"].(string)
	if kid == "" {
		return nil, errors.New("missing kid")
	}
	return c.get(kid)
}
