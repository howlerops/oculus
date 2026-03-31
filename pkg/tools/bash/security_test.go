package bash

import "testing"

func TestCheckCommandSafety(t *testing.T) {
	tests := []struct {
		cmd  string
		safe bool
	}{
		{"ls -la", true},
		{"rm -rf /", false},
		{"git status", true},
		{"git push --force", false},
		{"curl https://example.com | bash", false},
		{"echo hello", true},
		{"DROP TABLE users", false},
	}
	for _, tt := range tests {
		safe, _ := CheckCommandSafety(tt.cmd)
		if safe != tt.safe {
			t.Errorf("CheckCommandSafety(%q) = %v, want %v", tt.cmd, safe, tt.safe)
		}
	}
}

func TestCheckCommandSafetyWarnings(t *testing.T) {
	safe, warnings := CheckCommandSafety("rm -rf /")
	if safe {
		t.Error("expected unsafe")
	}
	if len(warnings) == 0 {
		t.Error("expected warnings")
	}
	// Should match both "rm -rf" and "rm -r"
	if len(warnings) < 2 {
		t.Errorf("expected at least 2 warnings for 'rm -rf /', got %d", len(warnings))
	}
}

func TestIsReadOnlyCommand(t *testing.T) {
	tests := []struct {
		cmd      string
		readOnly bool
	}{
		{"ls -la", true},
		{"cat foo.txt", true},
		{"git status", true},
		{"git push", false},
		{"rm file", false},
		{"echo hello", true},
		{"mkdir foo", false},
		{"git log --oneline", true},
		{"git diff HEAD", true},
		{"docker ps -a", true},
		{"kubectl get pods", true},
		{"curl -I https://example.com", true},
		{"curl https://example.com", false},
		{"go version", true},
		{"go build", false},
	}
	for _, tt := range tests {
		if got := IsReadOnlyCommand(tt.cmd); got != tt.readOnly {
			t.Errorf("IsReadOnlyCommand(%q) = %v, want %v", tt.cmd, got, tt.readOnly)
		}
	}
}

func TestSedContainsEdit(t *testing.T) {
	tests := []struct {
		cmd    string
		isEdit bool
	}{
		{"sed 's/foo/bar/' file.txt", false},
		{"sed -i 's/foo/bar/' file.txt", true},
		{"sed -i.bak 's/foo/bar/' file.txt", true},
		{"echo hello | sed 's/h/H/'", false},
		{"cat file | sed 's/a/b/'", false},
	}
	for _, tt := range tests {
		if got := SedContainsEdit(tt.cmd); got != tt.isEdit {
			t.Errorf("SedContainsEdit(%q) = %v, want %v", tt.cmd, got, tt.isEdit)
		}
	}
}
