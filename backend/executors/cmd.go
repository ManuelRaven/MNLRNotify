package executors

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
)

type CmdExecutor struct{}

func init() {
	RegisterExecutor("cmd", &CmdExecutor{})
}

func (e *CmdExecutor) Execute(mep *MessageExecutorParameters) string {
	message := mep.Message
	script := mep.Script

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Replace {{message}} placeholder if it exists
	script = strings.ReplaceAll(script, "{{message}}", message)

	// Split the script into command and arguments
	parts := strings.Fields(script)
	if len(parts) == 0 {
		return "Error: Empty command"
	}

	// Create command
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

	// Set up output buffers
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute command
	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "Error: Command execution timed out"
		}
		return "Error: " + err.Error() + "\n" + stderr.String()
	}

	result := stdout.String()
	if result == "" {
		return message
	}

	return strings.TrimSpace(result)
}
