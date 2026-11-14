# Live stream route failover

Each regional source reports name, revision, priority, observation tick, TTL
and health. Reports for the same `(region, name)` are revisioned; old or
duplicate reports must not undo a newer health state. At tick `now`, a report
is fresh only while `now < seenAt + TTL`. Route selection is isolated by
region and considers only fresh healthy sources. Highest priority wins, then
lexicographically smaller name. A caller's sticky source may remain selected
only while it is still healthy and no more than two priority points below the
best alternative; it cannot defeat an expired or unhealthy report.

`Work` counts route candidates examined by selections. A region with a few
current sources must not scan every old report from all regions.
