package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pavel-one/day-of-the-bot/internal/domain"
)

// PersonOfTheDayRepositoryImpl реализует PersonOfTheDayRepository
type PersonOfTheDayRepositoryImpl struct {
	db *Database
}

// NewPersonOfTheDayRepository создает новый экземпляр PersonOfTheDayRepository
func NewPersonOfTheDayRepository(db *Database) PersonOfTheDayRepository {
	return &PersonOfTheDayRepositoryImpl{db: db}
}

// Set устанавливает человека дня
func (r *PersonOfTheDayRepositoryImpl) Set(userID, chatID int64, date time.Time) error {
	dateStr := date.Format("2006-01-02")

	query := r.db.psql.Replace("person_of_the_day").
		Columns("user_id", "chat_id", "date").
		Values(userID, chatID, dateStr)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.conn.Exec(sqlStr, args...)
	if err != nil {
		return fmt.Errorf("failed to set person of the day: %w", err)
	}

	return nil
}

// GetByDate возвращает человека дня на указанную дату
func (r *PersonOfTheDayRepositoryImpl) GetByDate(chatID int64, date time.Time) (*domain.User, error) {
	dateStr := date.Format("2006-01-02")

	query := r.db.psql.Select("u.id", "u.username", "u.first_name", "u.last_name", "u.chat_id", "u.created_at").
		From("users u").
		Join("person_of_the_day p ON u.id = p.user_id AND u.chat_id = p.chat_id").
		Where(squirrel.Eq{"p.chat_id": chatID, "p.date": dateStr})

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
		return nil, fmt.Errorf("failed to get today's person: %w", err)
	}

	if username.Valid {
		user.Username = username.String
	}
	if lastName.Valid {
		user.LastName = lastName.String
	}

	return &user, nil
}

// GetUserStats возвращает статистику пользователей
func (r *PersonOfTheDayRepositoryImpl) GetUserStats(chatID int64) ([]domain.UserStats, error) {
	query := r.db.psql.Select(
		"u.id", "u.username", "u.first_name", "u.last_name", "u.chat_id", "u.created_at",
		"COALESCE(COUNT(p.id), 0) as count",
	).
		From("users u").
		LeftJoin("person_of_the_day p ON u.id = p.user_id AND u.chat_id = p.chat_id").
		Where(squirrel.Eq{"u.chat_id": chatID}).
		GroupBy("u.id", "u.username", "u.first_name", "u.last_name", "u.chat_id", "u.created_at").
		OrderBy("count DESC", "u.first_name")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.conn.Query(sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Ошибка закрытия rows: %v\n", err)
		}
	}()

	var stats []domain.UserStats
	for rows.Next() {
		var userStat domain.UserStats
		var username sql.NullString
		var lastName sql.NullString

		err := rows.Scan(
			&userStat.User.ID,
			&username,
			&userStat.User.FirstName,
			&lastName,
			&userStat.User.ChatID,
			&userStat.User.CreatedAt,
			&userStat.Count,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user stats: %w", err)
		}

		if username.Valid {
			userStat.User.Username = username.String
		}
		if lastName.Valid {
			userStat.User.LastName = lastName.String
		}

		stats = append(stats, userStat)
	}

	return stats, nil
}
