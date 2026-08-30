package submission

import (
	"context"
	"testing"
)

func TestCreateRejectsEmptyContent(t *testing.T) {
	u := NewCreateUsecase(nil)
	if _, e := u.Execute(context.Background(), "u", "a", CreateInput{Content: "  "}); e == nil {
		t.Fatal("expected validation error")
	}
}
