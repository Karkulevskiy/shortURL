package redirect_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/redirect"
	"url-shortener/internal/http-server/handlers/url/redirect/mocks"
	"url-shortener/internal/lib/logger/slogdiscard"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetURL(t *testing.T) {
	cases := []struct {
		name      string
		url       string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			url:   "https://yandex.ru",
			alias: "fkshfs",
		},
		{
			name:      "Bad request",
			alias:     "",
			mockError: errors.New("bad request"),
		},
		{
			name:      "Internal error",
			alias:     "aaaaaa",
			url:       "https://yandex.ru",
			mockError: errors.New("internal error"),
		},
	}

	for _, rt := range cases {
		rt := rt
		t.Run(rt.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := mocks.NewURLGetter(t)

			if rt.mockError == nil && rt.respError == "" {
				urlGetterMock.On("GetURL", mock.AnythingOfType("string")).
					Return(rt.url, rt.mockError).
					Once()
			} else if rt.alias == "" || rt.respError == "" {
				urlGetterMock.On("GetURL", mock.AnythingOfType("string")).
					Return(rt.url, rt.mockError).
					Once()
			}

			handler := redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock, nil)

			req, err := http.NewRequest(http.MethodGet, "/"+rt.alias, nil)

			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			defer rr.Result().Body.Close()

			body, err := io.ReadAll(rr.Result().Body)

			require.NoError(t, err)
			require.NotEmpty(t, body)
			require.Equal(t, rr.Code, http.StatusOK)
		})
	}
}
