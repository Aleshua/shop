package repositories

import (
	"context"

	d "auth/internal/domain"
)

type IUserRepository interface {
	// Возвращает пользователя по id.
	// Если не найден возвращается ошибка.
	Find(ctx context.Context, userId int32) (d.User, error)

	// Возвращает пользователя по email.
	// Если не найден возвращается ошибка.
	FindByEmail(ctx context.Context, email string) (d.User, error)

	// Возвращает id созданного пользователя.
	// Возвращает ошибку если такой email занят.
	Save(ctx context.Context, user d.User) (int32, error)

	// Если пользователь уже создан и у него не подтверждена почта то поменяет пароль и не вернёт ошибку.
	SaveOrUpdateUnverified(ctx context.Context, user d.User) (int32, error)

	// Устанавливает текующее время в поле email_verified_at
	MarkEmailVerified(ctx context.Context, userId int32) error

	// Не вернёт ошибку если есть у пользователя. Заменит на новый.
	SaveOrUpdateConfirmCode(ctx context.Context, code d.ConfirmCode) error

	// Находит данные подтверждающего кода по id пользователя.
	// Возвращает ошибку если не найдёт.
	FindConfirmCode(ctx context.Context, userId int32) (d.ConfirmCode, error)

	// Обновляет сам код, а также время и попытки
	UpdateConfirmCode(ctx context.Context, code d.ConfirmCode) error

	// Удаляет код подтверждения для пользователя
	DeleteConfirmCode(ctx context.Context, userId int32) error
}
