package storage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
)

const MaxUploadBytes = 50 << 20

type FileKind string

const (
	FileKindImage    FileKind = "image"
	FileKindAudio    FileKind = "audio"
	FileKindVideo    FileKind = "video"
	FileKindPDF      FileKind = "pdf"
	FileKindDocument FileKind = "doc"
)

var ErrUnsupportedFileType = errors.New("unsupported file type")

var acceptedTypes = map[string]struct {
	kind FileKind
	ext  string
}{
	"audio/webm":         {FileKindAudio, ".webm"},
	"audio/ogg":          {FileKindAudio, ".ogg"},
	"audio/mpeg":         {FileKindAudio, ".mp3"},
	"audio/mp4":          {FileKindAudio, ".m4a"},
	"audio/wav":          {FileKindAudio, ".wav"},
	"audio/x-wav":        {FileKindAudio, ".wav"},
	"audio/wave":         {FileKindAudio, ".wav"},
	"image/png":          {FileKindImage, ".png"},
	"image/jpeg":         {FileKindImage, ".jpg"},
	"image/jpg":          {FileKindImage, ".jpg"},
	"image/webp":         {FileKindImage, ".webp"},
	"image/gif":          {FileKindImage, ".gif"},
	"video/mp4":          {FileKindVideo, ".mp4"},
	"video/webm":         {FileKindVideo, ".webm"},
	"video/ogg":          {FileKindVideo, ".ogv"},
	"video/quicktime":    {FileKindVideo, ".mov"},
	"application/pdf":    {FileKindPDF, ".pdf"},
	"application/msword": {FileKindDocument, ".doc"},
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {FileKindDocument, ".docx"},
	"application/vnd.ms-excel": {FileKindDocument, ".xls"},
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         {FileKindDocument, ".xlsx"},
	"application/vnd.ms-powerpoint":                                             {FileKindDocument, ".ppt"},
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": {FileKindDocument, ".pptx"},
	"text/plain": {FileKindDocument, ".txt"},
	"application/vnd.oasis.opendocument.text": {FileKindDocument, ".odt"},
}

func ClassifyContentType(contentType string) (FileKind, string, error) {
	base := strings.ToLower(strings.TrimSpace(contentType))
	if index := strings.Index(base, ";"); index >= 0 {
		base = strings.TrimSpace(base[:index])
	}
	entry, ok := acceptedTypes[base]
	if !ok {
		return "", "", fmt.Errorf("%w: %s", ErrUnsupportedFileType, contentType)
	}
	return entry.kind, entry.ext, nil
}

func ResolveContentType(contentType string, body []byte) (string, FileKind, string, error) {
	if kind, ext, err := ClassifyContentType(contentType); err == nil {
		if sniffedKind, _, sniffErr := ClassifyContentType(http.DetectContentType(body)); sniffErr == nil && conflictingKinds(kind, sniffedKind) {
			return "", "", "", fmt.Errorf("%w: declared %s but the file is %s", ErrUnsupportedFileType, contentType, sniffedKind)
		}
		return contentType, kind, ext, nil
	}

	sniffed := http.DetectContentType(body)
	kind, ext, err := ClassifyContentType(sniffed)
	if err != nil {
		return "", "", "", err
	}
	return sniffed, kind, ext, nil
}

func conflictingKinds(declared, sniffed FileKind) bool {
	if declared == sniffed {
		return false
	}
	mediaContainer := func(k FileKind) bool {
		return k == FileKindAudio || k == FileKindVideo
	}
	return !(mediaContainer(declared) && mediaContainer(sniffed))
}

func (s *S3Storage) UploadFile(ctx context.Context, folder, userID, filename, contentType string, body []byte) (string, error) {
	if len(body) == 0 {
		return "", errors.New("empty file")
	}
	if len(body) > MaxUploadBytes {
		return "", fmt.Errorf("file is larger than %d MiB", MaxUploadBytes>>20)
	}

	contentType, _, ext, err := ResolveContentType(contentType, body)
	if err != nil {
		return "", err
	}

	key := buildFileKey(folder, userID, ext)
	if err := s.putObject(ctx, key, contentType, body); err != nil {
		return "", err
	}
	return s.objectURL(key), nil
}

func buildFileKey(folder, userID, ext string) string {
	return path.Join("file", cleanPathPart(folder), cleanPathPart(userID), randomHex(16)+ext)
}
