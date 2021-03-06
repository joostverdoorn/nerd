package jwt

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nerdalize/nerd/nerd/conf"
	"github.com/nerdalize/nerd/nerd/utils"
)

type staticSession struct {
	ss *conf.SessionSnapshot
}

//Read returns a snapshot of the session file
func (s *staticSession) Read() (*conf.SessionSnapshot, error) {
	return s.ss, nil
}

//WriteJWT writes the jwt object to the session file
func (s *staticSession) WriteJWT(jwt, refreshToken string) error {
	return fmt.Errorf("not implemented")
}

//WriteOAuth writes the oauth object to the session file
func (s *staticSession) WriteOAuth(accessToken, refreshToken string, expiration time.Time, scope, tokenType string) error {
	return fmt.Errorf("not implemented")
}

//WriteProject writes the project object to the session file
func (s *staticSession) WriteProject(name, awsRegion string) error {
	return fmt.Errorf("not implemented")
}
func TestEnvProvider(t *testing.T) {
	key := testkey(t)
	pub, _ := key.Public().(*ecdsa.PublicKey)
	now := time.Now().Unix()
	refreshedClaims := &jwt.StandardClaims{
		ExpiresAt: now + minute*10,
	}
	refreshedToken := getToken(key, refreshedClaims, t)
	session := &staticSession{
		ss: &conf.SessionSnapshot{},
	}
	client := &tokenClient{
		token: refreshedToken,
	}

	prov := NewEnvProvider(pub, session, client)
	prov.ExpireWindow = 0

	t.Run("normal", func(t *testing.T) {
		claims := &jwt.StandardClaims{
			ExpiresAt: now + minute*5,
		}
		token := getToken(key, claims, t)
		os.Setenv("NERD_JWT", token)
		ret, err := prov.Retrieve()
		utils.OK(t, err)
		utils.Equals(t, token, ret)
	})

	t.Run("noToken", func(t *testing.T) {
		os.Setenv("NERD_JWT", "")
		ret, err := prov.Retrieve()
		utils.Assert(t, err != nil, "expected error because no token was set")
		utils.Assert(t, strings.Contains(err.Error(), "not set"), "expected error because no token was set", err)
		utils.Equals(t, "", ret)
	})

	claimsExp := &jwt.StandardClaims{
		ExpiresAt: now - minute*5,
	}
	tokenExp := getToken(key, claimsExp, t)
	t.Run("expired", func(t *testing.T) {
		os.Setenv("NERD_JWT", tokenExp)
		ret, err := prov.Retrieve()
		utils.Assert(t, err != nil, "expected token to be expired")
		utils.Assert(t, strings.Contains(err.Error(), "expired"), "expected token to be expired", err)
		utils.Equals(t, "", ret)
	})

	t.Run("refresh", func(t *testing.T) {
		os.Setenv("NERD_JWT", tokenExp)
		os.Setenv("NERD_JWT_REFRESH_TOKEN", "abc")
		ret, err := prov.Retrieve()
		utils.OK(t, err)
		utils.Equals(t, refreshedToken, ret)
		utils.Equals(t, refreshedClaims.ExpiresAt, prov.expiration.Unix())
	})
}
