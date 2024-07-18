package helpers

import (
	"strings"
)

type config struct {
	hooks *hooks
}

func (h *hooks) Config() *config {
	return &config{hooks: h}
}

func (c *config) Get() string {
	out, _ := GitExec(c.hooks.project.FullPath(), "config", "--get", "core.hooksPath")

	return strings.TrimSpace(string(out))
}

func (c *config) Set() error {
	_, err := GitExec(c.hooks.project.FullPath(), "config", "core.hooksPath", DIR)

	return err
}

func (c *config) Correct() bool {
	return c.Get() == DIR
}
