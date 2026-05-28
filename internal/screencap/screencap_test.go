package screencap

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCaptureCommand(t *testing.T) {
	outPath := filepath.Join(os.TempDir(), "moonpool_test_real_screenshot.png")

	// Find moonpool binary — use absolute path since WT opens in a different cwd
	moonpool, err := filepath.Abs(filepath.Join("..", "..", "moonpool.exe"))
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}
	if _, err := os.Stat(moonpool); err != nil {
		t.Skipf("moonpool.exe not built: %v", err)
	}

	fixture, err := filepath.Abs(filepath.Join("..", "..", "testdata", "fixtures", "headings.md"))
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}

	err = CaptureCommand(moonpool, []string{"-s", "dark", "-w", "80", fixture}, outPath, 2*time.Second)
	if err != nil {
		t.Fatalf("CaptureCommand failed: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	t.Logf("Screenshot: %s (%d KB)", outPath, info.Size()/1024)
}
