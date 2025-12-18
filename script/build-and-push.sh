#!/bin/bash

#####################################################
# Flow Link Server - 构建和推送脚本
# 功能：构建优化的 Docker 镜像并推送到阿里云镜像仓库
#
# 使用方法：
#   ./build-and-push.sh [选项] [版本号]
#
# 选项：
#   -h, --help      显示帮助信息
#   -b, --build     仅构建，不推送
#   -p, --push      仅推送（需要先构建）
#   -c, --clean     构建后清理本地镜像
#   -m, --multi     多平台构建 (amd64 + arm64)
#   --no-cache      禁用 Docker 构建缓存
#   --no-latest     不推送 latest 标签
#
# 示例：
#   ./build-and-push.sh                    # 交互式构建推送 1.0 版本
#   ./build-and-push.sh 2.0                # 构建推送 2.0 版本
#   ./build-and-push.sh -b 1.5             # 仅构建 1.5 版本
#   ./build-and-push.sh -m 2.0             # 多平台构建 2.0 版本
#   ./build-and-push.sh --no-cache 1.0     # 禁用缓存构建
#
# 环境变量：
#   DOCKER_REGISTRY_PASSWORD  镜像仓库密码（可选，不设置则交互输入）
#####################################################

set -e  # 遇到错误立即退出

# ============================================
# 颜色定义
# ============================================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# ============================================
# 镜像仓库配置
# ============================================
REGISTRY="crpi-7i3xk868tuahj8mk.cn-shenzhen.personal.cr.aliyuncs.com"
NAMESPACE="lv_public"
USERNAME="litevar"
PASSWORD="${DOCKER_REGISTRY_PASSWORD:-}"

# ============================================
# 默认配置
# ============================================
IMAGE_NAME="flow-link-server"
DEFAULT_VERSION="1.0"
VERSION=""
BUILD_ONLY=false
PUSH_ONLY=false
CLEAN_AFTER=false
MULTI_PLATFORM=false
NO_CACHE=false
PUSH_LATEST=true

# ============================================
# 工具函数
# ============================================
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

print_step() {
    echo -e "\n${CYAN}${BOLD}▶ $1${NC}"
}

# 显示帮助信息
show_help() {
    cat << EOF
${BOLD}Flow Link Server - 构建和推送脚本${NC}

${YELLOW}用法:${NC}
    ./build-and-push.sh [选项] [版本号]

${YELLOW}选项:${NC}
    -h, --help      显示此帮助信息
    -b, --build     仅构建镜像，不推送
    -p, --push      仅推送镜像（需要先构建）
    -c, --clean     构建推送后清理本地镜像
    -m, --multi     多平台构建 (linux/amd64 + linux/arm64)
    --no-cache      禁用 Docker 构建缓存
    --no-latest     不推送 latest 标签

${YELLOW}示例:${NC}
    ./build-and-push.sh                    # 交互式构建推送 ${DEFAULT_VERSION} 版本
    ./build-and-push.sh 2.0                # 构建推送 2.0 版本
    ./build-and-push.sh -b 1.5             # 仅构建 1.5 版本
    ./build-and-push.sh -m 2.0             # 多平台构建 2.0 版本
    ./build-and-push.sh -c --no-cache 1.0  # 禁用缓存构建并清理

${YELLOW}环境变量:${NC}
    DOCKER_REGISTRY_PASSWORD    镜像仓库密码（不设置则交互输入）

${YELLOW}镜像仓库:${NC}
    ${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}

EOF
    exit 0
}

# 解析命令行参数
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                ;;
            -b|--build)
                BUILD_ONLY=true
                shift
                ;;
            -p|--push)
                PUSH_ONLY=true
                shift
                ;;
            -c|--clean)
                CLEAN_AFTER=true
                shift
                ;;
            -m|--multi)
                MULTI_PLATFORM=true
                shift
                ;;
            --no-cache)
                NO_CACHE=true
                shift
                ;;
            --no-latest)
                PUSH_LATEST=false
                shift
                ;;
            -*)
                print_error "未知选项: $1"
                echo "使用 --help 查看帮助"
                exit 1
                ;;
            *)
                VERSION="$1"
                shift
                ;;
        esac
    done
    
    # 设置默认版本
    VERSION="${VERSION:-$DEFAULT_VERSION}"
}

