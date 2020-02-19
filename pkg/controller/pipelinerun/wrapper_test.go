package pipelinerun

import (
	"reflect"
	"testing"

	"github.com/bigkevmcd/commit-status-tracker/pkg/tracker"
	tb "github.com/bigkevmcd/commit-status-tracker/test/builder"
)

func TestFindCommitWithRepository(t *testing.T) {
	pipelineRun := wrap(tb.MakePipelineRunWithResources(
		tb.MakeGitResource("https://github.com/tektoncd/triggers", "master")))
	want := &tracker.Commit{
		RepoURL: "https://github.com/tektoncd/triggers",
		Ref:     "master",
	}

	r, err := pipelineRun.FindCommit()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r, want) {
		t.Fatalf("got %+v, want %+v", r, want)
	}
}
