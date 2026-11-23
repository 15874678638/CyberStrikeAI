package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestFrontendRoutes(t *testing.T) {
	t.Chdir("../..")
	router := gin.New()
	registerFrontendRoutes(router)

	for _, tc := range []struct {
		path, location string
		status         int
	}{
		{"/", "", http.StatusOK},
		{"/classic", "/", http.StatusFound},
		{"/api-docs", "/#/api-docs", http.StatusFound},
		{"/static/favicon.ico", "", http.StatusOK},
		{"/static/js/chat.js", "", http.StatusNotFound},
	} {
		t.Run(tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if w.Code != tc.status || w.Header().Get("Location") != tc.location {
				t.Fatalf("status %d, location %q", w.Code, w.Header().Get("Location"))
			}
			if tc.path == "/" && (!strings.Contains(w.Body.String(), `<div id="root">`) || w.Header().Get("Cache-Control") != "no-cache") {
				t.Fatal("root must serve the revalidated Semi entry point")
			}
		})
	}
}
