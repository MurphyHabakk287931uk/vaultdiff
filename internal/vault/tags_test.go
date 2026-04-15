package vault

import (
	"errors"
	"testing"
)

func TestTagsClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewTagsClient(nil, map[string]string{}, "")
}

func TestTagsClient_PanicsOnNilTags(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil tags")
		}
	}()
	mock := NewMockClient(map[string]map[string]string{})
	NewTagsClient(mock, nil, "")
}

func TestTagsClient_InjectsTagsWithDefaultPrefix(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "hunter2"},
	})
	client := NewTagsClient(mock, map[string]string{"env": "prod"}, "")

	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_PASS"] != "hunter2" {
		t.Errorf("expected DB_PASS=hunter2, got %q", got["DB_PASS"])
	}
	if got["_tag.env"] != "prod" {
		t.Errorf("expected _tag.env=prod, got %q", got["_tag.env"])
	}
}

func TestTagsClient_InjectsTagsWithCustomPrefix(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"KEY": "val"},
	})
	client := NewTagsClient(mock, map[string]string{"region": "us-east-1"}, "meta")

	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["meta.region"] != "us-east-1" {
		t.Errorf("expected meta.region=us-east-1, got %q", got["meta.region"])
	}
}

func TestTagsClient_DoesNotMutateOriginalTags(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {},
	})
	orig := map[string]string{"env": "staging"}
	client := NewTagsClient(mock, orig, "")

	orig["injected"] = "after"

	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["_tag.injected"]; ok {
		t.Error("mutation of original tags map should not affect client")
	}
}

func TestTagsClient_PropagatesError(t *testing.T) {
	expected := errors.New("vault unavailable")
	mock := NewMockClient(nil)
	mock.SetError("secret/broken", expected)

	client := NewTagsClient(mock, map[string]string{"env": "prod"}, "")
	_, err := client.ReadSecrets("secret/broken")
	if !errors.Is(err, expected) {
		t.Errorf("expected wrapped error %v, got %v", expected, err)
	}
}

func TestTagsClient_ImplementsSecretReader(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{})
	var _ SecretReader = NewTagsClient(mock, map[string]string{}, "")
}
