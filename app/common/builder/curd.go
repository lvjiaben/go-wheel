package builder

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/utils/datatype"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 全局 Schema 缓存
var schemaCache = &sync.Map{}

// ListApi 列表查询参数
type ListApi struct {
	Page      int    `form:"page" binding:"min=1"`      // 页码
	PageSize  int    `form:"page_size" binding:"min=1"` // 每页数量
	Pid       int    `form:"pid"`                       // 父级
	Filter    string `form:"filter"`                    // 筛选条件（JSON格式）
	Search    string `form:"search"`                    // 搜索关键词
	Parent    string `form:"parent"`                    // 父级（字符串）
	SortBy    string `form:"sort_by"`                   // 排序字段
	SortOrder string `form:"sort_order"`                // 排序方向（asc/desc）
}

type CRUDBuilder[T any] struct {
	db             *gorm.DB                                   // gorm（静态 DB，可选）
	dbProvider     func(*gin.Context) *gorm.DB                // DB Provider（动态获取 DB）
	ctx            *gin.Context                               // 当前请求上下文（通过 WithContext 设置）
	validFields    []string                                   // 模型字段列表
	searchFields   []string                                   // 允许模糊搜索的参数
	before         func(query interface{}, db *gorm.DB) error // curd执行前回调
	after          func(query interface{}, db *gorm.DB) error // curd执行后回调
	useTransaction bool                                       // 是否开启事务
	usePagination  bool                                       // 是否使用分页（List 方法）
	useFilter      bool                                       // 是否使用筛选器（List 方法）
}

// NewCRUDBuilder 创建 CRUD 构建器（使用静态 DB）
// 推荐使用 db.WithContext(ctx.Request.Context()) 创建带上下文的 DB Session
// 注意：不推荐在 Service 初始化时使用此方法，建议使用 NewCRUDBuilderWithProvider
func NewCRUDBuilder[T any](db *gorm.DB) *CRUDBuilder[T] {
	builder := &CRUDBuilder[T]{
		db:             db,
		useTransaction: true, // 默认使用事务
		usePagination:  true, // 默认使用分页
		useFilter:      true, // 默认使用筛选器
	}
	// 初始化有效字段列表
	builder.initValidFields()
	return builder
}

// NewCRUDBuilderWithProvider 创建 CRUD 构建器（使用 DB Provider）
// 推荐在 Service 初始化时使用此方法，每次请求会动态获取 DB Session
// provider: 返回 DB 的函数，会在每次 CRUD 操作时调用，自动绑定请求上下文
func NewCRUDBuilderWithProvider[T any](provider func(*gin.Context) *gorm.DB) *CRUDBuilder[T] {
	builder := &CRUDBuilder[T]{
		dbProvider:     provider,
		useTransaction: true, // 默认使用事务
		usePagination:  true, // 默认使用分页
		useFilter:      true, // 默认使用筛选器
	}
	// 初始化有效字段列表
	builder.initValidFields()
	return builder
}

// WithContext 设置请求上下文（链式操作）
func (b *CRUDBuilder[T]) WithContext(ctx *gin.Context) *CRUDBuilder[T] {
	b.ctx = ctx
	return b
}

// Session 创建一个新的 DB Session（推荐在每次请求时使用）
// 这样可以避免查询条件污染，并且可以绑定请求上下文
func (b *CRUDBuilder[T]) Session(db *gorm.DB) *CRUDBuilder[T] {
	newBuilder := &CRUDBuilder[T]{
		db:             db,
		dbProvider:     b.dbProvider,     // 复用 Provider
		ctx:            b.ctx,            // 复用上下文
		validFields:    b.validFields,    // 复用字段列表
		searchFields:   b.searchFields,   // 复用搜索字段
		before:         b.before,         // 复用回调
		after:          b.after,          // 复用回调
		useTransaction: b.useTransaction, // 复用配置
		usePagination:  b.usePagination,  // 复用配置
		useFilter:      b.useFilter,      // 复用配置
	}
	return newBuilder
}

