# Development with Docker

Start Docker Desktop with Linux containers, then run:

```sh
docker compose up --build -d
docker compose logs -f api
```

The API is available at http://localhost:8080 once compilation finishes.
The development image includes GCC, enables CGO for SQLite, and runs the API
under Delve with compiler optimizations and inlining disabled.

## Attach a debugger

In GoLand, create a **Go Remote** run configuration with host `localhost`
and port `2345`. Set a breakpoint and start that debug configuration.

In VS Code with the Go extension, add this configuration to `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Attach to Docker",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "host": "127.0.0.1",
      "port": 2345,
      "substitutePath": [
        { "from": "${workspaceFolder}", "to": "/src" }
      ]
    }
  ]
}
```

Source files are mounted into `/src`. After editing Go code, rebuild and
restart the debug session with:

```sh
docker compose restart api
```

Reconnect the debugger after the restart. This setup rebuilds on container
start; it does not automatically reload on file changes. Changes to the
Dockerfile require `docker compose up --build -d`.

SQLite uses `/data/api.db` in the `sqlite-data` volume. An existing local
`api.db` is not imported automatically. `docker compose down` preserves
the volume; `docker compose down -v` deletes it.

The debugger port is bound to localhost. The development container runs
as root with ptrace enabled for debugging. The final Dockerfile stage
still builds a small, non-root release image with `docker build -t restapi .`.
