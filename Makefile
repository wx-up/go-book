.PHONY: docker
docker:
	@rm gobook || true
	@go mod tidy
	@GOOS=linux GOARCH=amd64 go build -tags=k8s -o gobook .
	@docker rmi -f wxsatellite/go_book:v1.0.0
	@docker build -t wxsatellite/go_book:v1.0.0 .