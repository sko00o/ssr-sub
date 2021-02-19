FROM golang:1.14.7-buster AS builder
LABEL maintainer="Ming Chen"

ENV PACKAGE github.com/mingcheng/ssr-subscriber
ENV GOPROXY https://goproxy.cn,https://goproxy.io,direct
ENV BUILD_DIR ${GOPATH}/src/${PACKAGE}
ENV TARGET_DIR ${BUILD_DIR}

# Print go version
RUN echo "GOROOT is ${GOROOT}"
RUN echo "GOPATH is ${GOPATH}"

# Build
COPY . ${BUILD_DIR}
WORKDIR ${BUILD_DIR}
RUN make clean && \
	make build && \
  	mv ${TARGET_DIR}/ssr-subscriber /usr/bin/ssr-subscriber

# Stage2
#FROM alpine:3.11.6
FROM centos:8

# @from https://mirrors.ustc.edu.cn/help/alpine.html
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN dnf install ca-certificates -y

COPY --from=builder /usr/bin/ssr-subscriber /bin/ssr-subscriber

HEALTHCHECK --interval=60s --timeout=10s \
	CMD curl -fs http://localhost/last-check || exit 1

EXPOSE 80
ENTRYPOINT ["/bin/ssr-subscriber", "--http"]
