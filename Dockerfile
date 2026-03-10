FROM node:22
WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

FROM golang:1.25.3
ARG BUILD_VERSION
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=0 /app/frontend/dist /app/frontend/dist

ENV CGO_ENABLED=0
RUN go build -o /go/bin/app ./cmd/server
RUN go build -o /go/bin/migrate ./cmd/migrate

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /
COPY --from=1 /go/bin/app /usr/local/bin/app
COPY --from=1 /go/bin/migrate /usr/local/bin/migrate
COPY --from=1 /app/db/migrations /db/migrations
CMD ["app"]
