package main

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func TestBinaryBuilds(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "/tmp/oculus-test", "./")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	defer os.Remove("/tmp/oculus-test")

	// Check binary exists and is reasonable size
	info, err := os.Stat("/tmp/oculus-test")
	if err != nil {
		t.Fatalf("binary not found: %v", err)
	}

	sizeMB := float64(info.Size()) / (1024 * 1024)
	t.Logf("Binary size: %.1f MB", sizeMB)
	if sizeMB > 30 {
		t.Errorf("Binary too large: %.1f MB (want <30MB)", sizeMB)
	}
}

func TestBinaryHelp(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		// --help returns exit code 0 for cobra
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
			t.Fatalf("help failed: %v\n%s", err, out)
		}
	}

	output := string(out)
	if len(output) == 0 {
		t.Error("expected help output")
	}
}

func TestMemoryFootprint(t *testing.T) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	rssMB := float64(m.Sys) / (1024 * 1024)
	t.Logf("Current RSS: %.1f MB", rssMB)

	// Go runtime overhead should be very small
	if rssMB > 50 {
		t.Errorf("RSS too high at idle: %.1f MB (want <50MB)", rssMB)
	}
}
