package tests

import (
	ssov1 "github.com/EvgenyPrf/protos/gen/go/sso"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sso/tests/suite"
	"testing"
	"time"
)

const (
	emptyAppID        = 0
	appID             = 1
	appSecret         = "test-secret"
	passDefaultLength = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, s := suite.New(t)
	email := gofakeit.Email()
	password := randomFakePassword()

	respReg, err := s.AuthClient.Register(ctx,
		&ssov1.RegisterRequest{
			Email:    email,
			Password: password,
		},
	)

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := s.AuthClient.Login(ctx,
		&ssov1.LoginRequest{
			Email:    email,
			Password: password,
			AppId:    appID,
		},
	)

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	loginTime := time.Now()
	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(s.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLength)
}
