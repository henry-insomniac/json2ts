package generator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TypeGenerator 结构体用于生成 TypeScript 类型
type TypeGenerator struct {
	interfaces       []string // 存储生成的接口
	interfaceCounter int      // 接口计数器
}

// New 创建新的类型生成器
func New() *TypeGenerator {
	return &TypeGenerator{
		interfaces:       make([]string, 0),
		interfaceCounter: 0,
	}
}

// ConvertJSONToTypeScript 将JSON数据转换为TypeScript接口定义
func ConvertJSONToTypeScript(jsonData []byte) (string, error) {
	// 解析JSON
	var v map[string]interface{}
	if err := json.Unmarshal(jsonData, &v); err != nil {
		return "", err
	}
	
	// 生成TypeScript
	generator := New()
	generator.generateInterface("Root", v)
	
	// 组装输出
	var result strings.Builder
	for i, interfaceDef := range generator.interfaces {
		if i > 0 {
			result.WriteString("\n\n")
		}
		result.WriteString(interfaceDef)
	}
	result.WriteString("\n")
	
	return result.String(), nil
}

// generateInterface 生成接口定义
func (tg *TypeGenerator) generateInterface(name string, obj map[string]interface{}) string {
	var fields []string
	for key, value := range obj {
		tsType := tg.toTSType(value)
		fields = append(fields, fmt.Sprintf("  %s: %s;", key, tsType))
	}
	
	interfaceDef := fmt.Sprintf("export interface %s {\n%s\n}", name, strings.Join(fields, "\n"))
	tg.interfaces = append(tg.interfaces, interfaceDef)
	return name
}

// toTSType 改进的类型推断函数
func (tg *TypeGenerator) toTSType(val interface{}) string {
	switch v := val.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case nil:
		return "null"
	case []interface{}:
		return tg.handleArray(v)
	case map[string]interface{}:
		return tg.handleObject(v)
	default:
		return "any"
	}
}

// handleArray 处理数组类型推断
func (tg *TypeGenerator) handleArray(arr []interface{}) string {
	if len(arr) == 0 {
		return "any[]"
	}
	
	// 分析数组元素类型
	elementTypes := make(map[string]bool)
	for _, item := range arr {
		tsType := tg.toTSType(item)
		elementTypes[tsType] = true
	}
	
	// 如果所有元素类型相同
	if len(elementTypes) == 1 {
		for elemType := range elementTypes {
			return elemType + "[]"
		}
	}
	
	// 多种类型，生成联合类型
	types := make([]string, 0, len(elementTypes))
	for elemType := range elementTypes {
		types = append(types, elemType)
	}
	return fmt.Sprintf("(%s)[]", strings.Join(types, " | "))
}

// handleObject 处理对象类型推断
func (tg *TypeGenerator) handleObject(obj map[string]interface{}) string {
	// 生成接口名
	tg.interfaceCounter++
	interfaceName := fmt.Sprintf("Interface%d", tg.interfaceCounter)
	
	// 生成接口定义
	return tg.generateInterface(interfaceName, obj)
}