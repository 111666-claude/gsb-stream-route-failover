package streamroute

// Route is a regional live-stream source report.
type Route struct {
	Region   string
	Name     string
	Revision int
	Priority int
	SeenAt   int
	TTL      int
	Healthy  bool
}

// Router selects a live source for a region.
type Router struct {
	// latest holds only the newest report per (region, name), so
	// selections scan current sources instead of every old report.
	latest map[string]map[string]Route
	Work   int
}

func (r *Router) Update(report Route) {
	if r.latest == nil {
		r.latest = make(map[string]map[string]Route)
	}
	byName := r.latest[report.Region]
	if byName == nil {
		byName = make(map[string]Route)
		r.latest[report.Region] = byName
	}
	// An older or duplicate revision must not undo a newer health state.
	if cur, ok := byName[report.Name]; ok && cur.Revision >= report.Revision {
		return
	}
	byName[report.Name] = report
}

// Choose prefers healthy routes while allowing a small sticky margin.
func (r *Router) Choose(region string, now int, sticky string) string {
	best := Route{}
	found := false
	stickyFound := false
	var stickyRoute Route
	for _, route := range r.latest[region] {
		r.Work++
		// Fresh only while now < seenAt + TTL; the boundary is expired.
		if !route.Healthy || now >= route.SeenAt+route.TTL {
			continue
		}
		if route.Name == sticky {
			stickyFound = true
			stickyRoute = route
		}
		if !found || route.Priority > best.Priority ||
			(route.Priority == best.Priority && route.Name < best.Name) {
			best, found = route, true
		}
	}
	if !found {
		return ""
	}
	// Sticky wins only while healthy/fresh and within two priority
	// points of the best alternative.
	if stickyFound && best.Priority-stickyRoute.Priority <= 2 {
		return stickyRoute.Name
	}
	return best.Name
}
