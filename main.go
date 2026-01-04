package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/list"
)
 
var (
	vaultDir string
	docStyle = lipgloss.NewStyle().Margin(1,2)
)

func init() {
	homeDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting home Dir", err) 
	}

	vaultDir = fmt.Sprintf("%s/.notion", homeDir)
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	newFileInput textinput.Model
	createFileInputVisible bool
	currentFile *os.File
	noteTextArea textarea.Model
	list list.Model
	showingList bool
	err error
}

func (m model) Init() tea.Cmd {
	return nil
} 

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

		case tea.WindowSizeMsg: 
			h, v := docStyle.GetFrameSize()
			m.list.SetSize(msg.Width - h, msg.Height - v - 5)

		case tea.KeyMsg:
			switch msg.String() {
				case "ctrl+c", "q":
					return m, tea.Quit

				case "esc": 
					if m.createFileInputVisible {
						m.createFileInputVisible = false
					}

					if m.currentFile != nil {
						m.noteTextArea.SetValue("")
						m.currentFile = nil
					}

					if m.showingList {
						if m.list.FilterState() == list.Filtering {
							break
						}

						m.showingList = false
					}

					return m, nil

				case "ctrl+l":
					noteList := listFiles()
					m.list.SetItems(noteList)
					m.showingList = true
					return m, nil

				case "ctrl+n" :
					m.createFileInputVisible = true
					return m, nil

				case "ctrl+s":
					if m.currentFile == nil {
						break
					}
					if err := m.currentFile.Truncate(0); err != nil {
						fmt.Println("Can't save the file 😔")
						return m, nil
					}

					if _, err := m.currentFile.Seek(0, 0); err != nil {
						fmt.Println("Can not save the file :( ")
						return m, nil
					}

					if _, err := m.currentFile.WriteString(m.noteTextArea.Value()); err != nil {
						fmt.Println("Can not save the file :( ")
						return m, nil
					}

					if err := m.currentFile.Close(); err != nil {
						fmt.Println("Can not close file.")
					}

					m.currentFile = nil
					m.noteTextArea.SetValue("")
					
					return m, nil 

				case "enter" :
					if m.currentFile != nil {
						break
					}

					if m.showingList {
						item, ok := m.list.SelectedItem().(item)
						if ok {
							filepath := fmt.Sprintf("%s/%s", vaultDir, item.title)

							content, err := os.ReadFile(filepath)
							if err != nil {
								fmt.Println("Error reading file ", err)
								return m, nil
							}

							m.noteTextArea.SetValue(string(content))
							
							f, err := os.OpenFile(filepath, os.O_RDWR, 0644)
							if err != nil {
								fmt.Println("Got error in open file ", err)
							}

							m.currentFile = f
							m.showingList = false
						}

						return m, nil
					}

					// todo : create folder 
					filename := m.newFileInput.Value()
					if filename != "" {
						filepath := fmt.Sprintf("%s/%s.md", vaultDir, filename)

						if _, err := os.Stat(filepath); err == nil {
							return m , nil
						}

						f, err := os.Create(filepath)
						if err != nil {
							log.Fatalf("%s", err)
						}

						m.currentFile = f
						m.createFileInputVisible = false
						m.newFileInput.SetValue("")
						
					}

					return m, nil
			}

    }

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
	}

	if m.showingList {
		m.list, cmd = m.list.Update(msg)
	}

    return m, cmd
}

func (m model) View() string {

	var style = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    Width(100).
	Align(lipgloss.Center)

	welcome := style.Render("Welcome to Notion App 📔")

	view := ""
	if m.createFileInputVisible {
		view = m.newFileInput.View()
	}

	if m.currentFile != nil {
		view = m.noteTextArea.View()
	}

	if m.showingList {
		view = m.list.View()
	}

	help := "Ctrl+N: new file . Ctrl+L: list . Esc: back . Ctrl+S: save . Ctrl+Q: Quit"

	return fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help) 
}

func initializeModel() model {

	err := os.MkdirAll(vaultDir, 0750)
	if err != nil {
		log.Fatal(err)
	}

	// initialize new file input
	ti := textinput.New()
	ti.Placeholder = "What would you like to call it?"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 100

	// Initialize textarea
	ta := textarea.New()
	ta.Placeholder = "Write your note here..."
	ta.ShowLineNumbers = false
	ta.Focus()

	// list
	noteList := listFiles()
	finalList := list.New(noteList, list.NewDefaultDelegate(), 0 ,0)
	finalList.Title = "All Notes 📔"
	finalList.Styles.Title = lipgloss.NewStyle().
	Foreground(lipgloss.Color("16")).
	Background(lipgloss.Color("254")).
	Padding(0, 1)

	return model{
		newFileInput: ti,
		createFileInputVisible: false, 
		noteTextArea: ta,
		list: finalList,
		err:       nil,
	}
}

func main() {

	p := tea.NewProgram(initializeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alaas there's been an error: %v", err)
		os.Exit(1)
	}
}

func listFiles() []list.Item {
	items := make([]list.Item, 0)

	entries, err := os.ReadDir(vaultDir)
	if err != nil {
		fmt.Println("Error reading list", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				fmt.Println("Got error in file info ", err)
				continue
			}

			modeTime := info.ModTime().Format("2006-01-02 15:04")
			
			items = append(items, item{
				title: entry.Name(),
				desc: fmt.Sprintf("Modified : %s", modeTime),
			})
		}
	}

	return items
}