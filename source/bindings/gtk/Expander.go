package gtk

/*
#include <gtk/gtk.h>
*/
import "C"
import "unsafe"

type Expander struct {
	Widget
}

func NewExpander(label string) *Expander {

	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))

	ptr := C.gtk_expander_new(cLabel)

	return &Expander{
		Widget: Widget{widget: (*C.GtkWidget)(unsafe.Pointer(ptr))},
	}

}

func (e *Expander) SetChild(child unsafe.Pointer) {
	C.gtk_expander_set_child(
		(*C.GtkExpander)(unsafe.Pointer(e.widget)),
		(*C.GtkWidget)(child),
	)
}

func (e *Expander) SetExpanded(expanded bool) {
	if expanded {
		C.gtk_expander_set_expanded((*C.GtkExpander)(unsafe.Pointer(e.widget)), C.TRUE)
	} else {
		C.gtk_expander_set_expanded((*C.GtkExpander)(unsafe.Pointer(e.widget)), C.FALSE)
	}
}

func (e *Expander) SetLabel(label string) {

	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))

	C.gtk_expander_set_label((*C.GtkExpander)(unsafe.Pointer(e.widget)), cLabel)

}

func (e *Expander) SetLabelWidget(child unsafe.Pointer) {
	C.gtk_expander_set_label_widget(
		(*C.GtkExpander)(unsafe.Pointer(e.widget)),
		(*C.GtkWidget)(child),
	)
}
