## Running tests

```bash
go test -v $(go list ./... | grep -v /vendor/)
```

## Saving new dependencies

Uses [Godep](https://github.com/tools/godep).

```bash
go get -u <package>
godep save -t $(go list ./... | grep -v /vendor/)
```
