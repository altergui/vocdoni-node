package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var env = []string{
	"COMPOSE_DOCKER_CLI_BUILD=1",
	"DOCKER_BUILDKIT=1",
	"COMPOSE_INTERACTIVE_NO_CLI=1",
}

func dockerCompose(arg ...string) *exec.Cmd {
	//	arg = append([]string{"docker-compose"}, arg...)
	//	return exec.Command("/bin/sh", "-c", strings.Join(arg, " ")) //nolint:gosec
	return exec.Command("docker-compose", strings.Join(arg, " ")) //nolint:gosec

}

func main() {
	cmd := dockerCompose("version")
	cmd = exec.Command("env")
	cmd.Env = env
	fmt.Println(cmd)
	out, err := cmd.CombinedOutput()

	fmt.Printf("%s", out)
	if err != nil {
		log.Fatal("x", err)
	}

	cmd = exec.Command("docker", "compose", "up", "-d")
	cmd.Env = env
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", out)
}
