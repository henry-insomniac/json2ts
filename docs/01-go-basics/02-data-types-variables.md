# 🔢 Go 数据类型和变量详解

> 深入理解我们项目中用到的 Go 数据类型

## 🎯 基础数据类型

### 字符串（string）
```go
var name string = "json2ts"
title := "JSON转TypeScript工具"

// 字符串操作
fmt.Printf("项目名: %s\n", name)      // 格式化输出
fullName := name + "-cli"              // 字符串拼接
```

### 数值类型
```go
// 整数类型
var count int = 42                     // 通用整数
var id int64 = 123456789              // 64位整数

// 浮点数类型
var price float64 = 3.14159           // 双精度浮点数
var ratio float32 = 0.5               // 单精度浮点数
```

### 布尔类型
```go
var isActive bool = true
var hasError bool = false

// 布尔运算
if isActive && !hasError {
    fmt.Println("系统正常运行")
}
```

## 🔍 Go 中的 interface{} 详解

在我们的项目中，最重要的是理解 `interface{}`：

```go
// interface{} 是空接口，可以存储任何类型
var data interface{}

data = "hello"           // 字符串
data = 42               // 整数
data = 3.14             // 浮点数
data = []string{"a"}    // 切片
data = map[string]int{} // 映射
```

### JSON 解析中的 interface{}
```go
// JSON 解析结果
var jsonData map[string]interface{}

// JSON: {"name": "张三", "age": 25, "active": true}
// 解析后:
// jsonData["name"]   -> string
// jsonData["age"]    -> float64 (注意：JSON数字默认解析为float64)
// jsonData["active"] -> bool
```

## 📚 复合类型

### 切片（Slice）
切片是 Go 中的动态数组，我们项目中广泛使用：

```go
// 创建切片
var interfaces []string                    // 空切片
interfaces = make([]string, 0)            // 使用 make 创建
interfaces = append(interfaces, "新接口")  // 添加元素

// 在我们项目中的使用
type TypeGenerator struct {
    interfaces []string    // 存储生成的接口定义
}
```

### 映射（Map）
映射类似其他语言的字典或哈希表：

```go
// JSON 对象在 Go 中表示为 map[string]interface{}
jsonObj := map[string]interface{}{
    "name": "张三",
    "age":  25,
}

// 访问映射
name := jsonObj["name"].(string)  // 类型断言获取具体值
```

### 结构体（Struct）
我们项目的核心数据结构：

```go
type TypeGenerator struct {
    interfaces       []string  // 接口定义列表
    interfaceCounter int       // 接口计数器
}

// 创建结构体实例
generator := TypeGenerator{
    interfaces:       make([]string, 0),
    interfaceCounter: 0,
}
```

## 🎭 类型断言和类型转换

### 类型断言
用于从 `interface{}` 中提取具体类型：

```go
func analyzeValue(val interface{}) {
    // 方式1：类型断言 + 检查
    if str, ok := val.(string); ok {
        fmt.Printf("这是字符串: %s\n", str)
    }
    
    // 方式2：类型 switch（我们项目中使用的）
    switch v := val.(type) {
    case string:
        fmt.Printf("字符串: %s\n", v)
    case float64:
        fmt.Printf("数字: %.2f\n", v)
    case bool:
        fmt.Printf("布尔值: %t\n", v)
    case []interface{}:
        fmt.Printf("数组，长度: %d\n", len(v))
    case map[string]interface{}:
        fmt.Printf("对象，键数量: %d\n", len(v))
    default:
        fmt.Printf("未知类型: %T\n", v)
    }
}
```

### 在我们项目中的实际应用
```go
func (tg *TypeGenerator) toTSType(val interface{}) string {
    switch v := val.(type) {
    case string:
        return "string"                    // 字符串 → TS string
    case float64:
        return "number"                    // 数字 → TS number
    case bool:
        return "boolean"                   // 布尔 → TS boolean
    case nil:
        return "null"                      // null → TS null
    case []interface{}:
        return tg.handleArray(v)           // 数组 → 递归处理
    case map[string]interface{}:
        return tg.handleObject(v)          // 对象 → 生成接口
    default:
        return "any"                       // 其他 → TS any
    }
}
```

