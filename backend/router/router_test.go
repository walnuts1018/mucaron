package router

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/walnuts1018/mucaron/backend/config"
	"github.com/walnuts1018/mucaron/backend/router/handler"
)

func TestHealth(t *testing.T) {
	router, err := NewRouter(config.Config{LogLevel: slog.LevelDebug}, handler.Handler{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"status":"ok"}`, w.Body.String())
}
