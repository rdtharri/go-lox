# Install stringer
go install golang.org/x/tools/cmd/stringer@latest

# Generate
go generate ./...

# Build
go build .
