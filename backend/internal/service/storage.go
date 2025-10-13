package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type StorageService interface {
	SaveFile(fileHeader *multipart.FileHeader) (string, int64, error)
}

type localStorage struct {
	basePath string
}

func NewLocalStorage() StorageService {
	base := viper.GetString("uploads.path")
	if base == "" {
		base = "./uploads"
	}
	return &localStorage{basePath: base}
}

func (s *localStorage) SaveFile(fh *multipart.FileHeader) (string, int64, error) {
	src, err := fh.Open()
	if err != nil {
		return "", 0, err
	}
	defer src.Close()
	// check configured max size (if provided)
	max := viper.GetInt64("uploads.max_size")
	if max > 0 && fh.Size > 0 && fh.Size > max {
		return "", 0, fmt.Errorf("file too large")
	}
	// generate safe path
	now := time.Now()
	dir := filepath.Join(s.basePath, fmt.Sprintf("%04d", now.Year()), fmt.Sprintf("%02d", now.Month()), fmt.Sprintf("%02d", now.Day()))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", 0, err
	}
	// random filename
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", 0, err
	}
	ext := filepath.Ext(fh.Filename)
	name := hex.EncodeToString(b) + ext
	full := filepath.Join(dir, name)
	dst, err := os.Create(full)
	if err != nil {
		return "", 0, err
	}
	defer dst.Close()
	// read first bytes to detect content type
	buf := make([]byte, 512)
	nread, _ := io.ReadFull(src, buf)
	detected := http.DetectContentType(buf[:nread])
	// allowed types check
	allowed := viper.GetStringSlice("uploads.allowed_types")
	if len(allowed) > 0 {
		ok := false
		for _, a := range allowed {
			if a == detected {
				ok = true
				break
			}
			// support wildcard like image/*
			if strings.HasSuffix(a, "/*") {
				prefix := strings.TrimSuffix(a, "*")
				if strings.HasPrefix(detected, prefix) {
					ok = true
					break
				}
			}
		}
		if !ok {
			// cleanup
			dst.Close()
			os.Remove(full)
			return "", 0, fmt.Errorf("disallowed content type: %s", detected)
		}
	}

	// write initial buffer and then the rest
	wn, err := dst.Write(buf[:nread])
	if err != nil {
		os.Remove(full)
		return "", 0, err
	}
	total := int64(wn)
	// copy rest, but monitor max size if needed
	if max > 0 {
		// limited reader ensuring we don't read beyond max
		lr := &io.LimitedReader{R: src, N: max - total}
		m, err := io.Copy(dst, lr)
		total += m
		if err != nil && err != io.EOF {
			os.Remove(full)
			return "", 0, err
		}
		if lr.N == 0 {
			// exceeded
			dst.Close()
			os.Remove(full)
			return "", 0, fmt.Errorf("file too large")
		}
	} else {
		m, err := io.Copy(dst, src)
		if err != nil {
			os.Remove(full)
			return "", 0, err
		}
		total += m
	}
	// return path relative to base
	rel, err := filepath.Rel(s.basePath, full)
	if err != nil {
		rel = full
	}
	return rel, total, nil
}
