# Multi-stage build for dura.
# TODO: pin to a specific base image once we settle on a single distroless tag.

FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/dura ./cmd/dura

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/dura /usr/local/bin/dura
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/dura"]
