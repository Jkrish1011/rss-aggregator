## To Migrate the database

```
goose postgres postgres://postgres:postgres@localhost:5432/rssagg up
```

## Generate go code by analyzing the quries and schema mentioned in the sqlc.yaml file

```
sqlc generate
```
