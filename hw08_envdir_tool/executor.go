package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	command.Env = prepareEnv(env)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		log.Printf("error while running command: %s", err)
		return 1
	}

	return 0
}

func prepareEnv(env Environment) []string {
	envsCount := len(env)
	result := make([]string, 0, envsCount)
	used := make(map[string]struct{}, envsCount)

	for _, existsEnv := range os.Environ() {
		index := strings.Index(existsEnv, "=")
		if index < 1 {
			continue
		}

		key := existsEnv[:index]
		if value, ok := env[key]; ok {
			used[key] = struct{}{}

			if value.NeedRemove {
				continue
			}
			result = append(result, key+"="+value.Value)
			continue
		}

		result = append(result, existsEnv)
	}

	if len(used) == envsCount {
		return result
	}

	for key, value := range env {
		if _, ok := used[key]; ok {
			continue
		}

		if value.NeedRemove {
			continue
		}

		result = append(result, key+"="+value.Value)
	}

	return result
}
