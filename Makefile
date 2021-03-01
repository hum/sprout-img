.PHONY: prod
prod:
	@env GOOS=linux CGO_ENABLED=0 GOARCH=arm GOARM=5 \
	  go build -v \
	  -o build/prod cmd/sprout-img/main.go

.PHONY: dev
dev:
	@env CGO_ENABLED=0 \
	  go build -v \
	  -o build/dev cmd/sprout-img/main.go
