# 基础镜像
FROM golang:1.20-alpine3.18 as builder

# 设置工作目录
WORKDIR /app

# 将本地文件复制到容器中
COPY . .

ARG VERSION_PATH
ARG BUILD_BRANCH
ARG BUILD_COMMIT
ARG BUILD_TIME
ARG BUILD_GO_VERSION
ARG BUILD_PATH
ARG MAIN_FILE

ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on

RUN echo "VERSION_PATH=${VERSION_PATH}" \
    && echo "BUILD_PATH=${BUILD_PATH}" \
    && echo "MAIN_FILE=${MAIN_FILE}"

COPY ${BUILD_PATH}/config/config.yaml /tmp/config.yaml

#RUN echo "BUILD_BRANCH=${BUILD_BRANCH}" \
#    && echo "BUILD_COMMIT=${BUILD_COMMIT}" \
#    && echo "BUILD_TIME=${BUILD_TIME}" \
#    && echo "BUILD_GO_VERSION=${BUILD_GO_VERSION}" \
#    && echo "BUILD_PATH=${BUILD_PATH}" \
#    && echo "MAIN_FILE=${MAIN_FILE}"
RUN go mod tidy
RUN go build -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GitBranch=${BUILD_BRANCH}' -X '${VERSION_PATH}.GitCommit=${BUILD_COMMIT}' -X '${VERSION_PATH}.BuildTime=${BUILD_TIME}' -X '${VERSION_PATH}.GoVersion=${BUILD_GO_VERSION}'" -o /tmp/main ${BUILD_PATH}/${MAIN_FILE}

FROM alpine
COPY --from=builder /tmp/main .
COPY --from=builder /tmp/config.yaml ./config/config.yaml
ENTRYPOINT ["/main"]
