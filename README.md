# 🤖 问AI (Wen AI) CLI -- 有问题记不清楚了，问AI。

问AI是一个专为服务器运维和个人主机管理设计的CLI工具，通过集成AI能力，帮助用户快速查找和执行系统命令，提升运维效率。它能够智能解析用户需求，提供精准的命令建议和执行方案，是运维人员的得力助手。

## ✨ 功能特性

- 🤖 智能对话：支持与AI进行自然语言对话，快速查找应用命令
- 🔍 智能感知：智能感知当前工作环境，AI回答更准确
- 🖥️ 跨平台兼容：支持Linux、MacOS、Windows（arm、amd架构）等平台
- 🌍 多语言支持：内置国际化支持，提供多语言界面（目前支持中、英文)
- ⚙️ 配置管理：支持自定义配置，包括API密钥等设置
- 📝 日志记录：详细的日志记录，方便问题排查

## 📦 安装

### 📋 前置要求

- Go 1.22 或更高版本
- Git

### 📝 安装步骤

#### 方式1. 📦 二进制安装
```bash
# 一行安装（中文版）
# 默认安装最新版本
curl https://raw.githubusercontent.com/zdt1013/wen-ai-cli/main/install.sh | bash

# 指定版本和加速源安装
curl https://raw.githubusercontent.com/zdt1013/wen-ai-cli/main/install.sh | bash -s -- -v v0.1.0 -m ghproxy
```
```bash
# 分步骤安装
# 下载安装脚本
curl -o install.sh https://raw.githubusercontent.com/zdt1013/wen-ai-cli/main/install.sh

# 添加执行权限
chmod +x install.sh

# 运行安装脚本（默认安装最新版本）
sudo ./install.sh

# 或者指定版本安装
sudo ./install.sh -v v0.1.0

# 使用加速源安装（可选值：ghproxy, wgetla）
sudo ./install.sh -v v0.1.0 -m ghproxy
```

#### 方式2. 🚀 本地编译
1. 克隆仓库：
```bash
git clone https://github.com/zdt1013/wen-ai-cli.git
cd wen-ai-cli
```

2. 安装依赖：
```bash
go mod download
```

3. 编译项目：
```bash
go build
```

4. 将编译后的可执行文件添加到系统PATH中（可选）

## 🚀 使用方法

### ⌨️ 基本命令

```bash
# 单轮提问模式
> wen + 输入任意问题
```

![问AI使用示例](docs/example1.png)

### 💬 启动对话模式
```bash
> wen chat + 输入首轮问题
或
> wen chat <回车>
  输入任意问题
```
![问AI使用示例](docs/example2.png)

### 🛠️ 其他命令
```bash
# 查看配置
> wen config --help
# 查看帮助信息
> wen --help
```

### 🔧 配置说明

在使用问AI之前，你需要配置必要的参数，比如API密钥等。可以通过以下命令进行配置：

```bash
> wen config --apiKey YOUR_API_KEY --baseUrl YOUR_API_BASE --model YOUR_API_MODEL
或
> wen config -k YOUR_API_KEY -u YOUR_API_BASE -m YOUR_API_MODEL
```

## 📁 项目结构

```
wen-ai-cli/
├── action/     # 命令动作实现
├── assets/     # 静态资源
│   └── lang/   # 多语言包
├── cmd/        # 命令行定义
├── docs/       # 说明文档与附件
├── common/     # 公共组件
├── execute/    # 执行器
├── logger/     # 日志模块
├── model/      # 数据模型
├── setup/      # 初始化设置
├── validate/   # 验证器
├── wenai/      # 核心功能实现
├── main.go     # 程序入口
├── go.mod      # Go模块定义
└── go.sum      # Go模块依赖校验
```

## ⚠️ 已知问题
  * 非命令行相关问题解答，回复样式暂未优化

## 🔮 未来计划 
 * 手册模式：wen man。
 * 感知：运行宿主用户、非sudo用户，智能调整命令。
 * 感知：本机已安装命令、未安装命令。
 * 工具链（functioncall、mcp）支持。
 * 思考模型兼容。
 * 用户体系，方便即安即用，不用自己去配置ai参数。
 * 用户知识库、使用偏好，方便保存自己使用习惯，Wen AI更懂你。

## 📚 依赖开源项目
 * [urfave/cli](https://github.com/urfave/cli) - 命令行应用框架
 * [cloudwego/eino](https://github.com/cloudwego/eino) - 字节跳动开源大型语言模型（LLM）应用开发框架
 * [gookit/config](https://github.com/gookit/config) - 配置管理库
 * [gookit/i18n](https://github.com/gookit/i18n) - 国际化支持
 * [shirou/gopsutil](https://github.com/shirou/gopsutil) - 系统信息采集
 * [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) - OpenAI API 客户端
 * [gookit/slog](https://github.com/gookit/slog) - 日志库
 * [go-cmd/cmd](https://github.com/go-cmd/cmd) - 命令执行库
 * [manifoldco/promptui](https://github.com/manifoldco/promptui) - 交互式命令行界面

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request 来帮助改进这个项目。在提交代码之前，请确保：

1. 代码符合项目的编码规范
2. 添加了必要的测试
3. 更新了相关文档

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 📞 联系方式

如有任何问题或建议，请通过以下方式联系：

- 提交 Issue
