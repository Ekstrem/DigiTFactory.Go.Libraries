package result

import "github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/monads"

// ToResult преобразует значение и ошибку в Result-монаду.
func ToResult[T any](value T, err error) monads.Result[T, error] {
	if err != nil {
		return monads.NewFailure[T, error](err)
	}
	return monads.NewSuccess[T, error](value)
}
