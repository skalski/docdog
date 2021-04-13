cd src/

#build Win
GOOS=windows GOARCH=amd64 go build -o ../bin/win/docdog.exe docdog.go

#build OSX
GOOS=darwin GOARCH=amd64 go build -o ../bin/macos_amd/docdog docdog.go
GOOS=darwin GOARCH=amd64 go build -o ../bin/macos_arm/docdog docdog.go

#build *nix64
go build -o ../bin/linux/docdog docdog.go