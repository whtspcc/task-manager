package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func userInputInt(scanner *bufio.Scanner) (int, error) {
	scanner.Scan()
	input, err := strconv.Atoi(scanner.Text())
	if err != nil || input < 0 {
		return 0, err
	}
	return input, nil
}

func userInputString(scanner *bufio.Scanner) (string, error) {
	scanner.Scan()
	text := strings.TrimSpace(scanner.Text())
	if text == "" {
		return "", errors.New("Текст заметки не может быть пустым")
	}
	return text, nil
}

func createNoteCommand(scanner *bufio.Scanner) {
	notes := readNotes()
	i := getMaxID(notes)
	i++

	fmt.Println("Введите содержание заметки:")
	text, err := userInputString(scanner)
	if err != nil {
		fmt.Println("Ошибка ввода содержания заметки")
		return
	}

	fmt.Println("Введите дедлайн (в минутах)")
	countOfMinutes, err := userInputInt(scanner)
	if err != nil {
		fmt.Println("Ошибка ввода дедлайна")
		return
	}

	newNote := createNote(i, text, false, time.Now(), time.Now().Add(time.Minute*time.Duration(countOfMinutes)))
	writeNote(&newNote)
}

func markNoteCommand(scanner *bufio.Scanner) {
	fmt.Println("Введите айди заметки:")
	id, err := userInputInt(scanner)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		return
	}
	toggleNote(id)
}

func deleteNoteCommand(scanner *bufio.Scanner) {
	fmt.Println("Какую заметку удалить?")
	id, err := userInputInt(scanner)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		return
	}
	deleteNote(id)
}

func editNoteCommand(scanner *bufio.Scanner) {
	fmt.Println("Введите номер заметки, которую вы хотите отредактировать:")
	id, err := userInputInt(scanner)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		return
	}
	editNote(id)
}

func searchNotesCommand(scanner *bufio.Scanner) {
	fmt.Println("Введите ключевое слово:")
	userInputString(scanner)
	searchNotes(scanner.Text())
}

func sortNotesCommand(scanner *bufio.Scanner) {
	fmt.Println("1 - сортировать по айди\n2 - сортировать по статусу: ")
	scanner.Scan()
	if scanner.Text() == "1" {
		fmt.Println("Заметки отсортированы по номеру(айди)")
		sortByID()
	} else if scanner.Text() == "2" {
		fmt.Println("Заметки отсортированы по статусу")
		sortByStatus()
	} else {
		fmt.Println("Введите 1 или 2")
	}
}

func printStatsCommand() {
	fmt.Println("Ваша статистика:")
	printStats()
}

func run() {
	scanner := bufio.NewScanner(os.Stdin)

OuterLoop:
	for {
		fmt.Println("Выберите действие:")
		fmt.Println(`1. Создать заметку
2. Показать заметки
3. Отметить заметку выполненной / снять отметку
4. Удалить заметку
5. Вывести активные заметки
6. Редактировать заметку
7. Искать заметку по кодовому слову
8. Отсортировать заметки
9. Статистика
10. Показать просроченные заметки
0. Выйти`)
		scanner.Scan()
		switch {
		case scanner.Text() == "1":
			createNoteCommand(scanner)
		case scanner.Text() == "2":
			printNotes()
		case scanner.Text() == "3":
			markNoteCommand(scanner)
		case scanner.Text() == "4":
			deleteNoteCommand(scanner)
		case scanner.Text() == "5":
			printActiveNotes()
		case scanner.Text() == "6":
			editNoteCommand(scanner)
		case scanner.Text() == "7":
			searchNotesCommand(scanner)
		case scanner.Text() == "8":
			sortNotesCommand(scanner)
		case scanner.Text() == "9":
			printStatsCommand()
		case scanner.Text() == "10":
			printExpiredNotes()
		case scanner.Text() == "0":
			break OuterLoop
		default:
			fmt.Println("Введите число из списка...")
		}
	}
}

func main() {
	gui()
}
