
## Build the app

just type
```
go build
```
from the projects root to build the app

## Run the app
type 
```
./plangen
```
to run the app. This will start the web server on `localhost`'s port `8080`

## Using the app
As per requirements there is a single endpoint which only accepts `POST` requests.
Example request:
```
curl -i -X POST 'http://localhost:8080/generate-plan' -d '{"duration":2, "loanAmount":"500", "nominalRate":"4.0", "startDate":"2018-01-01T00:00:00Z"}'
```