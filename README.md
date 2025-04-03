# 问AI (Wen AI) CLI

问AI是一个专为服务器运维和个人主机管理设计的CLI工具，通过集成AI能力，帮助用户快速查找和执行系统命令，提升运维效率。它能够智能解析用户需求，提供精准的命令建议和执行方案，是运维人员的得力助手。

## 功能特性

- 🤖 智能对话：支持与AI进行自然语言对话，快速查找应用命令
- 🌍 多语言支持：内置国际化支持，提供多语言界面（目前支持中、英文)
- ⚙️ 配置管理：支持自定义配置，包括API密钥等设置
- 📝 日志记录：详细的日志记录，方便问题排查

## 安装

### 前置要求

- Go 1.22 或更高版本
- Git

### 安装步骤

1. 克隆仓库：
```bash
git clone https://github.com/yourusername/wen-ai-cli.git
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

## 使用方法

### 基本命令

```bash
# 单轮提问模式
> wen + 输入任意问题

# 启动对话模式
> wen chat + 输入首轮问题
或
> wen chat 回车
> 输入任意问题

# 配置设置
wen config

# 查看帮助信息
wen --help
```


### 配置说明

在使用问AI之前，你需要配置必要的参数，比如API密钥等。可以通过以下命令进行配置：

```bash
wen config set api-key YOUR_API_KEY
```

## 项目结构

```
wen-ai-cli/
├── action/     # 命令动作实现
├── cmd/        # 命令行定义
├── common/     # 公共组件
├── conf/       # 配置文件
├── execute/    # 执行器
├── logger/     # 日志模块
├── model/      # 数据模型
├── setup/      # 初始化设置
├── validate/   # 验证器
└── wenai/      # 核心功能实现
```

## 贡献指南

欢迎提交 Issue 和 Pull Request 来帮助改进这个项目。在提交代码之前，请确保：

1. 代码符合项目的编码规范
2. 添加了必要的测试
3. 更新了相关文档

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 联系方式

如有任何问题或建议，请通过以下方式联系：

- 提交 Issue
