FROM golang:1.15-alpine AS build

RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
 && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.31-r0/glibc-2.31-r0.apk \
 && apk add glibc-2.31-r0.apk \
 && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.31-r0/glibc-bin-2.31-r0.apk \
 && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.31-r0/glibc-i18n-2.31-r0.apk \
 && apk add glibc-bin-2.31-r0.apk glibc-i18n-2.31-r0.apk \
 && /usr/glibc-compat/bin/localedef -i en_US -f UTF-8 en_US.UTF-8 \
 && apk add protobuf-dev make \
 && echo 'done'

WORKDIR /opt/app

COPY go.* ./

RUN go mod download \
 && echo 'done'

COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY main.go ./
COPY Makefile ./

RUN CGO_ENABLED=0 make build test \
 && echo 'done'

FROM alpine AS run

WORKDIR /opt/app

COPY --from=build /opt/app/bin/backend ./

ENTRYPOINT ["./backend"]
