#!/bin/bash

# 显示帮助信息
show_help() {
    echo "使用方法: $0 [选项]"
    echo "选项:"
    echo "  -v, --version VERSION  指定要安装的版本号"
    echo "  -m, --mirror MIRROR    指定加速源 (可选值: ghproxy, wgetla, default)"
    echo "  -h, --help             显示帮助信息"
    exit 0
}

# 检查用户权限并执行命令
check_and_execute() {
    local cmd="$1"
    if [ -w "/usr/bin" ]; then
        # 用户有写入权限，直接执行
        eval "$cmd"
    else
        # 用户没有写入权限，使用sudo
        echo "需要sudo权限来安装到 /usr/bin 目录"
        sudo $cmd
    fi
}

# 获取系统信息
case "$(uname -s)" in
    Linux*)
        OS="linux"
        ;;
    Darwin*)
        OS="darwin"
        ;;
    *)
        echo "不支持的操作系统: $(uname -s)"
        exit 1
        ;;
esac

# 获取系统架构
case "$(uname -m)" in
    x86_64|amd64)
        ARCH="x86_64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "不支持的架构: $(uname -m)"
        exit 1
        ;;
esac

# 根据操作系统设置下载文件扩展名和安装路径
case $OS in
    linux|darwin)
        EXT=""
        INSTALL_PATH="/usr/bin/wen"
        ;;
    *)
        echo "不支持的操作系统: $OS"
        exit 1
        ;;
esac

# 如果没有指定版本号，则获取最新版本
if [ -z "$VERSION" ]; then
    echo "正在获取最新版本号..."
    VERSION=$(curl -fSsL https://api.github.com/repos/zdt1013/wen-ai-cli/releases/latest | grep -o '"tag_name": "[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$VERSION" ]; then
        echo "无法获取最新版本号"
        exit 1
    fi
    echo "将安装最新版本: $VERSION"
else
    echo "将安装指定版本: $VERSION"
fi

# 定义加速域名列表
declare -A MIRRORS=(
    ["ghproxy"]="https://ghproxy.net/"
    ["ghproxy2"]="https://ghproxy.imciel.com/"
    ["wgetla"]="https://wget.la/"
    ["default"]=""
)

# 解析命令行参数
MIRROR="default"
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -m|--mirror)
            if [[ -n "${MIRRORS[$2]}" ]]; then
                MIRROR="$2"
            else
                echo "无效的加速源: $2"
                echo "可用的加速源: ${!MIRRORS[*]}"
                exit 1
            fi
            shift 2
            ;;
        -h|--help)
            show_help
            ;;
        *)
            echo "未知选项: $1"
            show_help
            ;;
    esac
done

# 处理版本号，分别用于目录和文件名
RELEASE_VERSION=$VERSION
VERSION=${VERSION#v}

# 构建下载URL
GITHUB_URL="https://github.com/zdt1013/wen-ai-cli/releases/download/${RELEASE_VERSION}/wen-ai-cli_${VERSION}_${OS}_${ARCH}${EXT}"

# 确保目标目录存在
check_and_execute "mkdir -p $(dirname $INSTALL_PATH)"

# 根据选择的加速源下载
if [ "$MIRROR" == "default" ]; then
    echo "使用默认下载方式..."
    check_and_execute "curl -fSL -o $INSTALL_PATH $GITHUB_URL"
    
    if [ $? -ne 0 ]; then
        echo "默认下载失败，尝试使用加速源..."
        for mirror_name in "${!MIRRORS[@]}"; do
            if [ "$mirror_name" != "default" ]; then
                mirror_url="${MIRRORS[$mirror_name]}${GITHUB_URL}"
                echo "尝试使用加速源: $mirror_name ($mirror_url)"
                check_and_execute "curl -L -o $INSTALL_PATH $mirror_url"
                if [ $? -eq 0 ]; then
                    echo "使用加速源 $mirror_name 下载成功"
                    break
                fi
            fi
        done
    fi
else
    mirror_url="${MIRRORS[$MIRROR]}${GITHUB_URL}"
    echo "使用加速源 $MIRROR 下载: $mirror_url"
    check_and_execute "curl -fSL -o $INSTALL_PATH $mirror_url"
fi

if [ $? -ne 0 ]; then
    echo "所有下载方式均失败"
    exit 1
fi

# 设置执行权限
check_and_execute "chmod +x $INSTALL_PATH"

echo "安装完成！"
echo "wen-ai-cli 已安装到 $INSTALL_PATH"
