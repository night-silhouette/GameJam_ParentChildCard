#!/bin/bash

if [ -z "$1" ]; then
    echo "使用错误: 请提供版本号 (例如: ./deploy.sh v2)"
    exit 1
fi
VERSION=$1
APP_NAME="pcc_card_${VERSION}"
OUTPUT_BIN="../bin/${VERSION}/pcc_card_${VERSION}"
PM2="/www/server/nodejs/v22.12.0/lib/node_modules/pm2/bin/pm2"

CURRENT_NUM=${VERSION#v}
PREV_NUM=$((CURRENT_NUM - 1))
PREV_VERSION="v${PREV_NUM}"

echo "开始构建项目版本: ${VERSION}..."


echo "正在编译 Golang 代码..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "$OUTPUT_BIN" ./main/main.go
if [ $? -ne 0 ]; then
    echo "编译失败，请检查代码！"
    exit 1
fi

echo "编译成功: $OUTPUT_BIN"

ssh root@120.26.145.68 "${PM2} stop pcc_card_${PREV_VERSION}" || true
ssh root@120.26.145.68 "${PM2} stop pcc_card_${VERSION}" || true
scp -r ../bin/${VERSION} root@120.26.145.68:/root/pcc_servers
ssh root@120.26.145.68 "${PM2} start ./pcc_servers/${VERSION}/pcc_card_${VERSION}"



