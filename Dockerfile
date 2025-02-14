FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o api ./app/services/api

FROM gcr.io/distroless/base AS final

WORKDIR /app

COPY --from=builder /app/api .

COPY --from=builder /app/zarf/keys /app/zarf/keys

EXPOSE 8000

CMD ["/app/api"]

LABEL \
    org.opencontainers.image.title="Template api" \
    org.opencontainers.image.description="Template description" \
    org.opencontainers.image.version="1.0.0" \
    org.opencontainers.image.authors="Your Name <your.email@example.com>" \
    org.opencontainers.image.url="https://template.api.url" \
    org.opencontainers.image.documentation="https://template.api.docs" \
    org.opencontainers.image.source="https://github.com/your-repo" \
    org.opencontainers.image.vendor="Company Name" \
    org.opencontainers.image.licenses="MIT" \
    com.example.api.environment="dev" \