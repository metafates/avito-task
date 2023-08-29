package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("Valid config", func(t *testing.T) {
		const (
			port                     = "6969"
			postgresConnectionString = "postgresql://user:pass@localhost:1111/test"
		)

		t.Setenv("SERVER_PORT", port)
		t.Setenv("SERVER_DB_POSTGRES", postgresConnectionString)

		config, err := Load()
		assert.NoError(t, err)

		assert.Equal(t, config.Port, port)
		assert.Equal(t, config.DB.PostgresURI, postgresConnectionString)
	})

	t.Run("Empty config", func(t *testing.T) {
		t.Setenv("SERVER_PORT", "")
		t.Setenv("SERVER_DB_POSTGRES", "")

		_, err := Load()
		assert.Error(t, err)
	})
}
