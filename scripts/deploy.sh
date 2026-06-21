#!/bin/bash

# ===========================================
# HairCut 一键部署脚本
# ===========================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查运行环境依赖..."
    
    local deps=("docker" "docker-compose" "git")
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            log_error "缺少依赖: $dep"
            exit 1
        fi
    done
    
    log_info "✅ 环境依赖检查通过"
}

# 加载环境变量
load_env() {
    if [ -f ".env" ]; then
        export $(cat .env | grep -v '^#' | xargs)
        log_info "✅ 已加载 .env 环境变量"
    else
        log_warn "未找到 .env 文件，使用默认配置"
    fi
}

# 启动基础设施服务
start_infrastructure() {
    log_info "启动基础设施服务 (MySQL/Redis/ES/MongoDB/RabbitMQ/MinIO)..."
    docker-compose up -d
    sleep 10  # 等待服务启动
    log_info "✅ 基础设施服务已启动"
}

# 初始化数据库
init_database() {
    log_info "初始化数据库表结构..."
    chmod +x scripts/init-db.sh
    ./scripts/init-db.sh
}

# 构建后端服务
build_server() {
    log_info "构建后端 API 服务..."
    cd server
    go mod download
    go build -o ../bin/server ./cmd/server/main.go
    cd ..
    log_info "✅ 后端服务构建完成"
}

# 构建前端应用
build_frontends() {
    log_info "构建前端应用..."
    
    # 运营后台
    log_info "   - 构建运营后台..."
    cd apps/admin-web
    npm ci --registry=https://registry.npmmirror.com
    npm run build
    cd ../..
    
    # 官网
    log_info "   - 构建官网..."
    cd apps/official-website
    npm ci --registry=https://registry.npmmirror.com
    npm run build
    cd ../..
    
    # 用户端小程序
    log_info "   - 构建用户端小程序..."
    cd apps/user-miniapp
    npm ci --registry=https://registry.npmmirror.com
    npm run build:weapp
    cd ../..
    
    log_info "✅ 所有前端应用构建完成"
}

# Docker 部署模式
deploy_docker() {
    log_info "使用 Docker 模式部署生产环境..."
    
    docker-compose -f docker-compose.prod.yml build
    docker-compose -f docker-compose.prod.yml up -d
    
    log_info "✅ Docker 部署完成"
    log_info ""
    log_info "🌐 服务访问地址:"
    log_info "   - API服务: http://localhost:8080"
    log_info "   - 运营后台: http://localhost:80"
    log_info "   - 官网网站: http://localhost:8080"
    log_info "   - Kibana(日志): http://localhost:5601"
    log_info "   - MinIO控制台: http://localhost:9001"
    log_info "   - RabbitMQ管理: http://localhost:15672"
}

# 显示帮助信息
show_help() {
    echo "HairCut 连锁理发店数字化平台 - 部署工具"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  init         初始化项目（安装依赖+数据库初始化）"
    echo "  dev          启动开发环境"
    echo "  build        构建所有应用"
    echo "  deploy       Docker 生产部署"
    echo "  stop         停止所有服务"
    echo "  status       查看服务状态"
    echo "  logs          查看日志"
    echo "  help          显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 init      # 首次部署时初始化"
    echo "  $0 dev       # 启动本地开发环境"
    echo "  $0 deploy    # 生产环境Docker部署"
}

# 主程序
main() {
    case "${1:-help}" in
        init)
            check_dependencies
            load_env
            start_infrastructure
            init_database
            ;;
        dev)
            check_dependencies
            load_env
            start_infrastructure
            log_info "🚀 开发环境就绪！请手动启动各服务:"
            log_info "   - 后端: cd server && go run cmd/server/main.go"
            log_info "   - 后台: cd apps/admin-web && npm run dev"
            log_info "   - 小程序: cd apps/user-miniapp && npm run dev:weapp"
            log_info "   - 官网: cd apps/official-website && npm run dev"
            ;;
        build)
            build_server
            build_frontends
            ;;
        deploy)
            deploy_docker
            ;;
        stop)
            log_info "停止所有服务..."
            docker-compose down
            log_info "✅ 所有服务已停止"
            ;;
        status)
            docker-compose ps
            ;;
        logs)
            docker-compose logs -f --tail=100 ${2:-}
            ;;
        help|*)
            show_help
            ;;
    esac
}

main "$@"
