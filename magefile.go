//go:build mage
// +build mage

package main

import (
	"errors"
	"os"
	"os/exec"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/metafates/avito-task/log"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Logger.Info().Str("name", name).Strs("args", args).Msg("running command")
	return cmd.Run()
}

func goInstall(url string) error {
	return run("go", "install", url)
}

func inPath(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// Spin up docker containers and run tests
func Test() error {
	err := run("docker", "compose", "up", "-d", "--no-deps", "--build", "server")
	if err != nil {
		return err
	}

	err = run("docker", "compose", "down")
	if err != nil {
		return err
	}

	err = run("docker", "compose", "up", "-d")
	if err != nil {
		return err
	}

	waitDuration := 10 * time.Second
	log.Logger.Info().Dur("wait", waitDuration).Msg("waiting for all containers to start")
	time.Sleep(10 * time.Second)

	err = run("go", "test", "./...")
	if err != nil {
		return err
	}

	return run("docker", "compose", "down")
}

// Rebuild Dockerfile and start docker compose
func Docker() error {
	err := run("docker", "compose", "up", "-d", "--no-deps", "--build", "server")
	if err != nil {
		return err
	}

	err = run("docker", "compose", "down")
	if err != nil {
		return err
	}

	err = run("docker", "compose", "up")
	if err != nil {
		return err
	}

	return nil
}

// Run code generation
func Generate() error {
	mg.Deps(installGenerators)

	log.Logger.Info().Msg("running code generation")

	err := run("go", "generate", "./...")
	if err != nil {
		return err
	}

	err = run("oapi-codegen", "--config", "oapi.cfg.yaml", "openapi.yaml")
	if err != nil {
		return err
	}

	if _, err := os.Stat("Dockerfile"); errors.Is(err, os.ErrNotExist) {
		err = run("goctl", "docker", "-go", "main.go", "--tz", "Europe/Moscow")
		if err != nil {
			return err
		}
	}

	return nil
}

// Manage your deps, or running package managers.
func installGenerators() error {
	log.Logger.Info().Msg("installing dependencies")

	if !inPath("oapi-codegen") {
		err := goInstall("github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest")
		if err != nil {
			return err
		}
	}

	if !inPath("goctl") {
		err := goInstall("github.com/zeromicro/go-zero/tools/goctl@latest")
		if err != nil {
			return err
		}
	}

	return nil
}
