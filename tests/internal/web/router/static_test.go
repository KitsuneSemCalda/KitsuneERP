package router_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
	"KitsuneSemCalda/KitsuneERP/internal/web/router"
)

func TestRouter_Static(t *testing.T) {
	gest.Describe("Static Files").
		It("should serve static files from a directory", func(g *gest.T) {
			r := router.New()

			// Create a temp directory for static files
			tmpDir, err := os.MkdirTemp("", "static-test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			testFile := "hello.txt"
			testContent := "hello static"
			err = os.WriteFile(filepath.Join(tmpDir, testFile), []byte(testContent), 0644)
			if err != nil {
				t.Fatal(err)
			}

			r.Static("/static", tmpDir)

			req := httptest.NewRequest(http.MethodGet, "/static/"+testFile, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			g.Expect(w.Code).ToBe(http.StatusOK)
			g.Expect(w.Body.String()).ToBe(testContent)
		}).
		Run(t)
}
