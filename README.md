# Talon

Talon is a service written in Go that runs in the background, fetching your workout data from the Hevy API and storing it in a local SQLite database to be analysed by OpenClaw or displayed on Grafana dashboards.

tl;dr:
- Minimal footprint. Implemented with pure Go and SQLite.
- Fetches all workouts from the Hevy API and stores them in a local SQLite database. Rate-limited to 1 QPS of outbound traffic.
- Syncs your latest workouts every hour. Syncs your full workout history every 24 hours.

## Usage

```sh
docker run -d \
    --restart unless-stopped \
    -v talon_data:/data \
    -e HEVY_API_KEY=your_api_key_uuid tymzd/talon:latest
```

## Local Development

1. Ensure you have a `.env` file in the project root containing your API key:
   ```env
   HEVY_API_KEY=your_api_key_uuid
   ```
2. Start the daemon:
   ```sh
   go run .
   ```
3. Talon will initialise `talon_test.db` locally, inject the schema, enable WAL-mode, and immediately begin syncing your entire Hevy history.
