package dashboard

import (
	"github.com/rivo/tview"
)

func ui() {
	box := tview.NewBox().SetBorder(true).SetTitle("Sidearm")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}
