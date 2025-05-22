package utils

import (
	"encoding/json"
	"regexp"
	"strings"
)

// ParseHealthItemString 解析健康项目字符串为键值对
// 输入格式: "项目1:数值1, 项目2:数值2, 项目3:数值3"
// 输出: map[string]string{"项目1": "数值1", "项目2": "数值2", "项目3": "数值3"}
func ParseHealthItemString(itemString string) map[string]string {
	result := make(map[string]string)
	if itemString == "" {
		return result
	}

	// 正则表达式匹配 "key:value" 或 "key:" 格式
	re := regexp.MustCompile(`([^:,]+):([^,]*)(,|$)`)
	matches := re.FindAllStringSubmatch(itemString, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])
			if key != "" {
				result[key] = value
			}
		}
	}

	return result
}

// BuildHealthItemString 将键值对构建为健康项目字符串
// 输入: map[string]string{"项目1": "数值1", "项目2": "数值2", "项目3": "数值3"}
// 输出: "项目1:数值1, 项目2:数值2, 项目3:数值3"
func BuildHealthItemString(itemMap map[string]string) string {
	if len(itemMap) == 0 {
		return ""
	}

	var parts []string
	for key, value := range itemMap {
		parts = append(parts, key+":"+value)
	}

	return strings.Join(parts, ", ")
}

// ConvertToJSON 将健康项目字符串转换为JSON格式
// 输入格式: "项目1:数值1, 项目2:数值2, 项目3:数值3"
// 输出: JSON字符串
func ConvertToJSON(itemString string) (string, error) {
	itemMap := ParseHealthItemString(itemString)
	jsonBytes, err := json.Marshal(itemMap)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// ParseJSONToHealthItemString 将JSON格式转换为健康项目字符串
// 输入: JSON字符串
// 输出格式: "项目1:数值1, 项目2:数值2, 项目3:数值3"
func ParseJSONToHealthItemString(jsonString string) (string, error) {
	var itemMap map[string]string
	if err := json.Unmarshal([]byte(jsonString), &itemMap); err != nil {
		return "", err
	}
	return BuildHealthItemString(itemMap), nil
}