FROM golang:1.15.2-alpine as builder
# Force the go compiler to use modules
ENV GO111MODULE=on
# Create the user and group files to run unprivileged
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
RUN apk update && apk add --no-cache git ca-certificates tzdata
RUN mkdir /build
COPY . /build/
WORKDIR /build
# Import the code
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o twitchdata .
FROM scratch AS final
LABEL author="leNicDev"
# Import the time zone files
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# Import the user and group files
COPY --from=builder /user/group /user/passwd /etc/
# Import the CA certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the compiled go executable
COPY --from=builder /build/twitchdata /
WORKDIR /
# Run as unpriveleged
USER nobody:nobody
ENTRYPOINT ["/twitchdata"]