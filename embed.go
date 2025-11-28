package main

import "embed"

//go:embed configs/config.yaml
var ConfigFS embed.FS

//go:embed configs/i18n/*.yaml
var I18nFS embed.FS

//go:embed app/views/*.html
var ViewsFS embed.FS

