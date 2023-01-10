package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty cmd", func(t *testing.T) {
		returnCode := RunCmd([]string{}, Environment{})

		assert.Equal(t, 1, returnCode)
	})

	t.Run("command without arguments", func(t *testing.T) {
		returnCode := RunCmd([]string{"pwd"}, Environment{})

		assert.Equal(t, 0, returnCode)
	})

	t.Run("command with arguments", func(t *testing.T) {
		returnCode := RunCmd([]string{"ls", "-la"}, Environment{})

		assert.Equal(t, 0, returnCode)
	})

	t.Run("command with error", func(t *testing.T) {
		returnCode := RunCmd([]string{"notexists"}, Environment{})

		assert.Equal(t, 1, returnCode)
	})

	t.Run("command with wrong arguments", func(t *testing.T) {
		returnCode := RunCmd([]string{"pwd", "-s"}, Environment{})

		assert.Equal(t, 1, returnCode)
	})
}