## 🔄 变量声明的多种方式

### 完整声明
```go
var message string = "Hello, Go!"
var count int = 0
var isReady bool = false
```

### 类型推断
```go
var message = "Hello, Go!"    // Go 自动推断为 string
var count = 0                 // Go 自动推断为 int
var isReady = false           // Go 自动推断为 bool
```

### 短声明（最常用）
```go
message := "Hello, Go!"       // 只能在函数内使用
count := 0
isReady := false
```

### 批量声明
```go
var (
    name    string = "json2ts"
    version int    = 1
    debug   bool   = true
)
```

## 📊 JSON 数据类型映射表

理解 JSON 数据在 Go 中的类型映射非常重要：

| JSON 类型 | Go 类型 | TypeScript 类型 | 示例 |
|-----------|---------|-----------------|------|
| string | string | string | "hello" |
| number | float64 | number | 42, 3.14 |
| boolean | bool | boolean | true, false |
| null | nil | null | null |
| array | []interface{} | T[] | [1, 2, 3] |
| object | map[string]interface{} | interface | {"key": "value"} |

### 重要注意事项

1. **JSON 数字都是 float64**
```go
// JSON: {"age": 25}
// Go 中 age 是 float64，不是 int
age := jsonData["age"].(float64)
```

2. **数组元素可能是混合类型**
```go
// JSON: [1, "hello", true]
// Go: []interface{}{float64(1), "hello", true}
```

3. **嵌套对象需要递归处理**
```go
// JSON: {"user": {"name": "张三", "age": 25}}
user := jsonData["user"].(map[string]interface{})
name := user["name"].(string)
```

## 🎯 实际代码示例

### 处理复杂 JSON 结构
```go
// 示例 JSON
jsonStr := `{
    "name": "张三",
    "age": 25,
    "hobbies": ["读书", "编程"],
    "address": {
        "city": "北京",
        "zip": "100000"
    }
}`

// 解析到 Go 数据结构
var data map[string]interface{}
json.Unmarshal([]byte(jsonStr), &data)

// 访问数据
name := data["name"].(string)                              // "张三"
age := data["age"].(float64)                               // 25.0
hobbies := data["hobbies"].([]interface{})                 // ["读书", "编程"]
address := data["address"].(map[string]interface{})        // 嵌套对象
city := address["city"].(string)                           // "北京"
```

### 类型安全的处理方式
```go
func safeGetString(data map[string]interface{}, key string) (string, bool) {
    if value, exists := data[key]; exists {
        if str, ok := value.(string); ok {
            return str, true
        }
    }
    return "", false
}

// 使用
if name, ok := safeGetString(data, "name"); ok {
    fmt.Printf("姓名: %s\n", name)
} else {
    fmt.Println("姓名字段不存在或类型不正确")
}
```

## 💡 性能和内存考虑

### 切片的容量管理
```go
// 预分配容量，避免频繁扩容
interfaces := make([]string, 0, 10)  // 长度0，容量10

// 不好的方式：频繁扩容
var interfaces []string
for i := 0; i < 1000; i++ {
    interfaces = append(interfaces, fmt.Sprintf("interface%d", i))
}

// 好的方式：预分配
interfaces := make([]string, 0, 1000)
for i := 0; i < 1000; i++ {
    interfaces = append(interfaces, fmt.Sprintf("interface%d", i))
}
```

### 指针 vs 值
```go
// 大结构体使用指针传递
type LargeStruct struct {
    data [1000]int
}

// 值传递：复制整个结构体（慢）
func processValue(ls LargeStruct) {
    // 处理逻辑
}

// 指针传递：只传递地址（快）
func processPointer(ls *LargeStruct) {
    // 处理逻辑
}
```

## 🎯 小结

通过本章，你应该掌握了：

1. **Go 的基础数据类型**：string, int, float64, bool
2. **interface{} 的威力**：可以存储任何类型
3. **复合类型**：slice, map, struct
4. **类型断言**：从接口中提取具体类型
5. **JSON 映射关系**：JSON 数据如何在 Go 中表示
6. **性能考虑**：切片容量和指针使用

这些概念在我们的类型推断引擎中都有实际应用！

---

**下一章**：[函数和方法详解](03-functions-methods.md) 🚀