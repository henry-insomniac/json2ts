# 📚 切片和映射详解 - Go 的核心复合类型

> 深入理解我们项目中最重要的两种数据结构

## 🎯 切片（Slice）详解

### 什么是切片？
切片是 Go 中最重要的数据结构之一，它是对数组的抽象，提供了动态数组的功能。

```go
// 数组 vs 切片的区别
var array [5]string          // 数组：固定长度 5
var slice []string           // 切片：动态长度

// 在我们项目中的使用
type TypeGenerator struct {
    interfaces []string      // 切片：存储生成的接口定义
}
```

### 切片的内部结构
```go
// 切片内部实际上是一个结构体
type slice struct {
    ptr    *Element  // 指向底层数组的指针
    len    int       // 当前长度
    cap    int       // 容量（底层数组长度）
}
```

## 🔧 切片的创建方式

### 1. 使用字面量创建
```go
// 在我们项目中的例子
func main() {
    // 创建包含初始值的切片
    supportedTypes := []string{"string", "number", "boolean", "null"}
    
    // 空切片
    var interfaces []string
    interfaces = []string{}  // 等价于上面的声明
}
```

### 2. 使用 make 函数创建
```go
// make([]Type, length, capacity)
interfaces := make([]string, 0)       // 长度0，默认容量
interfaces = make([]string, 0, 10)    // 长度0，容量10
interfaces = make([]string, 5)        // 长度5，容量5，元素为零值

// 在我们的构造函数中
func NewTypeGenerator() *TypeGenerator {
    return &TypeGenerator{
        interfaces: make([]string, 0, 10),  // 预分配容量
    }
}
```

### 3. 从数组创建切片
```go
array := [...]string{"interface1", "interface2", "interface3"}
slice := array[:]        // 完整切片
slice = array[1:3]       // 索引1到2的切片
slice = array[:2]        // 从开始到索引1
slice = array[1:]        // 从索引1到结束
```

## ⚡ 切片操作详解

### append 操作
```go
func (tg *TypeGenerator) addInterface(definition string) {
    // append 是向切片添加元素的标准方法
    tg.interfaces = append(tg.interfaces, definition)
    
    // 添加多个元素
    newInterfaces := []string{"interface A", "interface B"}
    tg.interfaces = append(tg.interfaces, newInterfaces...)
}

// append 的工作原理
func demonstrateAppend() {
    slice := make([]string, 0, 2)  // 长度0，容量2
    fmt.Printf("初始: len=%d, cap=%d\n", len(slice), cap(slice))
    
    slice = append(slice, "第一个")
    fmt.Printf("添加1个: len=%d, cap=%d\n", len(slice), cap(slice))  // len=1, cap=2
    
    slice = append(slice, "第二个")
    fmt.Printf("添加2个: len=%d, cap=%d\n", len(slice), cap(slice))  // len=2, cap=2
    
    slice = append(slice, "第三个")  // 触发扩容
    fmt.Printf("扩容后: len=%d, cap=%d\n", len(slice), cap(slice))   // len=3, cap=4
}
```

### 切片遍历
```go
func (tg *TypeGenerator) printAllInterfaces() {
    // 方式1：使用索引
    for i := 0; i < len(tg.interfaces); i++ {
        fmt.Printf("%d: %s\n", i, tg.interfaces[i])
    }
    
    // 方式2：使用 range（推荐）
    for index, definition := range tg.interfaces {
        fmt.Printf("%d: %s\n", index, definition)
    }
    
    // 方式3：只要值，不要索引
    for _, definition := range tg.interfaces {
        fmt.Println(definition)
    }
    
    // 方式4：只要索引，不要值
    for index := range tg.interfaces {
        fmt.Printf("接口 %d\n", index)
    }
}
```

### 切片复制和删除
```go
func (tg *TypeGenerator) copySlice() []string {
    // 复制切片
    backup := make([]string, len(tg.interfaces))
    copy(backup, tg.interfaces)
    return backup
}

func (tg *TypeGenerator) removeInterface(index int) {
    if index < 0 || index >= len(tg.interfaces) {
        return  // 索引越界
    }
    
    // 删除指定索引的元素
    tg.interfaces = append(tg.interfaces[:index], tg.interfaces[index+1:]...)
}

func (tg *TypeGenerator) insertInterface(index int, definition string) {
    if index < 0 || index > len(tg.interfaces) {
        return  // 索引越界
    }
    
    // 在指定位置插入元素
    tg.interfaces = append(tg.interfaces[:index], 
        append([]string{definition}, tg.interfaces[index:]...)...)
}
```

