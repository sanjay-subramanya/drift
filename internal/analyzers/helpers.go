package analyzers

import (
	"github.com/sanjay-subramanya/drift/internal/git"
	"strings"
)

func DependencyHits(localFiles, upstreamFiles []string) []string {
	out := []string{}

	for _, local := range localFiles {
		imports, err := extractImports(local)
		if err != nil {
			continue
		}
		for _, imp := range imports {
			for _, up := range upstreamFiles {
				if imp == up {
					out = append(out, up)
				}
			}
		}
	}
	return unique(out)
}

func extractImports(file string) ([]string, error) {
	data, err := git.RunGit("show", "HEAD:"+file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(data, "\n")
	out := []string{}

	for _, l := range lines {
		l = strings.TrimSpace(l)

		// Go / Python / JS style heuristics
		if strings.HasPrefix(l, "import ") ||
			strings.HasPrefix(l, "from ") ||
			strings.Contains(l, "require(") {

			parts := strings.Fields(l)
			for _, p := range parts {
				if strings.Contains(p, "/") && strings.Contains(p, ".") {
					out = append(out, strings.Trim(p, "\"'"))
				}
			}
		}
	}
	return out, nil
}

func stringJoin(files []string) string {
	if len(files) == 0 {
		return "none"
	}
	return strings.Join(files, ", ")
}

func unique(in []string) []string {
	m := map[string]bool{}
	out := []string{}
	for _, x := range in {
		if !m[x] {
			m[x] = true
			out = append(out, x)
		}
	}
	return out
}