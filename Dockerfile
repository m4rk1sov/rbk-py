FROM golang:1.25 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /docgen ./cmd/server

FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=build /docgen /docgen
COPY templates ./templates
ENV ADDR=:8000 \
    STATIC_TOKEN=default_token \
    JWT_SECRET=dev_secret \
    TEMPLATE_DIR=/app/templates \
    SERVICE_CONTEXT_URL=/document-generator \
    PDF_CONVERTER_URL=http://gotenberg:3000
EXPOSE 8000
ENTRYPOINT ["/docgen"]