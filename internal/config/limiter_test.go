//go:build integration

package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/api/internal/config"
)

func TestNewLimiter(t *testing.T) {
	t.Parallel()
	client, clientErr := config.NewValkey(t.Context(), valkeyURL)
	require.NoError(t, clientErr)
	assert.NotNil(t, client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		limiter, limiterErr := config.NewLimiter(client)
		require.NoError(t, limiterErr)
		assert.NotNil(t, limiter)

		t.Run("Check", func(t *testing.T) {
			t.Parallel()
			for i := range 6 {
				result, resultErr := limiter.Login.Allow(t.Context(), "test@test.com")
				require.NoError(t, resultErr)
				if i < 5 {
					assert.True(t, result.Allowed)
				} else {
					assert.False(t, result.Allowed)
				}
			}
		})
	})
}
