//go:build integration

package service_test

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	authv1 "neupaneanish.com.np/api/internal/protobuf/auth/v1"
	"neupaneanish.com.np/api/internal/redis"
	"neupaneanish.com.np/api/internal/service"
)

func TestVerification(t *testing.T) {
	t.Parallel()

	data := &service.ForgetPasswordSession{
		Key:    rand.Text(),
		ExAt:   time.Now().Add(service.SessionExpiry),
		UserID: uuid.NewString(),
		Code:   "A1B2C3D4",
	}

	hSetErr := redis.HSet[service.ForgetPasswordSession](
		t.Context(),
		service.ForgetPasswordSessionPrefix,
		data,
		cfg.Client,
	)
	require.NoError(t, hSetErr)

	t.Run("Invalid session", func(t *testing.T) {
		t.Parallel()
		req := &authv1.VerificationRequest{Session: rand.Text(), Code: data.Code}

		response, responseErr := authServiceClient.Verification(t.Context(), req)
		require.Error(t, responseErr)
		assert.Nil(t, response)
		st, _ := status.FromError(responseErr)
		assert.Equal(t, codes.Aborted, st.Code())
	})

	t.Run("Valid session invalid code", func(t *testing.T) {
		t.Parallel()
		req := &authv1.VerificationRequest{Session: data.Key, Code: "12345678"}

		response, responseErr := authServiceClient.Verification(t.Context(), req)
		require.Error(t, responseErr)
		assert.Nil(t, response)
		st, _ := status.FromError(responseErr)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("Valid session and code", func(t *testing.T) {
		t.Parallel()
		req := &authv1.VerificationRequest{Session: data.Key, Code: data.Code}

		response, responseErr := authServiceClient.Verification(t.Context(), req)
		require.NoError(t, responseErr)
		assert.NotNil(t, response)
	})

	t.Run("Rate limiter", func(t *testing.T) {
		session := rand.Text()
		code := "1234ABCD"

		req := &authv1.VerificationRequest{Session: session, Code: code}

		t.Run("Allowed", func(t *testing.T) {
			for range 5 {
				response, responseErr := authServiceClient.Verification(t.Context(), req)
				require.Error(t, responseErr)
				assert.Nil(t, response)
				st, _ := status.FromError(responseErr)
				assert.Equal(t, codes.Aborted, st.Code())
			}
		})

		t.Run("Blocked", func(t *testing.T) {
			response, err := authServiceClient.Verification(t.Context(), req)
			require.Error(t, err)
			assert.Nil(t, response)
			st, _ := status.FromError(err)
			assert.Equal(t, codes.ResourceExhausted, st.Code())
		})
	})
}
