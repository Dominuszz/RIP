package repository

import (
	"errors"
	"fmt"
	"lab3/internal/app/ds"
	"lab3/internal/app/serializer"

	"gorm.io/gorm"
)

func (r *Repository) GetUserByID(id int) (ds.Users, error) {
	user := ds.Users{}
	if id <= 0 {
		return ds.Users{}, fmt.Errorf("неверный id: должен быть > 0")
	}

	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return ds.Users{}, err
	}
	return user, nil
}

func (r *Repository) GetUserByLogin(login string) (ds.Users, error) {
	user := ds.Users{}
	if login == "" {
		return ds.Users{}, errors.New("логин не может быть пустым")
	}

	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Users{}, fmt.Errorf("%w: пользователь с логином %s не найден", ErrNotFound, login)
		}
		return ds.Users{}, err
	}
	return user, nil
}

func (r *Repository) CreateUser(userJSON serializer.UserJSON) (ds.Users, error) {
	user := serializer.UserFromJSON(userJSON)
	if user.Login == "" {
		return ds.Users{}, errors.New("логин обязателен для заполнения")
	}
	if user.Password == "" {
		return ds.Users{}, errors.New("пароль обязателен для заполнения")
	}
	_, err := r.GetUserByLogin(user.Login)
	if err == nil {
		return ds.Users{}, fmt.Errorf("%w: пользователь с логином %s уже существует", ErrAlreadyExists, user.Login)
	} else if !errors.Is(err, ErrNotFound) {
		return ds.Users{}, err
	}
	err = r.db.Create(&user).Error
	if err != nil {
		return ds.Users{}, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return user, nil
}

func (r *Repository) SignIn(userJSON serializer.UserJSON) (ds.Users, error) {
	if userJSON.Login == "" {
		return ds.Users{}, errors.New("логин обязателен для заполнения")
	}
	if userJSON.Password == "" {
		return ds.Users{}, errors.New("пароль обязателен для заполнения")
	}

	user, err := r.GetUserByLogin(userJSON.Login)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ds.Users{}, errors.New("неверный логин или пароль")
		}
		return ds.Users{}, err
	}

	if user.Password != userJSON.Password {
		return ds.Users{}, errors.New("неверный логин или пароль")
	}

	r.SetUserID(int(user.ID))
	return user, nil
}

func (r *Repository) EditInfo(id int, userJSON serializer.UserJSON) (ds.Users, error) {
	if id <= 0 {
		return ds.Users{}, fmt.Errorf("неверный id пользователя")
	}

	currentUser, err := r.GetUserByID(id)
	if err != nil {
		return ds.Users{}, err
	}

	updates := serializer.UserFromJSON(userJSON)

	if updates.IsModerator && !currentUser.IsModerator {
		updates.IsModerator = false
	}

	err = r.db.Model(&currentUser).Updates(updates).Error
	if err != nil {
		return ds.Users{}, fmt.Errorf("ошибка при обновлении профиля: %w", err)
	}

	return r.GetUserByID(id)
}
