package users

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"

	m "go_rest/internal/models"
	r "go_rest/internal/repositories"
)

func TestUsers(t *testing.T) {
	db, err := r.SetupTestdb()
	if err != nil {
		t.Fatal(err)
	}
	u := NewUserRepository(db)
	t.Run("Create", func(t *testing.T) {
		user := m.User{
			Email:   "test_user@email.com",
			Pass:    "test_pass",
			IsAdmin: true,
		}
		err = u.Create(&user)
		if err != nil {
			t.Fatal(err)
		}
		if user.Id != 1 {
			t.Errorf("Ошибка создания user")
		}
		var userId int
		var email, pass string
		var isAdmin bool
		row := db.QueryRow("SELECT userid, email, pass, is_admin FROM users WHERE userid = ?", 1)
		err = row.Scan(&userId, &email, &pass, &isAdmin)
		if err != nil {
			t.Fatal(err)
		}
		if userId != 1 {
			t.Errorf("Error: got %v, want %v", userId, 1)
		}
		if email != "test_user@email.com" {
			t.Errorf("Error: got %v, want %v", email, "test_user@email.com")
		}

		if pass != "test_pass" {
			t.Errorf("Error: got %v, want %v", pass, "test_pass")
		}
		if isAdmin != true {
			t.Errorf("Error: got %v, want %v", isAdmin, true)
		}

	})

	t.Run("Get", func(t *testing.T) {
		stmt, err := db.Prepare("INSERT INTO users(email, pass, is_admin) VALUES(?, ?, ?)")
		if err != nil {
			t.Fatal(err)
		}

		_, err = stmt.Exec("test@email.com", "123", true)
		if err != nil {
			t.Fatal(err)
		}
		user, err := u.Get(2)
		if err != nil {
			t.Fatal(err)
		}
		if user.Id != 2 {
			t.Errorf("Error: got %v, want %v", user.Id, 2)
		}
		if user.Email != "test@email.com" {
			t.Errorf("Error: got %v, want %v", user.Email, "test@email.com")
		}
		if user.Pass != "123" {
			t.Errorf("Error: got %v, want %v", user.Pass, "123")
		}
		if user.IsAdmin != true {
			t.Errorf("Error: got %v, want %v", user.IsAdmin, true)
		}
		userbyEmail, err := u.GetbyEmail("test_user@email.com")
		if err != nil {
			t.Fatal(err)
		}
		if userbyEmail.Email != "test_user@email.com" {
			t.Errorf("Error: got %v, want %v", userbyEmail.Email, "test_user@email.com")
		}
		if userbyEmail.Id != 1 {
			t.Errorf("Error: got %v, want %v", userbyEmail.Id, 1)
		}

	})

}
