package router_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
	"KitsuneSemCalda/KitsuneERP/internal/web/router/middlewares"
)

func TestLoggingMiddleware(t *testing.T) {
	gest.Describe("Logging Middleware").
		It("should log request details correctly", func(g *gest.T) {
			// Capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr) // Reset to default

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("logged"))
			})

			loggingHandler := middlewares.Logging(handler)

			req := httptest.NewRequest(http.MethodPost, "/log-me", nil)
			w := httptest.NewRecorder()

			loggingHandler.ServeHTTP(w, req)

			g.Expect(w.Code).ToBe(http.StatusAccepted)
			g.Expect(w.Body.String()).ToBe("logged")

			logOutput := buf.String()
			g.Expect(logOutput).ToContain("202")
			g.Expect(logOutput).ToContain("POST")
			g.Expect(logOutput).ToContain("/log-me")
			g.Expect(logOutput).ToContain("6B")
		}).
		It("should log default 200 status when not explicitly set", func(g *gest.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok"))
			})

			loggingHandler := middlewares.Logging(handler)

			req := httptest.NewRequest(http.MethodGet, "/default", nil)
			w := httptest.NewRecorder()

			loggingHandler.ServeHTTP(w, req)

			g.Expect(w.Code).ToBe(http.StatusOK)

			logOutput := buf.String()
			g.Expect(logOutput).ToContain("200")
			g.Expect(logOutput).ToContain("GET")
			g.Expect(logOutput).ToContain("/default")
			g.Expect(strings.Contains(logOutput, "2B")).ToBeTrue()
		}).
		Run(t)
}
