FROM golang:1.17-alpine as build
WORKDIR /app

# Restore modules - Start
COPY go.mod ./
COPY go.sum ./
RUN go mod download
# Restore modules - End

COPY *.go ./
COPY utils/ utils/
COPY data/ data/
COPY handlers/ handlers/
RUN go build -o /url-shortener

FROM alpine:3.15
WORKDIR /
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "1000" \
    "nonroot"
USER nonroot:nonroot

COPY --from=build /url-shortener /url-shortener
EXPOSE 5000

ENTRYPOINT ["/url-shortener"]