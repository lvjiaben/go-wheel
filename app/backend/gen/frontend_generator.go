package gen

import (
	"fmt"
	"strings"
)

// FrontendGenerator 前端代码生成器
type FrontendGenerator struct {
	config *GenConfig
}

// NewFrontendGenerator 创建前端代码生成器
func NewFrontendGenerator(config *GenConfig) *FrontendGenerator {
	return &FrontendGenerator{config: config}
}

// Generate 生成前端代码
func (f *FrontendGenerator) Generate() FrontendCode {
	return FrontendCode{
		API:        f.GenerateAPI(),
		ListView:   f.GenerateListView(),
		DataTS:     f.GenerateDataTS(),
		FormVue:    f.GenerateFormVue(),
		LocaleZhCN: f.GenerateLocaleZhCN(),
		LocaleEnUS: f.GenerateLocaleEnUS(),
	}
}

// GenerateAPI 生成 API 文件
func (f *FrontendGenerator) GenerateAPI() string {
	var sb strings.Builder

	// 导入
	sb.WriteString("import { requestClient } from '#/api/request';\n\n")

	// 命名空间
	sb.WriteString(fmt.Sprintf("export namespace %sApi {\n", f.config.StructName))

	// 接口定义
	sb.WriteString(fmt.Sprintf("\t/** %s */\n", f.config.TableComment))
	sb.WriteString(fmt.Sprintf("\texport interface %s {\n", f.config.StructName))
	sb.WriteString("\t\t[key: string]: any;\n")

	// 字段定义
	for _, field := range f.config.Fields {
		if !field.ShowInTable && !field.ShowInForm {
			continue
		}

		comment := field.ColumnComment
		if comment == "" {
			comment = field.FieldName
		}

		tsType := f.convertGoTypeToTSType(field.FieldType)
		optional := ""
		if !field.IsRequired || field.ColumnName == "password" {
			optional = "?"
		}

		sb.WriteString(fmt.Sprintf("\t\t/** %s */\n", comment))
		sb.WriteString(fmt.Sprintf("\t\t%s%s: %s;\n", field.JsonTag, optional, tsType))
	}

	sb.WriteString("\t}\n\n")

	// 列表请求参数
	sb.WriteString("\t/** 列表请求参数 */\n")
	sb.WriteString("\texport interface ListParams {\n")
	sb.WriteString("\t\tpage: number;\n")
	sb.WriteString("\t\tpage_size: number;\n")
	sb.WriteString("\t\tsearch?: string;\n")
	sb.WriteString("\t\tfilter?: string;\n")
	sb.WriteString("\t\tsort_by?: string;\n")
	sb.WriteString("\t\tsort_order?: 'asc' | 'desc';\n")
	sb.WriteString("\t}\n\n")

	// 列表响应
	sb.WriteString("\t/** 列表响应 */\n")
	sb.WriteString("\texport interface ListResponse {\n")
	sb.WriteString(fmt.Sprintf("\t\tlist: %s[];\n", f.config.StructName))
	sb.WriteString("\t\ttotal: number;\n")
	sb.WriteString("\t\tpage: number;\n")
	sb.WriteString("\t\tlimit: number;\n")
	sb.WriteString("\t}\n\n")

	// Operate 参数
	if f.config.Methods.Operate {
		sb.WriteString("\t/** 操作字段参数 */\n")
		sb.WriteString("\texport interface OperateParams {\n")
		sb.WriteString("\t\tids?: number[];\n")
		sb.WriteString("\t\tfield: string;\n")
		sb.WriteString("\t\tvalue: number;\n")
		sb.WriteString("\t}\n\n")
	}

	sb.WriteString("}\n\n")

	// API 函数
	// 列表
	sb.WriteString("/**\n")
	sb.WriteString(fmt.Sprintf(" * 获取%s列表\n", f.config.TableComment))
	sb.WriteString(" */\n")
	sb.WriteString(fmt.Sprintf("async function get%sList(params: %sApi.ListParams) {\n",
		f.config.StructName, f.config.StructName))
	sb.WriteString(fmt.Sprintf("\treturn requestClient.get<%sApi.ListResponse>('/%s/list', { params });\n",
		f.config.StructName, f.config.ModuleName))
	sb.WriteString("}\n\n")

	// Create
	if f.config.Methods.Create {
		sb.WriteString("/**\n")
		sb.WriteString(fmt.Sprintf(" * 创建%s\n", f.config.TableComment))
		sb.WriteString(" */\n")
		sb.WriteString(fmt.Sprintf("async function create%s(\n", f.config.StructName))
		sb.WriteString(fmt.Sprintf("\tdata: Omit<%sApi.%s, 'id' | 'created_at' | 'updated_at'>,\n",
			f.config.StructName, f.config.StructName))
		sb.WriteString(") {\n")
		sb.WriteString(fmt.Sprintf("\treturn requestClient.post('/%s/create', data);\n", f.config.ModuleName))
		sb.WriteString("}\n\n")
	}

	// Update
	if f.config.Methods.Update {
		sb.WriteString("/**\n")
		sb.WriteString(fmt.Sprintf(" * 更新%s\n", f.config.TableComment))
		sb.WriteString(" */\n")
		sb.WriteString(fmt.Sprintf("async function update%s(\n", f.config.StructName))
		sb.WriteString(fmt.Sprintf("\tdata: Partial<%sApi.%s> & { id: number },\n",
			f.config.StructName, f.config.StructName))
		sb.WriteString(") {\n")
		sb.WriteString(fmt.Sprintf("\treturn requestClient.post('/%s/update', data);\n", f.config.ModuleName))
		sb.WriteString("}\n\n")
	}

	// Delete
	if f.config.Methods.Delete {
		sb.WriteString("/**\n")
		sb.WriteString(fmt.Sprintf(" * 删除%s\n", f.config.TableComment))
		sb.WriteString(" */\n")
		sb.WriteString(fmt.Sprintf("async function delete%s(data: any) {\n", f.config.StructName))
		sb.WriteString(fmt.Sprintf("\treturn requestClient.post('/%s/delete', data);\n", f.config.ModuleName))
		sb.WriteString("}\n\n")
	}

	// Operate
	if f.config.Methods.Operate {
		sb.WriteString("/**\n")
		sb.WriteString(fmt.Sprintf(" * 操作%s字段\n", f.config.TableComment))
		sb.WriteString(" */\n")
		sb.WriteString(fmt.Sprintf("async function operate%s(data: %sApi.OperateParams) {\n",
			f.config.StructName, f.config.StructName))
		sb.WriteString(fmt.Sprintf("\treturn requestClient.post('/%s/operate', data);\n", f.config.ModuleName))
		sb.WriteString("}\n\n")
	}

	// 导出
	sb.WriteString("export {\n")
	sb.WriteString(fmt.Sprintf("\tget%sList,\n", f.config.StructName))
	if f.config.Methods.Create {
		sb.WriteString(fmt.Sprintf("\tcreate%s,\n", f.config.StructName))
	}
	if f.config.Methods.Update {
		sb.WriteString(fmt.Sprintf("\tupdate%s,\n", f.config.StructName))
	}
	if f.config.Methods.Delete {
		sb.WriteString(fmt.Sprintf("\tdelete%s,\n", f.config.StructName))
	}
	if f.config.Methods.Operate {
		sb.WriteString(fmt.Sprintf("\toperate%s,\n", f.config.StructName))
	}
	sb.WriteString("};\n")

	return sb.String()
}

