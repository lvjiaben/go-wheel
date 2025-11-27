package gen

// GenConfig 代码生成配置
type GenConfig struct {
	// 基础配置
	TableName       string `json:"table_name"`        // 表名（如 user_info）
	TableComment    string `json:"table_comment"`     // 表注释（用于菜单名称）
	ModuleName      string `json:"module_name"`       // 模块名（如 user_info，用于目录和文件名）
	PackageName     string `json:"package_name"`      // 包名（如 user_info）
	StructName      string `json:"struct_name"`       // 结构体名（如 UserInfo）
	FrontendSrcPath string `json:"frontend_src_path"` // 前端 src 目录绝对路径

	// 方法生成配置
	Methods MethodConfig `json:"methods"` // 方法配置

	// 字段配置
	Fields []FieldConfig `json:"fields"` // 字段配置列表

	// 搜索配置
	SearchFields []string `json:"search_fields"` // 模糊搜索字段列表（用于 WithSearchFields）

	// Operate 配置
	OperateFields []string `json:"operate_fields"` // Operate 方法允许操作的字段列表

	// 排序配置
	DefaultSortField string `json:"default_sort_field"` // 默认排序字段（如 id, sort, weigh）
	DefaultSortOrder string `json:"default_sort_order"` // 默认排序类型（asc 或 desc）

	// 菜单配置
	MenuConfig MenuConfig `json:"menu_config"` // 菜单配置
}

// MethodConfig 方法生成配置
type MethodConfig struct {
	List    bool `json:"list"`    // 是否生成 List 方法（必选，始终为 true）
	Create  bool `json:"create"`  // 是否生成 Create 方法
	Update  bool `json:"update"`  // 是否生成 Update 方法
	Delete  bool `json:"delete"`  // 是否生成 Delete 方法
	Operate bool `json:"operate"` // 是否生成 Operate 方法
}

// FieldConfig 字段配置
type FieldConfig struct {
	// 数据库字段信息
	ColumnName      string `json:"column_name"`       // 数据库列名（如 user_name）
	ColumnType      string `json:"column_type"`       // 数据库类型（如 varchar(255)）
	ColumnComment   string `json:"column_comment"`    // 字段注释
	IsNullable      bool   `json:"is_nullable"`       // 是否可为空
	ColumnDefault   string `json:"column_default"`    // 默认值
	IsPrimaryKey    bool   `json:"is_primary_key"`    // 是否主键
	IsAutoIncrement bool   `json:"is_auto_increment"` // 是否自增

	// Go 字段信息
	FieldName string `json:"field_name"` // Go 字段名（如 UserName）
	FieldType string `json:"field_type"` // Go 类型（如 string）
	JsonTag   string `json:"json_tag"`   // JSON 标签（如 user_name）
	GormTag   string `json:"gorm_tag"`   // GORM 标签

	// 验证器配置
	InCreate      bool   `json:"in_create"`      // 是否在 Create 验证器中
	InUpdate      bool   `json:"in_update"`      // 是否在 Update 验证器中
	IsRequired    bool   `json:"is_required"`    // 是否必填
	ValidateRules string `json:"validate_rules"` // 验证规则（如 required,min=1,max=32）

	// 前端表格配置
	ShowInTable      bool   `json:"show_in_table"`      // 是否在表格中显示
	TableSortable    bool   `json:"table_sortable"`     // 是否支持表格排序
	TableSearchable  bool   `json:"table_searchable"`   // 是否支持表格搜索
	TableSearchType  string `json:"table_search_type"`  // 表格搜索类型（input/select/date_range等）
	TableDisplayType string `json:"table_display_type"` // 表格显示类型（text/tag/time/image/edit等）
	SearchFormType   string `json:"search_form_type"`   // 搜索表单组件类型（与 FormComponent 一致）

	// 前端表单配置
	ShowInForm         bool                   `json:"show_in_form"`         // 是否在表单中显示
	FormComponent      string                 `json:"form_component"`       // 表单组件类型（Input/Select/DatePicker等）
	FormComponentProps map[string]interface{} `json:"form_component_props"` // 表单组件参数
	ComponentProps     string                 `json:"component_props"`      // 组件配置（JSON 字符串）

	// 智能识别标记
	IsOperateField  bool `json:"is_operate_field"`  // 是否为 Operate 字段（status、is_xxx等）
	IsSortField     bool `json:"is_sort_field"`     // 是否为排序字段（sort、weigh等）
	IsTimeField     bool `json:"is_time_field"`     // 是否为时间字段（xxx_at、xxx_time等）
	IsRelationField bool `json:"is_relation_field"` // 是否为关联字段（xxx_id等）
	IsEnumField     bool `json:"is_enum_field"`     // 是否为枚举字段
	IsSetField      bool `json:"is_set_field"`      // 是否为集合字段
	IsTextField     bool `json:"is_text_field"`     // 是否为长文本字段
	IsBoolField     bool `json:"is_bool_field"`     // 是否为布尔字段
	IsImageField    bool `json:"is_image_field"`    // 是否为图片字段
	IsImagesField   bool `json:"is_images_field"`   // 是否为多图字段
}

