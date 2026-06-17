package config_test

import (
	"testing"

	"github.com/shv-ng/fastx/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestValidMigrations(t *testing.T) {

	assert.Equal(t, []string{"None"}, config.ValidMigrations("None"))
	assert.Equal(t, []string{"Alembic", "None"}, config.ValidMigrations("SQLModel"))
}