// getDB 获取 DB 实例（优先使用 dbProvider）
func (b *CRUDBuilder[T]) getDB() *gorm.DB {
	if b.dbProvider != nil && b.ctx != nil {
		// 使用 Provider 动态获取 DB
		return b.dbProvider(b.ctx)
	}
	// 降级使用静态 DB
	return b.db
}

// initValidFields 初始化有效字段列表
func (b *CRUDBuilder[T]) initValidFields() {
	var model T
	var namer schema.Namer

	// 如果有静态 DB，使用其 NamingStrategy
	if b.db != nil {
		namer = b.db.NamingStrategy
	}

	// 如果没有 NamingStrategy，使用默认的
	if namer == nil {
		namer = schema.NamingStrategy{}
	}

	parsedSchema, err := schema.Parse(&model, schemaCache, namer)
	if err == nil && parsedSchema != nil {
		for _, field := range parsedSchema.Fields {
			b.validFields = append(b.validFields, field.DBName)
		}
	}
}

// WithTransaction 设置是否使用事务
func (b *CRUDBuilder[T]) WithTransaction(use bool) *CRUDBuilder[T] {
	b.useTransaction = use
	return b
}

func (b *CRUDBuilder[T]) WithPagination(use bool) *CRUDBuilder[T] {
	b.usePagination = use
	return b
}

func (b *CRUDBuilder[T]) WithFilter(use bool) *CRUDBuilder[T] {
	b.useFilter = use
	return b
}

func (b *CRUDBuilder[T]) WithSearchFields(fields ...string) *CRUDBuilder[T] {
	b.searchFields = fields
	return b
}

// Before 设置操作前的回调
func (b *CRUDBuilder[T]) Before(fn func(query interface{}, db *gorm.DB) error) *CRUDBuilder[T] {
	b.before = fn
	return b
}

// After 设置操作后的回调
func (b *CRUDBuilder[T]) After(fn func(query interface{}, db *gorm.DB) error) *CRUDBuilder[T] {
	b.after = fn
	return b
}

// ValidateFields 验证字段是否在模型中定义
func (b *CRUDBuilder[T]) ValidateFields(data map[string]interface{}) error {
	if len(b.validFields) == 0 {
		return fmt.Errorf("builder.curd.model_fields_empty")
	}

	invalidFields := []string{}
	for field := range data {
		if !datatype.Contains(b.validFields, field) {
			invalidFields = append(invalidFields, field)
		}
	}

	if len(invalidFields) > 0 {
		return fmt.Errorf("builder.curd.invalid_field: %s", strings.Join(invalidFields, ", "))
	}

	return nil
}

// Create 创建记录
func (b *CRUDBuilder[T]) Create(data interface{}) (*T, error) {
	// 获取 DB 实例
	db := b.getDB()

	// 转换数据为模型
	model, err := b.convertToModel(data)
	if err != nil {
		return nil, err
	}

	// 执行创建操作
	executeCreate := func(tx *gorm.DB) error {
		// 执行创建前回调
		if b.before != nil {
			if err := b.before(model, tx); err != nil {
				return err
			}
		}

		// 创建记录
		if err := tx.Create(model).Error; err != nil {
			return fmt.Errorf("builder.curd.create_failed: %w", err)
		}

		// 执行创建后回调
		if b.after != nil {
			if err := b.after(model, tx); err != nil {
				return err
			}
		}

		return nil
	}

	// 根据配置决定是否使用事务
	if b.useTransaction {
		if err := db.Transaction(executeCreate); err != nil {
			return nil, err
		}
	} else {
		if err := executeCreate(db); err != nil {
			return nil, err
		}
	}

	return model, nil
}