# ============================================
# Docker 检查函数
# ============================================
check_docker() {
    print_info "检查 Docker 环境..."
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        print_error "Docker 守护进程未运行，请启动 Docker"
        exit 1
    fi
    
    print_success "Docker 环境正常"
}

# 检查多平台构建支持
check_buildx() {
    if [ "$MULTI_PLATFORM" = true ]; then
        print_info "检查多平台构建支持..."
        
        if ! docker buildx version &> /dev/null; then
            print_error "Docker Buildx 未安装，无法进行多平台构建"
            print_info "请升级 Docker 或安装 buildx 插件"
            exit 1
        fi
        
        # 检查或创建 builder
        if ! docker buildx inspect flow-link-builder &> /dev/null; then
            print_info "创建多平台构建器..."
            docker buildx create --name flow-link-builder --use --bootstrap
        else
            docker buildx use flow-link-builder
        fi
        
        print_success "多平台构建环境就绪"
    fi
}

# 检查是否已登录
check_docker_login() {
    if [ -f ~/.docker/config.json ]; then
        if grep -q "\"${REGISTRY}\"" ~/.docker/config.json 2>/dev/null; then
            return 0
        fi
    fi
    return 1
}

# 登录镜像仓库
docker_login() {
    if check_docker_login; then
        print_info "已登录到镜像仓库"
        return 0
    fi
    
    print_info "登录阿里云镜像仓库..."
    
    if [ -n "${PASSWORD}" ]; then
        echo "${PASSWORD}" | docker login --username="${USERNAME}" --password-stdin "${REGISTRY}"
    else
        print_warning "请输入镜像仓库密码："
        docker login --username="${USERNAME}" "${REGISTRY}"
    fi
    
    if [ $? -ne 0 ]; then
        print_error "登录失败"
        exit 1
    fi
    
    print_success "登录成功"
}

# ============================================
# 构建函数
# ============================================
build_image() {
    local full_image="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${VERSION}"
    local latest_image="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:latest"
    
    print_step "构建 Docker 镜像"
    print_info "版本: ${VERSION}"
    print_info "镜像: ${full_image}"
    
    # 切换到项目目录
    SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
    cd "${SCRIPT_DIR}" || exit 1
    
    # 构建参数
    local build_args=(
        "--file" "../Dockerfile"
        "--build-arg" "VERSION=${VERSION}"
        "--tag" "${full_image}"
    )
    
    # 是否添加 latest 标签
    if [ "$PUSH_LATEST" = true ]; then
        build_args+=("--tag" "${latest_image}")
    fi
    
    # 是否禁用缓存
    if [ "$NO_CACHE" = true ]; then
        build_args+=("--no-cache")
        print_info "已禁用构建缓存"
    fi
    
    # 多平台构建
    if [ "$MULTI_PLATFORM" = true ]; then
        print_info "多平台构建: linux/amd64, linux/arm64"
        build_args+=("--platform" "linux/amd64,linux/arm64")
        
        if [ "$BUILD_ONLY" = false ]; then
            build_args+=("--push")
            print_info "多平台构建将直接推送到仓库"
        else
            build_args+=("--load")
            print_warning "多平台构建仅本地加载时只支持当前架构"
        fi
        
        docker buildx build "${build_args[@]}" ..
    else
        # 单平台构建
        docker build "${build_args[@]}" ..
    fi
    
    if [ $? -ne 0 ]; then
        print_error "镜像构建失败"
        exit 1
    fi
    
    print_success "镜像构建成功"
    
    # 显示镜像大小
    if [ "$MULTI_PLATFORM" = false ]; then
        print_info "镜像大小:"
        docker images "${full_image}" --format "  {{.Repository}}:{{.Tag}} - {{.Size}}"
    fi
}

# 推送镜像
push_image() {
    local full_image="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${VERSION}"
    local latest_image="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:latest"
    
    # 多平台构建已经在构建时推送
    if [ "$MULTI_PLATFORM" = true ]; then
        print_info "多平台镜像已在构建时推送"
        return 0
    fi
    
    print_step "推送镜像到仓库"
    
    # 推送版本标签
    print_info "推送 ${full_image}..."
    docker push "${full_image}"
    if [ $? -ne 0 ]; then
        print_error "镜像推送失败"
        exit 1
    fi
    print_success "版本镜像推送成功"
    
    # 推送 latest 标签
    if [ "$PUSH_LATEST" = true ]; then
        print_info "推送 ${latest_image}..."
        docker push "${latest_image}"
        if [ $? -ne 0 ]; then
            print_warning "latest 标签推送失败，但版本镜像已成功"
        else
            print_success "latest 标签推送成功"
        fi
    fi
}

