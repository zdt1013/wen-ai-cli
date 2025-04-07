# 🤖 Wen AI CLI -- Ask AI when you can't remember something.

Wen AI is a CLI tool specifically designed for server operations and personal host management. By integrating AI capabilities, it helps users quickly find and execute system commands, improving operational efficiency. It can intelligently parse user requirements and provide accurate command suggestions and execution solutions, making it a powerful assistant for operations personnel.

## ✨ Features

- 🤖 Smart Dialogue: Supports natural language conversations with AI to quickly find application commands
- 🔍 Smart Perception: Intelligently perceives the current working environment for more accurate AI responses
- 🖥️ Cross-Platform Compatibility: Supports Linux, MacOS, Windows (arm, amd architectures) and other platforms
- 🌍 Multi-language Support: Built-in internationalization support, providing multi-language interfaces (currently supports Chinese and English)
- ⚙️ Configuration Management: Supports custom configurations, including API keys and other settings
- 📝 Logging: Detailed logging for easy troubleshooting

## 📦 Installation

### 📋 Prerequisites

- Go 1.22 or higher
- Git

### 📝 Installation Steps

1. Clone the repository:
```bash
git clone https://github.com/zdt1013/wen-ai-cli.git
cd wen-ai-cli
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
go build
```

4. Add the compiled executable to system PATH (optional)

## 🚀 Usage

### ⌨️ Basic Commands

```bash
# Single-round question mode
> wen + enter any question
```

![Wen AI Usage Example](docs/example1.png)

### 💬 Start Chat Mode
```bash
> wen chat + enter first question
or
> wen chat <enter>
  enter any question
```
![Wen AI Usage Example](docs/example2.png)

### 🛠️ Other Commands
```bash
# View configuration
> wen config --help
# View help information
> wen --help
```

### 🔧 Configuration Instructions

Before using Wen AI, you need to configure necessary parameters such as API keys. You can configure them using the following commands:

```bash
> wen config --apiKey YOUR_API_KEY --baseUrl YOUR_API_BASE --model YOUR_API_MODEL
or
> wen config -k YOUR_API_KEY -u YOUR_API_BASE -m YOUR_API_MODEL
```

## 📁 Project Structure

```
wen-ai-cli/
├── action/     # Command action implementation
├── assets/     # Static resources
│   └── lang/   # Language packages
├── cmd/        # Command line definitions
├── docs/       # Documentation and attachments
├── common/     # Common components
├── execute/    # Executor
├── logger/     # Logging module
├── model/      # Data models
├── setup/      # Initialization settings
├── validate/   # Validator
├── wenai/      # Core functionality implementation
├── main.go     # Program entry
├── go.mod      # Go module definition
└── go.sum      # Go module dependency verification
```

## ⚠️ Known Issues
  * Non-command line related question responses, reply style not yet optimized

## 🔮 Future Plans 
 * Manual Mode: wen man
 * Perception: Detect host user, non-sudo user, intelligently adjust commands
 * Perception: Detect installed and uninstalled commands on the local machine
 * Toolchain (functioncall, mcp) support
 * Thinking model compatibility
 * User system for easy installation and use without manual AI parameter configuration
 * User knowledge base and usage preferences to save personal usage habits, making Wen AI understand you better

## 📚 Dependencies
 * [urfave/cli](https://github.com/urfave/cli) - Command line application framework
 * [cloudwego/eino](https://github.com/cloudwego/eino) - ByteDance open source large language model (LLM) application development framework
 * [gookit/config](https://github.com/gookit/config) - Configuration management library
 * [gookit/i18n](https://github.com/gookit/i18n) - Internationalization support
 * [shirou/gopsutil](https://github.com/shirou/gopsutil) - System information collection
 * [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) - OpenAI API client
 * [gookit/slog](https://github.com/gookit/slog) - Logging library
 * [go-cmd/cmd](https://github.com/go-cmd/cmd) - Command execution library
 * [manifoldco/promptui](https://github.com/manifoldco/promptui) - Interactive command line interface

## 🤝 Contributing

Welcome to submit Issues and Pull Requests to help improve this project. Before submitting code, please ensure:

1. Code complies with project coding standards
2. Necessary tests are added
3. Relevant documentation is updated

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## 📞 Contact

For any questions or suggestions, please contact us through:

- Submit an Issue
