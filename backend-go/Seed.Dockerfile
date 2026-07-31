FROM golang:1.26-alpine AS build

WORKDIR /src/backend-go
COPY backend-go/go.mod backend-go/go.sum ./
RUN go mod download
COPY backend-go ./
COPY backend-stub/media /src/backend-stub/media
RUN CGO_ENABLED=0 go build -o /out/seed ./seed

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app
COPY --from=build /out/seed /usr/local/bin/seed
COPY --from=build /src/backend-stub/media /app/media
USER app
ENTRYPOINT ["seed"]
