# syntax = docker/dockerfile:1-experimental

FROM golang as build
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -o /out/trellotogithub .

FROM build as test
RUN go test

FROM scratch as bin
COPY --from=build /out/trellotogithub /
ENTRYPOINT ["/trellotogithub"]