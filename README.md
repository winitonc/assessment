# Winiton's assessment to Go Software Engineering Bootcamp

##Example command to start application
###Without Docker

```
> DATABASE_URL=postgres://postgres:password@localhost:5432/postgres?sslmode=disable PORT=2565 AUTHORIZATION=November\ 10,\ 2009 go run server.go
```

###With Docker

```
> docker build -t goapp .
> docker run -p 2565:2565 --env-file local.env goapp
```
