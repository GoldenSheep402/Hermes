package service

import categoryV1 "github.com/GoldenSheep402/Hermes/pkg/proto/category/v1"

type metaTemplatePresetTemplate struct {
	Key          string
	Label        string
	Type         string
	Required     bool
	Options      string
	SortOrder    int32
	DefaultValue string
}

type metaTemplatePresetDef struct {
	Value       string
	Label       string
	Description string
	Templates   []metaTemplatePresetTemplate
}

var commonMetaTemplatePresetDefs = []metaTemplatePresetTemplate{
	{Key: "language", Label: "语言", Type: "text", SortOrder: 10},
	{Key: "region", Label: "地区", Type: "text", SortOrder: 20},
	{Key: "year", Label: "年份", Type: "number", SortOrder: 30},
	{Key: "source", Label: "来源", Type: "text", SortOrder: 40},
}

var metaTemplatePresetDefs = []metaTemplatePresetDef{
	{
		Value:       "movie",
		Label:       "电影",
		Description: "电影模板字段",
		Templates: appendPresetTemplates(commonMetaTemplatePresetDefs,
			metaTemplatePresetTemplate{Key: "resolution", Label: "分辨率", Type: "select", Options: "2160p,1080p,720p,480p", SortOrder: 110},
			metaTemplatePresetTemplate{Key: "video_codec", Label: "视频编码", Type: "select", Options: "H.265,H.264,AV1,VP9,Other", SortOrder: 120},
			metaTemplatePresetTemplate{Key: "audio_codec", Label: "音频编码", Type: "text", SortOrder: 130},
			metaTemplatePresetTemplate{Key: "subtitle", Label: "字幕", Type: "text", SortOrder: 140},
			metaTemplatePresetTemplate{Key: "imdb_id", Label: "IMDb ID", Type: "text", SortOrder: 150},
			metaTemplatePresetTemplate{Key: "douban_id", Label: "豆瓣 ID", Type: "text", SortOrder: 160},
		),
	},
	{
		Value:       "anime",
		Label:       "动漫",
		Description: "动漫模板字段",
		Templates: appendPresetTemplates(commonMetaTemplatePresetDefs,
			metaTemplatePresetTemplate{Key: "season", Label: "季度", Type: "number", SortOrder: 210},
			metaTemplatePresetTemplate{Key: "episode", Label: "集数", Type: "text", SortOrder: 220},
			metaTemplatePresetTemplate{Key: "resolution", Label: "分辨率", Type: "select", Options: "2160p,1080p,720p", SortOrder: 230},
			metaTemplatePresetTemplate{Key: "fansub_group", Label: "字幕组", Type: "text", SortOrder: 240},
			metaTemplatePresetTemplate{Key: "audio_lang", Label: "配音语言", Type: "text", SortOrder: 250},
		),
	},
	{
		Value:       "tv",
		Label:       "电视剧",
		Description: "电视剧模板字段",
		Templates: appendPresetTemplates(commonMetaTemplatePresetDefs,
			metaTemplatePresetTemplate{Key: "season", Label: "季", Type: "number", SortOrder: 310},
			metaTemplatePresetTemplate{Key: "episode", Label: "集", Type: "text", SortOrder: 320},
			metaTemplatePresetTemplate{Key: "total_episodes", Label: "总集数", Type: "number", SortOrder: 330},
			metaTemplatePresetTemplate{Key: "resolution", Label: "分辨率", Type: "select", Options: "2160p,1080p,720p", SortOrder: 340},
			metaTemplatePresetTemplate{Key: "subtitle", Label: "字幕", Type: "text", SortOrder: 350},
		),
	},
	{
		Value:       "documentary",
		Label:       "纪录片",
		Description: "纪录片模板字段",
		Templates: appendPresetTemplates(commonMetaTemplatePresetDefs,
			metaTemplatePresetTemplate{Key: "topic", Label: "主题", Type: "text", SortOrder: 410},
			metaTemplatePresetTemplate{Key: "duration", Label: "时长", Type: "text", SortOrder: 420},
			metaTemplatePresetTemplate{Key: "resolution", Label: "分辨率", Type: "select", Options: "2160p,1080p,720p", SortOrder: 430},
			metaTemplatePresetTemplate{Key: "imdb_id", Label: "IMDb ID", Type: "text", SortOrder: 440},
			metaTemplatePresetTemplate{Key: "douban_id", Label: "豆瓣 ID", Type: "text", SortOrder: 450},
		),
	},
	{
		Value:       "software",
		Label:       "软件",
		Description: "软件模板字段",
		Templates: appendPresetTemplates(commonMetaTemplatePresetDefs,
			metaTemplatePresetTemplate{Key: "os", Label: "操作系统", Type: "select", Options: "Windows,macOS,Linux,Cross-platform,Other", SortOrder: 510},
			metaTemplatePresetTemplate{Key: "version", Label: "版本", Type: "text", Required: true, SortOrder: 520},
			metaTemplatePresetTemplate{Key: "license", Label: "授权方式", Type: "select", Options: "Open-source,Freeware,Commercial,Trial,Other", SortOrder: 530},
			metaTemplatePresetTemplate{Key: "official_website", Label: "官网", Type: "url", SortOrder: 540},
		),
	},
	{
		Value:       "ebook",
		Label:       "电子书",
		Description: "电子书模板字段",
		Templates: appendPresetTemplates(commonMetaTemplatePresetDefs,
			metaTemplatePresetTemplate{Key: "author", Label: "作者", Type: "text", Required: true, SortOrder: 610},
			metaTemplatePresetTemplate{Key: "publisher", Label: "出版社", Type: "text", SortOrder: 620},
			metaTemplatePresetTemplate{Key: "isbn", Label: "ISBN", Type: "text", SortOrder: 630},
			metaTemplatePresetTemplate{Key: "book_format", Label: "书籍格式", Type: "select", Options: "PDF,EPUB,MOBI,AZW3,DJVU,Other", SortOrder: 640},
			metaTemplatePresetTemplate{Key: "pages", Label: "页数", Type: "number", SortOrder: 650},
		),
	},
}

func appendPresetTemplates(base []metaTemplatePresetTemplate, extras ...metaTemplatePresetTemplate) []metaTemplatePresetTemplate {
	result := make([]metaTemplatePresetTemplate, 0, len(base)+len(extras))
	result = append(result, base...)
	result = append(result, extras...)
	return result
}

func listMetaTemplatePresetProtos() []*categoryV1.MetaTemplatePreset {
	result := make([]*categoryV1.MetaTemplatePreset, 0, len(metaTemplatePresetDefs))
	for _, preset := range metaTemplatePresetDefs {
		items := make([]*categoryV1.CategoryMetaTemplate, 0, len(preset.Templates))
		for _, tpl := range preset.Templates {
			items = append(items, &categoryV1.CategoryMetaTemplate{
				Key:          tpl.Key,
				Label:        tpl.Label,
				Type:         tpl.Type,
				Required:     tpl.Required,
				Options:      tpl.Options,
				SortOrder:    tpl.SortOrder,
				DefaultValue: tpl.DefaultValue,
			})
		}

		result = append(result, &categoryV1.MetaTemplatePreset{
			Value:       preset.Value,
			Label:       preset.Label,
			Description: preset.Description,
			Templates:   items,
		})
	}
	return result
}
