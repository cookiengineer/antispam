package gtk

/*
#include <gtk/gtk.h>
*/
import "C"
import "unsafe"

type Button struct {
	Widget
}

func NewButton(label string) *Button {

	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))

	ptr := C.gtk_button_new_with_label(cLabel)

	return &Button{
		Widget: Widget{widget: (*C.GtkWidget)(unsafe.Pointer(ptr))},
	}

}

func (b *Button) SetLabel(label string) {

	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))

	C.gtk_button_set_label((*C.GtkButton)(unsafe.Pointer(b.widget)), cLabel)

}

func (b *Button) OnClick(fn func()) {
	connectSignal(unsafe.Pointer(b.widget), "clicked", fn)
}
