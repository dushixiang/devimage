package backend

import (
	"bytes"
	"os"
	"path"
	"time"
)

const (
	Processing = "processing"
	Success    = "success"
	Error      = "error"
)

type FileInfo struct {
	Data       []byte `json:"-"`
	Path       string `json:"path"`
	File       string `json:"file"`
	MimeType   string `json:"mimeType"`
	Size       int64  `json:"size"`
	Compressed int64  `json:"compressed"`
	Ratio      int    `json:"ratio"`
	Status     string `json:"status"`
	Error      string `json:"error"`
	CreatedAt  time.Time
}

func (r *FileInfo) Compress(quality int, outputDir string) (int64, error) {
	var buffer bytes.Buffer

	if r.Path != "" {
		ext := path.Ext(r.Path)
		compress, err := NewImageCompress(ext)
		if err != nil {
			return 0, err
		}
		file, err := os.Open(r.Path)
		if err != nil {
			return 0, err
		}
		img, err := compress.Decode(file)
		if err != nil {
			return 0, err
		}
		buffer, err = compress.Encode(img, &Options{Quality: quality})
		if err != nil {
			return 0, err
		}
	} else {
		ext := path.Ext(r.File)
		compress, err := NewImageCompress(ext)
		if err != nil {
			return 0, err
		}
		img, err := compress.Decode(bytes.NewReader(r.Data))
		if err != nil {
			return 0, err
		}
		buffer, err = compress.Encode(img, &Options{Quality: quality})
		if err != nil {
			return 0, err
		}
	}

	var (
		data    = buffer.Bytes()
		dataLen = int64(len(data))
	)

	if dataLen > r.Size {
		// TODO
	}

	_, err := os.Stat(outputDir)
	if err != nil {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return 0, err
		}
	}

	newFile, err := os.OpenFile(outputDir+"/"+r.File, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}

	_, err = newFile.Write(data)
	if err != nil {
		return 0, err
	}
	return dataLen, nil
}
