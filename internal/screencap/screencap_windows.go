// Package screencap captures screenshots of console windows on Windows.
// It launches a command in a new console window, waits for it to render,
// captures the window content using PrintWindow, and returns the image.
package screencap

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"
)

// Ensure syscall is used (needed for NewCallback).
var _ = syscall.NewCallback

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	gdi32                = syscall.NewLazyDLL("gdi32.dll")
	procGetWindowRect    = user32.NewProc("GetWindowRect")
	procPrintWindow      = user32.NewProc("PrintWindow")
	procGetDC            = user32.NewProc("GetDC")
	procReleaseDC        = user32.NewProc("ReleaseDC")
	procCreateCompatDC   = gdi32.NewProc("CreateCompatibleDC")
	procCreateCompatBmp  = gdi32.NewProc("CreateCompatibleBitmap")
	procSelectObject     = gdi32.NewProc("SelectObject")
	procDeleteObject     = gdi32.NewProc("DeleteObject")
	procDeleteDC         = gdi32.NewProc("DeleteDC")
	procGetDIBits        = gdi32.NewProc("GetDIBits")
	procBitBlt           = gdi32.NewProc("BitBlt")
)

type rect struct {
	Left, Top, Right, Bottom int32
}

type bitmapInfoHeader struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

type bitmapInfo struct {
	Header bitmapInfoHeader
	Colors [1]uint32
}

// CaptureCommand launches a command in a new Windows Terminal window,
// waits for it to render, captures the window, and saves as PNG.
func CaptureCommand(program string, args []string, outPath string, waitTime time.Duration) error {
	// Snapshot existing windows before launch
	beforeWindows := enumerateVisibleWindows()

	// Build the wt.exe arguments: wt -w new --size 100,30 cmd /k <program> <args>
	wtArgs := []string{"-w", "new", "--size", "100,30", "--pos", "100,100", "cmd", "/k", program}
	wtArgs = append(wtArgs, args...)

	cmd := exec.Command("wt.exe", wtArgs...)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start wt: %w", err)
	}

	// Poll for the new window to appear (up to waitTime)
	var hwnd uintptr
	deadline := time.Now().Add(waitTime)
	for time.Now().Before(deadline) {
		time.Sleep(500 * time.Millisecond)
		afterWindows := enumerateVisibleWindows()
		hwnd = findNewWindow(beforeWindows, afterWindows)
		if hwnd != 0 {
			break
		}
	}

	if hwnd == 0 {
		return fmt.Errorf("could not find new terminal window")
	}

	// Wait for content to fully render. WT needs time to:
	// 1. Initialize the terminal session (extra slow on first launch)
	// 2. Run cmd.exe /k which launches moonpool
	// 3. moonpool renders and outputs ANSI
	// 4. WT's GPU renderer paints the content
	if !wtWarmedUp {
		// First WT launch — WT itself needs to start up
		time.Sleep(5 * time.Second)
		wtWarmedUp = true
	} else {
		time.Sleep(3 * time.Second)
	}

	// Capture with retry — if image is blank, wait and try again
	var img *image.RGBA
	for attempt := 0; attempt < 3; attempt++ {
		var err2 error
		img, err2 = captureWindow(hwnd)
		if err2 != nil {
			closeAndWait(hwnd)
			return fmt.Errorf("capture window: %w", err2)
		}
		if !isBlankImage(img) {
			break
		}
		// Image is blank — wait and retry
		time.Sleep(2 * time.Second)
	}

	// Save PNG
	f, err := os.Create(outPath)
	if err != nil {
		closeAndWait(hwnd)
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		closeAndWait(hwnd)
		return fmt.Errorf("encode png: %w", err)
	}

	// Close and wait for the window to actually disappear
	closeAndWait(hwnd)
	return nil
}

// closeAndWait sends WM_CLOSE and waits until the window disappears.
func closeAndWait(hwnd uintptr) {
	procPostMessage := user32.NewProc("PostMessageW")
	procPostMessage.Call(hwnd, 0x0010, 0, 0) // WM_CLOSE

	// Poll until window is gone (up to 3 seconds)
	for i := 0; i < 30; i++ {
		time.Sleep(100 * time.Millisecond)
		visible, _, _ := procIsWindowVisible.Call(hwnd)
		if visible == 0 {
			return
		}
	}
}

var (
	procEnumWindows     = user32.NewProc("EnumWindows")
	procIsWindowVisible = user32.NewProc("IsWindowVisible")
)

func enumerateVisibleWindows() map[uintptr]bool {
	windows := make(map[uintptr]bool)
	cb := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		visible, _, _ := procIsWindowVisible.Call(hwnd)
		if visible != 0 {
			windows[hwnd] = true
		}
		return 1
	})
	procEnumWindows.Call(cb, 0)
	return windows
}

