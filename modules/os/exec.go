package os

import (
	"bytes"
	exec_ "os/exec"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

var exec = values.MK_NATIVE_FN(func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if len(args) != 2 {
		return nil, &errors.RuntimeError{Message: "exec() takes exactly 2 arguments"}
	}
	if args[0].Type != shared.String {
		return nil, &errors.RuntimeError{Message: "exec() expects string as first argument"}
	}
	if args[1].Type != shared.Array {
		return nil, &errors.RuntimeError{Message: "exec() expects array as second argument"}
	}

	program := args[0].Value.(string)
	args_ := args[1].Value.([]shared.RuntimeValue)
	rawArgs := make([]string, len(args_))
	for i, arg := range args_ {
		if arg.Type != shared.String {
			return nil, &errors.RuntimeError{Message: "exec() expects array of strings as second argument"}
		}
		rawArgs[i] = arg.Value.(string)
	}

	cmd := exec_.Command(program, rawArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Run()

	stdoutVal := values.MK_STRING(stdout.String())
	stderrVal := values.MK_STRING(stderr.String())
	codeVal := values.MK_NUMBER(float64(cmd.ProcessState.ExitCode()))

	retVal := values.MK_OBJECT(map[string]*shared.RuntimeValue{
		"stdout": &stdoutVal,
		"stderr": &stderrVal,
		"code":   &codeVal,
	})
	return &retVal, nil
})
