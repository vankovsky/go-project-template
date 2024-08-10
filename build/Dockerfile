FROM golang:1.21.6-alpine3.19 as build

RUN apk add make git

ARG NAME
ARG VERSION

WORKDIR /app
COPY . .
RUN make build

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build ${NAME} /app/${NAME}

# ENTRYPOINT ["/bin/main"]
