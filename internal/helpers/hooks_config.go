package helpers

import (
	"path/filepath"
	"strings"
)

type config struct {
	runner string
	hooks  *hooks
}

func (h *hooks) Config() *config {
	return &config{hooks: h, runner: filepath.Join(DIR, RUNNER)}
}

func (c *config) Get() string {
	out, _ := GitExec(c.hooks.project.FullPath(), "config", "--get", KEY)

	return strings.TrimSpace(string(out))
}

func (c *config) Set() error {
	_, err := GitExec(c.hooks.project.FullPath(), "config", KEY, c.runner)

	return err
}

func (c *config) Correct() bool {
	return c.Get() == c.runner
}
