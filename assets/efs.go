package assets

import (
	"embed"
	"path/filepath"

	"github.com/gookit/goutil/fsutil"
)

//go:embed lang
var fs embed.FS

// CopyLangFiles 复制语言文件到指定目录
func CopyLangFiles(targetLangDir string) error {
	files, err := fs.ReadDir("lang")
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := filepath.Join("lang", file.Name())
		filePath = filepath.ToSlash(filePath)
		content, err := fs.ReadFile(filePath)
		if err != nil {
			return err
		}
		copyToFilePath := filepath.Join(targetLangDir, file.Name())
		fsutil.WriteFile(copyToFilePath, content, 0644)
	}
	return nil
}
