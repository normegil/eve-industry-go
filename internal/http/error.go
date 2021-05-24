package http

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type errorHandler struct {
	LogUserError bool
}

func (h errorHandler) handle(w http.ResponseWriter, err error) {
	httpError := &Error{}
	if !errors.As(err, httpError) {
		httpError = &Error{
			Code:   50000,
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}
	if httpError.Status%400 < 100 && h.LogUserError {
		log.Error().Err(err).Msg("user error")
	} else if httpError.Status%400 >= 100 {
		log.Error().Err(err).Msg("rest error")
	}
	resp := httpError.toResponse()
	resp.Error = err.Error()
	bytes, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		log.Error().Err(marshalErr).Interface("Error", resp).Msg("Could not marshal Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	if _, writeErr := w.Write(bytes); nil != writeErr {
		log.Error().Err(writeErr).Interface("Error", resp).Msg("Could not write response with Error")
	}
}

type errorResponse struct {
	Code   int       `json:"code"`
	Status int       `json:"status"`
	Error  string    `json:"error"`
	Time   time.Time `json:"time"`
}

type Error struct {
	Code   int
	Status int
	Err    error
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e Error) toResponse() errorResponse {
	return errorResponse{
		Code:   e.Code,
		Status: e.Status,
		Error:  e.Err.Error(),
		Time:   time.Now(),
	}
}
