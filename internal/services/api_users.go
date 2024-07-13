package services

import (
	"Authorization-Service/internal/domain/models"
	"Authorization-Service/internal/server/configs"
	"Authorization-Service/internal/server/redis"
	"encoding/json"
	"mime"
	"net/http"
	"strings"
)

func ContentTypeCheck(w http.ResponseWriter, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")           // получаем тип контента в запросе
	mediatype, _, err := mime.ParseMediaType(contentType) // парсинг полученных данных из запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // ошибка, если не получили запрос
		return err
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType) // ошибка, если пришел не json
		return err
	}
	return nil
}

func GetUsers() ([]models.ResponseUser, error) {
	url := configs.Url + configs.P
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var rus []models.ResponseUser
	if err := dec.Decode(&rus); err != nil {
		return nil, err
	}

	return rus, nil
}

func FindUser(login, password string, users []models.ResponseUser) bool {
	for _, user := range users {
		if user.Login == login && user.Password == password {
			return true
		}
	}
	return false
}

// RememberUser godoc
// @Summary Authorize a new user
// @Description Авторизовать нового пользователя
// @Tags Users
// @Accept  json
// @Produce  json
// @Param User body models.ResponseUser true "User must be logged in"
// @Success 201 {array} string "{"token", "user_id"}"
// @Failure 400 {object} string "{"message"}"
// @Failure 500 {object} string "{"message"}"
// @Router /v1/remember/ [post]
func RememberUser(w http.ResponseWriter, r *http.Request) {

	type RequestUser struct {
		Login    string `json:"login,omitempty"`
		Password string `json:"password,omitempty"`
	}

	type Response struct {
		Token string `json:"token",omitempty`
		Id    string `json:"user_id",omitempty`
	}

	err := ContentTypeCheck(w, r)

	dec := json.NewDecoder(r.Body) // декодируем тело запроса
	dec.DisallowUnknownFields()    // проверка, что получили то, что готовы принять (ругается на id, если оно есть в запросе)
	var ru RequestUser
	if err := dec.Decode(&ru); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // получаем декодированные данные и проверяем, что все ок
		return
	}

	if len(ru.Login) == 0 || len(ru.Password) == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user_data_id string
	rus, err := GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, user := range rus {
		if user.Login == ru.Login && user.Password == ru.Password { // data = Login/Password/Token
			user_data_id = (user.UserId).String()
		}
	}

	// Writing to cache
	data, err := redis.Client.Get(user_data_id).Result() // ищем юзера к кэше
	var token string = ""

	data_mas := strings.Split(data, "/")
	data_new := ru.Login + ru.Password
	if err == nil {
		if data_mas[1] == ru.Password {

			token, err = GenerateJWT(ru.Login) // Создаем JWT токен
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError) // Не получилось создать токен
				return
			}

			data_new = data_new + token

			err := redis.Client.Set(user_data_id, data_new, 0).Err() // Перезаписываем кэш
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError) // Не получилось перезаписать кэш
				return
			}

		} else {

			http.Error(w, err.Error(), http.StatusBadRequest) // Неправильный пароль
			return

		}

	} else {
		// Не нашли юзера
		// ищем юзера в аккаунт-сервисе

		found := FindUser(ru.Login, ru.Password, rus)
		if found {

			token, err = GenerateJWT(ru.Login) // Создаем JWT токен
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError) // Не получилось создать токен
				return
			}

			data_new = data_new + token

			err := redis.Client.Set(user_data_id, data_new, 0).Err() // Перезаписываем кэш
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError) // Не получилось перезаписать кэш
				return
			}

		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	js, err := json.Marshal(Response{Token: token, Id: user_data_id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}
