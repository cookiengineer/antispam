package gtk

/*
#include <gtk/gtk.h>
*/
import "C"
import "unsafe"

type StackSidebar struct {
	Widget
}

func NewStackSidebar() *StackSidebar {

	ptr := C.gtk_stack_sidebar_new()

	return &StackSidebar{
		Widget: Widget{widget: (*C.GtkWidget)(unsafe.Pointer(ptr))},
	}

}

func (s *StackSidebar) SetStack(stack *Stack) {

	C.gtk_stack_sidebar_set_stack(
		(*C.GtkStackSidebar)(unsafe.Pointer(s.widget)),
		(*C.GtkStack)(unsafe.Pointer(stack.widget)),
	)

}
