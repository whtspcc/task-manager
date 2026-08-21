package main

import (
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Sticker struct {
	widget.BaseWidget
	Note        Note
	Image       *canvas.Image
	X           float32
	Y           float32
	isAnimating bool
	onTapped    func()
	app         fyne.App
}

type CreateButton struct {
	widget.BaseWidget
	Image         *canvas.Image
	onTapped      func()
	app           fyne.App
	noteContainer *fyne.Container
}

func newSticker(app fyne.App, note Note, imagePath string, startX, startY float32) *Sticker {
	s := &Sticker{X: startX, Y: startY, app: app, Note: note}
	s.ExtendBaseWidget(s)

	s.Image = canvas.NewImageFromFile(imagePath)
	s.Image.FillMode = canvas.ImageFillContain

	s.Resize(fyne.NewSize(120, 120))
	s.Move(fyne.NewPos(s.X, s.Y))

	return s
}

func newCreateButton(app fyne.App) *CreateButton {
	cb := &CreateButton{
		Image: canvas.NewImageFromFile("createNote.png"),
		app:   app,
	}

	cb.Image.FillMode = canvas.ImageFillContain

	cb.ExtendBaseWidget(cb)

	cb.Resize(fyne.NewSize(80, 80))
	cb.Move(fyne.NewPos(50, 50))

	return cb
}

func openCreateNoteWindow(cb *CreateButton, noteContainer *fyne.Container) {
	wToCreateNote := cb.app.NewWindow("форма")
	wToCreateNote.Resize(fyne.NewSize(500, 700))

	bg := canvas.NewImageFromFile("bg.png")
	bg.FillMode = canvas.ImageFillStretch

	textEntry := widget.NewEntry()
	textEntry.SetPlaceHolder("Введите текст заметки")

	deadlineEntry := widget.NewEntry()
	deadlineEntry.SetPlaceHolder("Дедлайн замекти (кол-во минут)")

	create := widget.NewButton("Создать", func() {
		notes := readNotes()
		i := getMaxID(notes)
		i++
		countOfMinutes, _ := strconv.Atoi(deadlineEntry.Text)
		newNote := createNote(i, textEntry.Text, false, time.Now(), time.Now().Add(time.Minute*time.Duration(countOfMinutes)))

		writeNote(&newNote)

		sticker := newSticker(cb.app, newNote, "note1.png", rand.Float32()*700, rand.Float32()*400)
		noteContainer.Add(sticker)
		noteContainer.Refresh()

		wToCreateNote.Close()
	})

	atCreate := container.NewVBox(textEntry, deadlineEntry, create)

	wToCreateNote.SetContent(container.NewMax(bg, atCreate))
	wToCreateNote.Show()
}

func (cb *CreateButton) Tapped(ev *fyne.PointEvent) {
	if cb.onTapped != nil {
		cb.onTapped()
	}

	openCreateNoteWindow(cb, cb.noteContainer)
}

func (s *Sticker) Tapped(ev *fyne.PointEvent) {
	if s.onTapped != nil {
		s.onTapped()
	}

	fyne.Do(func() {
		ww := s.app.NewWindow("that's it")

		ww.Resize(fyne.NewSize(500, 700))

		bg := canvas.NewImageFromFile("bg.png")

		bg.FillMode = canvas.ImageFillStretch

		ww.SetContent(container.NewMax(bg))
		ww.Show()
	})
}

func (cb *CreateButton) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(cb.Image)
}

func (s *Sticker) CreateRenderer() fyne.WidgetRenderer {
	s.Image.Resize(s.Size())

	doneOnLabel := widget.NewLabel("[ ]")
	deadlineOnLabel := widget.NewLabel(s.Note.Deadline.Format(dateFormat))
	if s.Note.Done {
		doneOnLabel = widget.NewLabel("[x]")
	} else if s.Note.Deadline.Before(time.Now()) {
		doneOnLabel = widget.NewLabel("[!]")
	}

	textOnLabel := widget.NewLabel(s.Note.Text)
	textOnLabel.TextStyle = fyne.TextStyle{Italic: true}
	doneOnLabel.Move(fyne.NewPos(20, 30))
	textOnLabel.Move(fyne.NewPos(20, 50))
	deadlineOnLabel.Move(fyne.NewPos(20, 70))
	containerForSticker := container.NewWithoutLayout(s.Image, doneOnLabel, textOnLabel, deadlineOnLabel)

	return widget.NewSimpleRenderer(containerForSticker)
}

func (s *Sticker) Dragged(ev *fyne.DragEvent) {
	if !s.isAnimating {
		s.isAnimating = true
		go s.animateDrag()
	}
	newPos := s.Position().Add(ev.Dragged)
	s.Move(newPos)

	s.X = newPos.X
	s.Y = newPos.Y

	if s.X > 600 && s.Y > 320 {
		s.Hide()
		deleteNote(s.Note.ID)
	}
}

func (s *Sticker) animateDrag() {
	frames := []string{"note1.png", "note2.png", "note3.png", "note4.png", "note5.png", "note4.png", "note3.png", "note2.png"}
	currentFrame := 0

	for s.isAnimating {
		time.Sleep(130 * time.Millisecond)

		if !s.isAnimating {
			break
		}

		currentFrame = (currentFrame + 1) % len(frames)
		framePath := frames[currentFrame]

		fyne.DoAndWait(func() {
			if s.Image != nil {
				s.Image.File = framePath
				s.Image.Refresh()
				s.Refresh()
			}
		})
	}
}

func (s *Sticker) DragEnd() {
	s.isAnimating = false
	s.Image.File = "note1.png"
	s.Refresh()
}

func gui() {
	notes := readNotes()

	a := app.NewWithID(idFyne)
	_ = a.NewWindow("notes")

	bg := canvas.NewImageFromFile("desk.jpg")
	bg.FillMode = canvas.ImageFillStretch

	trashCan := canvas.NewImageFromFile("trash.png")
	trashCan.FillMode = canvas.ImageFillContain
	trashCan.Resize(fyne.NewSize(80, 80))
	trashCan.Move(fyne.NewPos(690, 380))

	stickers := []fyne.CanvasObject{}

	for _, note := range notes {
		sticker := newSticker(a, note, "note1.png", rand.Float32()*700, rand.Float32()*400)
		stickers = append(stickers, sticker)
	}

	noteContainer := container.NewWithoutLayout(stickers...)

	createButton := newCreateButton(a)
	createButton.noteContainer = noteContainer

	foreground := container.NewWithoutLayout(trashCan, createButton, noteContainer)

	drv, ok := a.Driver().(desktop.Driver)
	if ok {
		w := drv.CreateSplashWindow()

		w.SetContent(container.NewMax(bg, foreground))

		w.Resize(fyne.NewSize(800, 500))
		w.ShowAndRun()
	}
}
