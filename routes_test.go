package streamroute

import "testing"

func TestOneRoute(t *testing.T) {
	r := &Router{}
	r.Update(Route{Region: "eu", Name: "edge-a", Revision: 1, Priority: 5, SeenAt: 0, TTL: 10, Healthy: true})
	if got := r.Choose("eu", 5, ""); got != "edge-a" {
		t.Fatalf("got %q", got)
	}
}
