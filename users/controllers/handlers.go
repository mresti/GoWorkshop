package controllers

import (
	"fmt"
	"log"

	pbList "github.com/wizelineacademy/GoWorkshop/proto/list"
	pb "github.com/wizelineacademy/GoWorkshop/proto/users"
	"github.com/wizelineacademy/GoWorkshop/shared/config"
	"github.com/wizelineacademy/GoWorkshop/users/models"
	"golang.org/x/net/context"
)

// Service struct
type Service struct{}

// CreateUser implementation
func (s *Service) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (response *pb.CreateUserResponse, err error) {
	user := &models.User{
		Name:  in.Name,
		Email: in.Email,
	}

	appContext := config.NewContext()
	defer appContext.Close()

	c := appContext.DbCollection("users")
	repo := &models.UserRepository{
		C: c,
	}
	var userID string
	userID, err = repo.Create(user)

	response = new(pb.CreateUserResponse)
	if err == nil {
		log.Printf("[user.Create] New user ID: %s", userID)

		// Create initial item in todo list
		_, listErr := appContext.ListService.CreateItem(context.Background(), &pbList.CreateItemRequest{
			Message: "Welcome to Workshop!",
			UserId:  userID,
		})
		if listErr != nil {
			log.Printf("[user.Create] Cannot create item: %v", listErr)
		}

		response.Message = fmt.Sprintf("User created successfully, ID: %s", userID)
		response.Code = 200
	} else {
		response.Message = err.Error()
		response.Code = 500
	}

	return response, err
}

// GetAllUsers implementation
func (s *Service) GetAllUsers(ctx context.Context, in *pb.GetAllUsersRequest) (response *pb.GetAllUsersResponse, err error) {
	appContext := config.NewContext()
	defer appContext.Close()

	c := appContext.DbCollection("users")
	repo := &models.UserRepository{
		C: c,
	}
	users := repo.GetAll()

	pbUsers := []*pb.User{}
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:    user.Id.Hex(),
			Name:  user.Name,
			Email: user.Email,
		})
	}
	response = &pb.GetAllUsersResponse{
		Users: pbUsers,
		Code:  200,
	}

	return response, err
}
