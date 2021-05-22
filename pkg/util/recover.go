package util

import (
	"errors"
	"fmt"
	"runtime"
)

func ConvertRecoverToError(r interface{}) error {
	if r != nil {
		var msg string
		for i := 2; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			msg += fmt.Sprintf("%s:%d\n", file, line)
		}

		return errors.New(msg)
	}
	return nil
}
