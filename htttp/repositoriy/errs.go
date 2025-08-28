package repositoriy

import "errors"

var ThisNameIsExist = errors.New("Это имя уже используется, попробуйте другое!")

var ThisNameIsNotExist = errors.New("Это имя не найдено!")
