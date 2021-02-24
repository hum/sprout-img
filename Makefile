.PHONY: prod
prod:
	@env GOOS=linux CGO_ENABLED=0 GOARCH=arm GOARM=5 \
	  go build -v \
	  -o build/sprout-img cmd/sprout-img/main.go
