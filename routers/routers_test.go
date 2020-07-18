package routers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestInitRouter(t *testing.T) {
	t.Run("GET /health returns 'ok' with status 200", func(t *testing.T) {
		router := InitRouter()

		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, "ok", w.Body.String())
		assert.Equal(t, 200, w.Code)
	})

	t.Run("GET /feed/:username returns rss with status 200", func(t *testing.T) {
		testUsername := "kenchan0130"
		router := InitRouter()

		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/feed/%s", testUsername), nil)
		router.ServeHTTP(w, req)

		fp := gofeed.NewParser()
		feed, err := fp.Parse(bytes.NewReader(w.Body.Bytes()))

		if err != nil {
			assert.Fail(t, err.Error())
		}

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, fmt.Sprintf("@%s like feed | Twitter Like Feed", testUsername), feed.Title)
	})
}