## 🗺️ 映射（Map）详解

### 什么是映射？
映射是键值对的无序集合，类似其他语言中的字典、哈希表。

```go
// 在我们项目中，JSON 对象就是映射
var jsonObject map[string]interface{}

// JSON: {"name": "张三", "age": 25}
// 对应的 Go 映射:
jsonObject = map[string]interface{}{
    "name": "张三",
    "age":  float64(25),  // JSON 数字是 float64
}
```

## 🔧 映射的创建方式

### 1. 使用字面量创建
```go
// 创建并初始化
typeMapping := map[string]string{
    "string":  "string",
    "number":  "number", 
    "boolean": "boolean",
    "null":    "null",
}

// 空映射
var emptyMap map[string]interface{}  // nil 映射，不能直接赋值
emptyMap = map[string]interface{}{}  // 空但已初始化的映射
```

### 2. 使用 make 函数创建
```go
// make(map[KeyType]ValueType)
typeCache := make(map[string]string)
typeCache["string"] = "string"
typeCache["number"] = "number"

// 在我们的项目中缓存类型推断结果
type TypeGenerator struct {
    interfaces []string
    typeCache  map[string]string  // 缓存已推断的类型
}

func NewTypeGenerator() *TypeGenerator {
    return &TypeGenerator{
        interfaces: make([]string, 0),
        typeCache:  make(map[string]string),
    }
}
```

## ⚡ 映射操作详解

### 基本操作
```go
func (tg *TypeGenerator) demonstrateMapOperations() {
    // 创建映射
    typeMap := make(map[string]string)
    
    // 添加键值对
    typeMap["string"] = "string"
    typeMap["number"] = "number"
    
    // 读取值
    tsType := typeMap["string"]  // 获取值
    
    // 安全读取（检查键是否存在）
    if tsType, exists := typeMap["boolean"]; exists {
        fmt.Printf("找到类型: %s\n", tsType)
    } else {
        fmt.Println("类型不存在")
    }
    
    // 删除键值对
    delete(typeMap, "string")
    
    // 获取映射长度
    count := len(typeMap)
    fmt.Printf("映射中有 %d 个键值对\n", count)
}
```

### 遍历映射
```go
func (tg *TypeGenerator) analyzeJSONObject(obj map[string]interface{}) {
    fmt.Println("分析 JSON 对象:")
    
    // 遍历映射
    for key, value := range obj {
        tsType := tg.toTSType(value)
        fmt.Printf("字段 %s: %s -> %s\n", key, getGoType(value), tsType)
    }
}

// 辅助函数：获取 Go 类型名
func getGoType(value interface{}) string {
    switch value.(type) {
    case string:
        return "string"
    case float64:
        return "float64"
    case bool:
        return "bool"
    case nil:
        return "nil"
    case []interface{}:
        return "[]interface{}"
    case map[string]interface{}:
        return "map[string]interface{}"
    default:
        return "unknown"
    }
}
```

## 🎯 实际项目应用

### 处理复杂 JSON 结构
```go
func (tg *TypeGenerator) processNestedObject(data map[string]interface{}) {
    // 示例 JSON:
    // {
    //   "user": {
    //     "name": "张三",
    //     "hobbies": ["读书", "编程"],
    //     "settings": {
    //       "theme": "dark"
    //     }
    //   }
    // }
    
    for key, value := range data {
        switch v := value.(type) {
        case string:
            fmt.Printf("%s: string\n", key)
            
        case float64:
            fmt.Printf("%s: number\n", key)
            
        case []interface{}:
            fmt.Printf("%s: array with %d elements\n", key, len(v))
            // 分析数组元素类型
            tg.analyzeArrayTypes(v)
            
        case map[string]interface{}:
            fmt.Printf("%s: nested object with %d fields\n", key, len(v))
            // 递归处理嵌套对象
            tg.processNestedObject(v)
        }
    }
}

func (tg *TypeGenerator) analyzeArrayTypes(arr []interface{}) {
    typeCount := make(map[string]int)
    
    for _, item := range arr {
        goType := getGoType(item)
        typeCount[goType]++
    }
    
    fmt.Printf("  数组类型统计: ")
    for goType, count := range typeCount {
        fmt.Printf("%s(%d) ", goType, count)
    }
    fmt.Println()
}
```

