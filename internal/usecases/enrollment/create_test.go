package enrollment

import (
	"context"
	"testing"
)

func TestCreateRejectsMissingEmail(t *testing.T) {
	u := NewCreateUsecase(nil)
	if _, e := u.Execute(context.Background(), "u", false, "c", CreateInput{}); e == nil {
		t.Fatal("expected validation error")
	}
}
