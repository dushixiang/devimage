package main

import (
	"changeme/backend"
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	fileManager *backend.FileManager
	configStore *backend.ConfigStore
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		fileManager: &backend.FileManager{},
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	configStore, err := backend.NewConfigStore()
	if err != nil {
		return
	}
	a.configStore = configStore
	a.fileManager.OnChange = func(infos []backend.FileInfo) {
		runtime.EventsEmit(a.ctx, "CompressChange", infos)
	}
	a.fileManager.OnCompressed = func(infos []backend.FileInfo) {
		runtime.EventsEmit(a.ctx, "CompressChange", infos)
	}
}

func (a *App) ClearCompressed() {
	a.fileManager.ClearCompressed()
	a.fileManager.OnChange(a.fileManager.FileList())
}

func (a *App) GetConfig() map[string]interface{} {
	config, err := a.configStore.Load()
	if err != nil {
		return map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}
	}
	return map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    config,
	}
}

func (a *App) SetConfig(values map[string]interface{}) map[string]interface{} {
	var config = backend.Config{
		Quality:   int(values["quality"].(float64)),
		OutputDir: values["outputDir"].(string),
	}

	if err := a.configStore.Store(&config); err != nil {
		return map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}
	}

	return map[string]interface{}{
		"code":    200,
		"message": "success",
	}
}

// OpenFileDialog 打开对话框
func (a *App) OpenFileDialog() map[string]interface{} {
	selection, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Images (*.png;*.jpg)",
				Pattern:     "*.png;*.jpg",
			},
		},
	})
	println("selection", selection)
	if err != nil {
		return map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}
	}

	config, err := a.configStore.Load()
	if err != nil {
		return map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}
	}

	var (
		outputDir = config.OutputDir
		quality   = config.Quality
	)

	for _, s := range selection {
		a.fileManager.Add(s, quality, outputDir)
	}
	fileList := a.fileManager.FileList()
	return map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    fileList,
	}
}

func (a *App) AddFile(content string) map[string]interface{} {
	var file map[string]interface{}
	err := json.Unmarshal([]byte(content), &file)
	if err != nil {
		return map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}
	}

	config, err := a.configStore.Load()
	if err != nil {
		return map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}
	}

	var (
		outputDir = config.OutputDir
		quality   = config.Quality
	)

	fileName := file["file"].(string)
	size := int(file["size"].(float64))
	raw := file["raw"].(string)
	imageBase64 := strings.Split(raw, ",")[1]
	var data = make([]byte, size)
	n, err := base64.StdEncoding.Decode(data, []byte(imageBase64))
	if err != nil {
		return map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}
	}
	a.fileManager.AddFromFrontend(fileName, data[:n], quality, outputDir)

	return map[string]interface{}{
		"code":    200,
		"message": "success",
	}
}
