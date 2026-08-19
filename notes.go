package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Note struct {
	ID        int
	Text      string
	Done      bool
	CreatedAt time.Time
	Deadline  time.Time
}

func createNote(index int, text string, done bool, createdAt time.Time, deadline time.Time) Note {
	return Note{ID: index, Text: text, Done: done, CreatedAt: createdAt, Deadline: deadline}
}

func printNotes() {
	var notes = readNotes()

	for _, note := range notes {
		printNote(note)
	}
}

func printNote(note Note) {
	status := "[ ]"
	if note.Done {
		status = "[x]"
	} else if note.Deadline.Before(time.Now()) {
		status = "[!]"
	}

	fmt.Printf("%s %d. %s\nСоздано: %v.\nДедлайн: %v\n",
		status,
		note.ID,
		note.Text,
		note.CreatedAt.Format(dateFormat),
		note.Deadline.Format(dateFormat))
}

func toggleNote(id int) {
	var notes = readNotes()
	var found = false

	for i := range notes {
		if notes[i].ID == id {
			if !notes[i].Done {
				fmt.Println("Было [ ], ставлю [x]. Выполнено!\n")
				found = true
			} else {
				fmt.Println("Было [x], ставлю [ ]. Перестало быть выполненым?\n")
			}
			notes[i].Done = !notes[i].Done
			found = true
		}
	}

	if !found {
		fmt.Println("Такой заметки нет\n")
		return
	}

	writeNotes(notes)
}

func deleteNote(id int) {
	var notes = readNotes()
	var found = false

	for i := range notes {
		if notes[i].ID == id {
			fmt.Printf("Удаляю заметку номер %d...\n", id)
			notes = append(notes[:i], notes[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("Заметки с айди %d не нашлось\n", id)
		return
	}

	writeNotes(notes)
}

func editNote(id int) {
	var notes = readNotes()
	var found = false

	scanner := bufio.NewScanner(os.Stdin)

	for i := range notes {
		if notes[i].ID == id {
			fmt.Println("Введите новый текст:\n")
			scanner.Scan()
			notes[i].Text = scanner.Text()
			found = true
		}
	}

	if !found {
		fmt.Println("Такой заметки нет\n")
		return
	}

	writeNotes(notes)
}

func printActiveNotes() {
	var notes = readNotes()
	var found = false

	for _, note := range notes {
		if !note.Done {
			fmt.Printf("%d %v %s\n", note.ID, note.Done, note.Text)
			found = true
		}
	}

	if !found {
		fmt.Println("Все заметки уже выполнены\n")
	}
}

func searchNotes(query string) {
	var notes = readNotes()
	var found = false

	for _, note := range notes {
		if strings.Contains(strings.ToLower(note.Text), strings.ToLower(query)) {
			printNote(note)
			found = true
		}
	}

	if !found {
		fmt.Println("Ничего не найдено")
	}
}

func printStats() {
	var notes = readNotes()

	var general = 0
	var completed = 0
	var unfulfilled = 0

	for _, note := range notes {
		general++
		if note.Done {
			completed++
		} else {
			unfulfilled++
		}
	}
	fmt.Printf("\nВсего задач: %d\nВыполнено: %d \nОсталось: %d\n\n", general, completed, unfulfilled)
}

func sortByID() {
	var notes = readNotes()
	sort.Slice(notes, func(i, j int) bool { return notes[i].ID < notes[j].ID })

	for _, note := range notes {
		printNote(note)
	}
}

func sortByStatus() {
	var notes = readNotes()
	sort.Slice(notes, func(i, j int) bool { return notes[i].Done == true && notes[j].Done == false })

	for _, note := range notes {
		printNote(note)
	}
}

func printExpiredNotes() {
	notes := readNotes()
	var found = false

	for _, note := range notes {
		if !note.Done && note.Deadline.Before(time.Now()) {
			printNote(note)
			found = true
		}
	}

	if !found {
		fmt.Println("Просроченных заметок нет")
	}
}
