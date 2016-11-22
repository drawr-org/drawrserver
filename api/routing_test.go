package api

import (
	"testing"

	"github.com/pressly/chi"
)

// iterate over subroutes of routes and execute fn
func exploreSubroutes(t *testing.T, routes []chi.Route, fn func(chi.Route)) {
	for _, v := range routes {
		if sub := v.SubRoutes; sub != nil {
			fn(v)
			exploreSubroutes(t, sub.Routes(), fn)
		} else {
			fn(v)
		}
	}
}

// TestAPIRouting tests if the api.Routing() function works
func TestAPIRouting(t *testing.T) {
	var (
		expected = []string{"/sessions", "/:sessionID", "/users", "/:userID"}
		found    = []string{}
	)

	// mimick start of main()
	r := chi.NewRouter()
	Routing(r)

	exploreSubroutes(t, r.Routes(), func(v chi.Route) {
		for _, wants := range expected {
			if v.Pattern == wants {
				t.Log(v.Pattern)
				found = append(found, v.Pattern)
			}
		}
	})

	// fail if not all routes where found
	if len(expected) != len(found) {
		t.FailNow()
	}
}