// convertGoTypeToTSType 将 Go 类型转换为 TypeScript 类型
func (f *FrontendGenerator) convertGoTypeToTSType(goType string) string {
	switch goType {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "number"
	case "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "string":
		return "string"
	case "[]byte":
		return "string"
	default:
		return "any"
	}
}

// generateListViewTemplate 生成列表页面的 template 部分
func (f *FrontendGenerator) generateListViewTemplate(sb *strings.Builder) {
	sb.WriteString("<template>\n")
	sb.WriteString("\t<Page>\n")

	// Form Drawer 放在 Grid 前面
	sb.WriteString("\t\t<FormDrawer @success=\"onRefresh\" />\n")

	// 表格
	sb.WriteString("\t\t<Grid>\n")

	// 图片字段插槽
	for _, field := range f.config.Fields {
		if field.ShowInTable && field.IsImageField {
			sb.WriteString(fmt.Sprintf("\t\t\t<template #%s=\"{ row }\">\n", field.JsonTag))
			sb.WriteString("\t\t\t\t<Image\n")
			sb.WriteString(fmt.Sprintf("\t\t\t\t\t:src=\"row.%s\"\n", field.JsonTag))
			sb.WriteString("\t\t\t\t\t:width=\"30\"\n")
			sb.WriteString("\t\t\t\t\t:height=\"30\"\n")
			sb.WriteString("\t\t\t\t\t:preview=\"true\"\n")
			sb.WriteString("\t\t\t\t\tloading=\"lazy\"\n")
			sb.WriteString("\t\t\t\t/>\n")
			sb.WriteString("\t\t\t</template>\n")
		}
	}

	// 状态字段插槽（只有配置为 tag 类型且是 operate 字段才生成可点击的 tag）
	for _, field := range f.config.Fields {
		if field.ShowInTable && field.TableDisplayType == "tag" && field.IsOperateField {
			sb.WriteString(fmt.Sprintf("\t\t\t<template #%s=\"{ row }\">\n", field.JsonTag))
			sb.WriteString("\t\t\t\t<Popconfirm\n")
			sb.WriteString("\t\t\t\t\t:title=\"$t('common.confirmTip')\"\n")
			sb.WriteString(fmt.Sprintf("\t\t\t\t\t@confirm=\"operateApi(row.%s === 1 ? 0 : 1, row, '%s')\"\n",
				field.JsonTag, field.JsonTag))
			sb.WriteString("\t\t\t\t>\n")
			sb.WriteString(fmt.Sprintf("\t\t\t\t\t<Tag :color=\"row.%s === 1 ? 'success' : 'error'\" class=\"ml-2 text-center\">\n",
				field.JsonTag))
			sb.WriteString(fmt.Sprintf("\t\t\t\t\t\t{{ row.%s === 1 ? $t('common.enable') : $t('common.disable') }}\n",
				field.JsonTag))
			sb.WriteString("\t\t\t\t\t</Tag>\n")
			sb.WriteString("\t\t\t\t</Popconfirm>\n")
			sb.WriteString("\t\t\t</template>\n")
		} else if field.ShowInTable && field.TableDisplayType == "tag" && !field.IsOperateField {
			// 普通 tag 类型（不可点击）
			sb.WriteString(fmt.Sprintf("\t\t\t<template #%s=\"{ row }\">\n", field.JsonTag))
			sb.WriteString(fmt.Sprintf("\t\t\t\t<Tag>{{ row.%s }}</Tag>\n", field.JsonTag))
			sb.WriteString("\t\t\t</template>\n")
		}
	}

	// link 类型插槽（参数跳转）
	for _, field := range f.config.Fields {
		if field.ShowInTable && field.TableDisplayType == "link" {
			// 计算跳转路径：如果是 _id 结尾，取前缀作为模块名；否则用整个字段名
			fieldName := field.JsonTag
			moduleName := fieldName
			if strings.HasSuffix(fieldName, "_id") {
				moduleName = strings.TrimSuffix(fieldName, "_id")
			}

			sb.WriteString(fmt.Sprintf("\t\t\t<template #%s=\"{ row }\">\n", fieldName))
			sb.WriteString("\t\t\t\t<span\n")
			sb.WriteString("\t\t\t\t\tclass=\"text-blue-600 underline cursor-pointer hover:text-blue-800\"\n")
			sb.WriteString(fmt.Sprintf("\t\t\t\t\t@click=\"router.push({ path: '/%s/list', query: { id: row.%s } })\"\n",
				moduleName, fieldName))
			sb.WriteString("\t\t\t\t>\n")
			sb.WriteString(fmt.Sprintf("\t\t\t\t\t{{ row.%s }}\n", fieldName))
			sb.WriteString("\t\t\t\t</span>\n")
			sb.WriteString("\t\t\t</template>\n")
		}
	}

	// links 类型插槽（多参数跳转）
	for _, field := range f.config.Fields {
		if field.ShowInTable && field.TableDisplayType == "links" {
			// 计算跳转路径：如果是 _ids 结尾，取前缀作为模块名；否则用整个字段名
			fieldName := field.JsonTag
			moduleName := fieldName
			if strings.HasSuffix(fieldName, "_ids") {
				moduleName = strings.TrimSuffix(fieldName, "_ids")
			} else if strings.HasSuffix(fieldName, "_id") {
				moduleName = strings.TrimSuffix(fieldName, "_id")
			}

			sb.WriteString(fmt.Sprintf("\t\t\t<template #%s=\"{ row }\">\n", fieldName))
			sb.WriteString("\t\t\t\t<div class=\"flex flex-wrap gap-1\">\n")
			sb.WriteString(fmt.Sprintf("\t\t\t\t\t<span\n"))
			sb.WriteString(fmt.Sprintf("\t\t\t\t\t\tv-for=\"(id, index) in (row.%s || '').split(',')\"\n", fieldName))
			sb.WriteString("\t\t\t\t\t\t:key=\"index\"\n")
			sb.WriteString("\t\t\t\t\t\tclass=\"text-blue-600 underline cursor-pointer hover:text-blue-800\"\n")
			sb.WriteString(fmt.Sprintf("\t\t\t\t\t\t@click=\"router.push({ path: '/%s/list', query: { id: id.trim() } })\"\n",
				moduleName))
			sb.WriteString("\t\t\t\t\t>\n")
			sb.WriteString("\t\t\t\t\t\t{{ id.trim() }}\n")
			sb.WriteString("\t\t\t\t\t</span>\n")
			sb.WriteString("\t\t\t\t</div>\n")
			sb.WriteString("\t\t\t</template>\n")
		}
	}

	// 操作列插槽
	sb.WriteString("\t\t\t<template #operation=\"{ row }\">\n")
	sb.WriteString("\t\t\t\t<div class=\"flex gap-2 justify-center\">\n")
	sb.WriteString("\t\t\t\t\t<VbenButtonGroup border>\n")

	if f.config.Methods.Update {
		sb.WriteString("\t\t\t\t\t\t<VbenButton variant=\"icon\" @click=\"onEdit(row)\">\n")
		sb.WriteString("\t\t\t\t\t\t\t<IconifyIcon\n")
		sb.WriteString("\t\t\t\t\t\t\t\tclass=\"size-6 outline-none text-blue-600\"\n")
		sb.WriteString("\t\t\t\t\t\t\t\ticon=\"ant-design:edit-outlined\"\n")
		sb.WriteString("\t\t\t\t\t\t\t/>\n")
		sb.WriteString("\t\t\t\t\t\t</VbenButton>\n")
	}

	if f.config.Methods.Delete {
		sb.WriteString("\t\t\t\t\t\t<Popconfirm\n")
		sb.WriteString("\t\t\t\t\t\t\t:title=\"$t('ui.actionTitle.delete')\"\n")
		sb.WriteString("\t\t\t\t\t\t\t@confirm=\"onDelete(row)\"\n")
		sb.WriteString("\t\t\t\t\t\t>\n")
		sb.WriteString("\t\t\t\t\t\t\t<VbenButton variant=\"icon\">\n")
		sb.WriteString("\t\t\t\t\t\t\t\t<X class=\"size-6 text-red-600\" />\n")
		sb.WriteString("\t\t\t\t\t\t\t</VbenButton>\n")
		sb.WriteString("\t\t\t\t\t\t</Popconfirm>\n")
	}

	sb.WriteString("\t\t\t\t\t</VbenButtonGroup>\n")
	sb.WriteString("\t\t\t\t</div>\n")
	sb.WriteString("\t\t\t</template>\n")

	// toolbar-actions 插槽
	sb.WriteString("\t\t\t<template #toolbar-actions>\n")

	// 创建按钮
	if f.config.Methods.Create {
		sb.WriteString("\t\t\t\t<VbenButton size=\"sm\" @click=\"onCreate\" variant=\"outline\" class=\"mr-3\">\n")
		sb.WriteString("\t\t\t\t\t<Plus class=\"size-3\" />\n")
		sb.WriteString("\t\t\t\t\t{{ $t('common.create') }}\n")
		sb.WriteString("\t\t\t\t</VbenButton>\n")
	}

	// 批量操作按钮组
	sb.WriteString("\t\t\t\t<VbenButtonGroup v-show=\"selectedCount > 0\" border>\n")

	// 批量删除
	if f.config.Methods.Delete {
		sb.WriteString("\t\t\t\t\t<Popconfirm\n")
		sb.WriteString("\t\t\t\t\t\t:title=\"$t('common.confirmDeleteSelected', [selectedCount])\"\n")
		sb.WriteString("\t\t\t\t\t\t:ok-text=\"$t('common.confirm')\"\n")
		sb.WriteString("\t\t\t\t\t\t:cancel-text=\"$t('common.cancel')\"\n")
		sb.WriteString("\t\t\t\t\t\t@confirm=\"onDelete()\"\n")
		sb.WriteString("\t\t\t\t\t>\n")
		sb.WriteString("\t\t\t\t\t\t<VbenButton variant=\"outline\">\n")
		sb.WriteString("\t\t\t\t\t\t\t<Trash class=\"size-3\" />\n")
		sb.WriteString("\t\t\t\t\t\t\t{{ $t('common.delete') }}\n")
		sb.WriteString("\t\t\t\t\t\t</VbenButton>\n")
		sb.WriteString("\t\t\t\t\t</Popconfirm>\n")
	}

	// 批量启用/禁用
	if f.config.Methods.Operate {
		sb.WriteString("\t\t\t\t\t<Popconfirm\n")
		sb.WriteString("\t\t\t\t\t\t:title=\"$t('common.confirmEnableSelected', [selectedCount])\"\n")
		sb.WriteString("\t\t\t\t\t\t:ok-text=\"$t('common.confirm')\"\n")
		sb.WriteString("\t\t\t\t\t\t:cancel-text=\"$t('common.cancel')\"\n")
		sb.WriteString("\t\t\t\t\t\t@confirm=\"operateApi(1, undefined, 'status')\"\n")
		sb.WriteString("\t\t\t\t\t>\n")
		sb.WriteString("\t\t\t\t\t\t<VbenButton variant=\"outline\">\n")
		sb.WriteString("\t\t\t\t\t\t\t<Check class=\"size-3\" />\n")
		sb.WriteString("\t\t\t\t\t\t\t{{ $t('common.enable') }}\n")
		sb.WriteString("\t\t\t\t\t\t</VbenButton>\n")
		sb.WriteString("\t\t\t\t\t</Popconfirm>\n")
		sb.WriteString("\t\t\t\t\t<Popconfirm\n")
		sb.WriteString("\t\t\t\t\t\t:title=\"$t('common.confirmDisableSelected', [selectedCount])\"\n")
		sb.WriteString("\t\t\t\t\t\t:ok-text=\"$t('common.confirm')\"\n")
		sb.WriteString("\t\t\t\t\t\t:cancel-text=\"$t('common.cancel')\"\n")
		sb.WriteString("\t\t\t\t\t\t@confirm=\"operateApi(0, undefined, 'status')\"\n")
		sb.WriteString("\t\t\t\t\t>\n")
		sb.WriteString("\t\t\t\t\t\t<VbenButton variant=\"outline\">\n")
		sb.WriteString("\t\t\t\t\t\t\t<X class=\"size-3\" />\n")
		sb.WriteString("\t\t\t\t\t\t\t{{ $t('common.disable') }}\n")
		sb.WriteString("\t\t\t\t\t\t</VbenButton>\n")
		sb.WriteString("\t\t\t\t\t</Popconfirm>\n")
	}

	sb.WriteString("\t\t\t\t</VbenButtonGroup>\n")
	sb.WriteString("\t\t\t</template>\n")

	// toolbar-tools 插槽 - 搜索框
	sb.WriteString("\t\t\t<template #toolbar-tools>\n")
	sb.WriteString("\t\t\t\t<InputSearch @search=\"onSearch\" @clear=\"onClearSearch\" allow-clear :placeholder=\"$t('common.fuzzySearch')\" />\n")
	sb.WriteString("\t\t\t</template>\n")

	sb.WriteString("\t\t</Grid>\n")
	sb.WriteString("\t</Page>\n")
	sb.WriteString("</template>\n")
}

