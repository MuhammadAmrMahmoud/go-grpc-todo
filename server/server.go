package main

import (
	"context"
	"log"
	"net"

	"github.com/MuhammadAmrMahmoud/grpc-todo-app/models"
	pb "github.com/MuhammadAmrMahmoud/grpc-todo-app/pb"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedTodoServiceServer
	db *gorm.DB
}

func (s *server) AddTodo(ctx context.Context, req *pb.AddTodoRequest) (*pb.AddTodoResponse, error) {
	todo := models.Todo{
		Title:       req.Title,
		Description: req.Description,
	}
	err := s.db.Create(&todo).Error
	if err != nil {
		return nil, err
	}

	return &pb.AddTodoResponse{Todo: &pb.Todo{
		Id:          int32(todo.ID),
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}}, nil
}

func (s *server) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	var todos []models.Todo
	err := s.db.Find(&todos).Error
	if err != nil {
		return nil, err
	}

	var pbTodos []*pb.Todo
	for _, t := range todos {
		pbTodos = append(pbTodos, &pb.Todo{
			Id:          int32(t.ID),
			Title:       t.Title,
			Description: t.Description,
			Completed:   t.Completed,
		})
	}

	return &pb.ListTodosResponse{Todos: pbTodos}, nil
}

func (s *server) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoResponse, error) {
	var todo models.Todo
	if err := s.db.First(&todo, req.Id).Error; err != nil {
		return nil, err
	}

	return &pb.GetTodoResponse{Todo: &pb.Todo{
		Id:          int32(todo.ID),
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}}, nil
}

func (s *server) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.UpdateTodoResponse, error) {
	var todo models.Todo
	if err := s.db.First(&todo, req.Id).Error; err != nil {
		return nil, err
	}

	todo.Title = req.Title
	todo.Description = req.Description
	todo.Completed = req.Completed

	if err := s.db.Save(&todo).Error; err != nil {
		return nil, err
	}

	return &pb.UpdateTodoResponse{Todo: &pb.Todo{
		Id:          int32(todo.ID),
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}}, nil
}

func (s *server) UpdateTodoStatus(ctx context.Context, req *pb.UpdateTodoStatusRequest) (*pb.UpdateTodoStatusResponse, error) {
	var todo models.Todo
	if err := s.db.First(&todo, req.Id).Error; err != nil {
		return nil, err
	}

	todo.Completed = req.Completed

	if err := s.db.Save(&todo).Error; err != nil {
		return nil, err
	}

	return &pb.UpdateTodoStatusResponse{Todo: &pb.Todo{
		Id:        int32(todo.ID),
		Completed: todo.Completed,
	}}, nil
}

func (s *server) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	if err := s.db.Delete(&models.Todo{}, req.Id).Error; err != nil {
		return nil, err
	}

	return &pb.DeleteTodoResponse{Success: true}, nil
}

func main() {
	db := models.InitDB()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(grpcServer, &server{db: db})

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
