package validate

import (
	"strings"

	"github.com/lvjiaben/go-wheel/app/genertate/gorm"
	"github.com/lvjiaben/go-wheel/pkg/actions"
	"github.com/lvjiaben/go-wheel/pkg/global"
	"github.com/lvjiaben/go-wheel/pkg/initialize"
	"github.com/lvjiaben/go-wheel/pkg/utils/file"
)

// 初始化需要传入一个model
func Genertate(TableName string, PackageName string, Path string, Cover bool) {
	tableNames := strings.Split(TableName, ",")
	initialize.ViperLoad()
	if global.DB != nil {
		db, _ := global.DB.DB()
		defer db.Close()
	}
	db := global.DB
	tables := gorm.GetTables(db, tableNames, global.CONFIG.Mysql.Dbname)
	for _, table := range tables {
		fields := gorm.GetFields(db, table.Name)
		generateValidate(PackageName, Path, Cover, table, fields)
	}
}

func getRequire(null string) string {
	if null == "YES" {
		return "-"
	} else {
		return "required"
	}
}

// 生成Model
func generateValidate(PackageName string, Path string, Cover bool, table gorm.Table, fields []gorm.Field) {
	var builder strings.Builder
	builder.WriteString("package " + PackageName + "\n\n")
	list := []string{"Create", "Update", "Delete", "Sort"}
	for _, item := range list {
		builder.WriteString("type " + actions.Marshal(table.Name) + item + " struct {\n")
		for _, field := range fields {
			fieldName := field.Field
			if (item == "Create" || item == "Update") && !actions.IsStringInSlice(fieldName, []string{"created_at", "create_time", "updated_at", "update_time", "deleted_at", "delete_time"}) {
				if item == "Create" && field.Key == "PRI" {
					continue
				}
				builder.WriteString("\t" + actions.Marshal(fieldName) + "\t" + gorm.GetFiledType(field) + "\t" +
					"`" + "json:\"" + fieldName + "\" binding:\"" + getRequire(field.Null) + "\" msg:\"" + gorm.GetFieldZh(field) + "有误\"" + "`\n")
			}
			if item == "Delete" && field.Key == "PRI" {
				builder.WriteString("\t" + actions.Marshal(fieldName) + "\t" + gorm.GetFiledType(field) + "\t" +
					"`" + "json:\"" + fieldName + "\" binding:\"" + getRequire(field.Null) + "\" msg:\"" + gorm.GetFieldZh(field) + "有误\"" + "`\n")
			}
		}
		builder.WriteString("}\n\n")
	}
	// 文件生成
	fileName := Path + table.Name + ".go"
	file.MakeFile(Path, fileName, builder.String(), Cover)
}
