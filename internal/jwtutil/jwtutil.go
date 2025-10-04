package jwtutil

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Decode(token string) (map[string]any, map[string]any, error) {
	parser := jwt.NewParser()
	t, _, err := parser.ParseUnverified(token, jwt.MapClaims{})
	if err != nil { return nil, nil, err }
	h := t.Header
	p, ok := t.Claims.(jwt.MapClaims)
	if !ok { return nil, nil, errors.New("invalid claims") }
	return h, p, nil
}

func VerifyHMAC(token string, secret []byte) error {
	_, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %T", t.Method)
		}
		return secret, nil
	})
	return err
}

func VerifyWithJWKS(token, jwksURL string) error {
	set, err := fetchJWKS(jwksURL)
	if err != nil { return err }
	_, err = jwt.Parse(token, func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		k := set.Key(kid)
		if k == nil { return nil, fmt.Errorf("no JWK for kid=%s", kid) }
		switch k.Kty {
		case "RSA":
			return k.RSA()
		case "EC":
			return k.ECDSA()
		default:
			return nil, fmt.Errorf("unsupported kty=%s", k.Kty)
		}
	})
	return err
}

func SignHMAC(payloadJSON []byte, secret []byte, alg string) (string, error) {
	var claims jwt.MapClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil { return "", err }
	var method jwt.SigningMethod = jwt.SigningMethodHS256
	switch alg {
	case "HS256": method = jwt.SigningMethodHS256
	case "HS384": method = jwt.SigningMethodHS384
	case "HS512": method = jwt.SigningMethodHS512
	default: return "", fmt.Errorf("unsupported alg %s for HMAC", alg)
	}
	t := jwt.NewWithClaims(method, claims)
	return t.SignedString(secret)
}

/* ---- Minimal JWKS ---- */

type jwkSet struct { Keys []jwk `json:"keys"` }
type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	// RSA
	N string `json:"n"`
	E string `json:"e"`
	// EC
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

func (s *jwkSet) Key(kid string) *jwk {
	for i := range s.Keys {
		if s.Keys[i].Kid == kid { return &s.Keys[i] }
	}
	// if kid empty, return first
	if kid == "" && len(s.Keys) > 0 { return &s.Keys[0] }
	return nil
}

func (k *jwk) RSA() (*rsa.PublicKey, error) {
	nb, err := b64url(k.N); if err != nil { return nil, err }
	eb, err := b64url(k.E); if err != nil { return nil, err }
	e := big.NewInt(0)
	if len(eb) < 8 {
		// small exponent like 65537
		var ei uint64
		for _, b := range eb { ei = (ei << 8) | uint64(b) }
		return &rsa.PublicKey{ N: big.NewInt(0).SetBytes(nb), E: int(ei) }, nil
	}
	// big exponent
	e.SetBytes(eb)
	return &rsa.PublicKey{ N: big.NewInt(0).SetBytes(nb), E: int(e.Int64()) }, nil
}

func (k *jwk) ECDSA() (*ecdsa.PublicKey, error) {
	xb, err := b64url(k.X); if err != nil { return nil, err }
	yb, err := b64url(k.Y); if err != nil { return nil, err }
	var curve elliptic.Curve
	switch k.Crv {
	case "P-256": curve = elliptic.P256()
	case "P-384": curve = elliptic.P384()
	case "P-521": curve = elliptic.P521()
	default: return nil, fmt.Errorf("unsupported EC curve %s", k.Crv)
	}
	x := big.NewInt(0).SetBytes(xb)
	y := big.NewInt(0).SetBytes(yb)
	return &ecdsa.PublicKey{ Curve: curve, X: x, Y: y }, nil
}

func fetchJWKS(url string) (*jwkSet, error) {
	res, err := http.Get(url)
	if err != nil { return nil, err }
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var set jwkSet
	if err := json.Unmarshal(body, &set); err != nil { return nil, err }
	if len(set.Keys) == 0 { return nil, errors.New("empty JWKS") }
	return &set, nil
}

func b64url(s string) ([]byte, error) {
	// padless base64url
	return base64.RawURLEncoding.DecodeString(s)
}

// helper to avoid unused import of encoding/binary (keep here if needed later)
var _ = binary.LittleEndian
