FROM golang:1.13-alpine as builder

# Download certificates
RUN apk --update add ca-certificates

# Create appuser.
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /code
COPY go.mod go.sum /code/
RUN go mod download

COPY . /code
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o /whale-cleaner cmd/cleaner/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /whale-cleaner /whale-cleaner
USER appuser:appuser
ENTRYPOINT ["/whale-cleaner"]