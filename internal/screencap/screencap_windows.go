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

	// Poll for the new window to appear
	var hwnd uintptr
	deadline := time.Now().Add(waitTime)
	for time.Now().Before(deadline) {
		time.Sleep(300 * time.Millisecond)
		afterWindows := enumerateVisibleWindows()
		hwnd = findNewWindow(beforeWindows, afterWindows)
		if hwnd != 0 {
			time.Sleep(1 * time.Second) // let content fully render
			break
		}
	}

	if hwnd == 0 {
		return fmt.Errorf("could not find new terminal window")
	}

	// Capture the window
	img, err := captureWindow(hwnd)
	if err != nil {
		return fmt.Errorf("capture window: %w", err)
	}

	// Save PNG
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("encode png: %w", err)
	}

	// Close the WT window by sending WM_CLOSE
	procPostMessage := user32.NewProc("PostMessageW")
	procPostMessage.Call(hwnd, 0x0010, 0, 0) // WM_CLOSE

	return nil
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

func captureWindow(hwnd uintptr) (*image.RGBA, error) {
	var r rect
	procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))

	width := int(r.Right - r.Left)
	height := int(r.Bottom - r.Top)
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid window dimensions: %dx%d", width, height)
	}

	// Get window DC
	hdcWindow, _, _ := procGetDC.Call(hwnd)
	defer procReleaseDC.Call(hwnd, hdcWindow)

	// Create compatible DC and bitmap
	hdcMem, _, _ := procCreateCompatDC.Call(hdcWindow)
	defer procDeleteDC.Call(hdcMem)

	hBitmap, _, _ := procCreateCompatBmp.Call(hdcWindow, uintptr(width), uintptr(height))
	defer procDeleteObject.Call(hBitmap)

	procSelectObject.Call(hdcMem, hBitmap)

	// PrintWindow with PW_RENDERFULLCONTENT (0x2)
	ret, _, _ := procPrintWindow.Call(hwnd, hdcMem, 0x2)
	if ret == 0 {
		// Fallback: BitBlt from screen
		const SRCCOPY = 0x00CC0020
		procBitBlt.Call(hdcMem, 0, 0, uintptr(width), uintptr(height),
			hdcWindow, 0, 0, SRCCOPY)
	}

	// Read pixels from bitmap
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

	return img, nil
}
