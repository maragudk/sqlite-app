package http_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	apphttp "github.com/maragudk/litefs-app/http"
)

func TestRedirectRegion(t *testing.T) {
	mux := chi.NewMux()
	mux.
		With(apphttp.RedirectRegion("foo")).
		Get("/", func(w http.ResponseWriter, r *http.Request) {})

	t.Run("sets HTTP header if region query param set and different to passed region", func(t *testing.T) {
		code, headers, _ := makeGetRequest(t, mux, "/?region=bar")
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "region=bar", headers.Get("fly-replay"))
	})

	t.Run("does not set HTTP header if region query param not set", func(t *testing.T) {
		code, headers, _ := makeGetRequest(t, mux, "/")
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "", headers.Get("fly-replay"))
	})

	t.Run("does not set HTTP header if region query param set to same as passed region", func(t *testing.T) {
		code, headers, _ := makeGetRequest(t, mux, "/?region=foo")
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "", headers.Get("fly-replay"))
	})
}

type mockPrimaryGetter struct {
	primary string
}

func (m *mockPrimaryGetter) GetPrimary() (string, error) {
	return m.primary, nil
}

func TestRedirectToPrimary(t *testing.T) {
	mux := chi.NewMux()
	mux.Use(apphttp.RedirectToPrimary(&mockPrimaryGetter{"foo"}, nil))
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {})

	t.Run("redirects post requests", func(t *testing.T) {
		code, headers, _ := makePostRequest(t, mux, "/", nil)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "instance=foo", headers.Get("fly-replay"))
	})

	t.Run("ignores get requests", func(t *testing.T) {
		code, headers, _ := makeGetRequest(t, mux, "/")
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "", headers.Get("fly-replay"))
	})

	t.Run("ignores requests with region set in URL query params", func(t *testing.T) {
		code, headers, _ := makePostRequest(t, mux, "/?region=bar", nil)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "", headers.Get("fly-replay"))
	})

	t.Run("does not redirect if primary is empty", func(t *testing.T) {
		mux := chi.NewMux()
		mux.
			With(apphttp.RedirectToPrimary(&mockPrimaryGetter{""}, nil)).
			Post("/", func(w http.ResponseWriter, r *http.Request) {})

		code, headers, _ := makePostRequest(t, mux, "/", nil)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "", headers.Get("fly-replay"))
	})
}

func makeGetRequest(t *testing.T, h http.Handler, target string) (int, http.Header, string) {
	return makeRequest(t, h, http.MethodGet, target, nil)
}

// makePostRequest and return the status code, response header, and the body.
func makePostRequest(t *testing.T, h http.Handler, target string, body io.Reader) (int, http.Header, string) {
	return makeRequest(t, h, http.MethodPost, target, body)
}

func makeRequest(t *testing.T, h http.Handler, method, target string, body io.Reader) (int, http.Header, string) {
	req := httptest.NewRequest(method, target, body)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	result := res.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	return result.StatusCode, result.Header, strings.TrimSpace(string(bodyBytes))
}
