# syntax = docker/dockerfile:1-experimental

FROM mcr.microsoft.com/vscode/devcontainers/go as base
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* ./
RUN go mod download
COPY . .

FROM base as test
RUN go test

FROM base as build
RUN --mount=type=cache,target=/root/.cache/go-build go build -o /out/trellotogithub .

FROM scratch 
COPY --from=build /out/trellotogithub /
ENTRYPOINT ["/trellotogithub"]