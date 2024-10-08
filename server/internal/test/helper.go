package test

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"gatekeeper/internal/server"
	"gatekeeper/pkg/jwt_provider"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func SendTestRequest(t *testing.T, s server.Server, method, path string, headers map[string]string, body any) *httptest.ResponseRecorder {
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(method, path, bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res := httptest.NewRecorder()
	s.EchoInst.ServeHTTP(res, req)
	return res
}

func GenerateWalletAddress(t *testing.T) (string, *ecdsa.PrivateKey) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	publicKey := privateKey.Public()
	address := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	return address.Hex(), privateKey
}

func GenerateProofToken(t *testing.T, jwtProvider jwt_provider.Provider, walletAddress string, expiredAt time.Time) string {
	proofToken, err := jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   walletAddress,
		ExpiresAt: &jwt.NumericDate{Time: expiredAt},
	})
	require.NoError(t, err)
	return proofToken
}

func ReadBody[T any](t *testing.T, bodyBuffer *bytes.Buffer) T {
	bodyBytes, err := io.ReadAll(bodyBuffer)
	require.NoError(t, err)

	var body T
	require.NoError(t, json.Unmarshal(bodyBytes, &body))

	return body
}
