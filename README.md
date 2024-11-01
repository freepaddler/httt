### Simple http test tool

Basically is used to check reverse proxy and balancers setup. Returns all headers and request information. Returns hostname in `X-Httt-Host` header.

#### Environmental variables
+ `PORT` port to bind, default `8080`
+ `HOST` overrides `os.hostname`
+ `WITH_BODY` returns request body as string, var is set when non-empty

#### Run
+ `docker run -p 8080:8080 freepaddler/httt`
+ `docker run -p 8080:8080 -e WITH_BODY=1 freepaddler/httt`