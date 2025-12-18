#!/bin/bash
# 更新服务脚本 - 拉取新镜像并重启

set -e

REGISTRY="crpi-7i3xk868tuahj8mk.cn-shenzhen.personal.cr.aliyuncs.com"
NAMESPACE="lv_public"
IMAGE_NAME="flow-link-server"
VERSION="${1:-latest}"

echo "=== 更新 Flow Link Server ==="
echo "版本: ${VERSION}"
echo ""

# 拉取新镜像
echo "📥 拉取新镜像..."
docker pull "${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${VERSION}"

# 重启服务
echo "🔄 重启服务..."
VERSION=${VERSION} docker-compose -f docker-compose.prod.yaml --env-file .env.production up -d --force-recreate web

# 清理旧镜像
echo "🧹 清理旧镜像..."
docker image prune -f

# 检查状态
echo ""
echo "📊 服务状态:"
docker-compose -f docker-compose.prod.yaml ps

echo ""
echo "✅ 更新完成！"
echo ""
echo "查看日志: docker-compose -f docker-compose.prod.yaml logs -f web"

