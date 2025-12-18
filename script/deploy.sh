#!/bin/bash
# 首次部署脚本 - 在服务器上运行

set -e

REGISTRY="crpi-7i3xk868tuahj8mk.cn-shenzhen.personal.cr.aliyuncs.com"
NAMESPACE="lv_public"
IMAGE_NAME="flow-link-server"
VERSION="${1:-latest}"

echo "=== Flow Link Server 部署 ==="
echo "版本: ${VERSION}"
echo ""

# 检查环境配置
if [ ! -f .env.production ]; then
    echo "❌ 请先创建 .env.production 文件"
    echo "   参考: cp .env.production.example .env.production"
    exit 1
fi

# 登录镜像仓库
echo "📦 登录镜像仓库..."
docker login "${REGISTRY}"

# 拉取镜像
echo "📥 拉取镜像..."
docker pull "${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${VERSION}"

# 启动服务
echo "🚀 启动服务..."
VERSION=${VERSION} docker-compose -f docker-compose.prod.yaml --env-file .env.production up -d

# 检查状态
echo ""
echo "📊 服务状态:"
docker-compose -f docker-compose.prod.yaml ps

echo ""
echo "✅ 部署完成！"
echo ""
echo "常用命令:"
echo "  查看日志: docker-compose -f docker-compose.prod.yaml logs -f web"
echo "  查看状态: docker-compose -f docker-compose.prod.yaml ps"
echo "  停止服务: docker-compose -f docker-compose.prod.yaml down"

