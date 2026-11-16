package gtk

/*
#include <gtk/gtk.h>
*/
import "C"
import "unsafe"

type ProgressBar struct {
	Widget
}

func NewProgressBar() *ProgressBar {

	ptr := C.gtk_progress_bar_new()

	return &ProgressBar{
		Widget: Widget{widget: (*C.GtkWidget)(unsafe.Pointer(ptr))},
	}

}

func (p *ProgressBar) SetFraction(fraction float64) {
	C.gtk_progress_bar_set_fraction((*C.GtkProgressBar)(unsafe.Pointer(p.widget)), C.double(fraction))
}

func (p *ProgressBar) Pulse() {
	C.gtk_progress_bar_pulse((*C.GtkProgressBar)(unsafe.Pointer(p.widget)))
}

func (p *ProgressBar) SetShowText(show bool) {
	if show {
		C.gtk_progress_bar_set_show_text((*C.GtkProgressBar)(unsafe.Pointer(p.widget)), C.TRUE)
	} else {
		C.gtk_progress_bar_set_show_text((*C.GtkProgressBar)(unsafe.Pointer(p.widget)), C.FALSE)
	}
}

func (p *ProgressBar) SetText(text string) {

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	C.gtk_progress_bar_set_text((*C.GtkProgressBar)(unsafe.Pointer(p.widget)), cText)

}
