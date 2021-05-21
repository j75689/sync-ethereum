FROM golang:1.15 as builder

WORKDIR /build
COPY . /build

ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN make tools
RUN go generate
RUN go build -mod=readonly -v -o server main.go

FROM alpine:3

RUN apk add --no-cache ca-certificates
RUN mkdir -p /server
WORKDIR /server

COPY --from=builder /build/server /server/server_bin

# Create appuser.
ENV USER=appuser
ENV UID=1001

RUN adduser \
--disabled-password \
--gecos "application user" \
--no-create-home \
--uid "${UID}" \
"${USER}"

RUN chown appuser:appuser /server
RUN chown appuser:appuser /server/*
USER appuser:appuser

ENTRYPOINT ["/server/server_bin"]