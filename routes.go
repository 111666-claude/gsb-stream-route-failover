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
	reports []Route
	Work    int
}

func (r *Router) Update(report Route) { r.reports = append(r.reports, report) }

// Choose prefers healthy routes while allowing a small sticky margin.
func (r *Router) Choose(region string, now int, sticky string) string {
	best := Route{}
	found := false
	stickyFound := false
	var stickyRoute Route
	for _, route := range r.reports {
		r.Work++
		if !route.Healthy || now > route.SeenAt+route.TTL {
			continue
		}
		if route.Name == sticky {
			stickyFound = true
			stickyRoute = route
		}
		if !found || route.Priority > best.Priority {
			best, found = route, true
		}
	}
	if stickyFound {
		return stickyRoute.Name
	}
	if !found {
		return ""
	}
	return best.Name
}
