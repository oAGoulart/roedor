FROM golang:1.13 AS build
WORKDIR /go/src/roedor
COPY . .
SHELL ["/bin/bash", "-c"]
RUN go get -d -v ./...
RUN go build -o ./dist/roedor -v ./cmd/roedor/

FROM python:3
WORKDIR /app
COPY --from=build /go/src/roedor/dist/ .
COPY --from=build /go/src/roedor/requirements.txt .
RUN pip install --no-cache-dir -r "requirements.txt"
