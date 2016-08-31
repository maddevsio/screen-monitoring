## Installation

1. Install Go 1.7.x or greater, git, setup `$GOPATH`, and `PATH=$PATH:$GOPATH/bin`

2. Run the server
    ```
    cd $GOPATH/src/github.com/maddevsio/screen-monitoring/agents/ping_agent
    go build
    ./ping_agent -httpAddr=7066 -dashboardURL=http://localhost:8080/dashboard/v1/register -targetHost=github.com
    ```

## Env usage
```
export PORT=8090
export DASHBOARD_URL="http://localhost:8080/dashboard/v1/register"
export TARGET_HOST="github.com"
```

## Flag usage
```
Usage of ./ping_agent:
  -dashboardURL string
       	Dashboard service URL (default "http://localhost:8080/dashboard/v1/register")
  -httpAddr string
       	HTTP listen address (default ":8090")
  -targetHost string
       	Target hostname and port (default "google.com")
```
