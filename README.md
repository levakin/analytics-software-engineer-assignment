# Test assignment for Analytics Software Engineer position

This repo contains GitHub event data for 1 hour.

CLI application outputs:

- Top 10 active users sorted by amount of PRs created and commits pushed
- Top 10 repositories sorted by amount of commits pushed
- Top 10 repositories sorted by amount of watch events

Requirements:

- Readable, well-structured code
- Tests
- Structured, meaningful commits
- Some instructions on how to run the solution

## Instructions to run the app

You can use the application as follows:

```shell
go run cmd/ghanalytics/main.go top-users -n 10 -p ./data.tar.gz
go run cmd/ghanalytics/main.go top-repos-by-commits -n 10 -p ./data.tar.gz
go run cmd/ghanalytics/main.go top-repos-by-watch-events -n 10 -p ./data.tar.gz
```

Or you can install the application and use it as binary.

```shell
go install ./cmd/ghanalytics
ghanalytics top-users -n 10 -p ./data.tar.gz
```

To see help instructions:

```shell
ghanalytics
```
