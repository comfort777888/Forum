run:
	go run ./cmd/web/main.go

build:
	go build -o forum cmd/web/main.go

dbuild: 
	docker image build -t f-image .

drun:
	docker container run -p 9090:9090 -d --name f-container f-image

dstop:
	docker stop f-container

drm: 
	docker rm f-container

drim:
	docker rmi f-image

dclear:
	docker system prune -a
