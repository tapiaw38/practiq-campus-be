package storage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
)

// MaxUploadBytes caps a single upload. Scanned PDFs and short lecture videos
// are the large cases; anything past this is likely a mistake.
const MaxUploadBytes = 50 << 20 // 50 MiB

// FileKind groups accepted content types into the buckets the product talks
// about, so the UI can say "documento" instead of listing MIME types.
type FileKind string

const (
	FileKindImage    FileKind = "image"
	FileKindAudio    FileKind = "audio"
	FileKindVideo    FileKind = "video"
	FileKindPDF      FileKind = "pdf"
	FileKindDocument FileKind = "doc"
)

var ErrUnsupportedFileType = errors.New("unsupported file type")

// acceptedTypes is a whitelist: an unknown content type is rejected rather
// than stored, so the bucket never receives arbitrary binaries.
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

// ClassifyContentType reports the bucket a content type belongs to. Parameters
// such as "audio/webm;codecs=opus" are stripped before matching.
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

// ResolveContentType returns the content type an upload will actually be
// stored with: the declared one when the whitelist accepts it *and* the bytes
// do not say otherwise, else the sniffed one.
//
// http.DetectContentType only recognizes some formats; when it cannot tell
// (application/octet-stream and friends) the declaration stands, since
// rejecting there would break every format Go cannot sniff.
func ResolveContentType(contentType string, body []byte) (string, FileKind, string, error) {
	if kind, ext, err := ClassifyContentType(contentType); err == nil {
		if sniffedKind, _, sniffErr := ClassifyContentType(http.DetectContentType(body)); sniffErr == nil && conflictingKinds(kind, sniffedKind) {
			return "", "", "", fmt.Errorf("%w: declared %s but the file is %s", ErrUnsupportedFileType, contentType, sniffedKind)
		}
		return contentType, kind, ext, nil
	}
	// Some browsers send an empty or generic type; sniff the bytes before
	// rejecting the upload.
	sniffed := http.DetectContentType(body)
	kind, ext, err := ClassifyContentType(sniffed)
	if err != nil {
		return "", "", "", err
	}
	return sniffed, kind, ext, nil
}

// conflictingKinds reports a declaration the bytes contradict.
//
// Audio and video are never treated as a conflict: WebM and Ogg are the same
// container either way, so an audio file is sniffed as video/webm and
// rejecting it would break legitimate uploads. Everything else that Go can
// identify has to agree.
func conflictingKinds(declared, sniffed FileKind) bool {
	if declared == sniffed {
		return false
	}
	mediaContainer := func(k FileKind) bool {
		return k == FileKindAudio || k == FileKindVideo
	}
	return !(mediaContainer(declared) && mediaContainer(sniffed))
}

// UploadFile stores a file and returns its canonical bucket URL. contentType is
// trusted only after passing the whitelist; when the client sends nothing
// usable the bytes are sniffed instead.
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

	// The extension always comes from the verified content type, never from the
	// client-supplied filename — a hostile name cannot traverse or collide.
	key := buildFileKey(folder, userID, ext)
	if err := s.putObject(ctx, key, contentType, body); err != nil {
		return "", err
	}
	return s.objectURL(key), nil
}

func buildFileKey(folder, userID, ext string) string {
	return path.Join("file", cleanPathPart(folder), cleanPathPart(userID), randomHex(16)+ext)
}
