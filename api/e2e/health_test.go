//go:build e2e

package e2e

import (
	"net/http"
	"testing"

	"github.com/geek-teru/go-react-ec/api/e2e/helper"
)

// DB を共有して各テストで TRUNCATE するので、t.Parallel() は付けない

func TestHealth(t *testing.T) {
	e := helper.Setup(t, sqlDB, baseURL)

	e.GET("/api/health").
		Expect().
		Status(http.StatusOK).
		JSON().Object().HasValue("status", "ok")
}

func TestHealthDB(t *testing.T) {
	e := helper.Setup(t, sqlDB, baseURL)

	e.GET("/api/health/db").
		Expect().
		Status(http.StatusOK).
		JSON().Object().HasValue("status", "ok")
}
