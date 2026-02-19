package errorsx

import "net/http"

type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func HTTPStatus(err error) int {
	e := ToError(err)
	switch e.Kind {
	case KindValidation:
		return http.StatusBadRequest
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