### 类型推断结果缓存
```go
func (tg *TypeGenerator) toTSTypeWithCache(value interface{}) string {
    // 为复杂类型生成缓存键
    cacheKey := tg.generateCacheKey(value)
    
    // 检查缓存
    if cachedType, exists := tg.typeCache[cacheKey]; exists {
        return cachedType
    }
    
    // 计算类型并缓存
    tsType := tg.computeTSType(value)
    tg.typeCache[cacheKey] = tsType
    
    return tsType
}

func (tg *TypeGenerator) generateCacheKey(value interface{}) string {
    switch v := value.(type) {
    case string:
        return fmt.Sprintf("str:%s", v)
    case float64:
        return fmt.Sprintf("num:%.2f", v)
    case []interface{}:
        return fmt.Sprintf("arr:len%d", len(v))
    case map[string]interface{}:
        keys := make([]string, 0, len(v))
        for key := range v {
            keys = append(keys, key)
        }
        return fmt.Sprintf("obj:%v", keys)
    default:
        return fmt.Sprintf("other:%T", v)
    }
}
```

## 🚀 性能优化技巧

### 切片性能优化
```go
// 不好的方式：频繁扩容
func generateManyInterfacesBad() []string {
    var interfaces []string  // 初始容量为0
    for i := 0; i < 1000; i++ {
        interfaces = append(interfaces, fmt.Sprintf("interface%d", i))
        // 每次append可能触发扩容，复制所有元素
    }
    return interfaces
}

// 好的方式：预分配容量
func generateManyInterfacesGood() []string {
    interfaces := make([]string, 0, 1000)  // 预分配容量
    for i := 0; i < 1000; i++ {
        interfaces = append(interfaces, fmt.Sprintf("interface%d", i))
        // 不会触发扩容，性能更好
    }
    return interfaces
}

// 最好的方式：直接设置长度
func generateManyInterfacesBest() []string {
    interfaces := make([]string, 1000)  // 预分配长度
    for i := 0; i < 1000; i++ {
        interfaces[i] = fmt.Sprintf("interface%d", i)
        // 直接赋值，无需append
    }
    return interfaces
}
```

### 映射性能优化
```go
// 预分配映射容量（Go 1.15+）
func createLargeMap() map[string]string {
    // 如果知道大致大小，预分配可以减少哈希表重建
    typeMap := make(map[string]string, 100)
    
    // 批量添加数据...
    return typeMap
}

// 使用字符串池减少内存分配
type StringPool struct {
    pool map[string]string
}

func NewStringPool() *StringPool {
    return &StringPool{
        pool: make(map[string]string),
    }
}

func (sp *StringPool) Get(s string) string {
    if cached, exists := sp.pool[s]; exists {
        return cached
    }
    sp.pool[s] = s
    return s
}
```

## 🎨 高级用法和模式

### 切片作为栈和队列
```go
type InterfaceStack struct {
    items []string
}

func (s *InterfaceStack) Push(item string) {
    s.items = append(s.items, item)
}

func (s *InterfaceStack) Pop() (string, bool) {
    if len(s.items) == 0 {
        return "", false
    }
    
    index := len(s.items) - 1
    item := s.items[index]
    s.items = s.items[:index]  // 移除最后一个元素
    return item, true
}

// 队列实现
type InterfaceQueue struct {
    items []string
}

func (q *InterfaceQueue) Enqueue(item string) {
    q.items = append(q.items, item)
}

func (q *InterfaceQueue) Dequeue() (string, bool) {
    if len(q.items) == 0 {
        return "", false
    }
    
    item := q.items[0]
    q.items = q.items[1:]  // 移除第一个元素
    return item, true
}
```

### 映射的高级模式
```go
// 嵌套映射处理 JSON 路径
type JSONPath map[string]interface{}

func (jp JSONPath) GetValue(path string) (interface{}, bool) {
    keys := strings.Split(path, ".")
    current := map[string]interface{}(jp)
    
    for _, key := range keys {
        if value, exists := current[key]; exists {
            if nextMap, ok := value.(map[string]interface{}); ok {
                current = nextMap
            } else {
                // 最后一层，返回值
                return value, true
            }
        } else {
            return nil, false
        }
    }
    
    return current, true
}

// 使用示例
jsonData := JSONPath{
    "user": map[string]interface{}{
        "profile": map[string]interface{}{
            "name": "张三",
        },
    },
}

name, exists := jsonData.GetValue("user.profile.name")
if exists {
    fmt.Printf("用户名: %s\n", name)  // 输出: 用户名: 张三
}
```

## ⚠️ 常见陷阱和注意事项

