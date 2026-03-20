package skills

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/openclaw-coding/skill-seed/embedfs"
)

// Loader Skills 模板加载器
type Loader struct {
	templates map[string]*template.Template
}

// NewLoader 创建 Skills 模板加载器
func NewLoader() *Loader {
	return &Loader{
		templates: make(map[string]*template.Template),
	}
}

// Load 加载指定名称的 Skills 模板
func (l *Loader) Load(name string) (*template.Template, error) {
	if cached, ok := l.templates[name]; ok {
		return cached, nil
	}

	// 从内嵌的文件系统加载模板
	// 例如: templates/skills/SKILL.md.tmpl
	path := "templates/skills/" + strings.ToUpper(name) + ".md.tmpl"
	data, err := embedfs.FS.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(name).Parse(string(data))
	if err != nil {
		return nil, err
	}

	l.templates[name] = tmpl
	return tmpl, nil
}

// LoadReference 加载 references 目录下的模板
func (l *Loader) LoadReference(category, name string) (*template.Template, error) {
	key := category + "/" + name

	if cached, ok := l.templates[key]; ok {
		return cached, nil
	}

	// 例如: templates/skills/references/naming-patterns/file-naming.md.tmpl
	path := "templates/skills/references/" + category + "-patterns/" + name + ".md.tmpl"
	data, err := embedfs.FS.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(key).Parse(string(data))
	if err != nil {
		return nil, err
	}

	l.templates[key] = tmpl
	return tmpl, nil
}

// Render 渲染指定名称的模板
func (l *Loader) Render(name string, data interface{}) (string, error) {
	tmpl, err := l.Load(name)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// RenderReference 渲染 references 模板
func (l *Loader) RenderReference(category, name string, data interface{}) (string, error) {
	tmpl, err := l.LoadReference(category, name)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Clear 清除缓存
func (l *Loader) Clear() {
	l.templates = make(map[string]*template.Template)
}
