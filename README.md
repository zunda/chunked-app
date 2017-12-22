# chunked-app
Trying to serve with both of Transfer-Encoding: chunked and Content-Length

## Run locally
```
$ go run main.go
```

and access http://localhost:3000/

## Run on Heroku

```
$ heroku create
$ git push heroku master
$ heroku open
```
