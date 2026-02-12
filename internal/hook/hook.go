package hook

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Run executes the shell command cmd in dir, injecting PROJJ_HOOK_CONFIG as an
// environment variable containing the JSON-serialised config value.
func Run(cmd, dir string, cfg any) error {
	cfgJSON, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal hook config: %w", err)
	}

	c := exec.Command("sh", "-c", cmd)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = append(os.Environ(), "PROJJ_HOOK_CONFIG="+string(cfgJSON))
	return c.Run()
}
