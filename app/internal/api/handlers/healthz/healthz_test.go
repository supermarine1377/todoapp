package healthz_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/supermarine1377/todoapp/app/internal/api/handlers/healthz"
)

func TestHealthz(t *testing.T) {
	t.Run("TestHealthz", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, healthz.Healthz(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}

		var buff bytes.Buffer
		_, _ = io.Copy(&buff, rec.Body)

		res := struct {
			Status string `json:"status"`
		}{}
		if err := json.Unmarshal(buff.Bytes(), &res); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "healthy", res.Status)
	})
}
