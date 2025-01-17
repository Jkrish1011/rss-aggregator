## To Migrate the database

```
goose postgres postgres://postgres:postgres@localhost:5432/rssagg up
```

```
goose postgres postgres://postgres:postgres@localhost:5432/rssagg down
```

## Generate go code by analyzing the quries and schema mentioned in the sqlc.yaml file

```
sqlc generate
```

## To download the dependencies into the local

```
go mod vendor
```

This will create a local copy of the external package dependencies used in the project
