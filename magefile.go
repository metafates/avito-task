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

func inPath(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Logger.Info().Str("name", name).Strs("args", args).Msg("running command")
	if err := cmd.Run(); err != nil {
		log.Logger.Fatal().Err(err).Str("command", name).Send()
	}
}

func cmd(name string, args ...string) func(...string) {
	return func(extraArgs ...string) {
		run(name, append(args, extraArgs...)...)
	}
}

// Spin up docker containers and run tests
func Test() {
	compose := cmd("docker", "compose")

	compose("down")

	compose("-f", "docker-compose-aux.yml", "up", "-d")

	waitDuration := 5 * time.Second
	log.Logger.Info().Dur("wait", waitDuration).Msg("waiting for all containers to start")
	time.Sleep(10 * time.Second)

	run("go", "test", "./...")
	compose("down")
}

// Start the server
func Run() {
	run("go", "run", ".")
}

type Docker mg.Namespace

// Rebuild Dockerfile and start docker compose
func (Docker) All() {
	compose := cmd("docker", "compose")

	compose("up", "-d", "--no-deps", "--build", "server")
	compose("down")
	compose("up")
}

// Start docker compose only with auxiliary containers (database, web ui) without the server itself
func (Docker) Dev() {
	compose := cmd("docker", "compose")

	compose("down")
	compose("-f", "docker-compose-aux.yml", "up")
}

// Run code generation
func Generate() {
	mg.Deps(installGenerators)

	log.Logger.Info().Msg("running code generation")

	run("go", "generate", "./...")

	run("oapi-codegen", "--config", "oapi.cfg.yaml", "openapi.yaml")

	if _, err := os.Stat("Dockerfile"); errors.Is(err, os.ErrNotExist) {
		run("goctl", "docker", "-go", "main.go", "--tz", "Europe/Moscow")
	}
}

// Manage your deps, or running package managers.
func installGenerators() {
	log.Logger.Info().Msg("installing dependencies")

	install := cmd("go", "install")
	if !inPath("oapi-codegen") {
		install("github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest")
	}

	if !inPath("goctl") {
		install("github.com/zeromicro/go-zero/tools/goctl@latest")
	}
}
