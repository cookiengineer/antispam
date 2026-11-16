package gtk

/*
#include <gtk/gtk.h>
*/
import "C"
import "unsafe"

type ListBox struct {
	Widget
}

func NewListBox() *ListBox {

	ptr := C.gtk_list_box_new()

	C.gtk_list_box_set_selection_mode((*C.GtkListBox)(unsafe.Pointer(ptr)), C.GTK_SELECTION_SINGLE)

	return &ListBox{
		Widget: Widget{widget: (*C.GtkWidget)(unsafe.Pointer(ptr))},
	}

}

func (l *ListBox) Append(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	label := C.gtk_label_new(cText)
	C.gtk_widget_set_visible(label, C.TRUE)
	C.gtk_list_box_append((*C.GtkListBox)(unsafe.Pointer(l.widget)), label)
}

func (l *ListBox) AppendRow(row *ListBoxRow) {
	C.gtk_list_box_append(
		(*C.GtkListBox)(unsafe.Pointer(l.widget)),
		row.widget,
	)
}

func (l *ListBox) Clear() {
	box := (*C.GtkListBox)(unsafe.Pointer(l.widget))

	row := C.gtk_widget_get_first_child(l.widget)
	for row != nil {
		next := C.gtk_widget_get_next_sibling(row)
		C.gtk_list_box_remove(box, row)
		row = next
	}
}

func (l *ListBox) OnRowActivated(fn func()) {
	connectListBoxSignal(unsafe.Pointer(l.widget), fn)
}

func (l *ListBox) OnRowActivatedIndex(fn func(int)) {
	connectListBoxIndexSignal(unsafe.Pointer(l.widget), fn)
}

func (l *ListBox) SelectRow(index int) {

	row := C.gtk_list_box_get_row_at_index((*C.GtkListBox)(unsafe.Pointer(l.widget)), C.int(index))

	if row != nil {
		C.gtk_list_box_select_row((*C.GtkListBox)(unsafe.Pointer(l.widget)), row)
	}

}
