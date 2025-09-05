# 📘 Go 语言入门基础

> 为完全没有接触过 Go 语言的新手准备的基础教程

## 🌟 Go 语言是什么？

Go（也称 Golang）是 Google 开发的开源编程语言，特别适合：
- **后端服务器开发** 
- **命令行工具**（我们的 json2ts 就是这种）
- **系统编程**
- **云原生应用**

## 🔑 Go 语言的核心特点

### 1. 静态类型 + 简洁语法
```go
// Go 中声明变量的方式
var name string = "json2ts"    // 完整声明
var age = 25                   // 类型推断
title := "工具"                 // 短声明（最常用）
```

### 2. 强大的类型系统
```go
// 基础类型
var text string = "hello"      // 字符串
var count int = 42             // 整数
var price float64 = 3.14       // 浮点数
var isActive bool = true       // 布尔值
```

### 3. 指针但不复杂
```go
// 指针：存储变量的内存地址
var x int = 10
var p *int = &x    // p 是指向 x 的指针
fmt.Println(*p)    // *p 获取指针指向的值，输出：10
```

## 📦 包和导入

### 什么是包（Package）？
**包**是 Go 代码组织的基本单位，类似其他语言的"模块"或"命名空间"：

```go
package main           // 声明这个文件属于 main 包

import (              // 导入其他包
    "fmt"             // 格式化输出包
    "os"              // 操作系统接口包
    "encoding/json"   // JSON 处理包
)
```

### 包的规则
- `package main` - 表示这是可执行程序
- `func main()` - 程序的入口函数
- 每个 `.go` 文件都必须声明属于哪个包

## 🏗️ 函数基础

### 函数定义语法
```go
// 基本格式：func 函数名(参数) 返回类型 { 函数体 }
func greet(name string) string {
    return "Hello, " + name
}

// 多返回值（Go 的特色功能）
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("除数不能为零")
    }
    return a / b, nil
}
```

### 在我们的项目中
```go
// 我们的类型推断函数
func toTSType(val interface{}) string {
    // 接收任意类型，返回字符串
    switch val.(type) {
    case string:
        return "string"
    case float64:
        return "number"
    }
}
```

## 🔄 控制结构

### if 语句
```go
// Go 的 if 语句很干净，不需要括号
if len(os.Args) < 2 {
    fmt.Println("Usage: json2ts <file>")
    return
}

// if 语句可以包含初始化
if err := json.Unmarshal(data, &v); err != nil {
    // 处理错误
}
```

### for 循环
```go
// Go 只有 for 循环，但很灵活

// 传统 for 循环
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// 遍历切片/数组
for index, value := range slice {
    fmt.Printf("索引: %d, 值: %v\n", index, value)
}

// 遍历字典
for key, value := range map {
    fmt.Printf("键: %s, 值: %v\n", key, value)
}
```

### switch 语句
```go
// Go 的 switch 不需要 break
switch val.(type) {
case string:
    return "string"
case float64:
    return "number"
case bool:
    return "boolean"
default:
    return "any"
}
```

## ❗ 错误处理

### Go 的错误处理哲学
Go 不使用异常（try/catch），而是将错误作为返回值：

```go
// 标准模式：返回结果和错误
data, err := os.ReadFile("file.json")
if err != nil {
    // 处理错误
    fmt.Println("读取文件失败:", err)
    return
}
// 使用 data
```

### 为什么这样设计？
- **显式错误处理**：强制开发者考虑错误情况
- **代码更清晰**：一眼就能看出哪里可能出错
- **性能更好**：不需要异常栈的开销

## 🧩 接口（interface{}）

### 空接口的魔法
```go
// interface{} 表示"任意类型"
var anything interface{}
anything = "hello"     // 可以是字符串
anything = 42          // 可以是数字
anything = []int{1,2}  // 可以是切片
```

### 类型断言
```go
// 检查接口中的具体类型
switch v := val.(type) {
case string:
    fmt.Println("这是字符串:", v)
case int:
    fmt.Println("这是整数:", v)
case []interface{}:
    fmt.Println("这是数组:", v)
}
```

## 💡 实际例子：在我们的项目中

```go
// 我们项目中的实际应用
func toTSType(val interface{}) string {
    // val 是空接口，可以接收任何类型
    switch v := val.(type) {          // 类型断言
    case string:                      // 如果是字符串
        return "string"
    case float64:                     // JSON 数字都是 float64
        return "number"
    case []interface{}:               // 数组类型
        // 递归处理数组元素
        return handleArray(v)
    default:                          // 其他情况
        return "any"
    }
}
```

## 🎯 小结

经过这一章，你应该理解了：

1. **Go 的基本语法**：变量声明、函数定义
2. **包系统**：如何组织和导入代码
3. **控制结构**：if、for、switch 的用法
4. **错误处理**：Go 独特的错误处理方式
5. **接口概念**：`interface{}` 的强大之处
6. **类型断言**：如何判断接口中的具体类型

这些概念在我们改进类型推断引擎时都有体现！

---

**下一章**：[数据类型和变量详解](02-data-types-variables.md) 🚀