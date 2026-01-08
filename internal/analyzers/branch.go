package analyzers

import (
	"github.com/sanjay-subramanya/drift/internal/core/model"
	"github.com/sanjay-subramanya/drift/internal/git"
	"github.com/sanjay-subramanya/drift/internal/workspace"
	"github.com/sanjay-subramanya/drift/internal/config"
	"fmt"
	"slices"
	"strings"
	"strconv"
)

func AnalyzeBranch(base string) ([]model.Drift, error) {
	if _, err := git.RunGit("rev-parse", "--verify", "--quiet", base); err != nil {
		return nil, fmt.Errorf("Base branch \"%s\" does not exist", base)
	}
	_, _ = git.RunGit("fetch", "origin")

	mergeBase, err := git.MergeBase("HEAD", base)
	if err != nil {
		return nil, err
	}

	behind, err := git.CommitsBehind("HEAD", base)
	if err != nil {
		return nil, err
	}

	if behind == 0 {
		return nil, nil
	}

	upstreamFiles, err := git.UpstreamFiles(mergeBase, base)
	if err != nil {
		return nil, err
	}

	localDirty, err := git.DirtyFiles()
	if err != nil {
		return nil, err
	}

	// localFilesOut, err := git.RunGit("ls-files")
	if err != nil {
		return nil, err
	}
	// localTracked := strings.Split(localFilesOut, "\n")

	ignores := config.LoadIgnoreFile()

	// Classification buckets (each file goes into only 1)
	var critical []string
	var high []string
	var low []string

	depHits := DependencyHits(localDirty, upstreamFiles)

	for _, f := range upstreamFiles {
		if config.IsIgnored(f, ignores) {
			continue
		}

		switch {
		// file exists locally (dirty)
		case slices.Contains(localDirty, f):
			critical = append(critical, f)

		// env / docker / deployment files
		case workspace.IsEnvFile(f) || workspace.IsDeploymentFile(f):
			high = append(high, f)

		// dependency you import
		case slices.Contains(depHits, f):
			high = append(high, f)

		// everything else
		default:
			low = append(low, f)
		}
	}

	var lines []string
	lines = append(lines, "branch behind by "+strconv.Itoa(behind)+" commits;")

	if len(critical) > 0 {
		lines = append(lines,
			"[CRITICAL] files YOU are editing changed upstream: "+stringJoin(critical))
	}
	if len(high) > 0 {
		lines = append(lines,
			"[HIGH] deployment / dependency files changed upstream: "+stringJoin(high))
	}
	if len(low) > 0 {
		lines = append(lines,
			"[LOW] other files changed upstream: "+stringJoin(low))
	}

	return []model.Drift{
		{
			Type:    model.DriftBranch,
			Summary: strings.Join(lines, "\n"),
		},
	}, nil
}
