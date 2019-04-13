FROM daocloud.io/golang:1.8.1 AS builder

COPY ./dep-linux-amd64 /usr/bin/dep

RUN chmod +x /usr/bin/dep

WORKDIR $GOPATH/src/delonakc.com/api

COPY Gopkg.toml Gopkg.lock ./

RUN dep ensure --vendor-only

COPY . ./

COPY ./config.yml /

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.EnvMode=prod" -a -installsuffix nocgo -o /app app.go

FROM scratch

COPY --from=builder /app ./
COPY --from=builder /config.yml ./

EXPOSE 8088

ENTRYPOINT ["./app"]