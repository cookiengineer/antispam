package gtk

/*
#include <gtk/gtk.h>
*/
import "C"
import "unsafe"

type Window struct {
	widget *C.GtkWidget
}

func NewWindow(app *Application) *Window {

	ptr := C.gtk_application_window_new(app.Ptr())

	return &Window{
		widget: (*C.GtkWidget)(unsafe.Pointer(ptr)),
	}

}

func (w *Window) AsPtr() unsafe.Pointer {
	return unsafe.Pointer(w.widget)
}

func (w *Window) SetTitle(title string) {

	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.gtk_window_set_title((*C.GtkWindow)(unsafe.Pointer(w.widget)), cTitle)

}

func (w *Window) SetDefaultSize(width, height int) {
	C.gtk_window_set_default_size((*C.GtkWindow)(unsafe.Pointer(w.widget)), C.int(width), C.int(height))
}

func (w *Window) SetChild(child unsafe.Pointer) {
	C.gtk_window_set_child((*C.GtkWindow)(unsafe.Pointer(w.widget)), (*C.GtkWidget)(child))
}

func (w *Window) Present() {
	C.gtk_window_present((*C.GtkWindow)(unsafe.Pointer(w.widget)))
}

func (w *Window) SetModal(modal bool) {

	if modal == true {
		C.gtk_window_set_modal((*C.GtkWindow)(unsafe.Pointer(w.widget)), C.TRUE)
	} else {
		C.gtk_window_set_modal((*C.GtkWindow)(unsafe.Pointer(w.widget)), C.FALSE)
	}

}

func (w *Window) SetTransientFor(parent *Window) {

	if parent != nil {
		C.gtk_window_set_transient_for(
			(*C.GtkWindow)(unsafe.Pointer(w.widget)),
			(*C.GtkWindow)(unsafe.Pointer(parent.widget)),
		)
	}

}

func (w *Window) Close() {
	C.gtk_window_destroy((*C.GtkWindow)(unsafe.Pointer(w.widget)))
}
