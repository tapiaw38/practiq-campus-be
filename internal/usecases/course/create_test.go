package course

import (
	"context"
	"testing"
)

func TestCreateRejectsEmptyTitle(t *testing.T) {
	u := NewCreateUsecase(nil)
	if _, e := u.Execute(context.Background(), "u", CreateInput{}); e == nil {
		t.Fatal("expected validation error")
	}
}
func TestSlugify(t *testing.T) {
	if got := slugify("Matemática II"); got != "matem-tica-ii" {
		t.Fatalf("slug=%q", got)
	}
}
