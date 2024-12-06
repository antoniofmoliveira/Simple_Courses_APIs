package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/antoniofmoliveira/courses/db/database"
	"github.com/antoniofmoliveira/courses/dto"
	"github.com/antoniofmoliveira/courses/entity"
	"github.com/antoniofmoliveira/courses/flatbuffersapi/fb"
	"github.com/go-chi/jwtauth"
	flatbuffers "github.com/google/flatbuffers/go"
)

type UserHandler struct {
	UserRepository database.UserRepositoryInterface
}

func NewUserHandler(userRepository database.UserRepositoryInterface) *UserHandler {
	return &UserHandler{
		UserRepository: userRepository,
	}
}

func (u *UserHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	if r.Header.Get("Accept") != octetStream {
		slog.Error("FindUser", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	id := r.PathValue("id")

	user, err := u.UserRepository.Find(id)
	if err != nil {
		slog.Error("FindUser", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	buf := userAsBytes(user)
	w.Write(*buf)

	slog.Info("FindUser", "msg", "user found", "id", id)
}

func userAsBytes(userOutput dto.UserOutputDto) *[]byte {
	bb := flatbuffers.NewBuilder(0)
	id := bb.CreateString(userOutput.ID)
	name := bb.CreateString(userOutput.Name)
	email := bb.CreateString(userOutput.Email)
	fb.UserOutputStart(bb)
	fb.UserOutputAddId(bb, id)
	fb.UserOutputAddName(bb, name)
	fb.UserOutputAddEmail(bb, email)
	fbUserOutput := fb.UserOutputEnd(bb)
	bb.Finish(fbUserOutput)
	buf := bb.FinishedBytes()
	return &buf
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	defer func() {
		if err := recover(); err != nil {
			slog.Info("createUser", "msg", "unexpected payload")
			slog.Error("createUser", "msg", err)
			sendFlatBufferMessage(w, err.(error).Error(), http.StatusInternalServerError)
		}
	}()

	if r.Header.Get("Content-Type") != octetStream {
		slog.Error("createUser", "msg", "invalid content type")
		sendFlatBufferMessage(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	if r.Header.Get("Accept") != octetStream {
		slog.Error("createUser", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("createUser", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fbUserInput := fb.GetRootAsUserInput(body, 0)

	userInputDto := dto.UserInputDto{
		Name:     string(fbUserInput.Name()),
		Email:    string(fbUserInput.Email()),
		Password: string(fbUserInput.Password()),
	}

	_, err = u.UserRepository.FindByEmail(userInputDto.Email)
	if err == nil {
		slog.Error("createUser", "msg", "user already exists")
		sendFlatBufferMessage(w, "user already exists", http.StatusConflict)
		return
	}

	entityUser, err := entity.NewUser(userInputDto.Name, userInputDto.Email, userInputDto.Password)
	if err != nil {
		slog.Error("createUser", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	userInputDto = dto.UserInputDto{
		ID:       entityUser.ID,
		Name:     entityUser.Name,
		Email:    entityUser.Email,
		Password: entityUser.Password,
	}

	user, err := u.UserRepository.Create(userInputDto)

	if err != nil {
		slog.Error("createUser", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	buf := userAsBytes(user)
	w.Write(*buf)

	slog.Info("createUser", "msg", "user created", "id", user.ID)
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	if r.Header.Get("Accept") != octetStream {
		slog.Error("UpdateUser", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	id := r.PathValue("id")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("UpdateUser", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fbUserInput := fb.GetRootAsUserInput(body, 0)

	userInputDto := dto.UserInputDto{
		ID:       id,
		Name:     string(fbUserInput.Name()),
		Email:    string(fbUserInput.Email()),
		Password: string(fbUserInput.Password()),
	}

	entityUser, err := entity.NewUser(userInputDto.Name, userInputDto.Email, userInputDto.Password)
	if err != nil {
		slog.Error("UpdateUser", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	userInputDto = dto.UserInputDto{
		ID:       id,
		Name:     entityUser.Name,
		Email:    entityUser.Email,
		Password: entityUser.Password,
	}

	err = u.UserRepository.Update(userInputDto)

	if err != nil {
		slog.Error("UpdateUser", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendFlatBufferMessage(w, "user updated", http.StatusOK)

	slog.Info("UpdateUser", "msg", "user updated", "id", id)
}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	if r.Header.Get("Accept") != octetStream {
		slog.Error("DeleteUser", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	id := r.PathValue("id")

	err := u.UserRepository.Delete(id)
	if err != nil {
		slog.Error("DeleteUser", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendFlatBufferMessage(w, "user deleted", http.StatusOK)

	slog.Info("DeleteUser", "msg", "user deleted", "id", id)
}

func (u *UserHandler) FindAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	if r.Header.Get("Accept") != octetStream {
		slog.Error("FindAllUsers", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	users, err := u.UserRepository.FindAll()
	if err != nil {
		slog.Error("FindAllUsers", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fbBuilder := flatbuffers.NewBuilder(0)

	elements := usersAsFlatBufferVector(fbBuilder, &users.Users)

	fb.UserListOutputStartElementsVector(fbBuilder, len(*elements))
	for i := len(*elements) - 1; i >= 0; i-- {
		fbBuilder.PrependUOffsetT((*elements)[i])
	}
	vec := fbBuilder.EndVector(len(*elements))

	fb.UserListOutputStart(fbBuilder)
	fb.UserListOutputAddElements(fbBuilder, vec)
	fbUsersOutput := fb.UserListOutputEnd(fbBuilder)
	fbBuilder.Finish(fbUsersOutput)

	w.WriteHeader(http.StatusOK)
	w.Write(fbBuilder.FinishedBytes())

	slog.Info("FindAllUsers", "msg", "users found")
}

func usersAsFlatBufferVector(fbBuilder *flatbuffers.Builder, users *[]dto.UserOutputDto) *[]flatbuffers.UOffsetT {
	var elements []flatbuffers.UOffsetT
	for _, user := range *users {
		id := fbBuilder.CreateString(user.ID)
		name := fbBuilder.CreateString(user.Name)
		email := fbBuilder.CreateString(user.Email)
		fb.UserOutputStart(fbBuilder)
		fb.UserOutputAddId(fbBuilder, id)
		fb.UserOutputAddName(fbBuilder, name)
		fb.UserOutputAddEmail(fbBuilder, email)
		fbUserOutput := fb.UserOutputEnd(fbBuilder)
		elements = append(elements, fbUserOutput)
	}
	return &elements
}

func (u *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	if r.Header.Get("Content-Type") != octetStream {
		slog.Error("GetJWT", "msg", "invalid content type")
		sendFlatBufferMessage(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	if r.Header.Get("Accept") != octetStream {
		slog.Error("GetJWT", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("GetJWT", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fbUserCredentials := fb.GetRootAsUserCredentials(body, 0)

	userCredentials := dto.UserInputDto{
		Email:    string(fbUserCredentials.Email()),
		Password: string(fbUserCredentials.Password()),
	}

	userFromDB, err := u.UserRepository.FindByEmail(userCredentials.Email)
	if err != nil {
		slog.Error("GetJWT", "msg", err)
		sendFlatBufferMessage(w, "invalid credentials", http.StatusInternalServerError)
		return
	}

	entityUser := entity.User{
		ID:       "",
		Name:     "",
		Email:    userFromDB.Email,
		Password: userFromDB.Password,
	}

	if !entityUser.ValidatePassword(userCredentials.Password) {
		slog.Error("GetJWT", "msg", err)
		sendFlatBufferMessage(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)
	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": entityUser.Email,
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	bb := flatbuffers.NewBuilder(0)
	fbToken := bb.CreateString(tokenString)
	fb.JWTTokenStart(bb)
	fb.JWTTokenAddToken(bb, fbToken)
	fbJWTToken := fb.JWTTokenEnd(bb)
	bb.Finish(fbJWTToken)

	w.WriteHeader(http.StatusOK)
	w.Write(bb.FinishedBytes())

	slog.Info("GetJWT", "msg", "jwt token generated", "email", userCredentials.Email)
}
