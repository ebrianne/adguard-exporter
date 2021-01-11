FROM golang:1.15-alpine3.12 as build
ARG TARGETOS
ARG TARGETARCH

WORKDIR /tmp/adguard_exporter

RUN echo "${TARGETOS}-${TARGETARCH}" && apk --no-cache add git alpine-sdk upx
COPY . .
RUN GO111MODULE=on go mod vendor
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags '-s -w' -o adguard_exporter ./
RUN upx -f --brute adguard_exporter

FROM scratch
LABEL name="adguard-exporter"

WORKDIR /root
COPY --from=build /tmp/adguard_exporter/adguard_exporter adguard_exporter

CMD ["./adguard_exporter"]