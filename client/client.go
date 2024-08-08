package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/MuhammadAmrMahmoud/grpc-todo-app/pb"
	"google.golang.org/grpc"
)

func main() {
	dsn := "localhost:8080"
	conn, err := grpc.Dial(dsn, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewTodoServiceClient(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Enter a command (add, list, get, update, update-status, delete, exit):")
		scanner.Scan()
		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) == 0 {
			continue
		}

		switch args[0] {

		case "add":
			if len(args) < 3 {
				fmt.Println("Usage: add <title> <description>	(Same Format No Spaces)")
				continue
			}
			addTodo(client, args[1], args[2])

		case "list":
			listTodos(client)

		case "get":
			if len(args) != 2 {
				fmt.Println("Usage: get <id>")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid ID:", args[1])
				continue
			}
			getTodoByID(client, int32(id))

		case "update":
			if len(args) < 5 {
				fmt.Println("Usage: update <id> <title> <description> 		(Same Format No Spaces)")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid ID:", args[1])
				continue
			}
			completed, err := strconv.ParseBool(args[4])
			if err != nil {
				fmt.Println("Invalid completed status:", args[4])
				continue
			}
			updateTodo(client, int32(id), args[2], args[3], completed)

		case "update-status":
			if len(args) < 3 {
				fmt.Println("Usage: update <id> <true-false>	(Same Format No Spaces)")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid ID:", args[1])
				continue
			}
			completed, err := strconv.ParseBool(args[2])
			if err != nil {
				fmt.Println("Invalid completed status:", args[2])
				continue
			}
			updateTodoStatus(client, int32(id), completed)

		case "delete":
			if len(args) != 2 {
				fmt.Println("Usage: delete <id>")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid ID:", args[1])
				continue
			}
			deleteTodoByID(client, int32(id))

		case "exit":
			return
		default:
			fmt.Println("Unknown command:", args[0])
		}
	}
}

func addTodo(client pb.TodoServiceClient, title, description string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.AddTodo(ctx, &pb.AddTodoRequest{Title: title, Description: description})
	if err != nil {
		log.Fatalf("could not add todo: %v", err)
	}
	fmt.Printf("Added Todo: %v\n", r.Todo)
}

func listTodos(client pb.TodoServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.ListTodos(ctx, &pb.ListTodosRequest{})
	if err != nil {
		log.Fatalf("could not list todos: %v", err)
	}
	fmt.Println("Todos:")
	for _, todo := range r.Todos {
		completionStatus := "Not Completed"
		if todo.Completed {
			completionStatus = "Completed"
		}
		fmt.Printf("ID: %d, Title: %s, Description: %s, Status: %s\n",
			todo.Id, todo.Title, todo.Description, completionStatus)
	}
}

func getTodoByID(client pb.TodoServiceClient, id int32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.GetTodo(ctx, &pb.GetTodoRequest{Id: id})
	if err != nil {
		log.Fatalf("could not get todo: %v", err)
	}
	fmt.Printf("Got Todo: %v\n", r.Todo)
}

func updateTodo(client pb.TodoServiceClient, id int32, title, description string, completed bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.UpdateTodo(ctx, &pb.UpdateTodoRequest{
		Id:          id,
		Title:       title,
		Description: description,
		Completed:   completed,
	})
	if err != nil {
		log.Fatalf("could not update todo: %v", err)
	}
	fmt.Printf("Updated Todo: %v\n", r.Todo)
}

func updateTodoStatus(client pb.TodoServiceClient, id int32, completed bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.UpdateTodoStatus(ctx, &pb.UpdateTodoStatusRequest{
		Id:        id,
		Completed: completed,
	})
	if err != nil {
		log.Fatalf("could not update todo status: %v", err)
	}
	fmt.Printf("Updated Todo status: %v\n", r.Todo)
}

func deleteTodoByID(client pb.TodoServiceClient, id int32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.DeleteTodo(ctx, &pb.DeleteTodoRequest{Id: id})
	if err != nil {
		log.Fatalf("could not delete todo: %v", err)
	}
	fmt.Printf("Deleted Todo: %v\n", r.Success)
}
