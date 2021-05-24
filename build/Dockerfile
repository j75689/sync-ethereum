# confluent kafka package need compile with cgo
FROM golang:1.16 as builder

WORKDIR /build
COPY . /build

ENV GO111MODULE=on
ENV GOOS=linux
RUN go get github.com/google/wire/cmd/wire
RUN go generate
RUN go build -mod=readonly -v -o server main.go

FROM golang:1.16

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