// MenuConfig 菜单配置
type MenuConfig struct {
	ParentMenuName string `json:"parent_menu_name"` // 父级菜单名称（如"代码生成"）
	MenuName       string `json:"menu_name"`        // 菜单名称（使用表注释或表名）
	MenuIcon       string `json:"menu_icon"`        // 菜单图标（可留空）
	MenuSort       int    `json:"menu_sort"`        // 菜单排序
}

// TableInfo 数据库表信息
type TableInfo struct {
	TableName    string       `json:"table_name"`    // 表名
	TableComment string       `json:"table_comment"` // 表注释
	Columns      []ColumnInfo `json:"columns"`       // 列信息
}

// ColumnInfo 数据库列信息
type ColumnInfo struct {
	ColumnName         string `json:"column_name"`              // 列名
	ColumnType         string `json:"column_type"`              // 列类型（如 varchar(255)）
	DataType           string `json:"data_type"`                // 数据类型（如 varchar）
	ColumnComment      string `json:"column_comment"`           // 列注释
	IsNullable         string `json:"is_nullable"`              // 是否可为空（YES/NO）
	ColumnDefault      string `json:"column_default"`           // 默认值
	ColumnKey          string `json:"column_key"`               // 键类型（PRI/UNI/MUL）
	Extra              string `json:"extra"`                    // 额外信息（如 auto_increment）
	CharacterMaxLength int    `json:"character_maximum_length"` // 字符最大长度
	NumericPrecision   int    `json:"numeric_precision"`        // 数字精度
	NumericScale       int    `json:"numeric_scale"`            // 数字小数位数
}

// GeneratedCode 生成的代码
type GeneratedCode struct {
	Backend  BackendCode  `json:"backend"`  // 后端代码
	Frontend FrontendCode `json:"frontend"` // 前端代码
	Menu     MenuCode     `json:"menu"`     // 菜单 SQL
}

// BackendCode 后端代码
type BackendCode struct {
	Controller string `json:"controller"` // Controller 代码
	Service    string `json:"service"`    // Service 代码
	Model      string `json:"model"`      // Model 代码
	Validate   string `json:"validate"`   // Validate 代码
	Route      string `json:"route"`      // 路由代码
}

// FrontendCode 前端代码
type FrontendCode struct {
	API        string `json:"api"`          // API 文件代码
	ListView   string `json:"list_view"`    // 列表页面代码
	DataTS     string `json:"data_ts"`      // data.ts 代码
	FormVue    string `json:"form_vue"`     // form.vue 代码
	LocaleZhCN string `json:"locale_zh_cn"` // 中文语言包
	LocaleEnUS string `json:"locale_en_us"` // 英文语言包
}

// MenuCode 菜单 SQL
type MenuCode struct {
	SQL string `json:"sql"` // 菜单 SQL 语句
}

// FilePath 生成的文件路径
type FilePath struct {
	Backend  BackendFilePath  `json:"backend"`  // 后端文件路径
	Frontend FrontendFilePath `json:"frontend"` // 前端文件路径
}

// BackendFilePath 后端文件路径
type BackendFilePath struct {
	Controller string `json:"controller"` // Controller 文件路径
	Service    string `json:"service"`    // Service 文件路径
	Model      string `json:"model"`      // Model 文件路径
	Validate   string `json:"validate"`   // Validate 文件路径
}

// FrontendFilePath 前端文件路径
type FrontendFilePath struct {
	API        string `json:"api"`          // API 文件路径
	ListView   string `json:"list_view"`    // 列表页面路径
	DataTS     string `json:"data_ts"`      // data.ts 路径
	FormVue    string `json:"form_vue"`     // form.vue 路径
	LocaleZhCN string `json:"locale_zh_cn"` // 中文语言包路径
	LocaleEnUS string `json:"locale_en_us"` // 英文语言包路径
}