// Update 更新记录
func (b *CRUDBuilder[T]) Update(id interface{}, data interface{}) (*T, error) {
	// 获取 DB 实例
	db := b.getDB()

	// 查询现有记录
	var model T
	if err := db.First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("builder.curd.record_not_found")
		}
		return nil, fmt.Errorf("builder.curd.query_record_failed: %w", err)
	}

	// 转换数据
	updateData, err := b.convertToUpdateData(data)
	if err != nil {
		return nil, err
	}

	// 执行更新操作
	executeUpdate := func(tx *gorm.DB) error {
		// 执行更新前回调
		if b.before != nil {
			if err := b.before(&model, tx); err != nil {
				return err
			}
		}

		// 更新记录
		if err := tx.Model(&model).Updates(updateData).Error; err != nil {
			return fmt.Errorf("builder.curd.update_failed: %w", err)
		}

		// 重新查询更新后的记录
		if err := tx.First(&model, id).Error; err != nil {
			return fmt.Errorf("builder.curd.query_updated_record_failed: %w", err)
		}

		// 执行更新后回调
		if b.after != nil {
			if err := b.after(&model, tx); err != nil {
				return err
			}
		}

		return nil
	}

	// 根据配置决定是否使用事务
	if b.useTransaction {
		if err := db.Transaction(executeUpdate); err != nil {
			return nil, err
		}
	} else {
		if err := executeUpdate(db); err != nil {
			return nil, err
		}
	}

	return &model, nil
}

// Delete 删除记录
func (b *CRUDBuilder[T]) Delete(ids interface{}) error {
	// 获取 DB 实例
	db := b.getDB()

	// 转换为 ID 列表
	idList := b.convertToIDList(ids)
	if len(idList) == 0 {
		return fmt.Errorf("builder.curd.no_ids_to_delete")
	}

	// 执行删除操作
	executeDelete := func(tx *gorm.DB) error {
		// 执行删除前回调
		if b.before != nil {
			if err := b.before(ids, tx); err != nil {
				return err
			}
		}

		// 删除记录
		var model T
		if err := tx.Delete(&model, idList).Error; err != nil {
			return fmt.Errorf("builder.curd.delete_failed: %w", err)
		}

		// 执行删除后回调
		if b.after != nil {
			if err := b.after(ids, tx); err != nil {
				return err
			}
		}

		return nil
	}

	// 根据配置决定是否使用事务
	if b.useTransaction {
		return db.Transaction(executeDelete)
	}
	return executeDelete(db)
}

// List 查询列表
func (b *CRUDBuilder[T]) List() map[string]interface{} {
	// 解析查询参数
	var query ListApi
	query.Page = 1
	query.PageSize = 20
	if pageStr := b.ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			query.Page = page
		}
	}
	if pageSizeStr := b.ctx.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			query.PageSize = pageSize
		}
	}
	query.Pid, _ = strconv.Atoi(b.ctx.Query("pid"))
	query.Filter = b.ctx.Query("filter")
	query.Search = b.ctx.Query("search")
	query.Parent = b.ctx.Query("parent")
	query.SortBy = b.ctx.Query("sort_by")
	query.SortOrder = b.ctx.Query("sort_order")
	if query.SortOrder != "" && query.SortOrder != "asc" && query.SortOrder != "desc" {
		query.SortOrder = "asc"
	}
	// 获取数据库对象（优先使用 dbProvider）
	db := b.getDB()
	// 应用筛选器（如果启用）
	db = b.listFilter(db, &query)
	// 执行查询前回调
	if b.before != nil {
		if err := b.before(&query, db); err != nil {
			return nil
		}
	}
	// 查询数据
	var list []T
	var total int64
	// 先统计总数（使用 Model 避免消费查询条件）
	countDB := db.Session(&gorm.Session{})
	countDB.Model(new(T)).Count(&total)
	// 应用默认分页（如果启用了分页）
	if b.usePagination {
		db = db.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).Order("id DESC")
	}
	// 执行查询后回调（可以继续链式操作 db）
	if b.after != nil {
		if err := b.after(&list, db); err != nil {
			return nil
		}
	}
	// 最后执行查询
	db.Find(&list)

	return map[string]interface{}{
		"list":      list,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	}
}

