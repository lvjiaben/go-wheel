package datatype

import (
	"reflect"
	"strings"
)

// ToTreeAssocMap 对任意数据切片结构体进行树形整理
func ToTreeAssocMap(data interface{}, idField, parentField, childrenField string) []map[string]interface{} {
	if idField == "" {
		idField = "id"
	}
	if parentField == "" {
		parentField = "parent_id"
	}
	if childrenField == "" {
		childrenField = "children"
	}

	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		return nil
	}

	// 使用切片保持顺序，同时用map快速查找
	var allNodes []map[string]interface{}
	nodeMap := make(map[interface{}]map[string]interface{})
	var rootNodes []map[string]interface{}

	// 第一遍：创建所有节点，保持原始顺序
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		idValue := getFieldValueByJSONTag(item, idField)
		if idValue == nil {
			continue
		}

		// 转换为map，使用JSON标签
		nodeMapData := structToMapByJSONTag(item)
		nodeMapData[childrenField] = []map[string]interface{}{}

		// 添加到切片中保持顺序
		allNodes = append(allNodes, nodeMapData)
		// 同时添加到map中用于快速查找
		nodeMap[idValue] = nodeMapData
	}

	// 第二遍：构建树形结构，按照原始顺序处理
	for _, node := range allNodes {
		parentIDValue := node[parentField]

		if parentIDValue != nil && parentIDValue != 0 {
			// 有父节点
			if parent, exists := nodeMap[parentIDValue]; exists {
				// 添加到父节点的Children中
				children := parent[childrenField].([]map[string]interface{})
				parent[childrenField] = append(children, node)
			} else {
				// 找不到父节点，作为根节点
				rootNodes = append(rootNodes, node)
			}
		} else {
			// 没有父节点，作为根节点
			rootNodes = append(rootNodes, node)
		}
	}

	return rootNodes
}

// 结构体转map的辅助函数，使用JSON标签
func structToMapByJSONTag(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		if field.CanInterface() {
			// 获取JSON标签
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag == "" {
				// 如果没有JSON标签，使用字段名
				jsonTag = fieldType.Name
			} else {
				// 处理JSON标签，去掉omitempty等选项
				if idx := strings.Index(jsonTag, ","); idx != -1 {
					jsonTag = jsonTag[:idx]
				}
			}

			result[jsonTag] = field.Interface()
		}
	}

	return result
}

// 通过JSON标签获取字段值
func getFieldValueByJSONTag(v reflect.Value, jsonTag string) interface{} {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// 获取JSON标签
		tag := fieldType.Tag.Get("json")
		if tag == "" {
			tag = fieldType.Name
		} else {
			// 处理JSON标签，去掉omitempty等选项
			if idx := strings.Index(tag, ","); idx != -1 {
				tag = tag[:idx]
			}
		}

		if tag == jsonTag && field.CanInterface() {
			return field.Interface()
		}
	}
	return nil
}
