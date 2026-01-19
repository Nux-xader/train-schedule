CGO_ENABLED=0 garble -tiny -literals -seed=random -debugdir=debug build -ldflags "-s -w -X main.SecretKey=$1" -o release/
