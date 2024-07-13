package telegrambot

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"

	cm "tg_bot/createMenu"
	db "tg_bot/db"
	gv "tg_bot/globalVar"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// Отправка сообщения с информацией о текущем элементе
func sendStudentsMessage(bot *tgbotapi.BotAPI, chatID int64, students []gv.Student, index int, messageID int) {
	currentStudent := students[index]
	text := fmt.Sprintf("%d/%d\n\n\nФИО студента: %s\n\nГруппа: %s", index+1, len(students), currentStudent.Fio, currentStudent.Group)
	editText := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editText.ReplyMarkup = cm.CreateInlineKeyboardMenuSearchStudent()
	bot.Send(editText)
}

// Обработка нажатий кнопок
func handleCallbackQueryStudent(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, students []gv.Student, currentIndex *int, chatID int64) {

	if callback.Data == "button_students_achievements" {

		achievements, err := db.AllAchievements(int64(students[*currentIndex].ChatID))
		if err != nil {
			log.Printf("Произошла ошибка при поиске всех элементов: %v", err)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
			bot.Send(msg)
		} else {
			if len(achievements) == 0 {
				msg := tgbotapi.NewMessage(chatID, "Ничего не найдено.")

				bot.Send(msg)
			}
			gv.UserAchievementIndex[chatID] = 0
			gv.UserSearchResults[chatID] = achievements
			msg := tgbotapi.NewMessage(chatID, "Результаты поиска:")
			msg.ReplyMarkup = cm.CreateInlineKeyboardMenuStudentAchievements()
			sentMsg, _ := bot.Send(msg)
			gv.UserMessageChatID[chatID] = sentMsg.MessageID
			sendAchievementMessage(bot, chatID, achievements, 0, sentMsg.MessageID)
			// Сбрасываем состояние и данные пользователя
			gv.UserStates[chatID] = ""
		}

	} else if callback.Data == "button_search_prev_t" {
		if *currentIndex > 0 {
			*currentIndex--
		}
	} else if callback.Data == "button_search_next_t" {
		if *currentIndex < len(students)-1 {
			*currentIndex++
		}
	}
	sendStudentsMessage(bot, callback.Message.Chat.ID, students, *currentIndex, callback.Message.MessageID)
}

// Отправка сообщения с информацией о текущем элементе
func sendStudentsMessageAchievements(bot *tgbotapi.BotAPI, chatID int64, data []gv.Data, index int, messageID int) {
	currentData := data[index]
	text := fmt.Sprintf("%d/%d\n\n\nФИО студента: %s\nГруппа: %s\n\nНазвание достижения: %s\nИмя файла: %s", index+1, len(data), currentData.Student.Fio, currentData.Student.Group, currentData.Achievement.Name, currentData.Achievement.Filename)
	editText := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editText.ReplyMarkup = cm.CreateInlineKeyboardMenuAchievements()
	bot.Send(editText)
}

// Отправка сообщения с информацией о текущем элементе
func sendStudentsMessageSerachAchievements(bot *tgbotapi.BotAPI, chatID int64, achievements []gv.Achievement, index int, messageID int) {
	currentAchievement := achievements[index]
	text := fmt.Sprintf("%d/%d\n\n\nНазвание: %s\n\nИмя файла: %s", index+1, len(achievements), currentAchievement.Name, currentAchievement.Filename)
	editText := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editText.ReplyMarkup = cm.CreateInlineKeyboardMenuStudentAchievements()
	bot.Send(editText)
}

// Обработка нажатий кнопок
func handleCallbackQueryAchievements(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, data []gv.Data, currentIndex *int) {

	if callback.Data == "button_upload_search" {
		if err := db.SendFile(bot, callback, data[*currentIndex].Student.ChatID, data[*currentIndex].Achievement.Name); err != nil {
			log.Printf("❌ Произошла ошибка при выгрузке файла: %v", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("%v", err))
			bot.Send(msg)
			return
		}
	} else if callback.Data == "button_search_prev_a" {
		if *currentIndex > 0 {
			*currentIndex--
		}
	} else if callback.Data == "button_search_next_a" {
		if *currentIndex < len(data)-1 {
			*currentIndex++
		}
	}
	sendStudentsMessageAchievements(bot, callback.Message.Chat.ID, data, *currentIndex, callback.Message.MessageID)
}

func handleCallbackQuerySearchAchievements(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, achievements []gv.Achievement, currentIndex *int) {

	if callback.Data == "button_upload_search_a" {
		log.Println(achievements[*currentIndex].ChatID)
		if err := db.SendFile(bot, callback, int(achievements[*currentIndex].ChatID), achievements[*currentIndex].Name); err != nil {
			log.Printf("❌ Произошла ошибка при выгрузке файла: %v", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("%v", err))
			bot.Send(msg)
			return
		}
	} else if callback.Data == "button_search_prev_sa" {
		log.Println(currentIndex)
		if *currentIndex > 0 {
			*currentIndex--
		}
	} else if callback.Data == "button_search_next_sa" {
		log.Println(currentIndex)
		if *currentIndex < len(achievements)-1 {
			*currentIndex++
		}
	}
	sendStudentsMessageSerachAchievements(bot, callback.Message.Chat.ID, achievements, *currentIndex, callback.Message.MessageID)
}

// Отправка сообщения с информацией о текущем элементе
func sendDataMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	fio, group, err := db.GetFIOAndGroup(chatID)
	if err != nil {
		log.Printf("❌ Не удалось получить личные данные пользователя: %v", err)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
		bot.Send(msg)
		return
	}
	text := fmt.Sprintf("Ваши личные данные \n\n\nФИО: %s\nГруппа: %s", fio, group)
	editText := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editText.ReplyMarkup = cm.CreateInlineKeyboardData()
	bot.Send(editText)
}