// listFilter 列表筛选器（内部方法）
func (b *CRUDBuilder[T]) listFilter(db *gorm.DB, query *ListApi) *gorm.DB {
	// 处理筛选条件
	if b.useFilter && query.Filter != "" {
		var filter map[string]interface{}
		if err := json.Unmarshal([]byte(query.Filter), &filter); err == nil {
			for column, value := range filter {
				if !datatype.Contains(b.validFields, column) {
					continue
				}
				switch v := value.(type) {
				case []interface{}:
					if len(v) == 0 {
						continue
					}
					if len(v) == 2 {
						db = db.Where(column+" BETWEEN ? AND ?", v[0], v[1])
					} else {
						values := make([]string, 0, len(v))
						for _, val := range v {
							values = append(values, fmt.Sprintf("%v", val))
						}
						db = db.Where(column+" IN (?)", values)
					}
				default:
					strValue := fmt.Sprintf("%v", value)
					if strValue == "" {
						continue
					}
					db = db.Where(column+" = ?", value)
				}
			}
		}
	}

	// 处理搜索关键词
	if query.Search != "" && len(b.searchFields) > 0 {
		var conditions []string
		var args []interface{}

		for _, field := range b.searchFields {
			if datatype.Contains(b.validFields, field) {
				conditions = append(conditions, field+" LIKE ?")
				args = append(args, "%"+query.Search+"%")
			}
		}

		if len(conditions) > 0 {
			db = db.Where("("+strings.Join(conditions, " OR ")+")", args...)
		}
	}

	// 处理父级筛选
	if query.Pid > 0 && datatype.Contains(b.validFields, "pid") {
		db = db.Where("pid = ?", query.Pid)
	}

	// 处理排序
	if query.SortBy != "" {
		// 验证排序字段是否安全（防止SQL注入）
		if datatype.Contains(b.validFields, query.SortBy) {
			sortOrder := "ASC"
			if query.SortOrder == "desc" {
				sortOrder = "DESC"
			}
			db = db.Order(query.SortBy + " " + sortOrder)
		}
	}

	return db
}

// convertToModel 将数据转换为模型
func (b *CRUDBuilder[T]) convertToModel(data interface{}) (*T, error) {
	// 如果已经是模型类型
	switch v := data.(type) {
	case *T:
		return v, nil
	case T:
		return &v, nil
	case map[string]interface{}:
		// 验证字段
		if err := b.ValidateFields(v); err != nil {
			return nil, err
		}
		// 转换 map 为结构体
		return b.mapToStruct(v)
	default:
		// 尝试将任意结构体转换为模型（支持 Gin 验证器结构体）
		return b.structToModel(data)
	}
}

// structToModel 将任意结构体转换为模型（支持 Gin 验证器结构体）
func (b *CRUDBuilder[T]) structToModel(data interface{}) (*T, error) {
	// 先将结构体转换为 map
	dataMap, err := b.structToMap(data)
	if err != nil {
		return nil, fmt.Errorf("builder.curd.unsupported_data_type: %T", data)
	}

	// 验证字段
	if err := b.ValidateFields(dataMap); err != nil {
		return nil, err
	}

	// 再将 map 转换为模型
	return b.mapToStruct(dataMap)
}

// convertToUpdateData 将数据转换为更新数据
func (b *CRUDBuilder[T]) convertToUpdateData(data interface{}) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		// 验证字段
		if err := b.ValidateFields(v); err != nil {
			return nil, err
		}
		return v, nil
	case *T, T:
		// 如果是结构体，转换为 map
		return b.structToMap(v)
	default:
		// 尝试将任意结构体转换为 map（支持 Gin 验证器结构体）
		dataMap, err := b.structToMap(data)
		if err != nil {
			return nil, fmt.Errorf("builder.curd.unsupported_data_type: %T", data)
		}
		// 验证字段
		if err := b.ValidateFields(dataMap); err != nil {
			return nil, err
		}
		return dataMap, nil
	}
}

// convertToIDList 将 ID 转换为列表
func (b *CRUDBuilder[T]) convertToIDList(ids interface{}) []interface{} {
	switch v := ids.(type) {
	case []interface{}:
		return v
	case []int:
		result := make([]interface{}, len(v))
		for i, id := range v {
			result[i] = id
		}
		return result
	case []int64:
		result := make([]interface{}, len(v))
		for i, id := range v {
			result[i] = id
		}
		return result
	case []string:
		result := make([]interface{}, len(v))
		for i, id := range v {
			result[i] = id
		}
		return result
	default:
		// 单个 ID
		return []interface{}{ids}
	}
}

