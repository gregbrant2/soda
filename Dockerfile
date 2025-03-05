# syntax=docker/dockerfile:1

# See: https://docs.docker.com/guides/golang/build-images/
FROM golang:1.23 AS base

FROM base AS builder
# Set the working directory
WORKDIR /src

# Copy the mod file and install deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY cmd ./cmd
COPY internal ./internal
RUN file="$(ls .)" && echo $file

# Build
RUN go build -o /soda cmd/soda/main.go

FROM base AS runner

WORKDIR /app

COPY --from=builder /soda ./

# Download and unpack go-migrate
# os=`uname -s`; arch=`uname -m`; 
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar xvz
RUN file="$(ls -l .)" && echo "$file"

COPY web ./web
COPY db ./db
COPY app/bin/run-soda.sh ./
RUN chmod +x ./run-soda.sh

EXPOSE 3030

CMD ["./run-soda.sh"]