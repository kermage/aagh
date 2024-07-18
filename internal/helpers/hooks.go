package helpers

import (
	"path/filepath"

	"github.com/kermage/GO-Mods/pathinfo"
)

type hooks struct {
	project pathinfo.PathInfo
}

func Hooks(path string) *hooks {
	pi := pathinfo.Get(path)

	return &hooks{project: pi}
}

func (h *hooks) Project() *pathinfo.PathInfo {
	return &h.project
}

func (h *hooks) Directory() *pathinfo.PathInfo {
	pi := pathinfo.Get(filepath.Join(h.project.FullPath(), DIR))

	return &pi
}

func (h *hooks) ValidRoot() bool {
	return h.Directory().Exists() && h.Config().Correct()
}
