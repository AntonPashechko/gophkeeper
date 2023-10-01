package models

import "fmt"

type NewDataDTO struct {
	Id   string `json:"id"`   //Идентификатор данных (любая фраза)
	Data string `json:"data"` //Зашифрованные данные в base64
}

func (m *NewDataDTO) Validate() error {
	if m.Id == `` {
		return fmt.Errorf("id required")
	}
	if m.Data == `` {
		return fmt.Errorf("data required")
	}

	return nil
}
