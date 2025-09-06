# ❗ Go 错误处理 - 优雅地处理异常情况

> Go 独特的错误处理哲学和实践方法

## 🎯 Go 错误处理的哲学

### 为什么 Go 不使用异常？

**其他语言的异常机制**：
```java
// Java 风格
try {
    String data = readFile("config.json");
    // 处理数据
} catch (IOException e) {
    // 处理异常
}
```

**Go 的错误处理**：
```go
// Go 风格
data, err := os.ReadFile("config.json")
if err != nil {
    // 处理错误
    return err
}
// 使用 data
```

### Go 方式的优势：

1. **显式错误处理**：强制开发者考虑错误情况
2. **代码路径清晰**：正常流程和错误流程分离明确
3. **性能更好**：没有异常栈的性能开销
4. **简单直接**：不需要复杂的异常层级

## 🔍 错误类型和创建

### 基本错误接口
```go
// Go 内置的 error 接口
type error interface {
    Error() string
}
```

### 创建错误的方式

#### 1. 使用 errors.New()
```go
import "errors"

func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("除数不能为零")
    }
    return a / b, nil
}
```

#### 2. 使用 fmt.Errorf()
```go
import "fmt"

func readConfig(filename string) error {
    if filename == "" {
        return fmt.Errorf("文件名不能为空")
    }
    
    if !fileExists(filename) {
        return fmt.Errorf("文件 %s 不存在", filename)
    }
    
    return nil
}
```

#### 3. 自定义错误类型
```go
// 自定义错误结构体
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("字段 %s 验证失败: %s", e.Field, e.Message)
}

// 使用自定义错误
func validateAge(age int) error {
    if age < 0 {
        return ValidationError{
            Field:   "age",
            Message: "年龄不能为负数",
        }
    }
    return nil
}
```

## 🛠️ 我们项目中的错误处理

### 文件读取错误处理
```go
func main() {
    // 检查命令行参数
    if len(os.Args) < 2 {
        fmt.Println("Usage: json2ts <file>")
        return  // 直接退出，不返回错误
    }
    
    // 读取文件
    jsonFile := os.Args[1]
    jsonData, err := os.ReadFile(jsonFile)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return  // 打印错误并退出
    }
    
    // JSON 解析
    var v map[string]interface{}
    if err := json.Unmarshal(jsonData, &v); err != nil {
        fmt.Println("Error parsing JSON:", err)
        return  // 打印错误并退出
    }
    
    // 继续处理...
}
```

### 改进的错误处理
让我们为项目添加更好的错误处理：

```go
// 定义项目特定的错误类型
type JSONError struct {
    Operation string
    Filename  string
    Err       error
}

func (e JSONError) Error() string {
    return fmt.Sprintf("%s 失败，文件: %s，错误: %v", e.Operation, e.Filename, e.Err)
}

// 包装错误的辅助函数
func wrapJSONError(op, filename string, err error) error {
    if err == nil {
        return nil
    }
    return JSONError{
        Operation: op,
        Filename:  filename,
        Err:       err,
    }
}

// 改进的主函数
func processJSON(filename string) error {
    // 读取文件
    jsonData, err := os.ReadFile(filename)
    if err != nil {
        return wrapJSONError("读取文件", filename, err)
    }
    
    // 解析 JSON
    var v map[string]interface{}
    if err := json.Unmarshal(jsonData, &v); err != nil {
        return wrapJSONError("解析JSON", filename, err)
    }
    
    // 生成 TypeScript
    generator := NewTypeGenerator()
    generator.generateInterface("Root", v)
    
    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: json2ts <file>")
        os.Exit(1)
    }
    
    if err := processJSON(os.Args[1]); err != nil {
        fmt.Fprintf(os.Stderr, "错误: %v\n", err)
        os.Exit(1)
    }
}
```

## 🔄 错误处理模式

### 1. 早期返回模式
```go
func processData(data []byte) (*Result, error) {
    // 验证输入
    if len(data) == 0 {
        return nil, errors.New("数据不能为空")
    }
    
    // 步骤1
    step1Result, err := step1(data)
    if err != nil {
        return nil, fmt.Errorf("步骤1失败: %w", err)
    }
    
    // 步骤2
    step2Result, err := step2(step1Result)
    if err != nil {
        return nil, fmt.Errorf("步骤2失败: %w", err)
    }
    
    return step2Result, nil
}
```

### 2. 错误聚合模式
```go
type ErrorList struct {
    errors []error
}

func (el *ErrorList) Add(err error) {
    if err != nil {
        el.errors = append(el.errors, err)
    }
}

func (el *ErrorList) Error() string {
    if len(el.errors) == 0 {
        return ""
    }
    
    var messages []string
    for _, err := range el.errors {
        messages = append(messages, err.Error())
    }
    return strings.Join(messages, "; ")
}

func (el *ErrorList) HasErrors() bool {
    return len(el.errors) > 0
}

// 使用错误聚合
func validateJSON(data map[string]interface{}) error {
    var errs ErrorList
    
    // 验证必需字段
    if _, exists := data["name"]; !exists {
        errs.Add(errors.New("缺少 name 字段"))
    }
    
    if _, exists := data["version"]; !exists {
        errs.Add(errors.New("缺少 version 字段"))
    }
    
    // 验证类型
    if name, exists := data["name"]; exists {
        if _, ok := name.(string); !ok {
            errs.Add(errors.New("name 字段必须是字符串"))
        }
    }
    
    if errs.HasErrors() {
        return &errs
    }
    
    return nil
}
```

