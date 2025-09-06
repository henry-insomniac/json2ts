package main

import (
	"fmt"
	"os"
	
	"github.com/henry-insomniac/json2ts/src/generator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: json2ts <file>")
		return
	}
	
	// 读取JSON数据
	jsonData, err := readJSONFile(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	
	// 解析并生成类型
	typeScript, err := convertJSONToTypeScript(jsonData)
	if err != nil {
		fmt.Println("Error converting:", err)
		return
	}
	
	// 输出结果
	fmt.Print(typeScript)
}

// readJSONFile 读取JSON文件
func readJSONFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// convertJSONToTypeScript 将JSON数据转换为TypeScript接口
func convertJSONToTypeScript(jsonData []byte) (string, error) {
	return generator.ConvertJSONToTypeScript(jsonData)
}