# 清理本地镜像
cleanup_images() {
    local full_image="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${VERSION}"
    local latest_image="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:latest"
    
    if [ "$CLEAN_AFTER" = true ]; then
        print_step "清理本地镜像"
        docker rmi "${full_image}" 2>/dev/null && print_info "已删除 ${full_image}"
        if [ "$PUSH_LATEST" = true ]; then
            docker rmi "${latest_image}" 2>/dev/null && print_info "已删除 ${latest_image}"
        fi
        print_success "清理完成"
    else
        echo ""
        read -p "是否清理本地构建的镜像？(y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker rmi "${full_image}" 2>/dev/null
            if [ "$PUSH_LATEST" = true ]; then
                docker rmi "${latest_image}" 2>/dev/null
            fi
            print_success "清理完成"
        fi
    fi
}

# 显示总结
show_summary() {
    local full_image="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${VERSION}"
    
    echo ""
    echo -e "${GREEN}════════════════════════════════════════════${NC}"
    echo -e "${GREEN}${BOLD}  ✓ 操作完成！${NC}"
    echo -e "${GREEN}════════════════════════════════════════════${NC}"
    echo ""
    echo -e "${BOLD}镜像信息:${NC}"
    echo "  版本标签: ${full_image}"
    if [ "$PUSH_LATEST" = true ]; then
        echo "  最新标签: ${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:latest"
    fi
    echo ""
    echo -e "${BOLD}拉取命令:${NC}"
    echo "  docker pull ${full_image}"
    echo ""
    echo -e "${BOLD}运行命令:${NC}"
    echo "  # 使用环境变量文件运行"
    echo "  docker run -d \\"
    echo "    --name flow-link-server \\"
    echo "    --env-file .env.production \\"
    echo "    -p 8080:8080 \\"
    echo "    ${full_image}"
    echo ""
    echo -e "${BOLD}或使用 docker-compose:${NC}"
    echo "  # 生产环境"
    echo "  VERSION=${VERSION} docker-compose -f docker-compose.prod.yaml up -d"
    echo ""
}

# ============================================
# 主流程
# ============================================
main() {
    # 解析参数
    parse_args "$@"
    
    local full_image="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${VERSION}"
    
    # 显示标题
    echo ""
    echo -e "${CYAN}════════════════════════════════════════════${NC}"
    echo -e "${CYAN}${BOLD}  Flow Link Server - 镜像构建工具${NC}"
    echo -e "${CYAN}════════════════════════════════════════════${NC}"
    echo ""
    echo -e "  版本:   ${BOLD}${VERSION}${NC}"
    echo -e "  镜像:   ${full_image}"
    echo -e "  模式:   $([ "$BUILD_ONLY" = true ] && echo "仅构建" || ([ "$PUSH_ONLY" = true ] && echo "仅推送" || echo "构建并推送"))"
    [ "$MULTI_PLATFORM" = true ] && echo -e "  平台:   linux/amd64, linux/arm64"
    [ "$NO_CACHE" = true ] && echo -e "  缓存:   ${YELLOW}禁用${NC}"
    echo ""
    
    # 确认操作
    if [ "$BUILD_ONLY" = false ] || [ "$PUSH_ONLY" = true ]; then
        print_warning "此操作将推送镜像到远程仓库"
    fi
    
    read -p "是否继续？(y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "操作已取消"
        exit 0
    fi
    
    # 执行检查
    check_docker
    check_buildx
    
    # 登录（推送时需要）
    if [ "$BUILD_ONLY" = false ]; then
        docker_login
    fi
    
    # 构建
    if [ "$PUSH_ONLY" = false ]; then
        build_image
    fi
    
    # 推送
    if [ "$BUILD_ONLY" = false ]; then
        push_image
    fi
    
    # 清理
    if [ "$PUSH_ONLY" = false ]; then
        cleanup_images
    fi
    
    # 总结
    show_summary
}

# 运行
main "$@"
