FROM golang:1.17-buster as build
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

FROM gcr.io/distroless/base-debian11
WORKDIR /

COPY --from=build /url-shortener /url-shortener
EXPOSE 5000
USER nonroot:nonroot
ENTRYPOINT ["/url-shortener"]