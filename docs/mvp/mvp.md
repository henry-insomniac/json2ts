# 技术思路（最小 MVP）
1.	读文件：用 os.ReadFile 或 os.Stdin 读取 JSON 内容。
2.	解析 JSON：用标准库 encoding/json → map[string]interface{}。
3.	类型推断：对 map 里的值做类型判断：
    •	float64 → number
	•	string → string
	•	bool → boolean
	•	nil → null
	•	[]any → 数组
	•	map[string]any → 对象（递归处理）
4.	输出 TS 代码：拼接字符串，最后写到 types.ts。

```json
{
  "id": 1,
  "name": "insomniac",
  "active": true
}
```

```TS
export interface Root {
  id: number;
  name: string;
  active: boolean;
}
```