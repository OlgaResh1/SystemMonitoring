FROM golang:1.22 as build

ENV BIN_FILE /opt/systemstatsdaemon/SystemStatsDaemon
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY . ${CODE_DIR}

RUN go mod download

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/service/*

FROM alpine:3.20

RUN apk add --no-cache iproute2

ENV BIN_FILE "/opt/systemstatsdaemon/SystemStatsDaemon"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/systemstatsdaemon/configd.toml
COPY ./configs/configd.toml ${CONFIG_FILE}

CMD ${BIN_FILE} --config ${CONFIG_FILE}
