package errorsx

import (
	"errors"
	"net/http"
	"testing"
)

func TestHTTPStatusByKind(t *testing.T) {
	cases := []struct {
		err    error
		status int
	}{
		{New(KindValidation, "INVALID", "invalid", nil), http.StatusBadRequest},
		{New(KindNotFound, "NOT_FOUND", "not found", nil), http.StatusNotFound},
		{New(KindConflict, "CONFLICT", "conflict", nil), http.StatusConflict},
		{errors.New("plain"), http.StatusInternalServerError},
	}

	for _, tc := range cases {
		if got := HTTPStatus(tc.err); got != tc.status {
			t.Fatalf("expected %d got %d", tc.status, got)
		}
	}
}
