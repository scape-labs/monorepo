# Multi-stage build for flow-ussd.
# TODO: pin to a specific base image once we settle on a single distroless tag.

FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/flow-ussd ./cmd/flow-ussd

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/flow-ussd /usr/local/bin/flow-ussd
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/flow-ussd"]
