package repository

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pavel-one/day-of-the-bot/internal/domain"
)

// UserRepositoryImpl реализует UserRepository
type UserRepositoryImpl struct {
	db *Database
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *Database) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Add добавляет или обновляет пользователя
func (r *UserRepositoryImpl) Add(user domain.User) error {
	query := r.db.psql.Replace("users").
		Columns("id", "username", "first_name", "last_name", "chat_id").
		Values(user.ID, user.Username, user.FirstName, user.LastName, user.ChatID)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.conn.Exec(sqlStr, args...)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}

	return nil
}

// GetByChatID возвращает всех пользователей в чате
func (r *UserRepositoryImpl) GetByChatID(chatID int64) ([]domain.User, error) {
	query := r.db.psql.Select("id", "username", "first_name", "last_name", "chat_id", "created_at").
		From("users").
		Where(squirrel.Eq{"chat_id": chatID}).
		OrderBy("first_name")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.conn.Query(sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat users: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			fmt.Printf("Ошибка закрытия rows: %v\n", err)
		}
	}()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		var username sql.NullString
		var lastName sql.NullString

		err := rows.Scan(&user.ID, &username, &user.FirstName, &lastName, &user.ChatID, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if username.Valid {
			user.Username = username.String
		}
		if lastName.Valid {
			user.LastName = lastName.String
		}

		users = append(users, user)
	}

	return users, nil
}

// GetByID возвращает пользователя по ID и chat ID
func (r *UserRepositoryImpl) GetByID(userID, chatID int64) (*domain.User, error) {
	query := r.db.psql.Select("id", "username", "first_name", "last_name", "chat_id", "created_at").
		From("users").
		Where(squirrel.Eq{"id": userID, "chat_id": chatID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := r.db.conn.QueryRow(sqlStr, args...)

	var user domain.User
	var username sql.NullString
	var lastName sql.NullString

	err = row.Scan(&user.ID, &username, &user.FirstName, &lastName, &user.ChatID, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if username.Valid {
		user.Username = username.String
	}
	if lastName.Valid {
		user.LastName = lastName.String
	}

	return &user, nil
}
