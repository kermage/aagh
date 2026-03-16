package helpers

const (
	NAME   = "aagh"
	DIR    = "." + NAME
	KEY    = "core.hooksPath"
	RUNNER = "_"
	SCRIPT = `echo "Hello from $0 script"
echo ""
`
	// File permissions
	PermExecutable = 0755
	PermReadWrite  = 0644
)

var VERSION = "(untracked)"
