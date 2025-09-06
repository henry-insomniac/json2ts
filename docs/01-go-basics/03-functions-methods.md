# 🔧 函数和方法详解

> 深入理解 Go 函数系统和我们项目中的核心函数设计

## 🎯 函数基础回顾

### 函数的基本语法
```go
func functionName(parameter type) returnType {
    // 函数体
    return value
}
```

在我们的项目中最核心的函数：
```go
func toTSType(val interface{}) string {
    // 接收任意类型，返回 TypeScript 类型字符串
    switch val.(type) {
    case string:
        return "string"
    case float64:
        return "number"
    }
}
```

## 🔄 多返回值：Go 的特色功能

### 为什么需要多返回值？
在我们的项目中，经常需要同时返回结果和错误信息：

```go
// 单返回值的问题
func readJSONFile(filename string) string {
    data, err := os.ReadFile(filename)
    if err != nil {
        // 错误怎么处理？只能panic或返回空字符串
        return ""
    }
    return string(data)
}

// Go 的多返回值解决方案
func readJSONFile(filename string) (string, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return "", err    // 返回空字符串和错误
    }
    return string(data), nil  // 返回数据和nil错误
}
```

### 在我们项目中的应用
```go
// 我们可以改进类型推断函数
func (tg *TypeGenerator) generateInterfaceDefinition(name string, obj map[string]interface{}) (string, error) {
    if len(obj) == 0 {
        return "", fmt.Errorf("空对象无法生成接口")
    }
    
    var fields []string
    for key, value := range obj {
        tsType := tg.toTSType(value)
        fields = append(fields, fmt.Sprintf("  %s: %s;", key, tsType))
    }
    
    interfaceDef := fmt.Sprintf("export interface %s {\n%s\n}", name, strings.Join(fields, "\n"))
    return interfaceDef, nil  // 成功返回接口定义和nil错误
}
```

## 🎭 函数类型和高阶函数

### 函数作为类型
Go 中函数也是一种类型，可以作为变量传递：

```go
// 定义函数类型
type TypeConverter func(interface{}) string

// 使用函数类型
func processValue(val interface{}, converter TypeConverter) string {
    return converter(val)
}

// 在我们项目中的应用
func main() {
    data := "hello"
    
    // 使用我们的类型转换函数
    result := processValue(data, func(val interface{}) string {
        switch val.(type) {
        case string:
            return "string"
        default:
            return "unknown"
        }
    })
    
    fmt.Println(result)  // 输出: string
}
```

### 闭包和函数工厂
```go
// 创建专门的类型转换器工厂
func createTypeConverter(defaultType string) func(interface{}) string {
    return func(val interface{}) string {
        switch v := val.(type) {
        case string:
            return "string"
        case float64:
            return "number"
        case bool:
            return "boolean"
        case nil:
            return "null"
        default:
            return defaultType  // 使用闭包中的默认类型
        }
    }
}

// 使用
converter := createTypeConverter("any")
tsType := converter(someValue)
```

## 🔍 defer 语句：优雅的清理

### defer 的基本用法
`defer` 语句会延迟执行，直到包含它的函数返回时才执行：

```go
func processJSONFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()  // 确保文件最终会被关闭
    
    // 读取和处理文件
    // 即使出现错误，文件也会被正确关闭
    return processFile(file)
}
```

### 在我们项目中的应用
```go
func (tg *TypeGenerator) generateFromFile(filename string) error {
    fmt.Printf("开始处理文件: %s\n", filename)
    defer fmt.Printf("完成处理文件: %s\n", filename)  // 最后执行
    
    jsonData, err := os.ReadFile(filename)
    if err != nil {
        return err  // 即使返回错误，defer 也会执行
    }
    
    var data map[string]interface{}
    if err := json.Unmarshal(jsonData, &data); err != nil {
        return err
    }
    
    tg.generateInterface("Root", data)
    return nil
}
```

## 🎯 错误处理模式

### 标准错误处理
```go
func (tg *TypeGenerator) handleArrayWithValidation(arr []interface{}) (string, error) {
    if len(arr) == 0 {
        return "any[]", nil  // 空数组返回 any[]
    }
    
    // 分析数组元素类型
    elementTypes := make(map[string]bool)
    for i, item := range arr {
        if item == nil {
            continue  // 跳过 null 值
        }
        
        tsType := tg.toTSType(item)
        if tsType == "" {
            return "", fmt.Errorf("无法处理数组索引 %d 的类型", i)
        }
        elementTypes[tsType] = true
    }
    
    // 生成联合类型
    if len(elementTypes) == 1 {
        for elemType := range elementTypes {
            return elemType + "[]", nil
        }
    }
    
    // 多种类型的数组
    types := make([]string, 0, len(elementTypes))
    for elemType := range elementTypes {
        types = append(types, elemType)
    }
    return fmt.Sprintf("(%s)[]", strings.Join(types, " | ")), nil
}
```

