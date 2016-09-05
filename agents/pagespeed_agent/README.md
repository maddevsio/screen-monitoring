## Installation

1. Install Go 1.7.x or greater, git, setup `$GOPATH`, and `PATH=$PATH:$GOPATH/bin`

2. Run the server
    ```
    cd $GOPATH/src/github.com/maddevsio/screen-monitoring/agents/pagespeed_agent
    go get -u -v .
    go build .
    ./pagespeed_agent -httpAddr=:8090 -dashboardURL=http://localhost:8080/dashboard/v1/register -targetURL=https://google.com
    ```
## Env usage
```
export PORT=8090
export DASHBOARD_URL="http://localhost:8080/dashboard/v1/register"
export TARGET_URL="https://google.com"
export API_KEY="{GOOGLE PAGE SPEED APIKEY}"
```

## Flag usage
```
Usage of ./pagespeed_agent:
  -dashboardURL string
       	Dashboard service URL (default "http://localhost:8080/dashboard/v1/register")
  -httpAddr string
       	HTTP listen address (default ":8090")
  -targetURL string
       	Target URL (default "https://google.com")
```
