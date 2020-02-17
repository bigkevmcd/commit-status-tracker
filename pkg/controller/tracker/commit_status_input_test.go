package tracker

import (
	"testing"
)

func TestAnnotationByName(t *testing.T) {
	nt := []struct {
		name        string
		annotations map[string]string
		want        string
	}{
		{"no annotations", map[string]string{}, "default"},
		{"no matching annotation", map[string]string{"testing": "testing"}, "default"},
		{"with matching annotation", map[string]string{StatusContextName: "test-lint"}, "test-lint"},
	}

	for _, tt := range nt {
		r := fakeObject{annotations: tt.annotations}
		if b := getAnnotationByName(r, StatusContextName, "default"); b != tt.want {
			t.Errorf("Context() %s got %v, want %v", tt.name, b, tt.want)
		}
	}
}

type fakeObject struct {
	annotations map[string]string
}

func (fo fakeObject) Annotations() map[string]string {
	return fo.annotations
}

func (fo fakeObject) RunState() State {
	return Pending
}