// Обработка сообщений
func handleMessageData(bot *tgbotapi.BotAPI, chatID int64) {

	if gv.UserStates[chatID] == "edit_fio_data" {
		err := db.EditFIO(chatID)
		if err != nil {
			log.Printf("❌ Не удалось изменить личные данные пользователя: %v", err)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Вы успешно поменяли ФИО на %s", gv.UserData[chatID]["FIO"]))
			bot.Send(msg)
		}
		// Сбрасываем состояние
		gv.UserStates[chatID] = ""
	} else if gv.UserStates[chatID] == "edit_group_data" {
		err := db.EditGroup(chatID)
		if err != nil {
			log.Printf("❌ Не удалось изменить личные данные пользователя: %v", err)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Вы успешно поменяли группу на %s", gv.UserData[chatID]["group"]))
			bot.Send(msg)
		}
	}
	sendDataMessage(bot, chatID, gv.UserMessageChatID[chatID])
}

// Отправка сообщения с информацией о текущем элементе
func sendAchievementMessage(bot *tgbotapi.BotAPI, chatID int64, achievements []gv.Achievement, index int, messageID int) {

	isStudent, ok, err := db.CheckStudentByChatID(chatID)

	if !ok || (err != nil) {
		log.Printf("❌ Произошла ошибка при подключении к БД: %v", err)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
		bot.Send(msg)
		return
	}

	if len(achievements) == 0 {
		text := "Нет доступных достижений."
		editText := tgbotapi.NewEditMessageText(chatID, messageID, text)
		bot.Send(editText)
		var inlineKeyboard = cm.CreateAchievementsMenuKeyboard()
		responseText := "Выберите дальнейшее действие:"
		msg := tgbotapi.NewMessage(chatID, responseText)
		msg.ReplyMarkup = inlineKeyboard
		bot.Send(msg)
		return
	}

	currentAchievement := achievements[index]
	text := fmt.Sprintf("%d/%d\n\n\nНазвание: %s\n\nИмя файла: %s", index+1, len(achievements), currentAchievement.Name, currentAchievement.Filename)
	editText := tgbotapi.NewEditMessageText(chatID, messageID, text)
	if isStudent {
		editText.ReplyMarkup = cm.CreateInlineKeyboard()
	} else {
		editText.ReplyMarkup = cm.CreateInlineKeyboardMenuStudentAchievements()
	}
	bot.Send(editText)
}

// Обработка нажатий кнопок
func handleCallbackQuery(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, achievements []gv.Achievement, currentIndex *int, chatID int64) {

	if callback.Data == "button_upload" {
		if err := db.SendFile(bot, callback, int(chatID), achievements[*currentIndex].Name); err != nil {
			log.Printf("❌ Произошла ошибка при выгрузке файла: %v", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("%v", err))
			bot.Send(msg)
			return
		}
	} else if callback.Data == "button_delete_ach" {
		// Удаление достижения из базы данных
		if err := db.DeleteRecord(callback.Message.Chat.ID, achievements[*currentIndex].Name); err != nil {
			log.Printf("❌ Произошла ошибка при удалении достижения: %v", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("%v", err))
			bot.Send(msg)
		} else {
			log.Printf("✅ Достижение из чата %d успешно удалено", callback.Message.Chat.ID)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("✅ Достижение '%s' успешно удалено", achievements[*currentIndex].Name))
			bot.Send(msg)

			// Обновление списка достижений после удаления
			achievements = append(achievements[:*currentIndex], achievements[*currentIndex+1:]...)
			if *currentIndex >= len(achievements) && *currentIndex > 0 {
				*currentIndex--
			}
			gv.UserSearchResults[chatID] = achievements
		}
	} else if callback.Data == "new_name_achievement" {
		if err := db.UpdateAchievementsInDatabase(chatID, achievements[*currentIndex].Name); err != nil {
			log.Printf("❌ Произошла ошибка при обновлении достижений: %v", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("Произошла ошибка при обновлении достижений: %v", err))
			bot.Send(msg)
		} else {
			log.Printf("✅ Достижения для чата %d успешно обновлены", chatID)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "✅ Достижения успешно обновлены")
			achievements[*currentIndex].Name = gv.UserData[chatID]["achievements"]
			bot.Send(msg)
		}
	} else if callback.Data == "button_search_prev" {
		if *currentIndex > 0 {
			*currentIndex--
		}
	} else if callback.Data == "button_search_next" {
		if *currentIndex < len(achievements)-1 {
			*currentIndex++
		}
	}
	sendAchievementMessage(bot, callback.Message.Chat.ID, achievements, *currentIndex, callback.Message.MessageID)
}

// Обработка сообщений
func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, achievements []gv.Achievement, currentIndex *int, chatID int64) {

	if gv.UserStates[chatID] == "edit_name_achievements" {
		if err := db.UpdateAchievementsInDatabase(chatID, achievements[*currentIndex].Name); err != nil {
			log.Printf("❌ Произошла ошибка при обновлении достижений: %v", err)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Произошла ошибка при обновлении достижений: %v", err))
			bot.Send(msg)
		} else {
			log.Printf("✅ Достижения для чата %d успешно обновлены", chatID)
			msg := tgbotapi.NewMessage(chatID, "✅ Достижение успешно обновлено")
			achievements[*currentIndex].Name = gv.UserData[chatID]["achievements"]
			bot.Send(msg)
		}
	} else if gv.UserStates[chatID] == "edit_file_achievements" {
		err := db.SaveFile(bot, update, chatID, false)
		if err != nil {
			log.Printf("❌ Не удалось сохранить файл: %v", err)
		} else {
			msg := tgbotapi.NewMessage(chatID, "✅ Файл успешно обновлен!")
			achievements[*currentIndex].Filename = update.Message.Document.FileName
			bot.Send(msg)
		}
	}
	sendAchievementMessage(bot, chatID, achievements, *currentIndex, gv.UserMessageChatID[chatID])
}

