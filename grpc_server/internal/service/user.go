package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/antoniofmoliveira/courses/db/database"
	"github.com/antoniofmoliveira/courses/dto"
	"github.com/antoniofmoliveira/courses/entity"
	"github.com/antoniofmoliveira/courses/grpcproto/pb"
	"github.com/go-chi/jwtauth"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	db database.UserRepositoryInterface
}

func NewUserService(db database.UserRepositoryInterface) *UserService {
	return &UserService{
		db: db,
	}
}

func (u *UserService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.User, error) {
	user := dto.UserInputDto{Name: in.Name, Email: in.Email, Password: in.Password}

	entityUser, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		return nil, err
	}
	user = dto.UserInputDto{
		ID:       entityUser.ID,
		Name:     entityUser.Name,
		Email:    entityUser.Email,
		Password: entityUser.Password,
	}

	userOutputDto, err := u.db.Create(user)
	if err != nil {
		return nil, err
	}
	return &pb.User{Id: userOutputDto.ID, Name: userOutputDto.Name, Email: userOutputDto.Email}, nil
}

func (u *UserService) GetUser(ctx context.Context, in *pb.UserGetRequest) (*pb.User, error) {
	userOutputDto, err := u.db.Find(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.User{Id: userOutputDto.ID, Name: userOutputDto.Name, Email: userOutputDto.Email}, nil
}

func (u *UserService) ListUsers(ctx context.Context, in *pb.Blank) (*pb.Users, error) {
	usersOutputDto, err := u.db.FindAll()
	if err != nil {
		return nil, err
	}
	pbUsers := []*pb.User{}
	for _, user := range usersOutputDto.Users {
		pbUsers = append(pbUsers, &pb.User{Id: user.ID, Name: user.Name, Email: user.Email})
	}
	return &pb.Users{Users: pbUsers}, nil
}

func (u *UserService) UpdateUser(ctx context.Context, in *pb.UserUpdateRequest) (*pb.User, error) {
	user := dto.UserInputDto{
		ID:       in.Id,
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	}
	err := u.db.Update(user)
	if err != nil {
		return nil, err
	}
	return &pb.User{Id: user.ID, Name: user.Name, Email: user.Email}, nil
}

func (u *UserService) DeleteUser(ctx context.Context, in *pb.UserDeleteRequest) (*pb.Response, error) {
	err := u.db.Delete(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Response{IsSuccess: true, Message: "User deleted successfully"}, nil
}

func (u *UserService) GetJWTToken(ctx context.Context, in *pb.UserForJWT) (*pb.JWTToken, error) {
	userCredentials := dto.UserInputDto{
		Email:    in.Email,
		Password: in.Password,
	}

	userFromDB, err := u.db.FindByEmail(userCredentials.Email)
	if err != nil {
		return nil, err
	}

	entityUser := entity.User{
		Email:    userFromDB.Email,
		Password: userFromDB.Password,
	}

	if !entityUser.ValidatePassword(userCredentials.Password) {
		slog.Error("GetJWT", "msg", err)
		return nil, err
	}

	jwt := ctx.Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := ctx.Value("jwtExpiresIn").(int)
	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": entityUser.Email,
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	return &pb.JWTToken{Token: tokenString}, nil
}
