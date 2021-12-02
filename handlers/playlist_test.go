package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"spoplaylist/Infrastructure/api/playlist"
	"spoplaylist/Infrastructure/api/source_music"
	"spoplaylist/use_cases"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/karlseguin/ccache/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func Test_GetHot100(t *testing.T) {
	cases := []struct {
		name            string
		date            string
		mockPlaylist    playlist.Iplaylist
		mockSourceMusic source_music.ISourceMusic
		cachePlaylist   func() *ccache.Cache
		expCode         int
		expBody         interface{}
	}{
		{
			name: "error",
			date: "2021-01-01",
			mockSourceMusic: source_music.MockSourceMusic{
				MockGetHot100Songs: func(txn *newrelic.Transaction, date string) ([]string, error) {
					return nil, errors.New("error")
				},
			},
			cachePlaylist: func() *ccache.Cache {
				return ccache.New(ccache.Configure())
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name: "ok-cache",
			date: "2021-12-01",
			cachePlaylist: func() *ccache.Cache {
				cache := ccache.New(ccache.Configure())
				cache.Set("2021-12-01", []string{"something", "something"}, time.Minute*5)
				return cache
			},
			expCode: http.StatusOK,
			expBody: []string{"something", "something"},
		},
		{
			name: "ok-api",
			date: "2021-12-01",
			mockSourceMusic: source_music.MockSourceMusic{
				MockGetHot100Songs: func(txn *newrelic.Transaction, date string) ([]string, error) {
					return []string{"something", "something"}, nil
				},
			},
			cachePlaylist: func() *ccache.Cache {
				cache := ccache.New(ccache.Configure())
				return cache
			},
			expCode: http.StatusOK,
			expBody: []string{"something", "something"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, engine := gin.CreateTestContext(httptest.NewRecorder())
			adminUseCase := use_cases.NewAdminPlaylist(c.mockPlaylist,
				c.mockSourceMusic,
				c.cachePlaylist())
			h := Handler{
				AdminPlaylist: adminUseCase,
			}
			engine.GET("/hot-100/:date", func(c *gin.Context) {
				h.GetHot100(c)
			})

			url := fmt.Sprintf("/hot-100/%s", c.date)
			res := httptest.NewRecorder()
			req := httptest.NewRequest("GET", url, nil)
			engine.ServeHTTP(res, req)

			assert.Equal(t, res.Code, c.expCode)
			if c.expBody != nil {
				marshalBody, _ := json.Marshal(c.expBody)
				assert.Equal(t, res.Body.String(), bytes.NewBuffer(marshalBody).String())
			}

		})
	}
}
