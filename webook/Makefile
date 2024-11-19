.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f neo007/webook:v0.0.2
	@docker build -t neo007/webook:v0.0.1 .