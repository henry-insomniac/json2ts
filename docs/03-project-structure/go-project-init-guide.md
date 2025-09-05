# Go项目初始化指南

## 项目结构概览
```
json2ts/
├── README.md
├── docs/
│   └── go-project-init-guide.md
├── src/
├── go.mod
├── go.sum
├── main.go
├── cmd/
│   └── json2ts/
│       └── main.go
├── pkg/
│   ├── converter/
│   │   └── json_to_ts.go
│   └── utils/
└── tests/
    └── converter_test.go
```

## 初始化步骤

### 1. 初始化Go模块
```bash
go mod init github.com/username/json2ts
```

### 2. 创建主程序入口
```bash
mkdir -p cmd/json2ts
touch cmd/json2ts/main.go
```

### 3. 创建核心包结构
```bash
mkdir -p pkg/converter
mkdir -p pkg/utils
mkdir -p tests
```

### 4. 创建基础文件
- `main.go` - 项目根目录的主入口（可选）
- `cmd/json2ts/main.go` - CLI工具入口
- `pkg/converter/json_to_ts.go` - 核心转换逻辑
- `tests/converter_test.go` - 单元测试

### 5. 添加常用依赖
```bash
# JSON处理
go get github.com/tidwall/gjson

# CLI工具框架
go get github.com/spf13/cobra

# 测试工具
go get github.com/stretchr/testify
```

## 项目配置文件

### .gitignore
```
# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
go.work

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db
```

### Makefile（可选）
```makefile
.PHONY: build test clean run

build:
	go build -o bin/json2ts cmd/json2ts/main.go

test:
	go test ./...

clean:
	rm -rf bin/

run:
	go run cmd/json2ts/main.go

install:
	go install cmd/json2ts/main.go
```

## 下一步计划
1. 执行go mod init命令
2. 创建目录结构
3. 编写基础的main.go文件
4. 实现核心转换逻辑
5. 添加测试用例