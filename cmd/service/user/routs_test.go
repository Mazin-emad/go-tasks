package user

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/Mazin-emad/todo-backend/types"
// 	"github.com/gorilla/mux"
// )

// func TestUserRoutesHandle(t *testing.T) {

// 	userStore := &mockUserStore{}
// 	userHandler := NewHandler(userStore)

// 	t.Run("should fial if the username is already taken", func(t *testing.T) {

// 		payload := types.RegisterUserPayload{
// 			UserName: "",
// 			Password: "",
// 		}

// 		marshal, _ := json.Marshal(payload)

// 		req,err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshal))
	
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		router := mux.NewRouter()
		
// 		router.HandleFunc("/register", userHandler.handleRegister)
// 		router.ServeHTTP(rr, req)

// 	if rr.Code != http.StatusBadRequest {
// 		t.Errorf("expected status %d but got %d", http.StatusBadRequest, rr.Code)
// 	}

// 	})

// }

// type mockUserStore struct {}

// func (m *mockUserStore) GetUserByUsername(username string) (*types.User, error) {
// 	return nil, fmt.Errorf("user not found")
// }

// func (m *mockUserStore) CreateUser(user *types.User) error {
// 	return nil
// }


// func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
// 	return nil, nil
// }

// func (m *mockUserStore) GetAllUsers() ([]*types.User, error) {
// 	return nil, nil
// }
