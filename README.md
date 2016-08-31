#Screen Monitor project

## Installation

1. Install Go 1.7.x or greater, git, setup `$GOPATH`, and `PATH=$PATH:$GOPATH/bin`
2. Install nodejs 5.8 or greater
3. Run `npm install`
4. Run `npm run build`
```
cd $GOPATH/src/github.com/maddevsio/screen-monitoring
go get -v && go build -v
./screen-monitoring
```

## Docker
```
docker build -t screen-monitoring .
docker run -p 8888:8080 -it --rm --name my_screen_monitoring screen-monitoring
```

###Example Dashboard API registration
```
curl -H "Content-Type: application/json" -X POST -d '{"id": "github_http_agent", "width": 200, "height": 122, "content": "<div style=\"border: 3px solid black;\"><p>github.com:443</p><h1>200 ms</h1></div>"}' http://localhost:8080/dashboard/v1/register
```

###React sources:
  * Install node.js version ```>= 5.8```
  * Execute: ```npm install```  
  * For Development:
    * Execute ```npm run build```
    * Created files ```index.html```, ```bundle-[hash].js```, ```bundle-[hash].js.map``` inside ```public```
  * For Production:
    * Execute ```npm run build-production```
    * Created files ```index.html```, ```bundle-[hash].js``` inside ```public```

###Running Go backend service:
  * Execute: ```go build``` then ```./screen-monitoring```
  * Open browser and use url ```http://localhost:8080/```
  * See dashboard page without widgets

###Agents:
  * Every agent should register using dashboard api ```/dashboard/v1/register```
  * Make ```POST``` request to ```/dashboard/v1/register``` with data, for example

  ```
  {
    "id": "github_http_agent",
    "width": 200,
    "height": 122,
    "content": "<div style=\"border: 3px solid black;\"><p>github.com:443</p><h1>200 ms</h1></div>"
  }
  ```

###Dashboard
  * After agents registration need refresh browser page for displaying all registered agents (because now simple solution realized)
