package tools

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/utils"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type (
	TemplateRenderer struct {
		Templates *template.Template
	}

	tplFile struct {
		Name string
		Path string
	}
)

// ParseGlob 自定义模版解析，扫描子目录
func ParseGlob(tpl *template.Template, dir string, pattern string) (t *template.Template, err error) {
	t = tpl
	fmt.Println("扫描模版目录：" + dir)
	files := allFiles(dir, pattern)
	for _, file := range files {
		fmt.Printf("挂载模板：%s\n", file.Path)
		b, err := ioutil.ReadFile(file.Path)
		if err != nil {
			return t, err
		}
		s := string(b)
		name := file.Name
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

// 目录下的所有文件
func allFiles(dir string, suffix string) (arr []*tplFile) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Println("扫描子目录：" + file.Name())
			arr = append(arr, allFiles(path.Join(dir, file.Name()), suffix)...)
		} else {
			ok, _ := filepath.Match(suffix, file.Name())
			if ok {
				pathName := path.Join(dir, file.Name())
				list := strings.Split(filepath.ToSlash(pathName), "/")
				if len(list) > 1 {
					list = list[1:]
				}
				// fmt.Println(pathName, list)
				name := strings.Join(list, "/")
				arr = append(arr, &tplFile{Name: name, Path: pathName})
			}
		}
	}

	return
}

// SiteInfo SiteInfo
type SiteInfo struct {
	Title string
	User  *models.UserModel
}

// SetUser SetUser
func (s *SiteInfo) SetUser(user *models.UserModel) string {
	s.User = user
	return ""
}

// TemplateFuncs 模板方法
var TemplateFuncs = template.FuncMap{
	"msg": func() string { return "hello this is a msg" },
	"strDefault": func(str string, def string) string {
		if str != "" {
			return str
		}
		return def
	},
	"timeFormat": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	},
	"time": func(t time.Time, layout string) string {
		return t.Format(layout)
	},
	"now": func() time.Time {
		return time.Now()
	},
	"tTime": func(t utils.JSONTime, layout string) string {
		return t.Format(layout)
	},
	"admin": func() map[string]interface{} {
		return map[string]interface{}{
			"name": "MD webSite",
		}
	}, "site": func(title string, args ...interface{}) *SiteInfo {
		return &SiteInfo{
			Title: title,
			User:  &models.UserModel{},
		}
	},
	"tFormatDate": func(t utils.JSONTime) string {
		return t.Format("2006-01-02")
	},
	"thisWeek": func() string {
		now := time.Now()
		delta := (time.Monday - now.Weekday())
		if delta > 0 {
			delta = -6
		}
		monday := now.AddDate(0, 0, int(delta))
		return monday.Format("2006-01-02")
	},
	"lastWeek": func() string {
		now := time.Now()
		delta := (time.Monday - now.Weekday())
		if delta > 0 {
			delta = -6
		}
		monday := now.AddDate(0, 0, int(delta-7))
		return monday.Format("2006-01-02")
	},
	"today": func() string {
		return time.Now().Format("2006-01-02")
	},
	"yesterday": func() string {
		return time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	},
	"thisMonth": func() string {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	},
	"lastMonth": func() string {
		now := time.Now()
		return time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	},
	"url": MakeURL,
	"strJoin": func(str string, args ...interface{}) string {
		return fmt.Sprintf(str, args...)
	},
}

// MakeURL 管理模版中的url
func MakeURL(_type string, args ...interface{}) string {
	urls := map[string]string{
		// 项目
		"projectDetail": "/projects/detail/%d",
		"projectAdd":    "/projects/add",
		"projectUpdate": "/projects/update/%d",
		"projectDel":    "/projects/del/", //ajax拼接id
		// 项目日志
		"projectAddLog":    "/projects/detail/%d/logs/add",
		"projectUpdateLog": "/project-logs/update/%d",
	}
	urlFormat, ok := urls[_type]
	if ok {
		return fmt.Sprintf(urlFormat, args...)
	}
	return "/error"
}
