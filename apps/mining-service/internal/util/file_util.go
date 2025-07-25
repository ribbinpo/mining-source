package util

import (
	"path/filepath"
	"strings"
)

func IsImageFile(path string) bool {
	// Get file extension
	ext := strings.ToLower(filepath.Ext(path))

	// List of common image extensions
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".svg":  true,
		".webp": true,
		".ico":  true,
		".tiff": true,
		".tif":  true,
	}

	return imageExtensions[ext]
}

func IsResourceFile(path string) bool {
	// Get file extension
	ext := strings.ToLower(filepath.Ext(path))

	// List of common resource file extensions
	resourceExtensions := map[string]bool{
		".css":   true,
		".js":    true,
		".json":  true,
		".xml":   true,
		".txt":   true,
		".pdf":   true,
		".doc":   true,
		".docx":  true,
		".xls":   true,
		".xlsx":  true,
		".ppt":   true,
		".pptx":  true,
		".zip":   true,
		".rar":   true,
		".7z":    true,
		".tar":   true,
		".gz":    true,
		".mp3":   true,
		".mp4":   true,
		".avi":   true,
		".mov":   true,
		".wmv":   true,
		".flv":   true,
		".webm":  true,
		".woff":  true,
		".woff2": true,
		".ttf":   true,
		".eot":   true,
		".otf":   true,
		".map":   true, // Source maps
		".min":   true, // Minified files
	}

	return resourceExtensions[ext]
}
