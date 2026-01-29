CGO_ENABLED=0 garble -tiny -literals -seed=random -debugdir=debug build -ldflags "-s -w -X main.SecretKey=1ed03b60ed1e94f7" -o release/
upx --ultra-brute ./train-schedule