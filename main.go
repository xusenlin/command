package command

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const DefaultTimeout = time.Minute

type Command struct {
	name   string
	args   []string
	envs   []string
	cancel context.CancelFunc
	ctx    context.Context
}

func New(name string) *Command {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	return &Command{
		name:   name,
		args:   []string{},
		envs:   []string{},
		cancel: cancel,
		ctx:    ctx,
	}
}

func NewTimeoutCmd(name string, timeout time.Duration) *Command {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &Command{
		name:   name,
		args:   []string{},
		envs:   []string{},
		cancel: cancel,
		ctx:    ctx,
	}
}

func (c *Command) run(dir string) ([]byte, error) {
	defer c.cancel()
	cmd := exec.CommandContext(c.ctx, c.name, c.args...)
	cmd.WaitDelay = 1
	cmd.Dir = dir
	return cmd.CombinedOutput()
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
func (c *Command) String() string {
	return fmt.Sprintf("%s %s %s", c.name, c.args, c.envs)
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
