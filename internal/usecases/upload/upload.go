package upload

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/storage"
)

const previewLinkTTL = time.Hour

type (
	Usecase interface {
		Execute(context.Context, Input) (*Output, apperrors.ApplicationError)
	}

	usecase struct {
		contextFactory appcontext.Factory
	}

	Input struct {
		UserID      string
		Folder      string
		Filename    string
		ContentType string
		Reader      io.Reader
		Size        int64
	}

	Output struct {
		Data FileData `json:"data"`
	}

	FileData struct {
		// URL is the canonical bucket URL — the value to persist. The bucket
		// is private, so it is not openable on its own.
		URL string `json:"url"`
		// PreviewURL is a short-lived signed URL for showing the file right
		// after upload. Never stored.
		PreviewURL  string `json:"preview_url,omitempty"`
		Filename    string `json:"filename"`
		ContentType string `json:"content_type"`
		Kind        string `json:"kind"`
		Size        int64  `json:"size"`
	}
)

func NewUsecase(contextFactory appcontext.Factory) Usecase {
	return &usecase{contextFactory: contextFactory}
}

func (u *usecase) Execute(ctx context.Context, input Input) (*Output, apperrors.ApplicationError) {
	app := u.contextFactory()

	if app.Storage == nil || !app.Storage.IsConfigured() {
		return nil, apperrors.NewApplicationError(mappings.UploadNotConfiguredError, nil)
	}
	if input.Size > storage.MaxUploadBytes {
		return nil, apperrors.NewApplicationError(mappings.UploadTooLargeError, nil)
	}

	// Read one byte past the cap so a lying Content-Length is still caught.
	body, err := io.ReadAll(io.LimitReader(input.Reader, storage.MaxUploadBytes+1))
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UploadError, err)
	}
	if len(body) > storage.MaxUploadBytes {
		return nil, apperrors.NewApplicationError(mappings.UploadTooLargeError, nil)
	}

	folder := input.Folder
	if folder == "" {
		folder = "materials"
	}

	contentType, kind, _, err := storage.ResolveContentType(input.ContentType, body)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UploadUnsupportedTypeError, err)
	}

	url, err := app.Storage.UploadFile(ctx, folder, input.UserID, input.Filename, contentType, body)
	if err != nil {
		if errors.Is(err, storage.ErrUnsupportedFileType) {
			return nil, apperrors.NewApplicationError(mappings.UploadUnsupportedTypeError, err)
		}
		log.Printf("[upload] failed user_id=%s filename=%q err=%v", input.UserID, input.Filename, err)
		return nil, apperrors.NewApplicationError(mappings.UploadError, err)
	}

	previewURL, _ := app.Storage.PresignGetURL(url, previewLinkTTL)

	return &Output{Data: FileData{
		URL:         url,
		PreviewURL:  previewURL,
		Filename:    input.Filename,
		ContentType: contentType,
		Kind:        string(kind),
		Size:        int64(len(body)),
	}}, nil
}
