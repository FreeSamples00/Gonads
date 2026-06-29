_default:
    @just --unsorted --list

build:
    go build

clean:
    go clean

[group('CI')]
check:
    go vet

[group('CI')]
test:
    go test

[group('CI')]
format:
    go fmt
