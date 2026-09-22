package internal_test

import (
	"bufio"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var frameworkImports = []string{
	"github.com/gin-gonic/",
	"github.com/jackc/pgx/",
	"github.com/segmentio/kafka-go",
}

func TestDependencyRule(t *testing.T) {
	internalPrefix := modulePath(t) + "/internal/"

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		from := filepath.ToSlash(filepath.Dir(path))
		for _, spec := range file.Imports {
			imp, _ := strconv.Unquote(spec.Path.Value)
			if reason := violation(from, imp, internalPrefix); reason != "" {
				t.Errorf("%s imports %s: %s", path, imp, reason)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDependencyRuleViolations(t *testing.T) {
	const prefix = "example.com/svc/internal/"

	cases := []struct {
		from, imp string
		allowed   bool
	}{
		{"domain/ticket", prefix + "domain/valueobjects", true},
		{"domain/ticket", prefix + "application/port/out", false},
		{"domain/ticket", prefix + "adapters/out/postgres", false},
		{"domain/ticket", "github.com/gin-gonic/gin", false},
		{"application/command", prefix + "application/port/out", true},
		{"application/command", prefix + "domain/ticket", true},
		{"application/command", prefix + "adapters/out/cache", false},
		{"application/port/out", "github.com/jackc/pgx/v5/pgxpool", false},
		{"adapters/in/http", prefix + "application/command", true},
		{"adapters/in/http", prefix + "adapters/in/http/dto", true},
		{"adapters/in/http", prefix + "adapters/out/cache", false},
		{"adapters/in/messaging", prefix + "adapters/in/http", false},
		{"adapters/out/postgres", "github.com/jackc/pgx/v5/pgxpool", true},
	}

	for _, c := range cases {
		got := violation(c.from, c.imp, prefix) == ""
		if got != c.allowed {
			t.Errorf("%s -> %s: allowed=%v, want %v", c.from, c.imp, got, c.allowed)
		}
	}
}

func violation(from, imp, internalPrefix string) string {
	layer := segment(from, 0)
	target, internal := strings.CutPrefix(imp, internalPrefix)

	if layer == "domain" || layer == "application" {
		for _, framework := range frameworkImports {
			if strings.HasPrefix(imp, framework) {
				return layer + " must not depend on frameworks or drivers"
			}
		}
	}

	if !internal {
		return ""
	}

	switch layer {
	case "domain":
		if segment(target, 0) != "domain" {
			return "domain must only depend on domain"
		}
	case "application":
		if segment(target, 0) == "adapters" {
			return "application must not depend on adapters"
		}
	case "adapters":
		if segment(target, 0) == "adapters" && adapterRoot(target) != adapterRoot(from) {
			return "an adapter must not depend on another adapter"
		}
	}
	return ""
}

func segment(path string, i int) string {
	parts := strings.Split(path, "/")
	if i < len(parts) {
		return parts[i]
	}
	return ""
}

func adapterRoot(path string) string {
	parts := strings.SplitN(path, "/", 4)
	return strings.Join(parts[:min(3, len(parts))], "/")
}

func modulePath(t *testing.T) string {
	t.Helper()

	f, err := os.Open("../go.mod")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if module, ok := strings.CutPrefix(scanner.Text(), "module "); ok {
			return strings.TrimSpace(module)
		}
	}
	t.Fatal("module directive not found in go.mod")
	return ""
}