// GenerateListView 生成列表页面
func (f *FrontendGenerator) GenerateListView() string {
	var sb strings.Builder

	// Script 部分
	sb.WriteString("<script lang=\"ts\" setup>\n")
	sb.WriteString("import type {\n")
	sb.WriteString("\tVxeTableGridOptions,\n")
	sb.WriteString("} from '#/adapter/vxe-table';\n")
	sb.WriteString("import { onMounted, ref } from 'vue';\n")
	sb.WriteString("import { useRoute, useRouter } from 'vue-router';\n\n")

	sb.WriteString("import { VbenButton, VbenButtonGroup, Page, useVbenDrawer } from '@vben/common-ui';\n")
	sb.WriteString("import { Plus, Trash, Check, X, IconifyIcon } from '@vben/icons';\n")
	sb.WriteString("import { $t } from '@vben/locales';\n\n")

	sb.WriteString("import { message, Popconfirm, Tag, Image, InputSearch } from 'ant-design-vue';\n\n")

	sb.WriteString("import { useVbenVxeGrid } from '#/adapter/vxe-table';\n")
	sb.WriteString("import {\n")
	if f.config.Methods.Delete {
		sb.WriteString(fmt.Sprintf("\tdelete%s,\n", f.config.StructName))
	}
	sb.WriteString(fmt.Sprintf("\tget%sList,\n", f.config.StructName))
	if f.config.Methods.Operate {
		sb.WriteString(fmt.Sprintf("\toperate%s,\n", f.config.StructName))
	}
	sb.WriteString(fmt.Sprintf("} from '#/api/%s';\n\n", f.config.ModuleName))

	sb.WriteString(fmt.Sprintf("import type { %sApi } from '#/api/%s';\n\n",
		f.config.StructName, f.config.ModuleName))

	sb.WriteString("import { useColumns, useSearchFormSchema } from './data';\n")
	sb.WriteString("import FormDrawerComponent from './modules/form.vue';\n\n")

	// 变量定义
	sb.WriteString("const route = useRoute();\n")
	sb.WriteString("const router = useRouter();\n\n")

	sb.WriteString("// 选中的行数\n")
	sb.WriteString("const selectedCount = ref(0);\n\n")

	sb.WriteString("// 模糊搜索内容\n")
	sb.WriteString("const searchValue = ref<string>(\"\");\n\n")

	// 方法定义
	sb.WriteString("// 获取选中的行数据（公共方法）\n")
	sb.WriteString(fmt.Sprintf("const getSelectedRecords = (): %sApi.%s[] => {\n",
		f.config.StructName, f.config.StructName))
	sb.WriteString(fmt.Sprintf("\treturn (gridApi.grid?.getCheckboxRecords() as %sApi.%s[]) || [];\n",
		f.config.StructName, f.config.StructName))
	sb.WriteString("}\n\n")

	sb.WriteString("// 更新选中行数量\n")
	sb.WriteString("const updateSelectedCount = () => {\n")
	sb.WriteString("\trequestAnimationFrame(() => {\n")
	sb.WriteString("\t\tselectedCount.value = getSelectedRecords().length;\n")
	sb.WriteString("\t});\n")
	sb.WriteString("}\n\n")

	// 搜索方法
	sb.WriteString("// 清空搜索\n")
	sb.WriteString("const onClearSearch = () => {\n")
	sb.WriteString("\tsearchValue.value = \"\";\n")
	sb.WriteString("\tonRefresh();\n")
	sb.WriteString("}\n\n")

	sb.WriteString("// 搜索\n")
	sb.WriteString("const onSearch = (content: string) => {\n")
	sb.WriteString("\tsearchValue.value = content;\n")
	sb.WriteString("\tonRefresh();\n")
	sb.WriteString("}\n\n")

	// Form Drawer
	sb.WriteString("const [FormDrawer, formDrawerApi] = useVbenDrawer({\n")
	sb.WriteString("\tconnectedComponent: FormDrawerComponent,\n")
	sb.WriteString("\tdestroyOnClose: true,\n")
	sb.WriteString("});\n\n")

	// Grid 配置
	sb.WriteString("const [Grid, gridApi] = useVbenVxeGrid({\n")
	sb.WriteString("\tformOptions: {\n")
	sb.WriteString("\t\tcollapsed: true,\n")
	sb.WriteString("\t\tschema: useSearchFormSchema(),\n")
	sb.WriteString("\t\tshowCollapseButton: true,\n")
	sb.WriteString("\t\tsubmitOnChange: false,\n")
	sb.WriteString("\t\tsubmitOnEnter: false,\n")
	sb.WriteString("\t},\n")
	sb.WriteString("\tgridOptions: {\n")
	sb.WriteString("\t\tcolumns: useColumns(),\n")
	sb.WriteString("\t\tkeepSource: true,\n")
	sb.WriteString("\t\tpagerConfig: {\n")
	sb.WriteString("\t\t\tenabled: true,\n")
	sb.WriteString("\t\t},\n")
	sb.WriteString("\t\tproxyConfig: {\n")
	sb.WriteString("\t\t\tajax: {\n")
	sb.WriteString("\t\t\t\tquery: async ({ page, sort }, formValues) => {\n")
	sb.WriteString(fmt.Sprintf("\t\t\t\t\tconst params: %sApi.ListParams = {\n", f.config.StructName))
	sb.WriteString("\t\t\t\t\t\tpage: page.currentPage,\n")
	sb.WriteString("\t\t\t\t\t\tpage_size: page.pageSize,\n")

	// 默认排序字段
	defaultSortField := f.config.DefaultSortField
	if defaultSortField == "" {
		defaultSortField = "id"
	}
	sb.WriteString(fmt.Sprintf("\t\t\t\t\t\tsort_by: sort.field ? sort.field : '%s',\n", defaultSortField))

	// 默认排序类型
	defaultSortOrder := f.config.DefaultSortOrder
	if defaultSortOrder == "" {
		defaultSortOrder = "desc"
	}
	if defaultSortOrder == "desc" {
		sb.WriteString("\t\t\t\t\t\tsort_order: sort.order === 'asc' ? 'asc' : 'desc',\n")
	} else {
		sb.WriteString("\t\t\t\t\t\tsort_order: sort.order === 'desc' ? 'desc' : 'asc',\n")
	}

	sb.WriteString("\t\t\t\t\t\tfilter: JSON.stringify(formValues),\n")
	sb.WriteString("\t\t\t\t\t\tsearch: searchValue.value || undefined,\n")
	sb.WriteString("\t\t\t\t\t};\n")
	sb.WriteString(fmt.Sprintf("\t\t\t\t\treturn await get%sList(params);\n", f.config.StructName))
	sb.WriteString("\t\t\t\t},\n")
	sb.WriteString("\t\t\t},\n")
	sb.WriteString("\t\t\tsort: true,\n")
	sb.WriteString("\t\t},\n")
	sb.WriteString("\t\trowConfig: {\n")
	sb.WriteString("\t\t\tkeyField: 'id',\n")
	sb.WriteString("\t\t\tisHover: true,\n")
	sb.WriteString("\t\t},\n")
	sb.WriteString("\t\tsortConfig: {\n")
	sb.WriteString("\t\t\tremote: true,\n")
	sb.WriteString("\t\t\ttrigger: 'cell',\n")
	sb.WriteString("\t\t},\n")
	sb.WriteString("\t\tstripe: false,\n")
	sb.WriteString("\t\ttoolbarConfig: {\n")
	sb.WriteString("\t\t\tcustom: true,\n")
	sb.WriteString("\t\t\trefresh: true,\n")
	sb.WriteString("\t\t\tzoom: true,\n")
	sb.WriteString("\t\t},\n")
	sb.WriteString("\t\tcheckboxConfig: {\n")
	sb.WriteString("\t\t\ttrigger: 'cell',\n")
	sb.WriteString("\t\t\thighlight: true,\n")
	sb.WriteString("\t\t\treserve: false,\n")
	sb.WriteString("\t\t},\n")

	// 编辑配置
	hasEditableFields := false
	for _, field := range f.config.Fields {
		if field.TableDisplayType == "editable" {
			hasEditableFields = true
			break
		}
	}

	if hasEditableFields {
		sb.WriteString("\t\teditConfig: {\n")
		sb.WriteString("\t\t\ttrigger: 'dblclick',\n")
		sb.WriteString("\t\t\tmode: 'cell',\n")
		sb.WriteString("\t\t\tshowStatus: true,\n")
		sb.WriteString("\t\t},\n")
	}

	sb.WriteString("\t} as VxeTableGridOptions,\n")
	sb.WriteString("\tgridEvents: {\n")
	sb.WriteString("\t\tcheckboxChange: updateSelectedCount,\n")
	sb.WriteString("\t\tcheckboxAll: updateSelectedCount,\n")

	if hasEditableFields && f.config.Methods.Operate {
		sb.WriteString("\t\teditClosed: async ({ row, column }: { row: any; column: any }) => {\n")

		// 生成可编辑字段的判断
		var editableFields []string
		for _, field := range f.config.Fields {
			if field.TableDisplayType == "editable" {
				editableFields = append(editableFields, field.JsonTag)
			}
		}

		if len(editableFields) > 0 {
			sb.WriteString(fmt.Sprintf("\t\t\tif (%s) {\n",
				"column.field === '"+strings.Join(editableFields, "' || column.field === '")+"'"))
			sb.WriteString("\t\t\t\tconst newValue = row[column.field];\n")
			sb.WriteString("\t\t\t\tawait operateApi(newValue, row, column.field);\n")
			sb.WriteString("\t\t\t}\n")
		}

		sb.WriteString("\t\t},\n")
	}

	sb.WriteString("\t},\n")
	sb.WriteString("});\n\n")

	// Operate 方法
	if f.config.Methods.Operate {
		sb.WriteString("// 通用字段操作方法\n")
		sb.WriteString(fmt.Sprintf("const operateApi = async (value: number, row?: %sApi.%s, field: string = 'status') => {\n",
			f.config.StructName, f.config.StructName))
		sb.WriteString("\tlet ids: number[] = [];\n\n")
		sb.WriteString("\tif (row) {\n")
		sb.WriteString("\t\tids = [row.id!];\n")
		sb.WriteString("\t} else {\n")
		sb.WriteString("\t\tconst selectRecords = getSelectedRecords();\n")
		sb.WriteString("\t\tif (selectRecords.length === 0) {\n")
		sb.WriteString("\t\t\tmessage.warning($t('common.tableSelectTip'));\n")
		sb.WriteString("\t\t\treturn;\n")
		sb.WriteString("\t\t}\n")
		sb.WriteString("\t\tids = selectRecords.map((item) => item.id!);\n")
		sb.WriteString("\t}\n\n")
		sb.WriteString("\ttry {\n")
		sb.WriteString(fmt.Sprintf("\t\tawait operate%s({\n", f.config.StructName))
		sb.WriteString("\t\t\tids: ids,\n")
		sb.WriteString("\t\t\tfield,\n")
		sb.WriteString("\t\t\tvalue,\n")
		sb.WriteString("\t\t});\n")
		sb.WriteString("\t\tmessage.success($t('common.success'));\n")
		sb.WriteString("\t\tonRefresh();\n")
		sb.WriteString("\t\treturn true;\n")
		sb.WriteString("\t} catch (error) {\n")
		sb.WriteString("\t\tmessage.error('error');\n")
		sb.WriteString("\t\treturn false;\n")
		sb.WriteString("\t}\n")
		sb.WriteString("}\n\n")
	}

	// 其他方法
	sb.WriteString("const onRefresh = () => {\n")
	sb.WriteString("\tgridApi.grid?.clearCheckboxRow();\n")
	sb.WriteString("\tgridApi.query();\n")
	sb.WriteString("\tselectedCount.value = 0;\n")
	sb.WriteString("}\n\n")

	sb.WriteString(fmt.Sprintf("const onEdit = (row: %sApi.%s) => {\n", f.config.StructName, f.config.StructName))
	sb.WriteString("\tformDrawerApi.setData(row).open();\n")
	sb.WriteString("}\n\n")

	sb.WriteString("const onCreate = () => {\n")
	sb.WriteString("\tformDrawerApi.setData({}).open();\n")
	sb.WriteString("}\n\n")

	if f.config.Methods.Delete {
		sb.WriteString(fmt.Sprintf("const onDelete = (row?: %sApi.%s) => {\n", f.config.StructName, f.config.StructName))
		sb.WriteString("\tlet ids: number[] = [];\n\n")
		sb.WriteString("\tif (row) {\n")
		sb.WriteString("\t\tids = [row.id!];\n")
		sb.WriteString("\t} else {\n")
		sb.WriteString("\t\tconst selectRecords = getSelectedRecords();\n")
		sb.WriteString("\t\tif (selectRecords.length === 0) {\n")
		sb.WriteString("\t\t\tmessage.warning($t('common.tableSelectTip'));\n")
		sb.WriteString("\t\t\treturn;\n")
		sb.WriteString("\t\t}\n")
		sb.WriteString("\t\tids = selectRecords.map((item) => item.id!);\n")
		sb.WriteString("\t}\n\n")
		sb.WriteString("\tconst hideLoading = message.loading({\n")
		sb.WriteString("\t\tcontent: $t('ui.actionMessage.deleting'),\n")
		sb.WriteString("\t\tduration: 0,\n")
		sb.WriteString("\t\tkey: 'action_process_msg',\n")
		sb.WriteString("\t});\n")
		sb.WriteString(fmt.Sprintf("\tdelete%s({ ids: ids })\n", f.config.StructName))
		sb.WriteString("\t\t.then(() => {\n")
		sb.WriteString("\t\t\tmessage.success({\n")
		sb.WriteString("\t\t\t\tcontent: $t('ui.actionMessage.deleteSuccess'),\n")
		sb.WriteString("\t\t\t\tkey: 'action_process_msg',\n")
		sb.WriteString("\t\t\t});\n")
		sb.WriteString("\t\t\tonRefresh();\n")
		sb.WriteString("\t\t})\n")
		sb.WriteString("\t\t.catch(() => {\n")
		sb.WriteString("\t\t\thideLoading();\n")
		sb.WriteString("\t\t});\n")
		sb.WriteString("}\n")
	}

	// 生成 onMounted 逻辑：从 URL 参数填充搜索表单
	sb.WriteString("\n// 获取搜索表单的字段名列表\n")
	sb.WriteString("const searchFormSchema = useSearchFormSchema();\n")
	sb.WriteString("const searchFormFields = searchFormSchema\n")
	sb.WriteString("\t? searchFormSchema.map((item) => item.fieldName).filter((name): name is string => !!name)\n")
	sb.WriteString("\t: [];\n\n")

	sb.WriteString("// 页面加载时从 URL 参数填充搜索表单\n")
	sb.WriteString("onMounted(() => {\n")
	sb.WriteString("\tconst query = route.query;\n")
	sb.WriteString("\tif (Object.keys(query).length > 0) {\n")
	sb.WriteString("\t\tconst formValues: Record<string, any> = {};\n")
	sb.WriteString("\t\tfor (const [key, value] of Object.entries(query)) {\n")
	sb.WriteString("\t\t\t// 只填充存在于搜索表单中的字段\n")
	sb.WriteString("\t\t\tif (searchFormFields.includes(key) && value) {\n")
	sb.WriteString("\t\t\t\tformValues[key] = value;\n")
	sb.WriteString("\t\t\t}\n")
	sb.WriteString("\t\t}\n")
	sb.WriteString("\t\tif (Object.keys(formValues).length > 0) {\n")
	sb.WriteString("\t\t\t// 设置表单值并触发搜索\n")
	sb.WriteString("\t\t\tgridApi.formApi.setValues(formValues);\n")
	sb.WriteString("\t\t}\n")
	sb.WriteString("\t}\n")
	sb.WriteString("});\n")

	sb.WriteString("</script>\n\n")

	// 继续生成 template 部分...
	f.generateListViewTemplate(&sb)

	return sb.String()
}

