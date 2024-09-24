# myGoTemplate

## 目录结构

```plaintext
myGoTemplate/
├── cmd/
│   ├── app1/
│   │   └── main.go
│   └── app2/
│       └── main.go
├── internal/
│   ├── app1/
│   │   ├── handler/
│   │   ├── service/
│   │   └── repository/
│   └── app2/
│       ├── handler/
│       ├── service/
│       └── repository/
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── logger/
│   │   └── logger.go
│   └── utils/
│       └── utils.go
├── go.mod
└── README.md
```

## 目录结构解释

- **cmd/**：每个子目录代表一个可执行的应用程序入口，`main.go` 包含应用的启动代码。
- **internal/**：包含项目的内部代码，不希望被外部依赖。
  - **app1/** 和 **app2/**：分别对应不同的应用程序逻辑。
    - **handler/**：负责处理请求和响应。
    - **service/**：业务逻辑层，处理具体的业务规则。
    - **repository/**：数据访问层，负责与数据库或外部服务交互。
- **pkg/**：公共库，可被本项目和其他项目引用。
  - **config/**：配置管理，加载和解析配置文件。
  - **logger/**：日志库，统一日志记录格式和输出。
  - **utils/**：通用工具函数。
- **go.mod**：Go 模块依赖管理文件。
- **README.md**：项目说明文档。

## 最佳实践说明

- **模块化结构**：将代码按功能拆分为不同的模块，如 `handler`、`service`、`repository`，实现关注点分离。
- **使用 `internal`**：将内部实现放在 `internal` 目录下，防止被外部包引用，保持实现的私有性。
- **`pkg` 目录**：存放公共的可重用代码，如配置、日志等，可被本项目和其他项目使用。
- **配置管理**：使用环境变量或配置文件，通过 `pkg/config` 统一管理，方便不同环境的配置。
- **日志处理**：通过 `pkg/logger` 统一日志输出，方便日志级别控制和格式统一。
- **多应用支持**：在 `cmd/` 目录下为每个应用程序提供独立的入口，便于管理和部署。

## 如何使用

1. **克隆模板**：将该模板框架克隆到本地，或者根据该结构创建新的项目。
2. **修改模块名称**：将 `myGoTemplate` 改为你的模块名称，并在 `go.mod` 中更新。或使用脚本 `shell/init.sh` 进行修改。
3. **开发业务逻辑**：在 `internal/app1` 或 `internal/app2` 中编写具体的业务逻辑。
4. **运行应用程序**：
   ```bash
   go run cmd/app1/main.go
   ```
5. **构建可执行文件**：
   ```bash
   make app1
   make app2
   ```

## ref

https://github.com/songquanpeng/gin-template
