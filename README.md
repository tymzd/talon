# Talon

Talon is a service written in Go that runs in the background, fetching your workout data from the Hevy API and storing it in a local SQLite database to be analysed by OpenClaw or displayed on Grafana dashboards.

Minimal footprint. Implemented with pure Go and SQLite.

## How to Run

1. Ensure you have a `.env` file in the project root containing your API key:
   ```env
   HEVY_API_KEY=your_api_key_uuid
   ```
2. Start the daemon:
   ```sh
   go run .
   ```
3. Talon will initialise `talon_test.db` locally, inject the schema, enable WAL-mode, and immediately begin syncing your entire Hevy history.

## Roadmap & TODOs

- [ ] **Daily Full Sync**: Implement the `FullSync` background loop to scrape from the beginning of time daily, permanently capturing any edge cases or silent upstream deletions not caught by incremental syncs.
- [ ] **REST Endpoints**: Eventually, support an RPC/HTTP endpoint to manually force a refresh without restarting the binary.
- [ ] **Dockerisation**: Containerise exactly into a multi-stage `FROM scratch` Docker image for incredibly lightweight and reproducible deployments to a Raspberry Pi homelab.
