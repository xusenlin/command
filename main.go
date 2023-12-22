package command

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"time"
)

const DefaultTimeout = time.Minute

type Command struct {
	name    string
	args    []string
	envs    []string
	timeout time.Duration
	ctx     context.Context
}

func New(name string) *Command {
	return &Command{
		name:    name,
		args:    []string{},
		envs:    []string{},
		timeout: DefaultTimeout,
	}
}

func (c *Command) run(dir string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.name, c.args...)
	cmd.Dir = dir
	cmd.Env = c.envs
	return cmd.CombinedOutput()
}

func (c *Command) SetTimeout(timeout time.Duration) *Command {
	c.timeout = timeout
	return c
}

func (c *Command) AddArgs(args ...string) *Command {
	c.args = append(c.args, args...)
	return c
}

func (c *Command) AddEnvs(envs ...string) *Command {
	c.envs = append(c.envs, envs...)
	return c
}

func (c *Command) Run() (string, error) {
	o, err := c.run("")
	return string(o), err
}
func (c *Command) RunInDir(dir string) (string, error) {

	if !IsDir(dir) {
		return "", errors.New("the running directory does not exist")
	}
	o, err := c.run(dir)
	return string(o), err
}

func IsDir(name string) bool {
	if info, err := os.Stat(name); err == nil {
		return info.IsDir()
	}
	return false
}
