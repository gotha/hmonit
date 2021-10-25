# HMonit

Service that monitors the health of other services.

On predefined interval hmonit will ping the service and record its state.

hmonit expects your services to have `/__health` endpoint and to return JSON like this:

```json
{"ok": "true"}
```

depending on their status

## Configure

you need to create `services.json` file, here is an [example](./services.example.json).


## Build

```sh
go build
```

or 

```sh
docker build .
```

## Use 

Open your browser and open [http://localhost:8080/](http://localhost:8080/)
you should see the dashboard

if you send header `Accept: application/json` you are going to get the results in JSON format.
