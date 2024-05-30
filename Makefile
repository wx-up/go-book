.PHONY: docker
docker:
	@rm gobook || true
	@go mod tidy
	@GOOS=linux GOARCH=amd64 go build -tags=k8s -o gobook .
	@docker rmi -f wxsatellite/go_book:v1.0.0
	@docker build -t wxsatellite/go_book:v1.0.0 .

.PHONY: wrk-login
wrk-login:
	wrk -t2 -c500 -d30s -s ./script/wrk/login.lua http://localhost:8080/users/login


.PHONY: buf
buf:
	buf generate api/proto