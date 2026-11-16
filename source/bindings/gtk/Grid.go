package gtk

/*
#include <gtk/gtk.h>
*/
import "C"
import "unsafe"

type Grid struct {
	Widget
}

func NewGrid() *Grid {

	ptr := C.gtk_grid_new()

	return &Grid{
		Widget: Widget{widget: (*C.GtkWidget)(unsafe.Pointer(ptr))},
	}

}

func (g *Grid) Attach(child unsafe.Pointer, column int, row int, width int, height int) {
	C.gtk_grid_attach(
		(*C.GtkGrid)(unsafe.Pointer(g.widget)),
		(*C.GtkWidget)(child),
		C.int(column),
		C.int(row),
		C.int(width),
		C.int(height),
	)
}

func (g *Grid) SetColumnSpacing(spacing uint) {
	C.gtk_grid_set_column_spacing((*C.GtkGrid)(unsafe.Pointer(g.widget)), C.guint(spacing))
}

func (g *Grid) SetRowSpacing(spacing uint) {
	C.gtk_grid_set_row_spacing((*C.GtkGrid)(unsafe.Pointer(g.widget)), C.guint(spacing))
}
