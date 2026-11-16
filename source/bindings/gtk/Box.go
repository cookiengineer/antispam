package gtk

/*
#include <gtk/gtk.h>
*/
import "C"
import "unsafe"

const (
	OrientationHorizontal = 0
	OrientationVertical   = 1
)

type Box struct {
	Widget
}

func NewBox(orientation int, spacing int) *Box {

	ptr := C.gtk_box_new(C.GtkOrientation(orientation), C.int(spacing))

	return &Box{
		Widget: Widget{widget: (*C.GtkWidget)(unsafe.Pointer(ptr))},
	}

}

func (b *Box) Append(child unsafe.Pointer) {
	C.gtk_box_append((*C.GtkBox)(unsafe.Pointer(b.widget)), (*C.GtkWidget)(child))
}

func (b *Box) SetSpacing(spacing int) {
	C.gtk_box_set_spacing((*C.GtkBox)(unsafe.Pointer(b.widget)), C.int(spacing))
}

func (b *Box) Clear() {

	child := C.gtk_widget_get_first_child(b.widget)

	for child != nil {
		next := C.gtk_widget_get_next_sibling(child)
		C.gtk_box_remove((*C.GtkBox)(unsafe.Pointer(b.widget)), child)
		child = next
	}

}
