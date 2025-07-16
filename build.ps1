$env:CGO_ENABLED=1
go build -tags cgo -buildmode=c-shared -o libclash.dll ./sys