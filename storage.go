package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func writeNote(note *Note) {
	notes := readNotes()
	notes = append(notes, *note)
	writeNotes(notes)
}

func readNotes() []Note {
	data, err := os.ReadFile(notesFile)
	if err != nil {
		fmt.Println("Не удалось открыть файл", err)
		return nil
	}

	if len(data) == 0 {
		return []Note{}
	}

	notes := []Note{}

	err = json.Unmarshal(data, &notes)
	if err != nil {
		fmt.Println("Ошибка сериализации: %v\n", err)
		return nil
	}

	return notes
}

func getMaxID(notes []Note) int {
	var maxID = 0

	for _, note := range notes {
		if note.ID > maxID {
			maxID = note.ID
		}
	}
	return maxID
}

func writeNotes(notes []Note) {
	jsonData, err := json.MarshalIndent(notes, "", " ")
	if err != nil {
		fmt.Println("Ошибка при сериализации: %v", err)
		return
	}

	_ = os.WriteFile("notes.json", jsonData, 0644)
}
