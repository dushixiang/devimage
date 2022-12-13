package backend

import (
	"os"
	"sort"
	"sync"
	"time"
)

type FileManager struct {
	files        sync.Map
	OnCompressed func(infos []FileInfo)
	OnChange     func(infos []FileInfo)
}

func (r *FileManager) FileList() []FileInfo {
	var fileList []FileInfo
	r.files.Range(func(key, value any) bool {
		fileList = append(fileList, value.(FileInfo))
		return true
	})
	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].CreatedAt.After(fileList[j].CreatedAt)
	})
	return fileList
}

func (r *FileManager) ClearCompressed() {
	r.files.Range(func(key, value any) bool {
		info := value.(FileInfo)
		if info.Status != Processing {
			r.files.Delete(key)
		}
		return true
	})
}

func (r *FileManager) Add(path string, quality int, outputDir string) {
	fileInfo := FileInfo{
		Path:       path,
		File:       "",
		Size:       0,
		Compressed: 0,
		Ratio:      0,
		Status:     Processing,
		Error:      "",
		CreatedAt:  time.Now(),
	}

	file, err := os.Stat(path)
	if err != nil {
		fileInfo.Status = Error
		fileInfo.Error = err.Error()

		r.files.Store(path, fileInfo)
		r.onStatusChange()
	} else {
		fileInfo.File = file.Name()
		fileInfo.Size = file.Size()

		r.files.Store(path, fileInfo)
		r.onStatusChange()

		go func() {
			compressed, err := fileInfo.Compress(quality, outputDir)
			if err != nil {
				fileInfo.Status = Error
				fileInfo.Error = err.Error()
			} else {
				fileInfo.Status = Success
				fileInfo.Compressed = int64(compressed)
				ratio := fileInfo.Compressed * 100 / fileInfo.Size
				fileInfo.Ratio = int(ratio)
			}

			r.files.Store(path, fileInfo)
			r.onStatusChange()

			var completed = true
			r.files.Range(func(key, value any) bool {
				fi := value.(FileInfo)
				if fi.Status == Processing {
					completed = false
				}
				return true
			})
			if completed {
				r.onCompressed()
			}
		}()
	}

	r.files.Store(path, fileInfo)
	return
}

func (r *FileManager) onStatusChange() {
	if r.OnChange != nil {
		r.OnChange(r.FileList())
	}
}

func (r *FileManager) onCompressed() {
	if r.OnCompressed != nil {
		r.OnCompressed(r.FileList())
	}
}

func (r *FileManager) AddFromFrontend(name string, raw []byte, quality int, outputDir string) {
	fileInfo := FileInfo{
		Data:       raw,
		Path:       "",
		File:       name,
		Size:       int64(len(raw)),
		Compressed: 0,
		Ratio:      0,
		Status:     Processing,
		Error:      "",
		CreatedAt:  time.Now(),
	}

	r.files.Store(name, fileInfo)
	r.onStatusChange()

	go func() {
		compressed, err := fileInfo.Compress(quality, outputDir)
		if err != nil {
			println("err", err.Error())
			fileInfo.Status = Error
			fileInfo.Error = err.Error()
		} else {
			fileInfo.Status = Success
			fileInfo.Compressed = int64(compressed)
			ratio := fileInfo.Compressed * 100 / fileInfo.Size
			fileInfo.Ratio = int(ratio)
		}

		r.files.Store(name, fileInfo)
		r.files.Store(name, fileInfo)

		var completed = true
		r.files.Range(func(key, value any) bool {
			fi := value.(FileInfo)
			if fi.Status == Processing {
				completed = false
			}
			return true
		})
		if completed {
			r.onCompressed()
		}
	}()
}
