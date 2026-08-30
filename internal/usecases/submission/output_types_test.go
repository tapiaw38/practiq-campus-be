package submission

import (
	"context"
	"testing"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/storage"
)

// fakeStorage owns only the objects it was told about, so a test can express
// "this file belongs to that student" without an S3 config.
type fakeStorage struct {
	ownedBy map[string]string
}

func (fakeStorage) IsConfigured() bool { return true }

func (fakeStorage) UploadFile(context.Context, string, string, string, string, []byte) (string, error) {
	return "", nil
}

func (fakeStorage) PresignGetURL(rawURL string, _ time.Duration) (string, bool) {
	return rawURL + "?signed=1", true
}

func (f fakeStorage) OwnsFileURL(rawURL, folder, userID string) bool {
	return folder == SubmissionsFolder && f.ownedBy[rawURL] == userID
}

func contextWith(s storage.Storage) *appcontext.Context {
	return &appcontext.Context{Storage: s}
}

func TestWithAttachmentsSignsTheSubmittersOwnFiles(t *testing.T) {
	const mine = "https://bucket.s3.us-east-1.amazonaws.com/file/submissions/student-1/abc-informe.pdf"
	store := fakeStorage{ownedBy: map[string]string{mine: "student-1"}}

	data := SubmissionData{
		UserID:      "student-1",
		Content:     "Ahí va el TP.\n\nArchivo adjunto: informe.pdf\n" + mine,
		Attachments: make([]AttachmentData, 0),
	}

	got := withAttachments(contextWith(store), data)

	if len(got.Attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(got.Attachments))
	}
	if got.Attachments[0].Filename != "informe.pdf" {
		t.Errorf("filename = %q, want the name from the marker", got.Attachments[0].Filename)
	}
	if got.Attachments[0].ViewURL != mine+"?signed=1" {
		t.Errorf("view_url = %q, want the presigned URL", got.Attachments[0].ViewURL)
	}
	if got.Attachments[0].URL != mine {
		t.Errorf("url = %q, want the stored value", got.Attachments[0].URL)
	}
}

// The submission body is free text, so a student can type any URL into it.
// Signing one that belongs to somebody else would make the form a way to read
// other people's files.
func TestWithAttachmentsIgnoresFilesOwnedBySomeoneElse(t *testing.T) {
	const theirs = "https://bucket.s3.us-east-1.amazonaws.com/file/submissions/student-2/secreto.pdf"
	store := fakeStorage{ownedBy: map[string]string{theirs: "student-2"}}

	data := SubmissionData{
		UserID:      "student-1",
		Content:     "Archivo adjunto: secreto.pdf\n" + theirs,
		Attachments: make([]AttachmentData, 0),
	}

	if got := withAttachments(contextWith(store), data); len(got.Attachments) != 0 {
		t.Fatalf("expected no attachments, got %d", len(got.Attachments))
	}
}

func TestWithAttachmentsFallsBackToTheObjectName(t *testing.T) {
	const mine = "https://bucket.s3.us-east-1.amazonaws.com/file/submissions/student-1/abc-informe%20final.pdf"
	store := fakeStorage{ownedBy: map[string]string{mine: "student-1"}}

	data := SubmissionData{
		UserID:      "student-1",
		Content:     "Sin marcador, solo el link: " + mine,
		Attachments: make([]AttachmentData, 0),
	}

	got := withAttachments(contextWith(store), data)

	if len(got.Attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(got.Attachments))
	}
	if got.Attachments[0].Filename != "abc-informe final.pdf" {
		t.Errorf("filename = %q, want the unescaped object name", got.Attachments[0].Filename)
	}
}

func TestWithAttachmentsLeavesPlainTextAlone(t *testing.T) {
	store := fakeStorage{ownedBy: map[string]string{}}

	data := SubmissionData{
		UserID:      "student-1",
		Content:     "Resolví los cinco ejercicios a mano.",
		Attachments: make([]AttachmentData, 0),
	}

	if got := withAttachments(contextWith(store), data); len(got.Attachments) != 0 {
		t.Fatalf("expected no attachments, got %d", len(got.Attachments))
	}
}

// Storage is a no-op in local development; a submission must still serialize.
func TestWithAttachmentsWithoutStorage(t *testing.T) {
	data := SubmissionData{UserID: "student-1", Content: "algo", Attachments: make([]AttachmentData, 0)}

	if got := withAttachments(&appcontext.Context{}, data); len(got.Attachments) != 0 {
		t.Fatalf("expected no attachments, got %d", len(got.Attachments))
	}
}
