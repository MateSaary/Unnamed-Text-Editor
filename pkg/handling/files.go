package handling

import (
	"io/ioutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

//opens a file dialog and loads the selected file's content into the editor.
func OpenFile(window fyne.Window, editor *widget.Entry) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if reader == nil { //if no file is selected
			return
		}
		defer reader.Close()

		data, err := ioutil.ReadAll(reader)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		editor.SetText(string(data))
	}, window)
}

//opens a file dialog and saves the editor's content to the selected file.
func SaveFile(window fyne.Window, editor *widget.Entry) {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		_, err = writer.Write([]byte(editor.Text)) //converts to bytes and writes to file
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
	}, window)
}

//clears the editor's content.
func ClearEditor(editor *widget.Entry) {
	editor.SetText("")
}
