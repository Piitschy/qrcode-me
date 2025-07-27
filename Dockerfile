FROM golang:1.24-alpine AS build

WORKDIR /app 

COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o ./pocketbase ./main.go

FROM alpine:latest AS production

ENV ADMIN_EMAIL
ENV ADMIN_PASSWORD

WORKDIR /pb

COPY --from=build /app/pocketbase /pb/pocketbase

EXPOSE 8080

# start PocketBase
ENTRYPOINT ["/pb/pocketbase", "serve", "--http=0.0.0.0:8080"]