### 错误包装和链式错误
```go
import "fmt"

func (tg *TypeGenerator) processComplexObject(obj map[string]interface{}) error {
    for key, value := range obj {
        if err := tg.validateField(key, value); err != nil {
            // 包装错误，提供更多上下文
            return fmt.Errorf("处理字段 '%s' 时出错: %w", key, err)
        }
    }
    return nil
}

func (tg *TypeGenerator) validateField(key string, value interface{}) error {
    if key == "" {
        return fmt.Errorf("字段名不能为空")
    }
    
    if value == nil {
        return fmt.Errorf("字段值不能为 nil")
    }
    
    return nil
}
```

## 🚀 函数性能优化

### 避免重复计算
```go
// 不好的方式：每次都重新计算
func (tg *TypeGenerator) toTSType(val interface{}) string {
    switch v := val.(type) {
    case []interface{}:
        if len(v) == 0 {
            return "any[]"
        }
        
        // 每次都重新分析所有元素
        elementTypes := make(map[string]bool)
        for _, item := range v {
            elementTypes[tg.toTSType(item)] = true
        }
        // ... 处理逻辑
    }
    return "any"
}

// 好的方式：使用缓存
type TypeGenerator struct {
    interfaces  []string
    typeCache   map[string]string  // 添加类型缓存
}

func (tg *TypeGenerator) toTSTypeWithCache(val interface{}) string {
    // 为复杂类型生成哈希作为缓存键
    cacheKey := generateCacheKey(val)
    if cachedType, exists := tg.typeCache[cacheKey]; exists {
        return cachedType
    }
    
    result := tg.computeTSType(val)
    tg.typeCache[cacheKey] = result
    return result
}
```

### 预分配切片容量
```go
// 不好的方式
func (tg *TypeGenerator) generateFields(obj map[string]interface{}) []string {
    var fields []string  // 初始容量为0，会多次扩容
    for key, value := range obj {
        tsType := tg.toTSType(value)
        fields = append(fields, fmt.Sprintf("%s: %s", key, tsType))
    }
    return fields
}

// 好的方式
func (tg *TypeGenerator) generateFields(obj map[string]interface{}) []string {
    fields := make([]string, 0, len(obj))  // 预分配容量
    for key, value := range obj {
        tsType := tg.toTSType(value)
        fields = append(fields, fmt.Sprintf("%s: %s", key, tsType))
    }
    return fields
}
```

## 🎨 函数设计模式

### 选项模式（Options Pattern）
```go
// 配置结构
type GeneratorConfig struct {
    PrettyPrint    bool
    SortFields     bool
    UseOptional    bool
    InterfacePrefix string
}

// 选项函数类型
type GeneratorOption func(*GeneratorConfig)

// 选项构造函数
func WithPrettyPrint() GeneratorOption {
    return func(config *GeneratorConfig) {
        config.PrettyPrint = true
    }
}

func WithSortedFields() GeneratorOption {
    return func(config *GeneratorConfig) {
        config.SortFields = true
    }
}

func WithInterfacePrefix(prefix string) GeneratorOption {
    return func(config *GeneratorConfig) {
        config.InterfacePrefix = prefix
    }
}

// 构造函数使用选项
func NewTypeGeneratorWithOptions(options ...GeneratorOption) *TypeGenerator {
    config := &GeneratorConfig{
        PrettyPrint:     false,
        SortFields:      false,
        UseOptional:     false,
        InterfacePrefix: "I",
    }
    
    // 应用所有选项
    for _, option := range options {
        option(config)
    }
    
    return &TypeGenerator{
        interfaces: make([]string, 0),
        config:     config,
    }
}

// 使用方式
generator := NewTypeGeneratorWithOptions(
    WithPrettyPrint(),
    WithSortedFields(),
    WithInterfacePrefix("Type"),
)
```

