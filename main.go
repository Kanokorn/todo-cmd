package main

import (
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
	Use:  "add todo item",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.OpenFile("todo.csv", os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalln("error open file:", err)
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
			Description: args[0],
			Done:        "No",
		})

		file, err = os.OpenFile("todo.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatalln("error open file:", err)
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

			fmt.Println(row)
			data = append(data, row)
		}
		w.WriteAll(data)

		fmt.Printf("added %q\n", args[0])
	},
}

var filterTodoCmd = &cobra.Command{
	Use: "filter todo item",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.OpenFile("todo.csv", os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalln("error open file:", err)
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

		if done {
			for _, todo := range todos {
				if todo.Done == "Yes" {
					fmt.Println(todo.ID, todo.Description)
				}
			}
		} else {
			for _, todo := range todos {
				if todo.Done == "No" {
					fmt.Println(todo.ID, todo.Description)
				}
			}
		}
	},
}

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var done bool

func main() {
	rootCmd.AddCommand(addTodoCmd)

	filterTodoCmd.Flags().BoolVarP(&done, "done", "d", true, "filter done todo")
	rootCmd.AddCommand(filterTodoCmd)
	rootCmd.Execute()
}
