package taskrun

import (
	"reflect"
	"testing"

	"github.com/bigkevmcd/commit-status-tracker/pkg/tracker"
	tb "github.com/bigkevmcd/commit-status-tracker/test/builder"
)

func TestFindCommitWithRepository(t *testing.T) {
	pipelineRun := wrap(tb.MakeTaskRunWithInputResources(
		tb.MakeGitResource("https://github.com/tektoncd/triggers", "master")))

	r, err := pipelineRun.FindCommit()
	if err != nil {
		t.Fatal(err)
	}
	want := &tracker.Commit{
		RepoURL: "https://github.com/tektoncd/triggers",
		Ref:     "master",
	}
	if !reflect.DeepEqual(r, want) {
		t.Fatalf("got %+v, want %+v", r, want)
	}
}

func TestRunState(t *testing.T) {
	t.Skip()
}

func TestAnnotations(t *testing.T) {
	t.Skip()
}
