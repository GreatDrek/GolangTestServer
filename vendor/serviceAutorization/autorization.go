package serviceAutorization

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/json"

	"bytes"
	"errors"
	"log"
	"time"
)

type LogginDataClient struct {
	Id    int    `json:"-"`
	Email string `json:"emailClient"`
	Key   []byte `json:"key"`
	Salt  []byte `json:"-"`
}

func Autorization(requestType byte, data []byte, db *sql.DB) (*LogginDataClient, error) {
	var logginData *LogginDataClient
	err := json.Unmarshal(data, &logginData)
	if err != nil {
		return logginData, err
	}

	if logginData == nil {
		return logginData, errors.New("Error logginData 101")
	}

	err = logickLoggin(requestType, logginData, db)
	if err != nil {
		return logginData, err
	}

	return logginData, nil
}

func logickLoggin(requestType byte, ld *LogginDataClient, db *sql.DB) error {
	var retunrError error

	switch requestType {
	case 100: // Запрос проверки логина пароля

		// Проверяем наличие не нулевых данных в сообщении
		if ld.Email == "" || len(ld.Key) == 0 {
			retunrError = errors.New("error email or key")
			break
		} else {
			// Идем в БД для проверки данных
			log.Println("Start BD", time.Now().String())
			infoClient, err := checkUser(*ld, db)
			if err != nil {
				retunrError = err
				break
			}
			log.Println("Stop BD", time.Now().String())
			// Если бд возвращает пустого клиента, говорим что логин и пароль не верны
			if infoClient == nil {
				retunrError = errors.New("dont accaunt")
				break
			} else {
				// Если клиент с таким логином есть, проверяем пароль

				hash, err := hashSum(ld.Key, infoClient.Salt)
				if err != nil {
					retunrError = err
					break
				}

				// Проверям ключи, если они не верны говорим что не верный пароль
				if bytes.Equal(hash, infoClient.Key) == false {
					retunrError = errors.New("error key")
					break
				}

				ld.Id = infoClient.Id
				log.Println("Connect")
			}
			break
		}
		break

	case 101: // Запрос регистрации

		// Проверяем что бы логин не был нулевым
		if ld.Email == "" {
			retunrError = errors.New("null email")
			break
		} else {
			// Создаем нового клиента и добавляем его в БД

			emailOaut, err := returnEmail(ld.Email)
			if err != nil {
				retunrError = errors.New("nil email content")
				break
			}

			if emailOaut == "" {
				retunrError = errors.New("nil email content")
				break
			} else {
				ld.Email = emailOaut
			}

			//Проверяем есть ли в БД пользователь с таким логином
			_infoClient, err := checkUser(*ld, db)
			if err != nil {
				retunrError = err
				break
			}

			// Генерируем ключ для него
			newKey, err := randGenerate(64)
			if err != nil {
				retunrError = err
				break
			}

			// Генерируем соль для ключа
			newSalt, err := randGenerate(8)
			if err != nil {
				retunrError = err
				break
			}

			ld.Key = newKey

			hash, err := hashSum(newKey, newSalt)
			if err != nil {
				retunrError = err
				break
			}

			infoClient := &LogginDataClient{Email: ld.Email, Key: hash, Salt: newSalt}

			// Если пользователя нет, то регистрируем его
			if _infoClient == nil {
				// Добавляем в бд нового пользователь
				id, _ := addUser(infoClient, db)
				if id == 0 {
					retunrError = errors.New("error add bd")
					break
				}

				log.Println("Regestration")
				ld.Id = id

				break
			} else {
				// Обновляем в бд пользователя
				err = updateUser(infoClient, db)
				if err != nil {
					retunrError = err
					break
				}
				// Такой аккаунт уже зарегистрирован
				log.Println("Re Regestration")
				ld.Id = _infoClient.Id
				break
			}
		}
		break

	default:
		retunrError = errors.New("default")
		break
	}
	return retunrError
}

func randGenerate(lenght byte) ([]byte, error) {
	b := make([]byte, lenght)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, err
}

func hashSum(key []byte, salt []byte) ([]byte, error) {
	h := sha256.New()
	key0 := make([]byte, 0)
	key0 = append(key0, key...)
	key0 = append(key0, salt...)
	_, err := h.Write(key0)
	if err != nil {
		return nil, err
	}

	generateKey := h.Sum(nil)
	return generateKey, err
}
