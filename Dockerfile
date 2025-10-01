# Build stage
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=0.0.0
ARG GIT_SHA=dev
ARG BUILD_TIME=unknown
ARG CHANNEL=dev
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o /out/hello .

# Runtime
FROM gcr.io/distroless/base-debian12
ENV VERSION=${VERSION}
ENV GIT_SHA=${GIT_SHA}
ENV BUILD_TIME=${BUILD_TIME}
ENV CHANNEL=${CHANNEL}
EXPOSE 8080
COPY --from=build /out/hello /hello
USER 65532:65532
ENTRYPOINT ["/hello"]
