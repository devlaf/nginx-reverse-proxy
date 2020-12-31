FROM golang:alpine as setup-tool-builder
WORKDIR /app
COPY nginx-setup-tool .
RUN go get -d -v
RUN go generate
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"


FROM nginx:alpine

RUN apk add --update \
        bash \
        certbot \
        certbot-nginx \
        openssl

ADD nginx-conf /etc/nginx

WORKDIR /app
COPY --from=setup-tool-builder /app/nginx-setup-tool .
COPY nginx-setup-tool/config.json .

