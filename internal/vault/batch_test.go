package vault

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestBatchClient_EmptyPaths(t *testing.T) {
	client := NewBatchClient(NewMockClient(nil, nil), 2)
	results := client.ReadAll(context.Background(), nil)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestBatchClient_SinglePath(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"key": "val"},
	}, nil)
	client := NewBatchClient(mock, 2)
	results := client.ReadAll(context.Background(), []string{"secret/a"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].Secrets["key"] != "val" {
		t.Errorf("expected val, got %q", results[0].Secrets["key"])
	}
}

func TestBatchClient_MultiplePaths_OrderPreserved(t *testing.T) {
	data := map[string]map[string]string{}
	paths := make([]string, 10)
	for i := 0; i < 10; i++ {
		p := fmt.Sprintf("secret/path%d", i)
		paths[i] = p
		data[p] = map[string]string{"idx": fmt.Sprintf("%d", i)}
	}
	client := NewBatchClient(NewMockClient(data, nil), 3)
	results := client.ReadAll(context.Background(), paths)
	for i, r := range results {
		expected := fmt.Sprintf("%d", i)
		if r.Secrets["idx"] != expected {
			t.Errorf("result[%d]: expected idx=%s, got %s", i, expected, r.Secrets["idx"])
		}
	}
}

func TestBatchClient_PropagatesError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	mock := NewMockClient(nil, map[string]error{"secret/bad": sentinel})
	client := NewBatchClient(mock, 2)
	results := client.ReadAll(context.Background(), []string{"secret/bad"})
	if !errors.Is(results[0].Err, sentinel) {
		t.Errorf("expected sentinel error, got %v", results[0].Err)
	}
}

func TestBatchClient_DefaultConcurrency(t *testing.T) {
	c := NewBatchClient(NewMockClient(nil, nil), 0)
	if c.concurrency != 4 {
		t.Errorf("expected default concurrency 4, got %d", c.concurrency)
	}
}

func TestBatchClient_NilInnerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil inner client")
		}
	}()
	NewBatchClient(nil, 2)
}

func TestBatchClient_PathsMatchResults(t *testing.T) {
	data := map[string]map[string]string{
		"secret/x": {"k": "v1"},
		"secret/y": {"k": "v2"},
	}
	client := NewBatchClient(NewMockClient(data, nil), 2)
	paths := []string{"secret/x", "secret/y"}
	results := client.ReadAll(context.Background(), paths)
	for i, r := range results {
		if r.Path != paths[i] {
			t.Errorf("result[%d]: expected path %q, got %q", i, paths[i], r.Path)
		}
	}
}

func TestMergeResults_Success(t *testing.T) {
	results := []BatchResult{
		{Path: "secret/a", Secrets: map[string]string{"x": "1"}},
		{Path: "secret/b", Secrets: map[string]string{"y": "2"}},
	}
	merged, err := MergeResults(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if merged["secret/a:x"] != "1" || merged["secret/b:y"] != "2" {
		t.Errorf("unexpected merged map: %v", merged)
	}
}

func TestMergeResults_ErrorAborts(t *testing.T) {
	sentinel := errors.New("read error")
	results := []BatchResult{
		{Path: "secret/a", Secrets: map[string]string{"x": "1"}},
		{Path: "secret/b", Err: sentinel},
	}
	_, err := MergeResults(results)
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}
