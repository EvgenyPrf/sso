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

func TestRegisterLogin_Register_DuplicateRegistration(t *testing.T) {
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

	respReg, err = s.AuthClient.Register(ctx,
		&ssov1.RegisterRequest{
			Email:    email,
			Password: password,
		},
	)

	assert.Empty(t, respReg.GetUserId())
	assert.Error(t, err)
	assert.ErrorContains(t, err, "User already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, s := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password: value is required",
		},
		{
			name:        "Register with empty email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "email: value is required",
		},
		{
			name:        "Register with both empty",
			email:       "",
			password:    "",
			expectedErr: "email: value is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.AuthClient.Register(ctx,
				&ssov1.RegisterRequest{
					Email:    tt.email,
					Password: tt.password,
				},
			)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, s := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Login with empty password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "password: value is required",
		},
		{
			name:        "Login with empty email",
			email:       "",
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "email: value is required",
		},
		{
			name:        "Login with both empty",
			email:       "",
			password:    "",
			appID:       appID,
			expectedErr: "email: value is required",
		},
		{
			name:        "Login with non-matching password",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "Invalid credentials",
		},
		{
			name:        "Login without app ID",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       emptyAppID,
			expectedErr: "app_id: value is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.AuthClient.Register(ctx,
				&ssov1.RegisterRequest{
					Email:    gofakeit.Email(),
					Password: randomFakePassword(),
				},
			)

			require.NoError(t, err)

			_, err = s.AuthClient.Login(ctx,
				&ssov1.LoginRequest{
					Email:    tt.email,
					Password: tt.password,
					AppId:    tt.appID,
				},
			)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLength)
}