### 切片陷阱
```go
func slicePitfalls() {
    // 陷阱1：nil 切片 vs 空切片
    var nilSlice []string         // nil 切片
    emptySlice := []string{}      // 空切片
    
    fmt.Printf("nil切片: %v, 长度: %d\n", nilSlice == nil, len(nilSlice))      // true, 0
    fmt.Printf("空切片: %v, 长度: %d\n", emptySlice == nil, len(emptySlice))   // false, 0
    
    // 陷阱2：切片共享底层数组
    original := []string{"a", "b", "c", "d"}
    slice1 := original[1:3]  // ["b", "c"]
    slice2 := original[2:4]  // ["c", "d"]
    
    slice1[1] = "modified"   // 修改 slice1 影响了 slice2
    fmt.Printf("slice2: %v\n", slice2)  // ["modified", "d"]
    
    // 安全做法：复制切片
    safeCopy := make([]string, len(slice1))
    copy(safeCopy, slice1)
}
```

### 映射陷阱
```go
func mapPitfalls() {
    // 陷阱1：nil 映射不能写入
    var nilMap map[string]string  // nil 映射
    // nilMap["key"] = "value"    // 运行时panic!
    
    // 正确做法
    nilMap = make(map[string]string)
    nilMap["key"] = "value"  // OK
    
    // 陷阱2：映射不是线程安全的
    // 如果多个 goroutine 同时读写映射，需要使用 sync.Map 或加锁
    
    // 陷阱3：映射遍历顺序是随机的
    typeMap := map[string]string{
        "a": "first",
        "b": "second", 
        "c": "third",
    }
    
    for key := range typeMap {
        fmt.Printf("键: %s\n", key)  // 每次运行顺序可能不同
    }
}
```

## 🎯 实际代码示例：改进的类型生成器

```go
type AdvancedTypeGenerator struct {
    interfaces      []string              // 生成的接口定义
    typeCache       map[string]string     // 类型缓存
    interfaceNames  map[string]bool       // 已使用的接口名
    nameCounter     map[string]int        // 名称计数器
}

func NewAdvancedTypeGenerator() *AdvancedTypeGenerator {
    return &AdvancedTypeGenerator{
        interfaces:     make([]string, 0, 20),
        typeCache:      make(map[string]string, 50),
        interfaceNames: make(map[string]bool, 20),
        nameCounter:    make(map[string]int),
    }
}

func (atg *AdvancedTypeGenerator) generateUniqueInterfaceName(baseName string) string {
    if !atg.interfaceNames[baseName] {
        atg.interfaceNames[baseName] = true
        return baseName
    }
    
    // 名称已存在，生成唯一名称
    counter := atg.nameCounter[baseName]
    counter++
    atg.nameCounter[baseName] = counter
    
    uniqueName := fmt.Sprintf("%s%d", baseName, counter)
    atg.interfaceNames[uniqueName] = true
    return uniqueName
}

func (atg *AdvancedTypeGenerator) handleComplexArray(arr []interface{}) string {
    if len(arr) == 0 {
        return "any[]"
    }
    
    // 使用映射统计类型分布
    typeStats := make(map[string]int)
    
    for _, item := range arr {
        tsType := atg.toTSType(item)
        typeStats[tsType]++
    }
    
    // 根据类型分布决定如何处理
    if len(typeStats) == 1 {
        // 单一类型数组
        for tsType := range typeStats {
            return tsType + "[]"
        }
    }
    
    // 混合类型数组，生成联合类型
    types := make([]string, 0, len(typeStats))
    for tsType := range typeStats {
        types = append(types, tsType)
    }
    
    return fmt.Sprintf("(%s)[]", strings.Join(types, " | "))
}
```

## 💡 小结

通过本章，你应该掌握了：

### 切片核心概念
1. **内部结构**：指针、长度、容量的三元组
2. **创建方式**：字面量、make、从数组切取
3. **操作方法**：append、copy、遍历、删除、插入
4. **性能优化**：预分配容量，避免频繁扩容

### 映射核心概念  
1. **键值对存储**：无序的关联数组
2. **创建和初始化**：字面量和 make 函数
3. **安全操作**：存在性检查，避免nil映射写入
4. **遍历和处理**：range 循环，类型断言

### 项目应用
1. **JSON 处理**：映射表示对象，切片存储接口
2. **缓存优化**：使用映射缓存类型推断结果
3. **批量处理**：切片存储和处理多个数据项
4. **性能调优**：预分配容量，减少内存分配

这些知识为我们构建高效的类型推断引擎提供了坚实基础！

---

**下一章**：[错误处理](06-error-handling.md) 🚀