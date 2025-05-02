package validate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lvjiaben/go-wheel/pkg/utils"
)

// Generate 生成验证器
func Generate(structName string) error {
	// 获取结构体名称
	name := utils.ToSnakeCase(structName)

	// 生成文件内容
	content := fmt.Sprintf(`package validate

// %s 验证器
type %s struct {
	// TODO: 添加字段
}
`, structName, structName)

	// 写入文件
	filename := filepath.Join("app/backend/validate", name+".go")
	return os.WriteFile(filename, []byte(content), 0644)
}
