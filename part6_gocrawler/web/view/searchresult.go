package view

import (
	"crawler/web/model"
	"html/template"
	"io"
)

type SearchResultView struct {
	template *template.Template
}

// CreateSearchResultView 配置模板
func CreateSearchResultView(filename string) SearchResultView {
	return SearchResultView{
		template: template.Must(template.ParseFiles(filename)),
	}
}

// Render 渲染模板
func (v *SearchResultView) Render(w io.Writer, data model.SearchResult) error {
	return v.template.Execute(w, data)
}