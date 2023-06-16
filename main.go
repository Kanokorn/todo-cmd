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

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func main() {
	rootCmd.AddCommand(addTodoCmd)
	rootCmd.Execute()
}

// // go run main.go add "ล้างจาน"
// >> added "ล้างจาน"

// // go run main.go add "ซักผ้า"
// >> added "ซักผ้า"

// // go run main.go done
// >> please select item to mark as done
// 1. ล้างจาน
// 2. ซักผ้า
// >> 1
// >> "ล้างจาน" is done.

// // go run main.go filter --done
// 1. ล้างจาน

// // go run main.go filter --not-done
// 1. ซักผ้า
