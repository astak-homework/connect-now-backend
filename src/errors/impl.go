package errors

import "errors"

func (err HttpResponseError) Error() string {
	return err.originalError.Error()
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
