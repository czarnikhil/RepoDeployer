# Start from the official golang base image
FROM golang:latest


# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ARG GIT_REPOSITORY_URL

RUN git clone "$GIT_REPOSITORY_URL" /home/app

WORKDIR /home/app

COPY . /home/app/

RUN go mod tidy
RUN go mod vendor


EXPOSE 8080

# Command to run the executable
ENTRYPOINT [ "/home/app/main.go" ]
