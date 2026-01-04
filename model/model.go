package model

import (
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/list"
)

type model struct {
	newFileInput textinput.Model
	createFileInputVisible bool
	currentFile *os.File
	noteTextArea textarea.Model
	list list.Model
	showingList bool
	err error
}