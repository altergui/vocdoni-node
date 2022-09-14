package main

import (
	"fmt"
	"log"
	"os/exec"
)

var env = []string{
	"COMPOSE_DOCKER_CLI_BUILD=1",
	"DOCKER_BUILDKIT=1",
	"COMPOSE_INTERACTIVE_NO_CLI=1",
}

const compose = "docker-compose"

func main() {
	cmd := exec.Command("echo", compose, "up")
	cmd.Env = env
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", out)
}
