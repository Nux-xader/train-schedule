set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
garble -tiny -literals build -ldflags "-s -w" -o release/
