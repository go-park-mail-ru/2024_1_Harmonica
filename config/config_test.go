package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvFunctions(t *testing.T) {
	t.Run("Good test get env as string", func(t *testing.T) {
		os.Setenv("abc", "1")
		assert.Equal(t, "1", GetEnv("abc", ""))
	})
	t.Run("Good test get env as int", func(t *testing.T) {
		os.Setenv("abc", "1")
		assert.Equal(t, 1, GetEnvAsInt("abc", 0))
	})
	t.Run("Good test get env as bool", func(t *testing.T) {
		os.Setenv("abc", "True")
		assert.Equal(t, true, GetEnvAsBool("abc", false))
	})
	t.Run("Bad test get env as string", func(t *testing.T) {
		assert.Equal(t, "", GetEnv("NotEnvVar", ""))
	})
	t.Run("Bad test get env as int", func(t *testing.T) {
		assert.Equal(t, 0, GetEnvAsInt("NotEnvVar", 0))
	})
	t.Run("Bad test get env as bool", func(t *testing.T) {
		assert.Equal(t, true, GetEnvAsBool("NotEnvVar", true))
	})
	t.Run("Good new config rest", func(t *testing.T) {
		assert.Equal(t, Config{}, *New())
	})
}
