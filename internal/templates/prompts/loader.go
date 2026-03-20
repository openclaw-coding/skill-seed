package prompts

import (
	"bytes"
	"text/template"

	"github.com/openclaw-coding/skill-seed/embedfs"
)

// Loader 提示词模板加载器
type Loader struct {
	agentName string
	templates map[string]*template.Template
}

// NewLoader 创建提示词加载器
func NewLoader(agentName string) *Loader {
	return &Loader{
		agentName: agentName,
		templates: make(map[string]*template.Template),
	}
}

// Load 加载指定名称的提示词模板
func (l *Loader) Load(name string) error {
	// 从内嵌的文件系统加载模板
	// 例如: templates/prompts/claude/analyze.txt.tmpl
	path := "templates/prompts/" + l.agentName + "/" + name + ".txt.tmpl"
	data, err := embedfs.FS.ReadFile(path)
	if err != nil {
		return err
	}

	tmpl, err := template.New(name).Parse(string(data))
	if err != nil {
		return err
	}

	l.templates[name] = tmpl
	return nil
}

// Render 渲染指定名称的提示词模板
func (l *Loader) Render(name string, data interface{}) string {
	// 如果模板未加载，先加载
	if _, ok := l.templates[name]; !ok {
		if err := l.Load(name); err != nil {
			return ""
		}
	}

	tmpl := l.templates[name]
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return ""
	}

	return buf.String()
}

// Clear 清除缓存
func (l *Loader) Clear() {
	l.templates = make(map[string]*template.Template)
}

// Preload 预加载所有模板
func (l *Loader) Preload(names []string) error {
	for _, name := range names {
		if err := l.Load(name); err != nil {
			return err
		}
	}
	return nil
}
