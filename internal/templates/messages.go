package templates

import "fmt"

// Messages содержит все шаблоны сообщений бота
type Messages struct {
	// Общие сообщения
	BotGroupOnly   *MessageTemplate
	UnknownCommand *MessageTemplate
	ErrorOccurred  *MessageTemplate

	// Справка
	HelpText *MessageTemplate

	// Пидор дня
	PersonAlreadySelected *MessageTemplate
	PersonSelected        *MessageTemplate
	NoActiveUsers         *MessageTemplate
	PersonInfo            *MessageTemplate
	NoPersonSelectedToday *MessageTemplate

	// Статистика
	StatsHeader *MessageTemplate
	StatsEmpty  *MessageTemplate
	StatsEntry  *MessageTemplate
}

// NewMessages создает новый набор сообщений
func NewMessages() (*Messages, error) {
	messages := &Messages{}

	// Инициализируем все шаблоны
	templates := map[string]**MessageTemplate{
		// Общие сообщения
		"BotGroupOnly":   &messages.BotGroupOnly,
		"UnknownCommand": &messages.UnknownCommand,
		"ErrorOccurred":  &messages.ErrorOccurred,

		// Справка
		"HelpText": &messages.HelpText,

		// Пидор дня
		"PersonAlreadySelected": &messages.PersonAlreadySelected,
		"PersonSelected":        &messages.PersonSelected,
		"NoActiveUsers":         &messages.NoActiveUsers,
		"PersonInfo":            &messages.PersonInfo,
		"NoPersonSelectedToday": &messages.NoPersonSelectedToday,

		// Статистика
		"StatsHeader": &messages.StatsHeader,
		"StatsEmpty":  &messages.StatsEmpty,
		"StatsEntry":  &messages.StatsEntry,
	}

	// Шаблоны сообщений
	messageTemplates := map[string]string{
		"BotGroupOnly": "Этот бот работает только в группах!",

		"UnknownCommand": "Неизвестная команда. Используйте /help для списка команд.",

		"ErrorOccurred": "Произошла ошибка: {{error}}",

		"HelpText": `🎯 Добро пожаловать в бота "Пидор дня"!

Доступные команды:
/pidor - Выбрать пидора дня
/pidorstats - Показать статистику всех участников
/pidorinfo - Информация о сегодняшнем пидоре дня
/help - Показать эту справку

Бот работает только в группах и выбирает случайного участника из числа активных пользователей.`,

		"PersonAlreadySelected": `🎯 Пидор дня уже выбран!

👤 {{person}}`,

		"PersonSelected": `🎉 Пидор дня выбран!

🎯 {{person}}

Поздравляем! 🎊`,

		"NoActiveUsers": "В группе нет активных участников для выбора.",

		"PersonInfo": `ℹ️ Информация о сегодняшнем пидоре дня:

👤 {{person}}
📅 {{date}}`,

		"NoPersonSelectedToday": "Сегодня пидор дня еще не выбран. Используйте /pidor для выбора!",

		"StatsHeader": "📊 Статистика \"Пидор дня\":\n\n",

		"StatsEmpty": "В этой группе пока нет статистики.",

		"StatsEntry": "{{position}} {{person}} - {{count}} раз\n",
	}

	// Создаем шаблоны
	for name, templatePtr := range templates {
		templateStr, exists := messageTemplates[name]
		if !exists {
			return nil, fmt.Errorf("template %s not found", name)
		}

		template, err := NewTemplate(templateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to create template %s: %w", name, err)
		}

		*templatePtr = template
	}

	return messages, nil
}

// GetPositionEmoji возвращает эмодзи для позиции в статистике
func GetPositionEmoji(position int) string {
	switch position {
	case 1:
		return "🥇"
	case 2:
		return "🥈"
	case 3:
		return "🥉"
	default:
		return fmt.Sprintf("%d.", position)
	}
}
