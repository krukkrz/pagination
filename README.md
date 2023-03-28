# Pagination

This is an example project to show how pagination may be implemented. It exposes HTTP endpoint `/books` which accepts `limit` and `offset` parameters.

Example:
```bash
curl "localhost:8000/books?limit=10&offset=0"
```

[//]: # (- cursor next/previous)

[//]: # (- auto Incremental PK of the ID)

## Run
In order to start application just run in terminal:
```bash
make start
```
And to stop service:
```bash
make stop
```

## Test
In order to run all tests in the project run:
```bash
make test
```
## Please challenge me!
I love challenges :muscle: 

If you think you would change something in this code, or maybe you know a better approach in pagination (or any other part of code) - don't hesitate to contact me or even start an issue!