FROM --platform=$BUILDPLATFORM golang:alpine AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
ADD dist/frontend/* /app/dist/frontend/
ADD go.mod /app/
ADD go.sum /app/
ADD main.go /app/
ADD server/ /app/server/
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o webtimer .

FROM scratch
COPY --from=build /app/webtimer /webtimer
ENTRYPOINT ["/webtimer"]
