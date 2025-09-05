# 🏗️ 结构体和方法 - Go 的面向对象

> 本章解释我们在类型推断改进中用到的核心概念

## 🎯 什么是结构体（Struct）？

**结构体**是 Go 中组织数据的方式，类似其他语言中的"类"，但更简单：

```go
// 定义一个结构体
type Person struct {
    Name string    // 姓名字段
    Age  int       // 年龄字段
}

// 创建结构体实例
var p Person
p.Name = "张三"
p.Age = 25

// 或者使用字面量创建
p2 := Person{
    Name: "李四",
    Age:  30,
}
```

## 🔧 我们项目中的 TypeGenerator 结构体

让我们分析我们改进代码中使用的结构体：

```go
// TypeGenerator 结构体用于生成 TypeScript 类型
type TypeGenerator struct {
    interfaces       []string  // 存储生成的接口
    interfaceCounter int       // 接口计数器
}
```

### 为什么要用结构体？

**改进前**（只有函数）：
```go
func toTSType(val interface{}) string {
    // 问题：无法保存状态，不能生成多个接口
    return "simple type"
}
```

**改进后**（使用结构体）：
```go
type TypeGenerator struct {
    interfaces       []string  // 保存所有生成的接口
    interfaceCounter int       // 记录生成了几个接口
}
```

### 结构体的优势：
1. **状态保持**：可以记住之前生成的接口
2. **数据组织**：相关的数据放在一起
3. **方法绑定**：可以给结构体添加方法

## 🎭 什么是方法（Method）？

**方法**就是绑定到结构体上的函数：

```go
// 普通函数
func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

// 方法：绑定到 Person 结构体
func (p Person) greet() {
    fmt.Printf("Hello, I'm %s!\n", p.Name)
}

// 使用方法
person := Person{Name: "王五"}
person.greet()  // 输出：Hello, I'm 王五!
```

## 🔍 方法接收器详解

### 值接收器 vs 指针接收器

这是 Go 中非常重要的概念：

```go
// 值接收器：接收结构体的副本
func (p Person) getValue() string {
    return p.Name  // 只能读取，不能修改原始数据
}

// 指针接收器：接收结构体的地址
func (p *Person) setValue(name string) {
    p.Name = name  // 可以修改原始数据
}
```

### 在我们项目中的应用

```go
// 我们使用指针接收器，因为需要修改结构体的状态
func (tg *TypeGenerator) generateInterface(name string, obj map[string]interface{}) string {
    // tg 是指向 TypeGenerator 的指针
    // 我们需要修改 tg.interfaces 和 tg.interfaceCounter
    
    // 生成接口并添加到 interfaces 切片中
    interfaceDef := fmt.Sprintf("export interface %s {...}", name)
    tg.interfaces = append(tg.interfaces, interfaceDef)  // 修改结构体状态
    tg.interfaceCounter++                                 // 修改计数器
    
    return name
}
```

### 为什么使用指针接收器？

1. **修改原始数据**：我们需要向 `interfaces` 切片添加新接口
2. **性能考虑**：避免复制整个结构体
3. **一致性**：如果有一个方法使用指针接收器，通常所有方法都用指针接收器

## 🏭 构造函数模式

Go 没有构造函数，但我们可以创建"构造函数"：

```go
// NewTypeGenerator 创建新的类型生成器（构造函数模式）
func NewTypeGenerator() *TypeGenerator {
    return &TypeGenerator{
        interfaces:       make([]string, 0),  // 初始化空切片
        interfaceCounter: 0,                   // 计数器从 0 开始
    }
}

// 使用方式
generator := NewTypeGenerator()
```

### 构造函数的好处：
1. **初始化控制**：确保结构体正确初始化
2. **默认值设置**：为字段设置合理的默认值
3. **代码清晰**：一眼就能看出如何创建对象

## 🎨 方法链和流式接口

我们可以让方法返回结构体本身，实现链式调用：

```go
func (tg *TypeGenerator) addInterface(def string) *TypeGenerator {
    tg.interfaces = append(tg.interfaces, def)
    return tg  // 返回自身，支持链式调用
}

// 链式调用
generator.addInterface("interface A {}").addInterface("interface B {}")
```

## 🧩 实际代码分析

让我们分析我们项目中的完整方法：

```go
// handleObject 处理对象类型推断
func (tg *TypeGenerator) handleObject(obj map[string]interface{}) string {
    // 1. 生成接口名
    tg.interfaceCounter++
    interfaceName := fmt.Sprintf("Interface%d", tg.interfaceCounter)
    
    // 2. 生成接口定义
    return tg.generateInterface(interfaceName, obj)
}
```

**步骤解析**：
1. **修改计数器**：`tg.interfaceCounter++`
2. **生成唯一名称**：`Interface1`, `Interface2`, ...
3. **调用其他方法**：`tg.generateInterface(...)`
4. **返回接口名**：用于引用这个接口

## 🔄 方法之间的协作

```go
func (tg *TypeGenerator) toTSType(val interface{}) string {
    switch v := val.(type) {
    case map[string]interface{}:
        return tg.handleObject(v)     // 调用其他方法
    case []interface{}:
        return tg.handleArray(v)      // 调用其他方法
    default:
        return "simple type"
    }
}
```

### 方法协作的好处：
1. **职责分离**：每个方法只做一件事
2. **代码复用**：方法可以互相调用
3. **易于测试**：可以单独测试每个方法
4. **易于维护**：修改一个功能只需要改对应的方法

## 🎯 实践建议

### 1. 何时使用结构体？
- 需要**组织相关数据**时
- 需要**保持状态**时
- 需要给数据**添加行为**时

### 2. 何时使用指针接收器？
- 需要**修改结构体**时
- 结构体**比较大**时（性能考虑）
- 需要**一致性**时

### 3. 方法设计原则
- **单一职责**：每个方法只做一件事
- **清晰命名**：方法名要说明它的功能
- **合理返回**：返回调用者需要的信息

## 💡 与其他语言对比

| 概念 | Go | Java/C# | Python |
|------|----|---------| -------|
| 数据组织 | struct | class | class |
| 方法绑定 | 接收器 | 类方法 | 类方法 |
| 构造函数 | New函数 | constructor | __init__ |
| 继承 | 组合 | extends | extends |

## 🎯 小结

通过结构体和方法，我们实现了：

1. **状态保持**：TypeGenerator 记住生成的所有接口
2. **代码组织**：相关功能组织在一起
3. **递归处理**：方法可以调用其他方法处理复杂结构
4. **扩展性**：容易添加新的类型处理逻辑

这就是为什么我们的类型推断引擎变得更强大了！

---

**下一章**：[切片和映射详解](05-slices-maps.md) 🚀