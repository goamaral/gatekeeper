package server_test

import (
	"context"
	"gatekeeper/internal"
	"gatekeeper/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChallengeController_Issue(t *testing.T) {
	s := server.NewServer(internal.NewTestInjector(t))
	res := internal.SendTestRequest(t, s,
		http.MethodPost, "/v1/challenges/issue", server.ChallengeController_IssueRequest{WalletAddress: "WalletAddress"},
	)
	require.Equal(t, http.StatusOK, res.Code)
}

func TestChallengeController_Validate(t *testing.T) {
	walletAddressA, privateKeyA := internal.GenerateWalletAddress(t)
	challengeTokenA, err := server.GenerateChallengeToken()
	require.NoError(t, err)
	challengeA := server.ChallengeMessagePrefix + challengeTokenA
	challengeHashA := crypto.Keccak256Hash([]byte(challengeA)).Bytes()
	signatureA, err := crypto.Sign(challengeHashA, privateKeyA)
	require.NoError(t, err)

	_, privateKeyB := internal.GenerateWalletAddress(t)
	challengeTokenB, err := server.GenerateChallengeToken()
	require.NoError(t, err)
	challengeB := server.ChallengeMessagePrefix + challengeTokenB
	challengeHashB := crypto.Keccak256Hash([]byte(challengeB)).Bytes()
	signatureB, err := crypto.Sign(challengeHashB, privateKeyB)
	require.NoError(t, err)

	sendReq := func(t *testing.T, challenge, signature string, expiredAt time.Time) *httptest.ResponseRecorder {
		s := server.NewServer(internal.NewTestInjector(t))

		_, err = s.ChallengeCtrl.DB.ExecContext(context.Background(),
			"INSERT INTO challenges (wallet_address, token, expired_at) VALUES (?, ?, ?)",
			walletAddressA, challengeTokenA, expiredAt,
		)
		require.NoError(t, err)

		return internal.SendTestRequest(t, s,
			http.MethodPost, "/v1/challenges/validate", server.ChallengeController_ValidateRequest{
				Challenge: challenge,
				Signature: signature,
			},
		)
	}

	t.Run("Success", func(t *testing.T) {
		res := sendReq(t, challengeA, hexutil.Encode(signatureA), time.Now().UTC().Add(time.Minute))
		require.Equal(t, http.StatusNoContent, res.Code)
	})

	t.Run("ChallengeDoesNotExist", func(t *testing.T) {
		res := sendReq(t, challengeB, hexutil.Encode(signatureB), time.Now().UTC().Add(time.Minute))
		require.Equal(t, http.StatusUnprocessableEntity, res.Code)
		body := internal.ReadBody[server.ErrorResponse](t, res.Body)
		assert.Equal(t, server.MsgChallengeDoesNotExistOrExpired, body.Error)
	})

	t.Run("ChallengeExpired", func(t *testing.T) {
		res := sendReq(t, challengeA, hexutil.Encode(signatureA), time.Now().UTC().Add(-time.Minute))
		require.Equal(t, http.StatusUnprocessableEntity, res.Code)
		body := internal.ReadBody[server.ErrorResponse](t, res.Body)
		assert.Equal(t, server.MsgChallengeDoesNotExistOrExpired, body.Error)
	})

	t.Run("InvalidSignature", func(t *testing.T) {
		res := sendReq(t, challengeA, hexutil.Encode(signatureB), time.Now().UTC().Add(time.Minute))
		require.Equal(t, http.StatusUnprocessableEntity, res.Code)
		body := internal.ReadBody[server.ErrorResponse](t, res.Body)
		assert.Equal(t, server.MsgInvalidSignature, body.Error)
	})
}