// GenerateDataTS 生成 data.ts
func (f *FrontendGenerator) GenerateDataTS() string {
	var sb strings.Builder

	// 导入
	sb.WriteString("import type { VxeTableGridOptions } from '#/adapter/vxe-table';\n")
	sb.WriteString(fmt.Sprintf("import type { %sApi } from '#/api/%s';\n", f.config.StructName, f.config.ModuleName))
	sb.WriteString("import type { VbenFormProps } from '#/adapter/form';\n")
	sb.WriteString("import { $t } from '#/locales';\n")

	// 状态选项函数
	hasStatusField := false
	for _, field := range f.config.Fields {
		if field.IsOperateField && field.ShowInTable {
			hasStatusField = true
			break
		}
	}
	if hasStatusField {
		sb.WriteString(fmt.Sprintf("export function get%sStatusOptions() {\n", f.config.StructName))
		sb.WriteString("\treturn [\n")
		sb.WriteString("\t\t{ color: 'success', label: $t('common.enable'), value: 1 },\n")
		sb.WriteString("\t\t{ color: 'error', label: $t('common.disable'), value: 0 },\n")
		sb.WriteString("\t];\n")
		sb.WriteString("}\n\n")
	}

	// useSearchFormSchema 函数 - 放在 useColumns 之前
	sb.WriteString("export function useSearchFormSchema(): VbenFormProps['schema'] {\n")
	sb.WriteString("\treturn [\n")

	// 生成搜索表单配置
	hasSearchFields := false
	for _, field := range f.config.Fields {
		if !field.TableSearchable {
			continue
		}
		hasSearchFields = true

		sb.WriteString("\t\t{\n")
		sb.WriteString(fmt.Sprintf("\t\t\tcomponent: '%s',\n", f.getSearchComponent(field)))
		sb.WriteString(fmt.Sprintf("\t\t\tfieldName: '%s',\n", field.JsonTag))

		// 标题 - 使用国际化
		i18nKey := fmt.Sprintf("%s.%s", strings.ToLower(f.config.StructName), field.JsonTag)
		sb.WriteString(fmt.Sprintf("\t\t\tlabel: $t('%s'),\n", i18nKey))

		// 组件属性
		searchComponent := f.getSearchComponent(field)
		if searchComponent == "RangePicker" {
			sb.WriteString("\t\t\tcomponentProps: {\n")
			sb.WriteString("\t\t\t\tshowTime: true,\n")
			sb.WriteString("\t\t\t\tformat: 'YYYY-MM-DD HH:mm:ss',\n")
			sb.WriteString("\t\t\t\ttimestampFormat: true,\n")
			sb.WriteString("\t\t\t},\n")
		} else if searchComponent == "Switch" {
			// Switch 组件需要覆盖 w-full 样式并添加 allowClear
			sb.WriteString("\t\t\tcomponentProps: {\n")
			sb.WriteString("\t\t\t\tclass: '',\n")
			sb.WriteString("\t\t\t\tallowClear: true,\n")
			sb.WriteString("\t\t\t},\n")
		} else if searchComponent == "TableSelect" || searchComponent == "TableSelects" {
			// TableSelect 使用 config 属性配置
			multiple := "false"
			if searchComponent == "TableSelects" {
				multiple = "true"
			}
			// 从 FormComponentProps 读取配置
			api := fmt.Sprintf("/api/%s/list", strings.ToLower(field.ColumnName))
			labelField := "name"
			valueField := "id"
			if field.FormComponentProps != nil {
				if configData, ok := field.FormComponentProps["config"].(map[string]interface{}); ok {
					if v, ok := configData["api"].(string); ok && v != "" {
						api = v
					}
					if v, ok := configData["labelField"].(string); ok && v != "" {
						labelField = v
					}
					if v, ok := configData["valueField"].(string); ok && v != "" {
						valueField = v
					}
				}
			}
			sb.WriteString("\t\t\tcomponentProps: {\n")
			sb.WriteString(fmt.Sprintf("\t\t\t\tconfig: { api: '%s', labelField: '%s', valueField: '%s' },\n", api, labelField, valueField))
			sb.WriteString(fmt.Sprintf("\t\t\t\tmultiple: %s,\n", multiple))
			sb.WriteString("\t\t\t},\n")
		} else if (field.IsOperateField || field.IsBoolField) && searchComponent == "Select" {
			// status, is_xxx 等布尔类字段在搜索时使用 Select 并配置 options
			sb.WriteString("\t\t\tcomponentProps: {\n")
			sb.WriteString("\t\t\t\tallowClear: true,\n")
			sb.WriteString("\t\t\t\toptions: [\n")
			sb.WriteString("\t\t\t\t\t{ label: $t('common.all'), value: '' },\n")
			sb.WriteString("\t\t\t\t\t{ label: $t('common.yes'), value: '1' },\n")
			sb.WriteString("\t\t\t\t\t{ label: $t('common.no'), value: '0' },\n")
			sb.WriteString("\t\t\t\t],\n")
			sb.WriteString("\t\t\t},\n")
		} else if field.ComponentProps != "" {
			sb.WriteString(fmt.Sprintf("\t\t\tcomponentProps: %s,\n", field.ComponentProps))
		}

		sb.WriteString("\t\t},\n")
	}

	if !hasSearchFields {
		sb.WriteString("\t\t// 暂无搜索字段\n")
	}

	sb.WriteString("\t];\n")
	sb.WriteString("}\n\n")

	// useColumns 函数
	sb.WriteString(fmt.Sprintf("export function useColumns(): VxeTableGridOptions<%sApi.%s>['columns'] {\n",
		f.config.StructName, f.config.StructName))
	sb.WriteString("\treturn [\n")
	sb.WriteString("\t\t{\n")
	sb.WriteString("\t\t\ttype: 'checkbox',\n")
	sb.WriteString("\t\t\twidth: 50,\n")
	sb.WriteString("\t\t\talign: 'center',\n")
	sb.WriteString("\t\t\tfixed: 'left',\n")
	sb.WriteString("\t\t},\n")

	// 生成列配置
	for _, field := range f.config.Fields {
		if !field.ShowInTable {
			continue
		}

		sb.WriteString("\t\t{\n")

		// 对齐方式
		if field.IsPrimaryKey || field.IsTimeField || field.IsOperateField {
			sb.WriteString("\t\t\talign: 'center',\n")
		} else if field.FieldType == "float32" || field.FieldType == "float64" {
			sb.WriteString("\t\t\talign: 'right',\n")
		} else {
			sb.WriteString("\t\t\talign: 'left',\n")
		}

		// 字段名
		sb.WriteString(fmt.Sprintf("\t\t\tfield: '%s',\n", field.JsonTag))

		// 图片字段使用 slots
		if field.IsImageField {
			sb.WriteString(fmt.Sprintf("\t\t\tslots: { default: '%s' },\n", field.JsonTag))
		}

		// 标题 - 使用国际化
		i18nKey := fmt.Sprintf("%s.%s", strings.ToLower(f.config.StructName), field.JsonTag)
		sb.WriteString(fmt.Sprintf("\t\t\ttitle: $t('%s'),\n", i18nKey))

		// 宽度
		if field.IsPrimaryKey {
			sb.WriteString("\t\t\twidth: 80,\n")
		} else if field.IsImageField {
			sb.WriteString("\t\t\twidth: 50,\n")
		} else if field.IsTimeField {
			sb.WriteString("\t\t\twidth: 160,\n")
		} else if field.IsOperateField {
			sb.WriteString("\t\t\twidth: 80,\n")
		}

		// 排序
		if field.TableSortable {
			sb.WriteString("\t\t\tsortable: true,\n")
		}

		// 可编辑
		if field.TableDisplayType == "editable" {
			sb.WriteString("\t\t\teditRender: {\n")
			sb.WriteString("\t\t\t\tname: 'input',\n")
			sb.WriteString("\t\t\t\tattrs: {\n")
			sb.WriteString("\t\t\t\t\ttype: 'number',\n")
			sb.WriteString("\t\t\t\t\tplaceholder: '双击编辑',\n")
			sb.WriteString("\t\t\t\t},\n")
			sb.WriteString("\t\t\t},\n")
		}

		// tag 类型字段使用 slots
		if field.TableDisplayType == "tag" {
			sb.WriteString(fmt.Sprintf("\t\t\tslots: { default: '%s' },\n", field.JsonTag))
		}

		// link/links 类型字段使用 slots
		if field.TableDisplayType == "link" || field.TableDisplayType == "links" {
			sb.WriteString(fmt.Sprintf("\t\t\tslots: { default: '%s' },\n", field.JsonTag))
		}

		// 时间字段使用 formatter
		if field.IsTimeField {
			sb.WriteString("\t\t\tformatter: ({ cellValue }) => {\n")
			sb.WriteString("\t\t\t\treturn new Date(cellValue * 1000).toLocaleString();\n")
			sb.WriteString("\t\t\t},\n")
		}

		sb.WriteString("\t\t},\n")
	}

	// 操作列
	sb.WriteString("\t\t{\n")
	sb.WriteString("\t\t\tfield: 'operation',\n")
	sb.WriteString("\t\t\tfixed: 'right',\n")
	sb.WriteString("\t\t\tslots: { default: 'operation' },\n")
	sb.WriteString("\t\t\ttitle: $t('common.operation'),\n")
	sb.WriteString("\t\t\twidth: 100,\n")
	sb.WriteString("\t\t}\n")

	sb.WriteString("\t];\n")
	sb.WriteString("}\n")

	return sb.String()
}