func findNewWindow(before, after map[uintptr]bool) uintptr {
	for hwnd := range after {
		if !before[hwnd] {
			// Verify it has non-zero dimensions
			var r rect
			procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
			w := r.Right - r.Left
			h := r.Bottom - r.Top
			if w > 100 && h > 50 {
				return hwnd
			}
		}
	}
	return 0
}

var (
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
)

// wtWarmedUp tracks whether the first WT window has been launched.
// The first launch is much slower as WT itself needs to initialize.
var wtWarmedUp bool

// isBlankImage checks if an image is essentially a single solid color
// (indicating the terminal hasn't rendered content yet).
func isBlankImage(img *image.RGBA) bool {
	if img == nil {
		return true
	}
	bounds := img.Bounds()
	if bounds.Dx() < 10 || bounds.Dy() < 10 {
		return true
	}

	// Sample pixels from the middle of the image
	// If they're all the same color, the image is blank
	midY := bounds.Dy() / 2
	refR, refG, refB, _ := img.At(bounds.Dx()/4, midY).RGBA()
	uniqueColors := 0

	for x := bounds.Dx()/4; x < 3*bounds.Dx()/4; x += 10 {
		r, g, b, _ := img.At(x, midY).RGBA()
		if r != refR || g != refG || b != refB {
			uniqueColors++
			if uniqueColors > 3 {
				return false // has real content
			}
		}
	}
	return true // all sampled pixels are the same color
}

func captureWindow(hwnd uintptr) (*image.RGBA, error) {
	var r rect
	procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))

	width := int(r.Right - r.Left)
	height := int(r.Bottom - r.Top)
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid window dimensions: %dx%d", width, height)
	}

	// Bring window to foreground
	procSetForegroundWindow.Call(hwnd)
	time.Sleep(500 * time.Millisecond)

	// Try PrintWindow first (captures window directly, Z-order independent).
	// PW_RENDERFULLCONTENT (0x2) tells DWM to render the full content including
	// GPU-accelerated surfaces.
	img := captureViaPrintWindow(hwnd, width, height)
	if img != nil && !isBlankImage(img) {
		return img, nil
	}

	// Fallback: BitBlt from screen DC at window position.
	// Requires window to be visible and on top.
	return captureViaScreen(r, width, height)
}

func captureViaPrintWindow(hwnd uintptr, width, height int) *image.RGBA {
	hdcWindow, _, _ := procGetDC.Call(hwnd)
	if hdcWindow == 0 {
		return nil
	}
	defer procReleaseDC.Call(hwnd, hdcWindow)

	hdcMem, _, _ := procCreateCompatDC.Call(hdcWindow)
	defer procDeleteDC.Call(hdcMem)

	hBitmap, _, _ := procCreateCompatBmp.Call(hdcWindow, uintptr(width), uintptr(height))
	defer procDeleteObject.Call(hBitmap)

	procSelectObject.Call(hdcMem, hBitmap)

	// PW_RENDERFULLCONTENT = 0x2
	ret, _, _ := procPrintWindow.Call(hwnd, hdcMem, 0x2)
	if ret == 0 {
		return nil
	}

	return readBitmap(hdcMem, hBitmap, width, height)
}

func captureViaScreen(r rect, width, height int) (*image.RGBA, error) {
	hdcScreen, _, _ := procGetDC.Call(0)
	defer procReleaseDC.Call(0, hdcScreen)

	hdcMem, _, _ := procCreateCompatDC.Call(hdcScreen)
	defer procDeleteDC.Call(hdcMem)

	hBitmap, _, _ := procCreateCompatBmp.Call(hdcScreen, uintptr(width), uintptr(height))
	defer procDeleteObject.Call(hBitmap)

	procSelectObject.Call(hdcMem, hBitmap)

	const SRCCOPY = 0x00CC0020
	procBitBlt.Call(hdcMem, 0, 0, uintptr(width), uintptr(height),
		hdcScreen, uintptr(r.Left), uintptr(r.Top), SRCCOPY)

	return readBitmap(hdcMem, hBitmap, width, height), nil
}

func readBitmap(hdcMem, hBitmap uintptr, width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	bi := bitmapInfo{
		Header: bitmapInfoHeader{
			Size:        uint32(unsafe.Sizeof(bitmapInfoHeader{})),
			Width:       int32(width),
			Height:      -int32(height), // top-down
			Planes:      1,
			BitCount:    32,
			Compression: 0, // BI_RGB
		},
	}

	procGetDIBits.Call(
		hdcMem,
		hBitmap,
		0,
		uintptr(height),
		uintptr(unsafe.Pointer(&img.Pix[0])),
		uintptr(unsafe.Pointer(&bi)),
		0, // DIB_RGB_COLORS
	)

	// Swap BGRA to RGBA
	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i], img.Pix[i+2] = img.Pix[i+2], img.Pix[i]
	}

	return img
}
