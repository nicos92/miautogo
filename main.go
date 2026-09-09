package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procKeybdEvent          = user32.NewProc("keybd_event")
	procFindWindow          = user32.NewProc("FindWindowW")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
)

// keybd_event es más tolerante que SendInput
func keybdEvent(bVk byte, bScan byte, dwFlags uint32, dwExtraInfo uintptr) {
	procKeybdEvent.Call(uintptr(bVk), uintptr(bScan), uintptr(dwFlags), uintptr(dwExtraInfo))
}

// Constantes
const (
	KEYEVENTF_KEYDOWN = 0x0000
	KEYEVENTF_KEYUP   = 0x0002
	VK_RETURN         = 0x0D
	VK_LEFT           = 0x25
	VK_UP             = 0x26
	VK_RIGHT          = 0x27
	VK_DOWN           = 0x28
	VK_LWIN           = 0x5B
	VK_R              = 0x52
)

func sendKey(vk byte) {
	keybdEvent(vk, 0, KEYEVENTF_KEYDOWN, 0)
	time.Sleep(50 * time.Millisecond)
	keybdEvent(vk, 0, KEYEVENTF_KEYUP, 0)
}

// Escribe texto caracter por caracter (solo ASCII, sin acentos)
func typeText(s string) {
	for _, ch := range s {
		// Para mayúsculas
		if ch >= 'A' && ch <= 'Z' {
			keybdEvent(0x10, 0, KEYEVENTF_KEYDOWN, 0) // Shift down
			keybdEvent(byte(ch), 0, KEYEVENTF_KEYDOWN, 0)
			time.Sleep(20 * time.Millisecond)
			keybdEvent(byte(ch), 0, KEYEVENTF_KEYUP, 0)
			keybdEvent(0x10, 0, KEYEVENTF_KEYUP, 0)
		} else if ch >= 'a' && ch <= 'z' {
			// minúscula: VK es la misma, pero sin shift
			keybdEvent(byte(ch-32), 0, KEYEVENTF_KEYDOWN, 0)
			time.Sleep(20 * time.Millisecond)
			keybdEvent(byte(ch-32), 0, KEYEVENTF_KEYUP, 0)
		} else if ch >= '0' && ch <= '9' {
			keybdEvent(byte(ch), 0, KEYEVENTF_KEYDOWN, 0)
			time.Sleep(20 * time.Millisecond)
			keybdEvent(byte(ch), 0, KEYEVENTF_KEYUP, 0)
		} else if ch == ' ' {
			keybdEvent(0x20, 0, KEYEVENTF_KEYDOWN, 0)
			time.Sleep(20 * time.Millisecond)
			keybdEvent(0x20, 0, KEYEVENTF_KEYUP, 0)
		}
		// otros caracteres (punto, etc.) se omiten por simplicidad
		time.Sleep(30 * time.Millisecond)
	}
}

// Enfoca una ventana por su título (parcial)
func focusWindow(title string) bool {
	// Convertir título a UTF-16
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	hwnd, _, _ := procFindWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	if hwnd == 0 {
		return false
	}
	procSetForegroundWindow.Call(hwnd)
	time.Sleep(500 * time.Millisecond)
	return true
}

func main() {
	// Espera 3 segundos para que el usuario tenga tiempo de poner el foco en la ventana correcta
	fmt.Println("El script comenzará en 3 segundos...")
	time.Sleep(3 * time.Second)

	// Opcional: intentar enfocar una ventana específica (ej: "Bloc de notas")
	// focusWindow("Bloc de notas")

	// Abrir menú Inicio
	sendKey(VK_LWIN)
	time.Sleep(500 * time.Millisecond)

	// Escribir "notepad"
	typeText("notepad")
	time.Sleep(300 * time.Millisecond)
	sendKey(VK_RETURN)
	time.Sleep(1500 * time.Millisecond)

	// Escribir texto
	typeText("nicos")
	time.Sleep(300 * time.Millisecond)

	sendKey(VK_RETURN)
	// Escribir texto
	typeText("340480")
	time.Sleep(300 * time.Millisecond)

	sendKey(VK_RETURN)
	time.Sleep(150 * time.Millisecond)

	// Flechas
	for range 4 {
		sendKey(VK_DOWN)
		time.Sleep(150 * time.Millisecond)
	}
	for range 10 {

		sendKey(VK_RETURN)
	}

	fmt.Println("Flujo completado.")
}
