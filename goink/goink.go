package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func panic(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err.Error())
		os.Exit(1)
	}
}

func ioctl(a1, a2 uintptr, a3 unsafe.Pointer) error {
	_, _, errno := syscall.RawSyscall(syscall.SYS_IOCTL, a1, a2, uintptr(a3))
	if errno != 0 {
		return errno
	}
	return nil
}

func main() {
	// Open framebuffer
	fb0, err := os.OpenFile("/dev/fb0", os.O_RDWR|os.O_APPEND, os.ModeDevice)
	panic("OpenFile", err)

	// Get screen info
	var screen fb_var_screeninfo
	err = ioctl(fb0.Fd(), FBIOGET_VSCREENINFO, unsafe.Pointer(&screen))
	panic("vscreeninfo", err)

	// TODO: rotation_hack

	screensize := screen.Xres_virtual * screen.Yres_virtual * 2

	fb0map, err := syscall.Mmap(int(fb0.Fd()), 0, int(screensize), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	panic("mmap", err)

	for i := uint32(0); i < screensize; i++ {
		fb0map[i] = 0
	}

	// TODO: PNG

	var update mxcfb_update_data

	update = mxcfb_update_data{
		Temp:          TEMP_USE_AMBIENT,
		Update_marker: 0,
		Update_mode:   UPDATE_MODE_FULL,
		Update_region: mxcfb_rect{
			Height: 64,
			Width:  64,
		},
		Waveform_mode: WAVEFORM_MODE_AUTO,
	}

	err = ioctl(fb0.Fd(), MXCFB_SEND_UPDATE, unsafe.Pointer(&update))
	panic("mxcfb_send_update", err)
}
