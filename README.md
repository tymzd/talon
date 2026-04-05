# Talon

Talon is a service written in Go that runs in the background, fetching your
workout data from the Hevy API and storing it in a local SQLite database to be
analysed by OpenClaw or displayed on Grafana dashboards, for example.

## Syncing

There are two background loops that sync data from Hevy:

1. Continuous sync from last updated.
2. Daily sync from the beginning of time.

TODO: Eventually, I should support an RPC/HTTP endpoint to force refresh since the last date or something.
TODO: Should I be dockerising talon? For easier/reproducible deployments? I could potentially do a `FROM scratch` docker image.
