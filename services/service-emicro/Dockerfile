# Multi-stage build for emicro.
# TODO: pin to a specific base image once we settle on a single distroless tag.

FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/emicro ./cmd/emicro

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/emicro /usr/local/bin/emicro
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/emicro"]
