## Installation

1. Install Go 1.7.x or greater, git, setup `$GOPATH`, and `PATH=$PATH:$GOPATH/bin`

2. Run the server
    ```
    cd $GOPATH/src/github.com/maddevsio/screen-monitoring/agents/ahrefs
    go build
    ./ahrefs -dashboardURL=http://127.0.0.1:8080/dashboard/v1/register -httpAddr=:8090 -ahrefsEmail=email@mail.com -ahrefsPassword=password -ahrefsProject=project_name
    ```

## Env usage
``` 
export PORT=8090
export DASHBOARD_URL="http://localhost:8080/dashboard/v1/register"
export AHREFS_EMAIL = "email@mail.com"
export AHREFS_PASSWORD = "password"
export AHREFS_PROJECT  = "myproject.com"
```

## Flag usage
```
Usage of ./ping_agent:
  -dashboardURL string
       	Dashboard service URL (default "http://localhost:8080/dashboard/v1/register")
  -httpAddr string
       	HTTP listen address (default ":8090")
  -ahrefsEmail string
       	Email address of your ahrefs.com account. (default "email@mail.com")
  -ahrefsPassword string
       	Password. (default "password")
  -ahrefsProject string
       	Name of the project which data metrics you want to get. Be sure to use the exact name which is shown at ahrefs dahsboard. (default "myproject.com")

```
