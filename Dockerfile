FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
#executes a commands during the build process
RUN apk add build-base && go build -o forum cmd/web/main.go

FROM alpine:3.16
LABEL key="forumlabel"
#sets the working directory for commands run in the container
WORKDIR /app
#copies files from the host machine to the image
COPY --from=builder /app .
#declares the ports  that the container listens on at runtime
EXPOSE 9090
# default command that will be run when a container will be started from the image
CMD ["./forum"]