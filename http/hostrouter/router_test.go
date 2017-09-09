package hostrouter_test

import (
	"net/http"
	"testing"

	"github.com/prasannavl/go-gluons/http/hostrouter"
	"github.com/prasannavl/mchain"
)

func TestReplaceArrayItem(t *testing.T) {
	router := hostrouter.HostRouter{Threshold: 3}
	router.Set("host1", createHandler())
	router.Set("host2", createHandler())
	router.Set("host2", createHandler())

	items := router.Items.([]hostrouter.HostRouterItem)
	if len(items) != 2 {
		t.Fail()
	}
}

func TestSwitchToBackToArray(t *testing.T) {
	router := hostrouter.HostRouter{Threshold: 3}
	router.Set("host1", createHandler())
	router.Set("host2", createHandler())
	router.Set("host3", createHandler())
	router.Set("host4", createHandler())

	items := router.Items.(map[string]mchain.Handler)
	if len(items) != 4 {
		t.Fail()
	}
	router.Set("host3", nil)
	router.Set("host4", nil)

	items2 := router.Items.([]hostrouter.HostRouterItem)
	if len(items2) != 2 {
		t.Fail()
	}
}

func TestSwitchToMap(t *testing.T) {
	router := hostrouter.HostRouter{Threshold: 3}
	router.Set("host1", createHandler())
	router.Set("host2", createHandler())
	router.Set("host3", createHandler())
	router.Set("host4", createHandler())

	items := router.Items.(map[string]mchain.Handler)
	if len(items) != 4 {
		t.Fail()
	}
}

func createHandler() mchain.Handler {
	return mchain.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})
}