// getSearchComponent 获取搜索组件类型
func (f *FrontendGenerator) getSearchComponent(field FieldConfig) string {
	if field.SearchFormType != "" {
		return field.SearchFormType
	}
	// 时间字段默认使用 RangePicker
	if field.IsTimeField || strings.HasSuffix(field.ColumnName, "_at") || strings.HasSuffix(field.ColumnName, "_time") {
		return "RangePicker"
	}
	return "Input"
}

// GenerateFormVue 生成表单组件
func (f *FrontendGenerator) GenerateFormVue() string {
	var sb strings.Builder

	// Script 部分
	sb.WriteString("<script lang=\"ts\" setup>\n")
	sb.WriteString("import { computed, nextTick, reactive, ref } from 'vue';\n\n")
	sb.WriteString("import { useVbenDrawer, useVbenForm } from '@vben/common-ui';\n")
	sb.WriteString("import { $t } from '@vben/locales';\n\n")
	sb.WriteString("import { message } from 'ant-design-vue';\n")
	sb.WriteString("import { useBreakpoints } from '@vueuse/core';\n\n")
	sb.WriteString(fmt.Sprintf("import { create%s, update%s } from '#/api/%s';\n",
		f.config.StructName, f.config.StructName, f.config.ModuleName))
	sb.WriteString(fmt.Sprintf("import type { %sApi } from '#/api/%s';\n\n",
		f.config.StructName, f.config.ModuleName))

	// Emits
	sb.WriteString("const emit = defineEmits<{\n")
	sb.WriteString("\tsuccess: [];\n")
	sb.WriteString("}>();\n\n")

	// 响应式数据
	sb.WriteString(fmt.Sprintf("const rowData = ref<%sApi.%s>();\n", f.config.StructName, f.config.StructName))
	sb.WriteString("const loading = ref(false);\n")
	sb.WriteString("const isEdit = computed(() => !!rowData.value?.id);\n\n")

	// 响应式布局
	sb.WriteString("const breakpoints = useBreakpoints({ md: 768 });\n")
	sb.WriteString("const isHorizontal = breakpoints.greater('md');\n\n")

	// Form Schema - 使用 computed
	sb.WriteString("const formSchema = computed(() => [\n")

	// 生成表单字段
	for _, field := range f.config.Fields {
		if !field.ShowInForm || field.IsPrimaryKey || field.IsAutoIncrement {
			continue
		}

		// 标题 - 使用国际化
		i18nKey := fmt.Sprintf("%s.%s", strings.ToLower(f.config.StructName), field.JsonTag)

		sb.WriteString("\t{\n")
		sb.WriteString(fmt.Sprintf("\t\tfieldName: '%s',\n", field.JsonTag))
		sb.WriteString(fmt.Sprintf("\t\tlabel: $t('%s'),\n", i18nKey))

		// 组件类型
		component := field.FormComponent
		if component == "" {
			component = "Input"
		}
		sb.WriteString(fmt.Sprintf("\t\tcomponent: '%s',\n", component))

		// 必填
		if field.IsRequired {
			sb.WriteString("\t\trules: 'required',\n")
		}

		// 组件属性
		if field.ComponentProps != "" {
			// 使用用户配置的 componentProps
			sb.WriteString(fmt.Sprintf("\t\tcomponentProps: %s,\n", field.ComponentProps))
		} else if field.IsTextField && component == "Textarea" {
			sb.WriteString("\t\tcomponentProps: {\n")
			sb.WriteString("\t\t\trows: 4,\n")
			sb.WriteString("\t\t},\n")
		} else if component == "TableSelect" || component == "TableSelects" {
			// TableSelect 使用 config 属性配置
			multiple := "false"
			if component == "TableSelects" {
				multiple = "true"
			}
			// 从 FormComponentProps 读取配置
			api := fmt.Sprintf("/api/%s/list", strings.ToLower(field.ColumnName))
			labelField := "name"
			valueField := "id"
			if field.FormComponentProps != nil {
				if configData, ok := field.FormComponentProps["config"].(map[string]interface{}); ok {
					if v, ok := configData["api"].(string); ok && v != "" {
						api = v
					}
					if v, ok := configData["labelField"].(string); ok && v != "" {
						labelField = v
					}
					if v, ok := configData["valueField"].(string); ok && v != "" {
						valueField = v
					}
				}
			}
			sb.WriteString("\t\tcomponentProps: {\n")
			sb.WriteString(fmt.Sprintf("\t\t\tconfig: { api: '%s', labelField: '%s', valueField: '%s' },\n", api, labelField, valueField))
			sb.WriteString(fmt.Sprintf("\t\t\tmultiple: %s,\n", multiple))
			sb.WriteString("\t\t},\n")
		} else if component == "Switch" {
			// Switch 组件需要覆盖 w-full 样式
			sb.WriteString("\t\tcomponentProps: {\n")
			sb.WriteString("\t\t\tclass: '',\n")
			sb.WriteString("\t\t},\n")
		}

		sb.WriteString("\t},\n")
	}

	sb.WriteString("]);\n\n")

	// Form - 使用 reactive 包装
	sb.WriteString("const [Form, formApi] = useVbenForm(\n")
	sb.WriteString("\treactive({\n")
	sb.WriteString("\t\tcommonConfig: {\n")
	sb.WriteString("\t\t\tcomponentProps: {\n")
	sb.WriteString("\t\t\t\tclass: 'w-full',\n")
	sb.WriteString("\t\t\t},\n")
	sb.WriteString("\t\t},\n")
	sb.WriteString("\t\tschema: formSchema,\n")
	sb.WriteString("\t\tshowDefaultActions: false,\n")
	sb.WriteString("\t}),\n")
	sb.WriteString(");\n\n")

	// 提交方法
	sb.WriteString("const handleSubmit = async () => {\n")
	sb.WriteString("\tconst { valid } = await formApi.validate();\n")
	sb.WriteString("\tif (!valid) return;\n\n")
	sb.WriteString("\tconst values = await formApi.getValues();\n")
	sb.WriteString("\tloading.value = true;\n\n")
	sb.WriteString("\ttry {\n")
	sb.WriteString("\t\tif (isEdit.value && rowData.value?.id) {\n")
	sb.WriteString(fmt.Sprintf("\t\t\tawait update%s({ ...values, id: rowData.value.id });\n", f.config.StructName))
	sb.WriteString("\t\t} else {\n")
	sb.WriteString(fmt.Sprintf("\t\t\tawait create%s(values);\n", f.config.StructName))
	sb.WriteString("\t\t}\n")
	sb.WriteString("\t\tmessage.success($t('ui.actionMessage.operationSuccess'));\n")
	sb.WriteString("\t\tdrawerApi.close();\n")
	sb.WriteString("\t\temit('success');\n")
	sb.WriteString("\t} finally {\n")
	sb.WriteString("\t\tloading.value = false;\n")
	sb.WriteString("\t}\n")
	sb.WriteString("};\n\n")

	// Drawer - 使用正确的模式
	sb.WriteString("const [Drawer, drawerApi] = useVbenDrawer({\n")
	sb.WriteString("\tonConfirm: () => handleSubmit(),\n")
	sb.WriteString("\tasync onOpenChange(isOpen: boolean) {\n")
	sb.WriteString("\t\tif (isOpen) {\n")
	sb.WriteString(fmt.Sprintf("\t\t\tconst data = drawerApi.getData<%sApi.%s>();\n",
		f.config.StructName, f.config.StructName))
	sb.WriteString("\t\t\trowData.value = data;\n")
	sb.WriteString("\t\t\tdrawerApi.setState({\n")
	sb.WriteString("\t\t\t\ttitle: isEdit.value\n")
	sb.WriteString(fmt.Sprintf("\t\t\t\t\t? $t('common.edit') + ' ' + $t('%s.title')\n",
		strings.ToLower(f.config.StructName)))
	sb.WriteString(fmt.Sprintf("\t\t\t\t\t: $t('common.create') + ' ' + $t('%s.title'),\n",
		strings.ToLower(f.config.StructName)))
	sb.WriteString("\t\t\t});\n")
	sb.WriteString("\t\t\tawait nextTick();\n")
	sb.WriteString("\t\t\tif (isEdit.value && data) {\n")
	sb.WriteString("\t\t\t\tformApi.setValues({\n")

	// 设置表单值
	for _, field := range f.config.Fields {
		if !field.ShowInForm || field.IsPrimaryKey || field.IsAutoIncrement {
			continue
		}
		sb.WriteString(fmt.Sprintf("\t\t\t\t\t%s: data.%s,\n", field.JsonTag, field.JsonTag))
	}

	sb.WriteString("\t\t\t\t});\n")
	sb.WriteString("\t\t\t} else {\n")
	sb.WriteString("\t\t\t\tformApi.resetForm();\n")
	sb.WriteString("\t\t\t}\n")
	sb.WriteString("\t\t}\n")
	sb.WriteString("\t},\n")
	sb.WriteString("});\n")
	sb.WriteString("</script>\n\n")

	// Template 部分
	sb.WriteString("<template>\n")
	sb.WriteString("\t<Drawer class=\"w-[600px]\" :confirm-loading=\"loading\" :loading=\"loading\">\n")
	sb.WriteString("\t\t<Form class=\"mx-4\" :layout=\"isHorizontal ? 'horizontal' : 'vertical'\" />\n")
	sb.WriteString("\t</Drawer>\n")
	sb.WriteString("</template>\n")

	return sb.String()
}

