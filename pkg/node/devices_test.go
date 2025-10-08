package node

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNode_GetDevice(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		n := New("test")
		device, err := n.GetDevice("nonexistent")
		if status.Code(err) != codes.NotFound {
			t.Fatalf("expected NotFound error, got %v", err)
		}
		if device != nil {
			t.Fatalf("expected nil device, got %v", device)
		}
	})

	t.Run("found child", func(t *testing.T) {
		n := New("test")
		n.Announce("foo")
		got, err := n.GetDevice("foo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.GetName() != "foo" {
			t.Fatalf("expected device name 'foo', got '%s'", got.GetName())
		}
	})

	t.Run("found root", func(t *testing.T) {
		n := New("test")
		got, err := n.GetDevice("test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.GetName() != "test" {
			t.Fatalf("expected device name 'test', got '%s'", got.GetName())
		}
	})
	t.Run("found default", func(t *testing.T) {
		n := New("test")
		got, err := n.GetDevice("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.GetName() != "test" {
			t.Fatalf("expected device name 'test', got '%s'", got.GetName())
		}
	})
}
