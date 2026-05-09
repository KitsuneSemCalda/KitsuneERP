package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
	"KitsuneSemCalda/KitsuneERP/internal/web/router"
)

func TestRouter_Methods(t *testing.T) {
	r := router.New()

	methods := []struct {
		name    string
		fn      func(string, http.Handler)
		method  string
		pattern string
	}{
		{"Get", r.Get, http.MethodGet, "/get"},
		{"Post", r.Post, http.MethodPost, "/post"},
		{"Put", r.Put, http.MethodPut, "/put"},
		{"Delete", r.Delete, http.MethodDelete, "/delete"},
		{"Patch", r.Patch, http.MethodPatch, "/patch"},
	}

	for _, tt := range methods {
		gest.Describe(tt.name).
			It("should handle "+tt.method+" requests", func(g *gest.T) {
				tt.fn(tt.pattern, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte(tt.name))
				}))

				req := httptest.NewRequest(tt.method, tt.pattern, nil)
				w := httptest.NewRecorder()

				r.ServeHTTP(w, req)

				g.Expect(w.Code).ToBe(http.StatusOK)
				g.Expect(w.Body.String()).ToBe(tt.name)
			}).
			Run(t)
	}
}

func TestRouter_Group(t *testing.T) {
	gest.Describe("Router Groups").
		It("should handle nested groups and correct path prefixing", func(g *gest.T) {
			r := router.New()

			r.Group("/api", func(api *router.Router) {
				api.Get("/v1", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte("api v1"))
				}))

				api.Group("/admin", func(admin *router.Router) {
					admin.Get("/users", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.Write([]byte("admin users"))
					}))
				})
			})

			// Test /api/v1
			req1 := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
			w1 := httptest.NewRecorder()
			r.ServeHTTP(w1, req1)
			g.Expect(w1.Code).ToBe(http.StatusOK)
			g.Expect(w1.Body.String()).ToBe("api v1")

			// Test /api/admin/users
			req2 := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
			w2 := httptest.NewRecorder()
			r.ServeHTTP(w2, req2)
			g.Expect(w2.Code).ToBe(http.StatusOK)
			g.Expect(w2.Body.String()).ToBe("admin users")
		}).
		Run(t)
}