### 3. 错误包装和解包装（Go 1.13+）
```go
import "errors"

// 包装错误
func readUserConfig(userID int) error {
    configPath := fmt.Sprintf("/configs/user_%d.json", userID)
    
    data, err := os.ReadFile(configPath)
    if err != nil {
        // 使用 %w 包装错误，保留原始错误信息
        return fmt.Errorf("读取用户 %d 的配置失败: %w", userID, err)
    }
    
    return nil
}

// 检查特定错误类型
func handleConfigError(err error) {
    // 检查是否是文件不存在错误
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("配置文件不存在，使用默认配置")
        return
    }
    
    // 检查是否是权限错误
    var permErr *os.PathError
    if errors.As(err, &permErr) {
        fmt.Printf("权限不足: %v\n", permErr)
        return
    }
    
    fmt.Printf("未知错误: %v\n", err)
}
```

## 🎯 在我们的 TypeGenerator 中添加错误处理

让我们为类型生成器添加错误处理：

```go
// 为 TypeGenerator 添加错误处理
type TypeGenerationError struct {
    Operation string
    Path      string
    Err       error
}

func (e TypeGenerationError) Error() string {
    return fmt.Sprintf("类型生成失败 [%s] 路径: %s, 错误: %v", e.Operation, e.Path, e.Err)
}

// 改进的 TypeGenerator
type TypeGenerator struct {
    interfaces       []string
    interfaceCounter int
    maxDepth         int  // 防止无限递归
    currentDepth     int
}

func NewTypeGenerator() *TypeGenerator {
    return &TypeGenerator{
        interfaces:       make([]string, 0),
        interfaceCounter: 0,
        maxDepth:         10,  // 最大递归深度
        currentDepth:     0,
    }
}

// 添加错误检查的类型推断
func (tg *TypeGenerator) toTSType(val interface{}, path string) (string, error) {
    // 检查递归深度
    if tg.currentDepth > tg.maxDepth {
        return "", TypeGenerationError{
            Operation: "类型推断",
            Path:      path,
            Err:       errors.New("超出最大递归深度"),
        }
    }
    
    switch v := val.(type) {
    case string:
        return "string", nil
    case float64:
        return "number", nil
    case bool:
        return "boolean", nil
    case nil:
        return "null", nil
    case []interface{}:
        return tg.handleArray(v, path)
    case map[string]interface{}:
        return tg.handleObject(v, path)
    default:
        return "", TypeGenerationError{
            Operation: "类型推断",
            Path:      path,
            Err:       fmt.Errorf("不支持的类型: %T", v),
        }
    }
}

func (tg *TypeGenerator) handleArray(arr []interface{}, path string) (string, error) {
    if len(arr) == 0 {
        return "any[]", nil
    }
    
    tg.currentDepth++
    defer func() { tg.currentDepth-- }()
    
    elementTypes := make(map[string]bool)
    for i, item := range arr {
        itemPath := fmt.Sprintf("%s[%d]", path, i)
        tsType, err := tg.toTSType(item, itemPath)
        if err != nil {
            return "", err
        }
        elementTypes[tsType] = true
    }
    
    if len(elementTypes) == 1 {
        for elemType := range elementTypes {
            return elemType + "[]", nil
        }
    }
    
    types := make([]string, 0, len(elementTypes))
    for elemType := range elementTypes {
        types = append(types, elemType)
    }
    return fmt.Sprintf("(%s)[]", strings.Join(types, " | ")), nil
}

func (tg *TypeGenerator) handleObject(obj map[string]interface{}, path string) (string, error) {
    tg.currentDepth++
    defer func() { tg.currentDepth-- }()
    
    tg.interfaceCounter++
    interfaceName := fmt.Sprintf("Interface%d", tg.interfaceCounter)
    
    var fields []string
    for key, value := range obj {
        fieldPath := fmt.Sprintf("%s.%s", path, key)
        tsType, err := tg.toTSType(value, fieldPath)
        if err != nil {
            return "", err
        }
        fields = append(fields, fmt.Sprintf("  %s: %s;", key, tsType))
    }
    
    interfaceDef := fmt.Sprintf("export interface %s {\n%s\n}", interfaceName, strings.Join(fields, "\n"))
    tg.interfaces = append(tg.interfaces, interfaceDef)
    
    return interfaceName, nil
}
```

## 💡 错误处理最佳实践

### 1. 错误信息要清晰
```go
// 不好的错误信息
return errors.New("error")

// 好的错误信息
return fmt.Errorf("解析 JSON 文件 %s 失败，第 %d 行存在语法错误", filename, lineNum)
```

### 2. 错误要包含上下文
```go
// 添加调用上下文
func processUserData(userID int) error {
    if err := validateUser(userID); err != nil {
        return fmt.Errorf("处理用户数据失败，用户ID: %d, 错误: %w", userID, err)
    }
    return nil
}
```

### 3. 适当的错误处理级别
```go
// 库函数：返回错误，让调用者决定如何处理
func parseJSON(data []byte) (*Result, error) {
    // 解析逻辑
    return result, err
}

// 应用层：处理错误并记录日志
func handleRequest(w http.ResponseWriter, r *http.Request) {
    result, err := parseJSON(requestData)
    if err != nil {
        log.Printf("解析请求失败: %v", err)
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    // 继续处理
}

// 主函数：处理错误并退出
func main() {
    if err := run(); err != nil {
        log.Fatalf("程序运行失败: %v", err)
    }
}
```

## 🎯 小结

Go 的错误处理虽然看起来繁琐，但带来了：

1. **明确的错误边界**：每个可能失败的操作都必须检查
2. **清晰的代码逻辑**：正常流程和错误流程分离
3. **更好的程序健壮性**：强制考虑异常情况
4. **易于调试**：错误信息包含完整的上下文

掌握 Go 的错误处理是成为合格 Go 开发者的重要一步！

---

**下一章**：[包和模块管理](06-packages-modules.md) 🚀