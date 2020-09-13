$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o twitchdata .\main.go