### 流式处理模式
```go
type ProcessingPipeline struct {
    steps []func(interface{}) (interface{}, error)
}

func NewPipeline() *ProcessingPipeline {
    return &ProcessingPipeline{
        steps: make([]func(interface{}) (interface{}, error), 0),
    }
}

func (p *ProcessingPipeline) AddStep(step func(interface{}) (interface{}, error)) *ProcessingPipeline {
    p.steps = append(p.steps, step)
    return p
}

func (p *ProcessingPipeline) Process(input interface{}) (interface{}, error) {
    result := input
    for i, step := range p.steps {
        var err error
        result, err = step(result)
        if err != nil {
            return nil, fmt.Errorf("步骤 %d 处理失败: %w", i+1, err)
        }
    }
    return result, nil
}

// 使用示例
pipeline := NewPipeline().
    AddStep(validateJSON).
    AddStep(parseJSON).
    AddStep(normalizeTypes).
    AddStep(generateInterfaces)

result, err := pipeline.Process(jsonString)
```

## 🧪 单元测试友好的函数设计

### 纯函数设计
```go
// 好的设计：纯函数，易于测试
func convertJSONTypeToTS(jsonValue interface{}) string {
    switch v := jsonValue.(type) {
    case string:
        return "string"
    case float64:
        return "number"
    case bool:
        return "boolean"
    case nil:
        return "null"
    default:
        return "any"
    }
}

// 测试代码
func TestConvertJSONTypeToTS(t *testing.T) {
    tests := []struct {
        input    interface{}
        expected string
    }{
        {"hello", "string"},
        {42.0, "number"},
        {true, "boolean"},
        {nil, "null"},
    }
    
    for _, test := range tests {
        result := convertJSONTypeToTS(test.input)
        if result != test.expected {
            t.Errorf("输入 %v，期望 %s，得到 %s", test.input, test.expected, result)
        }
    }
}
```

### 依赖注入模式
```go
// 接口定义
type FileReader interface {
    ReadFile(filename string) ([]byte, error)
}

type OSFileReader struct{}

func (r OSFileReader) ReadFile(filename string) ([]byte, error) {
    return os.ReadFile(filename)
}

// 可测试的函数设计
func ProcessJSONFromFile(reader FileReader, filename string) ([]string, error) {
    data, err := reader.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    // 处理数据...
    return []string{"interface1", "interface2"}, nil
}

// 测试时使用模拟的 FileReader
type MockFileReader struct {
    data []byte
    err  error
}

func (m MockFileReader) ReadFile(filename string) ([]byte, error) {
    return m.data, m.err
}
```

## 🎯 实际应用示例

### 完整的类型推断函数
```go
func (tg *TypeGenerator) generateTypeScriptDefinition(jsonData []byte) (string, error) {
    // 1. 解析 JSON
    var data interface{}
    if err := json.Unmarshal(jsonData, &data); err != nil {
        return "", fmt.Errorf("JSON 解析失败: %w", err)
    }
    
    // 2. 重置生成器状态
    tg.interfaces = make([]string, 0)
    tg.interfaceCounter = 0
    
    // 3. 生成根接口
    rootType := tg.toTSType(data)
    
    // 4. 如果根类型是对象，生成主接口
    if rootType != "any" && rootType != "string" && rootType != "number" {
        mainInterface := fmt.Sprintf("export interface Root {\n  // 由 %s 类型表示\n}", rootType)
        tg.interfaces = append([]string{mainInterface}, tg.interfaces...)
    }
    
    // 5. 返回所有接口定义
    return strings.Join(tg.interfaces, "\n\n"), nil
}
```

## 💡 最佳实践总结

### 函数设计原则
1. **单一职责**：每个函数只做一件事
2. **纯函数优先**：相同输入总是产生相同输出
3. **错误处理**：明确的错误返回和处理
4. **性能考虑**：避免不必要的计算和内存分配

### 命名约定
- 函数名使用驼峰命名：`generateInterface`
- 布尔函数用 `is/has/can` 开头：`isValidType`
- 转换函数用 `to` 开头：`toTSType`
- 处理函数用 `handle/process` 开头：`handleArray`

## 🎯 小结

通过本章，你应该掌握了：

1. **多返回值**：Go 独特的错误处理方式
2. **函数类型**：函数作为一等公民的使用
3. **defer 语句**：优雅的资源清理
4. **错误处理**：标准化的错误处理模式
5. **性能优化**：缓存和预分配的技巧
6. **设计模式**：选项模式和流式处理
7. **测试友好**：纯函数和依赖注入

这些技能让我们能够写出更强大、更可维护的类型推断引擎！

---

**下一章**：[结构体和方法](04-structs-methods.md) 🚀