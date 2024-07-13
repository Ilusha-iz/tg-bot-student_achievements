package createmenu

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// CreateInlineKeyboardMenuSearchStudent создает меню для поиска студентов.
func CreateInlineKeyboardMenuSearchStudent() *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<", "button_search_prev_t"),
			tgbotapi.NewInlineKeyboardButtonData(">", "button_search_next_t"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏆 Перейти к достиженям", "button_students_achievements"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "button_back_teacher_search_menu"),
		),
	)
	return &keyboard
}

// CreateInlineKeyboardMenuAchievements создает меню для работы с достижениями.
func CreateInlineKeyboardMenuAchievements() *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<", "button_search_prev_a"),
			tgbotapi.NewInlineKeyboardButtonData("⬇️ Скачать", "button_upload_search"),
			tgbotapi.NewInlineKeyboardButtonData(">", "button_search_next_a"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Выйти в меню", "button_back_teacher_search_menu"),
		),
	)
	return &keyboard
}

// CreateInlineKeyboardMenuStudentAchievements создает меню для работы со студентскими достижениями.
func CreateInlineKeyboardMenuStudentAchievements() *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<", "button_search_prev_sa"),
			tgbotapi.NewInlineKeyboardButtonData("⬇️ Скачать", "button_upload_search_a"),
			tgbotapi.NewInlineKeyboardButtonData(">", "button_search_next_sa"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "button_cancel"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Выйти в меню", "button_back_teacher_search_menu"),
		),
	)
	return &keyboard
}

// CreateMenuKeyboard создает основное меню.
func CreateMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📇 Перейти к личным данным", "button_data"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➡️ Перейти к достижениям", "button_go_to_achievements"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить свой профиль", "button_delete_profile"),
		),
	)
}

// CreateTeacherMenu создает меню для преподавателя.
func CreateTeacherMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔍 Поиск данных", "button_teacher_search_menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬇️ Выгрузка данных", "button_download_menu"),
		),
	)
}

// CreateTeacherSearchMenu создает меню для поиска информации о студентах.
func CreateTeacherSearchMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📇 ФИО студента", "button_search_fio"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 Группе", "button_search_group"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏆 Достижению", "button_search_achievements"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "button_back_teacher_menu"),
		),
	)
}

// CreateTeacherUploadMenu создает меню для загрузки данных преподавателя.
func CreateTeacherUploadMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 По группе", "button_upload_group"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏆 По достижению", "button_upload_achievements"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗂 Выгрузить всё", "button_upload_all"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "button_back_teacher_menu"),
		),
	)
}

// CreateAchievementsMenuKeyboard создает меню для работы с достижениями.
func CreateAchievementsMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить достижения", "button_add"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Изменить достижения", "button_edit_achievements"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👀 Посмотреть достижения", "button_see"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔍 Найти достижения", "button_find"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить достижения", "button_delete_menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "button_back_to_menu"),
		),
	)
}

// CreateDeleteMenuKeyboard создает меню для удаления достижений.
func CreateDeleteMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить все достижения", "button_delete_all"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить конкретное достижение", "button_find"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "button_back_to_menu_achievements"),
		),
	)
}

// CreateAchievementAddedKeyboard создает меню после добавления достижения.
func CreateAchievementAddedKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить еще", "button_add"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "button_back_to_menu_achievements"),
		),
	)
}

// CreateInlineKeyboardData создает меню для изменения личных данных.
func CreateInlineKeyboardData() *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить ФИО", "button_data_edit_fio"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить группу", "button_data_edit_group"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Выход", "button_back_to_menu"),
		),
	)
	return &keyboard
}

// CreateInlineKeyboard создает общее меню с возможностью управления.
func CreateInlineKeyboard() *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<", "button_search_prev"),
			tgbotapi.NewInlineKeyboardButtonData(">", "button_search_next"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬇️ Скачать", "button_upload"),
			tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить", "button_delete_ach"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Изменить название достижения", "button_edit_name_achievements"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Изменить файл", "button_edit_file"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Выход", "button_back_to_menu_achievements"),
		),
	)
	return &keyboard
}
