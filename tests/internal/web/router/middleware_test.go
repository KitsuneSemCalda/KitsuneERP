package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
	"KitsuneSemCalda/KitsuneERP/internal/web/router"
)

func TestRouter_Middleware(t *testing.T) {
	gest.Describe("Middleware").
		It("should execute middlewares in the correct order", func(g *gest.T) {
			r := router.New()

			var order []string

			mw1 := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					order = append(order, "mw1-start")
					next.ServeHTTP(w, req)
					order = append(order, "mw1-end")
				})
			}

			mw2 := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					order = append(order, "mw2-start")
					next.ServeHTTP(w, req)
					order = append(order, "mw2-end")
				})
			}

			r.Use(mw1)
			r.Get("/mid", mw2(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				order = append(order, "handler")
				w.Write([]byte("done"))
			})))

			req := httptest.NewRequest(http.MethodGet, "/mid", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			g.Expect(w.Code).ToBe(http.StatusOK)
			g.Expect(w.Body.String()).ToBe("done")

			expectedOrder := []string{"mw1-start", "mw2-start", "handler", "mw2-end", "mw1-end"}
			g.Expect(order).ToEqual(expectedOrder)
		}).
		Run(t)
}

func TestRouter_GroupMiddleware(t *testing.T) {
	gest.Describe("Group Middleware").
		It("should isolate middleware between groups", func(g *gest.T) {
			r := router.New()

			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Header().Add("X-Global", "true")
					next.ServeHTTP(w, req)
				})
			})

			r.Group("/api", func(api *router.Router) {
				api.Use(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.Header().Add("X-Api", "true")
						next.ServeHTTP(w, req)
					})
				})

				api.Get("/test", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte("api test"))
				}))
			})

			r.Get("/root", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Write([]byte("root"))
			}))

			// Test /api/test - should have both headers
			req1 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
			w1 := httptest.NewRecorder()
			r.ServeHTTP(w1, req1)
			g.Expect(w1.Header().Get("X-Global")).ToBe("true")
			g.Expect(w1.Header().Get("X-Api")).ToBe("true")

			// Test /root - should only have X-Global
			req2 := httptest.NewRequest(http.MethodGet, "/root", nil)
			w2 := httptest.NewRecorder()
			r.ServeHTTP(w2, req2)
			g.Expect(w2.Header().Get("X-Global")).ToBe("true")
			g.Expect(w2.Header().Get("X-Api")).ToBe("")
		}).
		Run(t)
}
