 # syntax=docker/dockerfile:1
FROM golang:1.22.2-alpine AS base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . /src


FROM base AS build-app
COPY --from=base src/tests ./tests
RUN CGO_ENABLED=0 go build -o /bin/yadro_test ./cmd/main.go




FROM scratch AS app
COPY --from=build-app /bin/yadro_test /bin/app
COPY --from=build-app src/tests ./tests
ENTRYPOINT [ "./bin/app" ]

