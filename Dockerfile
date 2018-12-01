# iron/go:dev is the alpine image with the go tools added
FROM iron/go:dev
RUN apk add --update gcc musl-dev
WORKDIR /app
# Set an env var that matches your github repo name, replace treeder/dockergo here with your repo name
ENV SRC_DIR=/go/src/installer
# Add the source code:
ADD . $SRC_DIR
# Build it:
RUN cd $SRC_DIR; GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build  -o installer; cp installer /app/

FROM docker:dind
RUN apk add --no-cache \
		openssl curl
WORKDIR /root/
COPY --from=0 /go/src/installer/build.sh  build.sh
RUN chmod +x build.sh&&./build.sh
COPY --from=0 /app/installer .
RUN chmod +x installer
COPY --from=0 /go/src/installer/template /root/template
CMD ["./installer"]