// structToMap 将结构体转换为 map
func (b *CRUDBuilder[T]) structToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("builder.curd.unsupported_data_type")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// 获取列名
		columnName := b.getColumnName(field)
		if columnName == "" || columnName == "-" {
			continue
		}

		// 跳过零值（可选）
		if !fieldValue.IsZero() {
			result[columnName] = fieldValue.Interface()
		}
	}

	// 验证字段
	if err := b.ValidateFields(result); err != nil {
		return nil, err
	}

	return result, nil
}

// getColumnName 获取字段的列名
func (b *CRUDBuilder[T]) getColumnName(field reflect.StructField) string {
	// 优先使用 gorm 标签
	gormTag := field.Tag.Get("gorm")
	if gormTag != "" {
		tags := strings.Split(gormTag, ";")
		for _, tag := range tags {
			if strings.HasPrefix(tag, "column:") {
				return strings.TrimPrefix(tag, "column:")
			}
		}
		// 如果有 "-" 标签，跳过该字段
		if gormTag == "-" {
			return "-"
		}
	}

	// 使用字段名转换为蛇形命名
	return toSnakeCase(field.Name)
}

// mapToStruct 将 map 转换为结构体
func (b *CRUDBuilder[T]) mapToStruct(data map[string]interface{}) (*T, error) {
	var model T
	modelValue := reflect.ValueOf(&model).Elem()
	modelType := modelValue.Type()

	// 遍历结构体字段
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		// 获取 gorm 标签中的列名
		gormTag := field.Tag.Get("gorm")
		var columnName string

		// 解析 gorm 标签获取列名
		if gormTag != "" {
			tags := strings.Split(gormTag, ";")
			for _, tag := range tags {
				if strings.HasPrefix(tag, "column:") {
					columnName = strings.TrimPrefix(tag, "column:")
					break
				}
			}
		}

		// 如果没有 column 标签，使用字段名转换为蛇形命名
		if columnName == "" {
			columnName = toSnakeCase(field.Name)
		}

		// 从 data 中获取值并设置到结构体
		if value, ok := data[columnName]; ok && fieldValue.CanSet() {
			if err := setFieldValue(fieldValue, value); err != nil {
				return nil, fmt.Errorf("builder.curd.convert_failed: %s: %w", field.Name, err)
			}
		}
	}

	return &model, nil
}

// setFieldValue 设置字段值
func setFieldValue(field reflect.Value, value interface{}) error {
	if value == nil {
		return nil
	}

	valueReflect := reflect.ValueOf(value)
	fieldType := field.Type()

	// 类型匹配，直接设置
	if valueReflect.Type().AssignableTo(fieldType) {
		field.Set(valueReflect)
		return nil
	}

	// 类型转换
	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case int:
			field.SetInt(int64(v))
		case int8:
			field.SetInt(int64(v))
		case int16:
			field.SetInt(int64(v))
		case int32:
			field.SetInt(int64(v))
		case int64:
			field.SetInt(v)
		case float64:
			field.SetInt(int64(v))
		case string:
			// 尝试从字符串转换
			var intVal int64
			_, err := fmt.Sscanf(v, "%d", &intVal)
			if err != nil {
				return fmt.Errorf("builder.curd.cannot_convert_string_to_int: %s", v)
			}
			field.SetInt(intVal)
		default:
			return fmt.Errorf("builder.curd.cannot_convert_to_int: %T", value)
		}
	case reflect.String:
		field.SetString(fmt.Sprintf("%v", value))
	case reflect.Float32, reflect.Float64:
		switch v := value.(type) {
		case float32:
			field.SetFloat(float64(v))
		case float64:
			field.SetFloat(v)
		case int:
			field.SetFloat(float64(v))
		case int64:
			field.SetFloat(float64(v))
		default:
			return fmt.Errorf("builder.curd.cannot_convert_to_type: %T", value)
		}
	case reflect.Bool:
		switch v := value.(type) {
		case bool:
			field.SetBool(v)
		case int:
			field.SetBool(v != 0)
		default:
			return fmt.Errorf("builder.curd.cannot_convert_to_bool: %T", value)
		}
	default:
		return fmt.Errorf("builder.curd.unsupported_field_type: %s", fieldType.Kind())
	}

	return nil
}

// toSnakeCase 将驼峰命名转换为蛇形命名
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
