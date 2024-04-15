package errors

import (
	"net/http"

	"github.com/astak-homework/connect-now-backend/resources"
)

type HttpResponseError struct {
	resourceKey          resources.LocalizableMessageKey
	internalErrorMessage string
	originalError        error
	status               int
}

func NewBadRequest() HttpResponseError {
	return HttpResponseError{status: http.StatusBadRequest}
}

func NewNotFound() HttpResponseError {
	return HttpResponseError{status: http.StatusNotFound}
}

func NewInternal() HttpResponseError {
	return HttpResponseError{status: http.StatusInternalServerError}
}

func (err HttpResponseError) Err(original error) HttpResponseError {
	err.originalError = original
	return err
}

func (err HttpResponseError) Log(msg string) HttpResponseError {
	err.internalErrorMessage = msg
	return err
}

func (err HttpResponseError) Msg(key resources.LocalizableMessageKey) HttpResponseError {
	err.resourceKey = key
	return err
}

func (err HttpResponseError) InvalidData() HttpResponseError {
	return err.Msg(resources.MsgInvalidData)
}

func (err HttpResponseError) UserNotFound() HttpResponseError {
	return err.Msg(resources.MsgUserNotFound)
}

func (err HttpResponseError) ProfileNotFound() HttpResponseError {
	return err.Msg(resources.MsgProfileNotFound)
}

func (err HttpResponseError) IncorrectPassword() HttpResponseError {
	return err.Msg(resources.MsgIncorrectPassword)
}

func (err HttpResponseError) Internal() HttpResponseError {
	return err.Msg(resources.MsgInternalServerError)
}
