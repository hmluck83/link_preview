package linkpreview

import (
	"context"
	"testing"
)

func TestGetLinkPreview(t *testing.T) {
	// TODO: Write your test cases here

	lp, err := GetLinkPreview(context.Background(), "https://")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Title: %s", lp.Title)
	t.Logf("Description: %s", lp.Description)
	t.Logf("Image: %s", lp.Image)


}