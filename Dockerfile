FROM golang:1.22 as build
WORKDIR /src

COPY ./src /src
RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" -v -o /bin/app ./cmd/api/main.go

FROM scratch
COPY --from=build /bin/app /bin/app
COPY --from=build /src/lang/*.json /bin/lang/

WORKDIR /bin
ENTRYPOINT ["/bin/app"]