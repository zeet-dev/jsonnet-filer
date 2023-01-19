package jsonnet

import (
	"bytes"
	"context"
	"errors"

	"github.com/zeet-dev/jsonnet-filer/internal/sh"
)

func EvaluateFile(file string) (string, error) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exitCode, err := sh.Run(context.Background(), "jsonnet", func(o *sh.RunOptions) {
		o.Args = []string{
			file,
		}
		o.Stdout = &out
		o.Stderr = &errOut
	})

	if err != nil {
		return "", err
	}

	if exitCode != 0 {
		return "", errors.New(errOut.String())
	}

	return out.String(), nil
}
