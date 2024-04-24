package delete_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/lib/logger/slogdiscard"

	"github.com/stretchr/testify/require"
)

func TestDeleteURL(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		mockError error
	}{
		{
			name:  "Success",
			alias: "fkhrff",
		},
		{
			name:      "Invalid request",
			alias:     "",
			mockError: fmt.Errorf("invalid request"),
		},
		{
			name:      "Internal error",
			alias:     "khfske",
			mockError: fmt.Errorf("internal error"),
		},
	}

	for _, td := range cases {
		t.Run(td.name, func(t *testing.T) {
			td := td
			t.Parallel()

			urlDeleteMock := mocks.NewURLDeleter(t)

			urlDeleteMock.On("DeleteURL", td.alias).
				Return(td.mockError).
				Once()

			handler := delete.New(slogdiscard.NewDiscardLogger(), urlDeleteMock)

			req, err := http.NewRequest(http.MethodGet, "/"+td.alias, nil)

			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, td.mockError, rr.Result().Body.Close().Error())
		})
	}
}