// Функция для проверки наличия эмодзи в строке
func containsEmoji(s string) bool {
	emojiRegex := regexp.MustCompile(`[\x{1F600}-\x{1F6FF}\x{1F300}-\x{1F5FF}\x{1F900}-\x{1F9FF}\x{1F680}-\x{1F6FF}]`)
	return emojiRegex.MatchString(s)
}

func startMenu(bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Добро пожаловать в бота для сбора и каталогизации информации о достижениях студентов.\n\nДля начала работы бота необходимо выбрать, кем вы являетесь.")

	var inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👨‍🎓 Студент", "button_student"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👨‍🏫 Преподаватель", "button_teacher"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	sentMsg, _ := bot.Send(msg)
	gv.UserMessageChatID[chatId] = sentMsg.MessageID
}

func checkPassword(password string) bool {
	const correctPassword = "123"
	return password == correctPassword
}

func TelegramBot() {
	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatalf("Не удалось создать бота: %v", err)
	}
	// Устанавливаем время обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Получаем обновления от бота
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Не удалось получить обновления: %v", err)
	}

	for update := range updates {
		if (update.Message == nil) && (update.CallbackQuery == nil) {
			continue
		}

		if update.Message != nil {
			chatID := update.Message.Chat.ID

			if _, exists := gv.UserData[chatID]; !exists {
				gv.UserData[chatID] = make(map[string]string)
			}
			state := gv.UserStates[chatID]

			if state == "waiting_for_file" {
				err := db.SaveFile(bot, update, chatID, true)
				if err != nil {
					log.Printf("Не удалось сохранить файл: %v", err)
					continue
				} else {
					msg := tgbotapi.NewMessage(chatID, "Ваши достижения успешно сохранены!")
					msg.ReplyMarkup = cm.CreateAchievementAddedKeyboard()
					bot.Send(msg)
				}

				// Сбрасываем состояние и данные пользователя
				gv.UserStates[chatID] = ""

				continue
			} else if state == "waiting_for_new_file" {
				achievements, exists := gv.UserSearchResults[chatID]
				if !exists {
					responseText := "Результаты поиска не найдены."
					msg := tgbotapi.NewMessage(chatID, responseText)
					bot.Send(msg)
					continue
				}
				if index, exists := gv.UserAchievementIndex[chatID]; exists {
					gv.UserStates[chatID] = "edit_file_achievements"
					handleMessage(bot, update, achievements, &index, chatID)
					gv.UserAchievementIndex[chatID] = index
				} else {
					gv.UserAchievementIndex[chatID] = 0
					sendAchievementMessage(bot, chatID, achievements, 0, 0)
				}
				continue
			}

			// Проверяем, что от пользователя пришло именно текстовое сообщение
			if (reflect.TypeOf(update.Message.Text).Kind() == reflect.String) && (update.Message.Text != "") {

				if containsEmoji(update.Message.Text) {
					msg := tgbotapi.NewMessage(chatID, "Эмодзи не допускаются. Пожалуйста, используйте текстовые символы.")
					bot.Send(msg)
					continue
				}

				state := gv.UserStates[chatID]

				switch update.Message.Text {
				case "/start":
					gv.UserStates[chatID] = ""
					isStudent, ok, err := db.CheckStudentByChatID(chatID)
					if ok && (err == nil) {
						if isStudent {

							isRegistered, err := db.IsUserRegisteredStudent(chatID)
							if err != nil {
								responseText := fmt.Sprintf("Ошибка при проверке регистрации: %v", err)
								msg := tgbotapi.NewMessage(chatID, responseText)
								bot.Send(msg)
								continue
							} else if !isRegistered {
								startMenu(bot, chatID)
								continue
							} else {
								msg := tgbotapi.NewMessage(chatID, "Вы уже авторизированы.\nВыберите дальнейшее действие:")
								msg.ReplyMarkup = cm.CreateMenuKeyboard()
								bot.Send(msg)
								continue
							}
						} else {
							msg := tgbotapi.NewMessage(chatID, "Вы уже авторизированы.\nВыберите дальнейшее действие:")
							msg.ReplyMarkup = cm.CreateTeacherMenu()
							bot.Send(msg)
							continue
						}
					} else if err != nil {
						log.Printf("❌ Произошла ошибка при подключении к БД: %v", err)
						msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
						bot.Send(msg)
						continue
					}
					// Отправляем приветственное сообщение
					startMenu(bot, chatID)

				case "/help":
					gv.UserStates[chatID] = ""
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Доступные команды:\n/start - Перезапустить бота для регистрации\n/menu - Вызвать основное меню взаимодействий")
					bot.Send(msg)

				case "/menu":
					gv.UserStates[chatID] = ""
					isRegistered, err := db.IsUserRegistered(chatID)
					if err != nil {
						responseText := fmt.Sprintf("Ошибка при проверке регистрации: %v", err)
						msg := tgbotapi.NewMessage(chatID, responseText)
						bot.Send(msg)
						continue
					} else if !isRegistered {
						responseText := "Вы не зарегистрированы. Пожалуйста, зарегистрируйтесь."
						msg := tgbotapi.NewMessage(chatID, responseText)
						bot.Send(msg)
						continue
					} else {
						isStudent, ok, err := db.CheckStudentByChatID(chatID)
						if ok && (err == nil) {
							if isStudent {

								isRegistered, err := db.IsUserRegisteredStudent(chatID)
								if err != nil {
									responseText := fmt.Sprintf("Ошибка при проверке регистрации: %v", err)
									msg := tgbotapi.NewMessage(chatID, responseText)
									bot.Send(msg)
									continue
								} else if !isRegistered {
									responseText := "Вы не зарегистрированы. Пожалуйста, зарегистрируйтесь."
									msg := tgbotapi.NewMessage(chatID, responseText)
									bot.Send(msg)
									continue
								} else {
									msg := tgbotapi.NewMessage(chatID, "Выберите дальнейшее действие:")
									msg.ReplyMarkup = cm.CreateMenuKeyboard()
									bot.Send(msg)
									continue
								}
							} else {
								msg := tgbotapi.NewMessage(chatID, "Выберите дальнейшее действие:")
								msg.ReplyMarkup = cm.CreateTeacherMenu()
								bot.Send(msg)
								continue
							}
						} else if err != nil {
							log.Printf("❌ Произошла ошибка при подключении к БД: %v", err)
							msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
							bot.Send(msg)
							continue
						}
					}
					var inlineKeyboard = cm.CreateMenuKeyboard()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите дальнейшее действие:")
					msg.ReplyMarkup = inlineKeyboard
					bot.Send(msg)

				default:
					switch state {
					case "waiting_for_FIO":
						// Проверка, если первый символ ответа пользователя '/'
						if len(update.Message.Text) > 0 && update.Message.Text[0] == '/' {
							// Сообщаем пользователю, что ввод некорректен
							msg := tgbotapi.NewMessage(chatID, "Некорректный ввод. Пожалуйста, введите ваши ФИО без использования символа '/':")
							bot.Send(msg)
							continue
						}

						// Сохраняем ФИО и запрашиваем группу
						gv.UserData[chatID]["FIO"] = update.Message.Text
						gv.UserStates[chatID] = "waiting_for_group"
						msg := tgbotapi.NewMessage(chatID, "Введите вашу группу:")
						bot.Send(msg)
						continue
					case "waiting_for_new_FIO":
						// Проверка, если первый символ ответа пользователя '/'
						if len(update.Message.Text) > 0 && update.Message.Text[0] == '/' {
							// Сообщаем пользователю, что ввод некорректен
							msg := tgbotapi.NewMessage(chatID, "Некорректный ввод. Пожалуйста, введите ваши ФИО без использования символа '/':")
							bot.Send(msg)
							continue
						}
						gv.UserData[chatID]["FIO"] = update.Message.Text
						gv.UserStates[chatID] = "edit_fio_data"
						handleMessageData(bot, chatID)
					case "waiting_for_new_group":
						// Проверка, если первый символ ответа пользователя '/'
						if len(update.Message.Text) > 0 && update.Message.Text[0] == '/' {
							// Сообщаем пользователю, что ввод некорректен
							msg := tgbotapi.NewMessage(chatID, "Некорректный ввод. Пожалуйста, введите ваши ФИО без использования символа '/':")
							bot.Send(msg)
							continue
						}

						gv.UserData[chatID]["group"] = update.Message.Text
						gv.UserStates[chatID] = "edit_group_data"
						handleMessageData(bot, chatID)

					case "waiting_for_group":
						// Проверка, если первый символ ответа пользователя '/'
						if len(update.Message.Text) > 0 && update.Message.Text[0] == '/' {
							// Сообщаем пользователю, что ввод некорректен
							msg := tgbotapi.NewMessage(chatID, "Некорректный ввод. Пожалуйста, введите ваши ФИО без использования символа '/':")
							bot.Send(msg)
							continue
						}

						// Сохраняем группу и регистрируем пользователя
						gv.UserData[chatID]["group"] = update.Message.Text
						err := db.RegisterUser(chatID, gv.UserData[chatID]["FIO"], gv.UserData[chatID]["group"])
						if err != nil {
							log.Printf("❌ Не удалось зарегистрировать пользователя: %v", err)
							msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
							bot.Send(msg)
						} else {
							msg := tgbotapi.NewMessage(chatID, "Вы успешно зарегистрированы!")
							bot.Send(msg)
							// Отправляем меню после успешной регистрации
							msg = tgbotapi.NewMessage(chatID, "Выберите дальнейшее действие:")
							msg.ReplyMarkup = cm.CreateMenuKeyboard()
							bot.Send(msg)
						}
						// Сбрасываем состояние
						gv.UserStates[chatID] = ""

					case "waiting_for_achievements":
						// Проверка, если первый символ ответа пользователя '/'
						if len(update.Message.Text) > 0 && update.Message.Text[0] == '/' {
							// Сообщаем пользователю, что ввод некорректен
							msg := tgbotapi.NewMessage(chatID, "Некорректный ввод. Пожалуйста, введите ваши ФИО без использования символа '/':")
							bot.Send(msg)
							continue
						}
						// Сохраняем достижения и вставляем данные в базу данных
						gv.UserData[chatID]["achievements"] = update.Message.Text
						gv.UserStates[chatID] = "waiting_for_file"
						msg := tgbotapi.NewMessage(chatID, "Прикрепите файл, который потверждает ваше достижение\nФайл должен быть в формате pdf!\n")
						msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_cancel"),
							))
						bot.Send(msg)
						continue
					case "waiting_for_new_ach":
						// Проверка, если первый символ ответа пользователя '/'
						if len(update.Message.Text) > 0 && update.Message.Text[0] == '/' {
							// Сообщаем пользователю, что ввод некорректен
							msg := tgbotapi.NewMessage(chatID, "Некорректный ввод. Пожалуйста, введите ваши ФИО без использования символа '/':")
							bot.Send(msg)
							continue
						}
						// Сохраняем достижения и вставляем данные в базу данных
						gv.UserData[chatID]["achievements"] = update.Message.Text
						achievements, exists := gv.UserSearchResults[chatID]
						if !exists {
							responseText := "Результаты поиска не найдены."
							msg := tgbotapi.NewMessage(chatID, responseText)
							bot.Send(msg)
							continue
						}
						if index, exists := gv.UserAchievementIndex[chatID]; exists {
							gv.UserStates[chatID] = "edit_name_achievements"
							handleMessage(bot, update, achievements, &index, chatID)
							gv.UserAchievementIndex[chatID] = index
						} else {
							gv.UserAchievementIndex[chatID] = 0
							sendAchievementMessage(bot, chatID, achievements, 0, 0)
						}
						continue

					case "waiting_for_find_by_fio":
						//Поиск данных по ФИО
						query := update.Message.Text
						found, err := db.SearchDataByFIO(query)
						if err != nil {
							log.Printf("❌ Произошла ошибка при поиске студентов по ФИО: %v", err)
							msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
							bot.Send(msg)
						} else {
							if len(found) == 0 {
								msg := tgbotapi.NewMessage(chatID, "Ничего не найдено. \nВведите запрос снова либо нажмите '🚫 Отмена'")
								msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
									tgbotapi.NewInlineKeyboardRow(
										tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_back_teacher_search_menu"),
									))
								bot.Send(msg)
								continue
							}
							gv.UserStudentIndex[chatID] = 0
							gv.UserSearchStudents[chatID] = found
							msg := tgbotapi.NewMessage(chatID, "Результаты поиска:")
							msg.ReplyMarkup = cm.CreateInlineKeyboardMenuSearchStudent()
							sentMsg, _ := bot.Send(msg)
							gv.UserMessageChatID[chatID] = sentMsg.MessageID
							sendStudentsMessage(bot, chatID, found, 0, sentMsg.MessageID)
							// Сбрасываем состояние и данные пользователя
							gv.UserStates[chatID] = ""
						}

					case "waiting_for_find_by_achievements_upload":
						query := update.Message.Text
						found, err := db.SearchDataByAchievements(query)
						if err != nil {
							log.Printf("❌ Произошла ошибка при поиске студентов по ФИО: %v", err)
							msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
							bot.Send(msg)
						} else {
							if len(found) == 0 {
								msg := tgbotapi.NewMessage(chatID, "Ничего не найдено. \nВведите запрос снова либо нажмите '🚫 Отмена'")
								msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
									tgbotapi.NewInlineKeyboardRow(
										tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_download_menu"),
									))
								bot.Send(msg)
								continue
							}
							err = db.CreateExcelFile(found, int(chatID))
							if err != nil {
								log.Printf("❌ Произошла ошибка при создании таблицы: %v", err)
								msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
								bot.Send(msg)
								continue
							} else {
								err = db.SendExcelFile(bot, chatID)
								if err != nil {
									log.Printf("❌ Произошла ошибка при отправке таблицы: %v", err)
									msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
									bot.Send(msg)
									continue
								}
								gv.UserStates[chatID] = ""

							}
						}

					case "waiting_for_find_by_achievements":
						query := update.Message.Text
						found, err := db.SearchDataByAchievements(query)
						if err != nil {
							log.Printf("❌ Произошла ошибка при поиске студентов по ФИО: %v", err)
							msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
							bot.Send(msg)
						} else {
							if len(found) == 0 {
								msg := tgbotapi.NewMessage(chatID, "Ничего не найдено. \nВведите запрос снова либо нажмите '🚫 Отмена'")
								msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
									tgbotapi.NewInlineKeyboardRow(
										tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_back_teacher_search_menu"),
									))
								bot.Send(msg)
								continue
							}
							gv.UserStudentIndex[chatID] = 0
							gv.UserSearchResult[chatID] = found
							msg := tgbotapi.NewMessage(chatID, "Результаты поиска:")
							msg.ReplyMarkup = cm.CreateInlineKeyboardMenuAchievements()
							sentMsg, _ := bot.Send(msg)
							gv.UserMessageChatID[chatID] = sentMsg.MessageID
							sendStudentsMessageAchievements(bot, chatID, found, 0, sentMsg.MessageID)
							// Сбрасываем состояние и данные пользователя
							gv.UserStates[chatID] = ""
						}

					case "waiting_for_find_by_group_upload":
						query := update.Message.Text
						found, err := db.SearchDataByUsergroup(query)
						if err != nil {
							log.Printf("Произошла ошибка при поиске: %v", err)
							msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
							bot.Send(msg)
						} else {
							if len(found) == 0 {
								msg := tgbotapi.NewMessage(chatID, "Ничего не найдено. \nВведите запрос снова либо нажмите '🚫 Отмена'")
								msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
									tgbotapi.NewInlineKeyboardRow(
										tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_download_menu"),
									))
								bot.Send(msg)
								continue
							}
							err = db.CreateExcelFile(found, int(chatID))
							if err != nil {
								log.Printf("❌ Произошла ошибка при создании таблицы: %v", err)
								msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
								bot.Send(msg)
								continue
							} else {
								err = db.SendExcelFile(bot, chatID)
								if err != nil {
									log.Printf("❌ Произошла ошибка при отправке таблицы: %v", err)
									msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
									bot.Send(msg)
									continue
								}
								gv.UserStates[chatID] = ""
							}
						}

					case "waiting_for_password":
						password := update.Message.Text
						if checkPassword(password) {
							msg := tgbotapi.NewMessage(chatID, "✅ Вы успешно авторизовались")
							bot.Send(msg)

							err := db.AddStudent(chatID, false)
							if err != nil {
								responseText := fmt.Sprintf("Ошибка при выборе: %v", err)
								msg := tgbotapi.NewMessage(chatID, responseText)
								bot.Send(msg)
								gv.UserStates[chatID] = ""
								continue
							} else {
								deleteMsg := tgbotapi.NewDeleteMessage(chatID, gv.UserMessageChatID[chatID])
								bot.Send(deleteMsg)
								msg := tgbotapi.NewMessage(chatID, "Выберите дальнейшее действие:")
								msg.ReplyMarkup = cm.CreateTeacherMenu()
								bot.Send(msg)
								gv.UserStates[chatID] = ""
								delete(gv.UserMessageChatID, chatID)
								continue
							}
						} else {
							msg := tgbotapi.NewMessage(chatID, "Неыерный пароль.\n Повторите ввод")
							msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_back_start_menu"),
								))
							bot.Send(msg)
							continue
						}

					case "waiting_for_find_by_group":
						query := update.Message.Text
						found, err := db.SearchDataByGroup(query)
						if err != nil {
							log.Printf("Произошла ошибка при поиске: %v", err)
							msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
							bot.Send(msg)
						} else {
							if len(found) == 0 {
								msg := tgbotapi.NewMessage(chatID, "Ничего не найдено. \nВведите запрос снова либо нажмите '🚫 Отмена'")
								msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
									tgbotapi.NewInlineKeyboardRow(
										tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_back_teacher_search_menu"),
									))
								bot.Send(msg)
								continue
							}

							gv.UserStudentIndex[chatID] = 0
							gv.UserSearchStudents[chatID] = found
							msg := tgbotapi.NewMessage(chatID, "Результаты поиска:")
							msg.ReplyMarkup = cm.CreateInlineKeyboardMenuSearchStudent()
							sentMsg, _ := bot.Send(msg)
							gv.UserMessageChatID[chatID] = sentMsg.MessageID
							sendStudentsMessage(bot, chatID, found, 0, sentMsg.MessageID)
							// Сбрасываем состояние и данные пользователя
							gv.UserStates[chatID] = ""

						}
					case "waiting_for_find":
						// Поиск достижений
						query := update.Message.Text
						found, err := db.SearchAchievements(query, chatID)
						if err != nil {
							log.Printf("Произошла ошибка при поиске: %v", err)
							msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
							bot.Send(msg)
						} else {
							if len(found) == 0 {
								msg := tgbotapi.NewMessage(chatID, "Ничего не найдено. \nВведите запрос снова либо нажмите '🚫 Отмена'")
								msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
									tgbotapi.NewInlineKeyboardRow(
										tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_back_to_menu_achievements"),
									))
								bot.Send(msg)
								continue
							}
							gv.UserAchievementIndex[chatID] = 0
							gv.UserSearchResults[chatID] = found
							msg := tgbotapi.NewMessage(chatID, "Результаты поиска:")
							msg.ReplyMarkup = cm.CreateInlineKeyboard()
							sentMsg, _ := bot.Send(msg)
							gv.UserMessageChatID[chatID] = sentMsg.MessageID
							sendAchievementMessage(bot, chatID, found, 0, sentMsg.MessageID)
							// Сбрасываем состояние и данные пользователя
							gv.UserStates[chatID] = ""
						}

						// Сбрасываем состояние
						gv.UserStates[chatID] = ""

					default:
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, я не понимаю, о чем вы говорите")
						bot.Send(msg)
					}
				}
			} else {
				// Отправляем сообщение
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, используйте текстовые сообщения для команд и ввода данных.")
				bot.Send(msg)
			}
		}

		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID

			var responseText string

			switch update.CallbackQuery.Data {
			case "button_student":
				err := db.AddStudent(chatID, true)
				if err != nil {
					responseText = fmt.Sprintf("Ошибка при выборе: %v", err)
				} else {
					msg := tgbotapi.NewMessage(chatID, "Вам необходимо зарегистрироваться\nВведите ваше ФИО:")
					gv.UserStates[chatID] = "waiting_for_FIO"
					bot.Send(msg)
					deleteMsg := tgbotapi.NewDeleteMessage(chatID, update.CallbackQuery.Message.MessageID)
					bot.Send(deleteMsg)
					continue
				}
			case "button_teacher":
				responseText = "Введите пароль: "
				gv.UserStates[chatID] = "waiting_for_password"

			case "button_back_start_menu":
				gv.UserStates[chatID] = ""
				startMenu(bot, chatID)

			case "button_search_fio":
				responseText = "Поиск данных\nВведите фрагмент или полное ФИО: "
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_back_teacher_search_menu"),
					))
				gv.UserStates[chatID] = "waiting_for_find_by_fio"
				bot.Send(msg)
				continue

			case "button_search_group":
				responseText = "Поиск данных\nВведите группу:"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_back_teacher_search_menu"),
					))
				gv.UserStates[chatID] = "waiting_for_find_by_group"
				bot.Send(msg)
				continue

			case "button_search_achievements":
				responseText = "Поиск данных\nВведите фрагмент или полное название достижения:"
				gv.UserStates[chatID] = "waiting_for_find_by_achievements"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_back_teacher_search_menu"),
					))
				bot.Send(msg)
				continue

			case "button_upload_group":
				responseText = "Поиск данных\nВведите группу:"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_download_menu"),
					))
				gv.UserStates[chatID] = "waiting_for_find_by_group_upload"
				bot.Send(msg)
				continue

			case "button_upload_achievements":
				responseText = "Поиск данных\nВведите фрагмент или полное название достижения:"
				gv.UserStates[chatID] = "waiting_for_find_by_achievements_upload"

			case "button_upload_all":
				found, err := db.FetchData()
				if err != nil {
					log.Printf("❌ Произошла ошибка при поиске студентов по ФИО: %v", err)
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
					bot.Send(msg)
				} else {
					if len(found) == 0 {
						msg := tgbotapi.NewMessage(chatID, "Ничего не найдено. \nВведите запрос снова либо нажмите '🚫 Отмена'")
						msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_download_menu"),
							))
						bot.Send(msg)
						continue
					}
					err = db.CreateExcelFile(found, int(chatID))
					if err != nil {
						log.Printf("❌ Произошла ошибка при создании таблицы: %v", err)
						msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
						bot.Send(msg)
						continue
					} else {
						err = db.SendExcelFile(bot, chatID)
						if err != nil {
							log.Printf("❌ Произошла ошибка при отправке таблицы: %v", err)
							msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
							bot.Send(msg)
							continue
						}
						gv.UserStates[chatID] = ""
					}
				}

			case "button_register":
				isRegistered, err := db.IsUserRegistered(chatID)
				if err != nil {
					responseText = fmt.Sprintf("Ошибка при проверке регистрации: %v", err)
				} else if isRegistered {
					responseText = "Вы уже зарегистрированы."
				} else {
					responseText = "Для регистрации\nВведите ваше ФИО:"
					gv.UserStates[chatID] = "waiting_for_FIO"
					msg := tgbotapi.NewMessage(chatID, responseText)
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_cancel"),
						))
					bot.Send(msg)
					continue
				}
			case "button_continue":
				isRegistered, err := db.IsUserRegistered(chatID)
				if err != nil {
					responseText = fmt.Sprintf("Ошибка при проверке регистрации: %v", err)
				} else if !isRegistered {
					responseText = "Вы не зарегистрированы. Пожалуйста, зарегистрируйтесь."
				} else {
					// Создаем меню и отправляем его пользователю
					msg := tgbotapi.NewMessage(chatID, "Выберите дальнейшее действие:")
					msg.ReplyMarkup = cm.CreateMenuKeyboard()
					bot.Send(msg)
					continue
				}
			case "button_data":
				msg := tgbotapi.NewMessage(chatID, "Выберите дальнейшее действие:")
				msg.ReplyMarkup = cm.CreateInlineKeyboardData()
				sentMsg, _ := bot.Send(msg)
				gv.UserMessageChatID[chatID] = sentMsg.MessageID
				sendDataMessage(bot, chatID, sentMsg.MessageID)
				continue
			case "button_cancel":
				// Удаляем сообщение с кнопкой "Отмена"
				deleteMsg := tgbotapi.NewDeleteMessage(chatID, update.CallbackQuery.Message.MessageID)
				bot.Send(deleteMsg)
				gv.UserStates[chatID] = ""
				continue

			case "button_data_edit_fio":
				responseText = "Введите ФИО"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_cancel"),
					))
				gv.UserMessageChatID[chatID] = update.CallbackQuery.Message.MessageID
				gv.UserStates[chatID] = "waiting_for_new_FIO"
				bot.Send(msg)
				continue
			case "button_data_edit_group":
				responseText = "Введите группу"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_cancel"),
					))
				gv.UserStates[chatID] = "waiting_for_new_group"
				bot.Send(msg)
				continue
			case "button_add":
				responseText = "Добавление \nВведите ваше достижение:"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_cancel"),
					))
				bot.Send(msg)
				gv.UserStates[chatID] = "waiting_for_achievements"
				continue
			case "button_see":
				achievements, err := db.AllAchievements(chatID)
				if err != nil {
					log.Printf("Произошла ошибка при поиске всех элементов: %v", err)
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", err))
					bot.Send(msg)
				} else {
					if len(achievements) == 0 {
						msg := tgbotapi.NewMessage(chatID, "Ничего не найдено")
						bot.Send(msg)
						continue
					}
					gv.UserAchievementIndex[chatID] = 0
					gv.UserSearchResults[chatID] = achievements
					msg := tgbotapi.NewMessage(chatID, "Результаты поиска:")
					msg.ReplyMarkup = cm.CreateInlineKeyboard()
					sentMsg, _ := bot.Send(msg)
					gv.UserMessageChatID[chatID] = sentMsg.MessageID
					sendAchievementMessage(bot, chatID, achievements, 0, sentMsg.MessageID)
					// Сбрасываем состояние и данные пользователя
					gv.UserStates[chatID] = ""
				}
			case "button_find":
				responseText = "Поиск\n Введите фрагмент или полное наименования достижения:"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_back_to_menu_achievements"),
					))
				bot.Send(msg)
				gv.UserStates[chatID] = "waiting_for_find"
				continue
			case "button_edit_achievements":
				responseText = "Изменение достижения\n Введите фрагмент или полное наименования достижения для его поиска:"
				gv.UserStates[chatID] = "waiting_for_find"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_back_to_menu_achievements"),
					))
				bot.Send(msg)
				continue
			case "button_delete_menu":
				var inlineKeyboard = cm.CreateDeleteMenuKeyboard()
				responseText = "Выберите тип удаления: "
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = inlineKeyboard
				bot.Send(msg)
				continue
			case "button_delete_all":
				err := db.DeleteRecordsByChatIDFromFiles(chatID)
				var msg tgbotapi.MessageConfig
				if err != nil {
					msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Ошибка: %v", err))
					log.Printf("❌ Ошибка: %v", err)
				} else {
					msg = tgbotapi.NewMessage(chatID, "✅ Достижения успешно удалены")
				}
				bot.Send(msg)
			case "button_go_to_menu":
				var inlineKeyboard = cm.CreateTeacherMenu()
				responseText = "Выберите дальнейшее действие:"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = inlineKeyboard
				bot.Send(msg)
				continue
			case "button_teacher_search_menu":
				var inlineKeyboard = cm.CreateTeacherSearchMenu()
				responseText = "Выберите по какому параметру будете искать:"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = inlineKeyboard
				bot.Send(msg)
				continue
			case "button_download_menu":
				var inlineKeyboard = cm.CreateTeacherUploadMenu()
				responseText = "Выберите по какому параметру будете выгружать данные студентов:"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = inlineKeyboard
				bot.Send(msg)
				gv.UserStates[chatID] = ""
				continue

			case "button_back_teacher_menu":
				var inlineKeyboard = cm.CreateTeacherMenu()
				responseText = "Выберите дальнейшее действие:"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = inlineKeyboard
				bot.Send(msg)
				gv.UserStates[chatID] = ""
				continue

			case "button_back_teacher_search_menu":
				var inlineKeyboard = cm.CreateTeacherSearchMenu()
				responseText = "Выберите дальнейшее действие:"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = inlineKeyboard
				bot.Send(msg)
				gv.UserStates[chatID] = ""
				continue
			case "button_go_to_achievements":
				var inlineKeyboard = cm.CreateAchievementsMenuKeyboard()
				responseText = "Выберите дальнейшее действие:"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = inlineKeyboard
				bot.Send(msg)
				continue
			case "button_back_to_menu_achievements":
				var inlineKeyboard = cm.CreateAchievementsMenuKeyboard()
				responseText = "Выберите дальнейшее действие:"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = inlineKeyboard
				bot.Send(msg)
				gv.UserStates[chatID] = ""
				continue
			case "button_back_to_menu":
				var inlineKeyboard = cm.CreateMenuKeyboard()
				responseText = "Выберите дальнейшее действие:"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = inlineKeyboard
				bot.Send(msg)
				gv.UserStates[chatID] = ""
				continue
			case "button_edit_name_achievements":
				responseText = "Введите новое название достижения: "
				gv.UserStates[chatID] = "waiting_for_new_ach"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_cancel"),
					))
				bot.Send(msg)
				continue
			case "button_edit_file":
				responseText = "Прикрепите файл, который потверждает ваше достижение\nФайл должен быть в формате pdf!"
				gv.UserStates[chatID] = "waiting_for_new_file"
				msg := tgbotapi.NewMessage(chatID, responseText)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "button_cancel"),
					))
				bot.Send(msg)
				continue
			case "button_search_prev_t", "button_search_next_t", "button_students_achievements":
				students, exists := gv.UserSearchStudents[chatID]

				if !exists {
					responseText = "Результаты поиска не найдены."
					msg := tgbotapi.NewMessage(chatID, responseText)
					bot.Send(msg)
					continue
				}
				if index, exists := gv.UserStudentIndex[chatID]; exists {
					handleCallbackQueryStudent(bot, update.CallbackQuery, students, &index, chatID)
					gv.UserStudentIndex[chatID] = index
				} else {
					gv.UserStudentIndex[chatID] = 0
					sendStudentsMessage(bot, chatID, students, 0, 0)
				}
				continue

			case "button_delete_profile":
				err := db.DeleteRecordsByChatID(chatID)
				if err != nil {
					log.Printf("❌ Ошибка при удалении записей: %v\n", err)
				} else {
					log.Printf("✅ Записи успешно удалены.")
					startMenu(bot, chatID)
				}

			case "button_search_prev", "button_search_next", "button_upload", "button_delete_ach", "new_name_achievement":

				achievements, exists := gv.UserSearchResults[chatID]
				if !exists {
					responseText = "Результаты поиска не найдены."
					msg := tgbotapi.NewMessage(chatID, responseText)
					bot.Send(msg)
					continue
				}
				if index, exists := gv.UserAchievementIndex[chatID]; exists {
					handleCallbackQuery(bot, update.CallbackQuery, achievements, &index, chatID)
					gv.UserAchievementIndex[chatID] = index
				} else {
					gv.UserAchievementIndex[chatID] = 0
					sendAchievementMessage(bot, chatID, achievements, 0, 0)
				}
				continue
			case "button_search_prev_a", "button_search_next_a", "button_upload_search":
				data, exists := gv.UserSearchResult[chatID]
				if !exists {
					responseText = "Результаты поиска не найдены."
					msg := tgbotapi.NewMessage(chatID, responseText)
					bot.Send(msg)
					continue
				}
				if index, exists := gv.UserStudentIndex[chatID]; exists {
					handleCallbackQueryAchievements(bot, update.CallbackQuery, data, &index)
					gv.UserStudentIndex[chatID] = index
				} else {
					gv.UserStudentIndex[chatID] = 0
					sendStudentsMessageAchievements(bot, chatID, data, 0, 0)
				}
				continue
			case "button_search_prev_sa", "button_search_next_sa", "button_upload_search_a":
				log.Println(gv.UserAchievementIndex[chatID])
				data, exists := gv.UserSearchResults[chatID]
				if !exists {
					responseText = "Результаты поиска не найдены."
					msg := tgbotapi.NewMessage(chatID, responseText)
					bot.Send(msg)
					continue
				}
				if index, exists := gv.UserAchievementIndex[chatID]; exists {
					handleCallbackQuerySearchAchievements(bot, update.CallbackQuery, data, &index)
					gv.UserAchievementIndex[chatID] = index
				} else {
					gv.UserAchievementIndex[chatID] = 0
					sendAchievementMessage(bot, chatID, data, 0, 0)
				}
				continue

			default:
				responseText = "Неизвестная кнопка."
			}

			msg := tgbotapi.NewMessage(chatID, responseText)
			bot.Send(msg)
		}
	}
}
