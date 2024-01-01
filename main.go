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
	args   []string
	envs   []string
	cancel context.CancelFunc
	Cmd    *exec.Cmd
}

func New(name string) *Command {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	cmd := exec.CommandContext(ctx, name)
	cmd.WaitDelay = 1
	return &Command{
		args:   []string{},
		envs:   []string{},
		cancel: cancel,
		Cmd:    cmd,
	}
}

func NewTimeoutCmd(name string, timeout time.Duration) *Command {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	cmd := exec.CommandContext(ctx, name)
	cmd.WaitDelay = 1
	return &Command{
		args:   []string{},
		envs:   []string{},
		cancel: cancel,
		Cmd:    cmd,
	}
}

func (c *Command) run(dir string) ([]byte, error) {
	defer c.cancel()
	c.Cmd.Args = c.args
	c.Cmd.Dir = dir
	c.Cmd.Env = c.envs
	return c.Cmd.CombinedOutput()
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
	return c.Cmd.String()
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
