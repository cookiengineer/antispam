package gtk

/*
#include <gtk/gtk.h>
*/
import "C"
import "unsafe"

type Stack struct {
	Widget
}

func NewStack() *Stack {

	ptr := C.gtk_stack_new()

	return &Stack{
		Widget: Widget{widget: (*C.GtkWidget)(unsafe.Pointer(ptr))},
	}

}

func (s *Stack) AddTitled(child unsafe.Pointer, name string, title string) *Stack {

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	C.gtk_stack_add_titled(
		(*C.GtkStack)(unsafe.Pointer(s.widget)),
		(*C.GtkWidget)(child),
		cName,
		cTitle,
	)

	return s

}

func (s *Stack) SetVisibleChildName(name string) {

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	C.gtk_stack_set_visible_child_name(
		(*C.GtkStack)(unsafe.Pointer(s.widget)),
		cName,
	)

}