// GenerateLocaleZhCN 生成中文语言包
func (f *FrontendGenerator) GenerateLocaleZhCN() string {
	var sb strings.Builder

	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("  \"title\": \"%s\",\n", f.config.TableComment))
	sb.WriteString(fmt.Sprintf("  \"name\": \"%s\",\n", f.config.TableComment))

	// 生成字段翻译
	for i, field := range f.config.Fields {
		// 使用字段注释作为中文翻译
		comment := field.ColumnComment
		if comment == "" {
			comment = field.JsonTag
		}
		// id 字段特殊处理为 ID
		if strings.ToLower(field.JsonTag) == "id" && (comment == "id" || comment == "ID" || comment == "") {
			comment = "ID"
		}
		// 移除末尾的逗号（最后一个字段）
		if i == len(f.config.Fields)-1 {
			sb.WriteString(fmt.Sprintf("  \"%s\": \"%s\"\n", field.JsonTag, comment))
		} else {
			sb.WriteString(fmt.Sprintf("  \"%s\": \"%s\",\n", field.JsonTag, comment))
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

// GenerateLocaleEnUS 生成英文语言包
func (f *FrontendGenerator) GenerateLocaleEnUS() string {
	var sb strings.Builder

	sb.WriteString("{\n")
	// 使用表名作为英文标题
	sb.WriteString(fmt.Sprintf("  \"title\": \"%s\",\n", f.config.StructName))
	sb.WriteString(fmt.Sprintf("  \"name\": \"%s\",\n", f.config.StructName))

	// 生成字段翻译（英文使用字段名的驼峰转换）
	for i, field := range f.config.Fields {
		// 英文使用字段名的首字母大写形式
		englishName := toTitleCase(field.JsonTag)
		// 移除末尾的逗号（最后一个字段）
		if i == len(f.config.Fields)-1 {
			sb.WriteString(fmt.Sprintf("  \"%s\": \"%s\"\n", field.JsonTag, englishName))
		} else {
			sb.WriteString(fmt.Sprintf("  \"%s\": \"%s\",\n", field.JsonTag, englishName))
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

// toTitleCase 将 snake_case 转换为 Title Case
func toTitleCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		if len(word) > 0 {
			// id 特殊处理为 ID
			if strings.ToLower(word) == "id" {
				words[i] = "ID"
			} else {
				words[i] = strings.ToUpper(string(word[0])) + word[1:]
			}
		}
	}
	return strings.Join(words, " ")
}
