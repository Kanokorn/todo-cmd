package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type Todo struct {
	ID          int
	Description string
	Done        string
}

var addTodoCmd = &cobra.Command{
	Use:   "add",
	Short: "add todo item",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := addTodo(args[0])
		if err != nil {
			log.Fatalln("error add todo:", err)
		}
	},
}

var filterTodoCmd = &cobra.Command{
	Use:   "filter",
	Short: "filter todo item",
	Run: func(cmd *cobra.Command, args []string) {
		err := filterTodo()
		if err != nil {
			log.Fatalln("error filter todo:", err)
		}
	},
}

var markAsDoneCmd = &cobra.Command{
	Use:   "mark",
	Short: "mark todo item as done",
	Run: func(cmd *cobra.Command, args []string) {
		err := markAsDone()
		if err != nil {
			log.Fatalln("error mark as done:", err)
		}
	},
}

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Todo CLI")
	},
}

var done bool

func main() {
	filterTodoCmd.Flags().BoolVarP(&done, "done", "d", true, "filter done todo")
	rootCmd.AddCommand(filterTodoCmd)
	rootCmd.AddCommand(markAsDoneCmd)
	rootCmd.AddCommand(addTodoCmd)
	rootCmd.Execute()
}

func addTodo(description string) error {
	file, err := os.OpenFile("todo.csv", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	var todos []*Todo
	for _, record := range records {
		var todo Todo
		id, _ := strconv.Atoi(record[0])
		todo.ID = id
		todo.Description = record[1]
		todo.Done = record[2]

		todos = append(todos, &todo)
	}

	todos = append(todos, &Todo{
		ID:          len(todos) + 1,
		Description: description,
		Done:        "No",
	})

	file, err = os.OpenFile("todo.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	var data [][]string
	for _, todo := range todos {
		row := []string{
			strconv.Itoa(todo.ID),
			todo.Description,
			todo.Done,
		}

		data = append(data, row)
	}
	w.WriteAll(data)
	fmt.Printf("added %q\n", description)

	return nil
}

func filterTodo() error {
	file, err := os.OpenFile("todo.csv", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	doneTodos := make([]Todo, 0)
	notDoneTodos := make([]Todo, 0)
	for _, record := range records {
		var todo Todo
		id, _ := strconv.Atoi(record[0])
		todo.ID = id
		todo.Description = record[1]
		todo.Done = record[2]

		if todo.Done == "Yes" {
			doneTodos = append(doneTodos, todo)
		} else {
			notDoneTodos = append(notDoneTodos, todo)
		}
	}

	if done {
		if len(doneTodos) == 0 {
			fmt.Println("no todo item")
			return nil
		}

		for _, todo := range doneTodos {

			if todo.Done == "Yes" {
				fmt.Println(todo.ID, todo.Description)
			}
		}
	} else {
		if len(notDoneTodos) == 0 {
			fmt.Println("no todo item")
			return nil
		}

		for _, todo := range notDoneTodos {
			if todo.Done == "No" {
				fmt.Println(todo.ID, todo.Description)
			}
		}
	}

	return nil
}

func markAsDone() error {
	file, err := os.OpenFile("todo.csv", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	var todos []*Todo
	for _, record := range records {
		var todo Todo
		id, _ := strconv.Atoi(record[0])
		todo.ID = id
		todo.Description = record[1]
		todo.Done = record[2]

		todos = append(todos, &todo)
	}

	for _, todo := range todos {
		if todo.Done == "No" {
			id := strconv.Itoa(todo.ID)
			fmt.Printf("ID: %s, Desciption: %s\n", id, todo.Description)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Please select todo item (ID) to mark as done: ")
	scanner.Scan()

	for _, todo := range todos {
		if scanner.Text() == strconv.Itoa(todo.ID) {
			todo.Done = "Yes"
			fmt.Printf("mark '%s' as done\n", todo.Description)
		}
	}

	file, err = os.OpenFile("todo.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	var data [][]string
	for _, todo := range todos {
		row := []string{
			strconv.Itoa(todo.ID),
			todo.Description,
			todo.Done,
		}

		data = append(data, row)
	}
	w.WriteAll(data)

	return nil
}
