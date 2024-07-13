package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	gv "tg_bot/globalVar"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

// Получаем переменные окружения для подключения к базе данных
var (
	host     = os.Getenv("HOST")
	port     = os.Getenv("PORT")
	user     = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
	dbname   = os.Getenv("DBNAME")
	sslmode  = os.Getenv("SSLMODE")
	dbInfo   = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
)

// CreateTables создает необходимые таблицы в базе данных
func CreateTables() error {
	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	// Создаем таблицу users
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users
		(ID SERIAL PRIMARY KEY, 
		CHATID BIGINT UNIQUE, 
		FIO TEXT, 
		USERGROUP TEXT);`)
	if err != nil {
		return err
	}

	// Создаем таблицу files
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS files
		(ID SERIAL PRIMARY KEY, 
		CHATID BIGINT,  
		ACHIEVEMENTS TEXT,
		FILENAME TEXT,
		DATA BYTEA);`)
	if err != nil {
		return err
	}

	// Создаем таблицу students
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students
		(ID SERIAL PRIMARY KEY, 
		CHATID BIGINT UNIQUE,  
		ISSTUDENT BOOLEAN,
		DATA BYTEA);`)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRecordsByChatID удаляет записи, связанные с chatID, из всех таблиц
func DeleteRecordsByChatID(chatID int64) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	// Начало транзакции
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Удаление записей из таблицы files
	_, err = tx.Exec(`DELETE FROM files WHERE chatid = $1`, chatID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление записей из таблицы students
	_, err = tx.Exec(`DELETE FROM students WHERE chatid = $1`, chatID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление записей из таблицы users
	_, err = tx.Exec(`DELETE FROM users WHERE chatid = $1`, chatID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Завершение транзакции
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// AddStudent добавляет запись о студенте
func AddStudent(chatID int64, isStudent bool) error {
	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO students (CHATID, ISSTUDENT) VALUES ($1, $2)`
	_, err = db.Exec(query, chatID, isStudent)
	if err != nil {
		msg := fmt.Sprintf("❌ Невозможно добавить элемент в таблицу: %v", err)
		log.Println(msg)
		return fmt.Errorf(msg)
	}
	return nil
}

// CheckStudentByChatID проверяет, существует ли запись о студенте и возвращает значение ISSTUDENT
func CheckStudentByChatID(chatID int64) (bool, bool, error) {

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return false, false, err
	}
	defer db.Close()

	var isStudent bool
	query := `SELECT ISSTUDENT FROM students WHERE CHATID = $1`
	err = db.QueryRow(query, chatID).Scan(&isStudent)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, false, nil
		}
		return false, false, nil
	}
	return isStudent, true, nil
}

// SearchDataByAchievements ищет записи по подстроке достижения
func SearchDataByAchievements(substr string) ([]gv.Data, error) {
	var data []gv.Data

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return data, err
	}
	defer db.Close()

	searchPattern := "%" + substr + "%"

	// Выполнение запроса для получения данных о достижениях
	query := `SELECT chatid, achievements, filename FROM files WHERE achievements ILIKE $1`

	rows, err := db.Query(query, searchPattern)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	// Обработка результатов
	for rows.Next() {
		var d gv.Data
		err := rows.Scan(&d.Achievement.Id, &d.Achievement.Name, &d.Achievement.Filename)
		if err != nil {
			return data, err
		}

		// Выполнение запроса для получения данных о пользователях
		queryUser := `SELECT id, chatid, fio, usergroup FROM users WHERE chatid = $1`
		row := db.QueryRow(queryUser, d.Achievement.Id)

		var Student gv.Student
		err = row.Scan(&Student.Id, &Student.ChatID, &Student.Fio, &Student.Group)
		if err != nil {
			return data, err
		}

		d.Student = Student
		data = append(data, d)
	}

	// Проверка на ошибки при итерации по строкам
	if err = rows.Err(); err != nil {
		return data, err
	}

	return data, nil
}

// ByFio реализует интерфейс sort.Interface для []gv.Data на основе поля Fio
type ByFio []gv.Data

func (a ByFio) Len() int           { return len(a) }
func (a ByFio) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFio) Less(i, j int) bool { return a[i].Student.Fio < a[j].Student.Fio }

// SearchDataByUsergroup ищет записи по группе пользователя
func SearchDataByUsergroup(usergroup string) ([]gv.Data, error) {
	var data []gv.Data

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return data, err
	}
	defer db.Close()

	// Выполнение запроса для получения данных о пользователях
	queryUsers := `SELECT id, chatid, fio, usergroup FROM users WHERE usergroup = $1`
	rowsUsers, err := db.Query(queryUsers, usergroup)
	if err != nil {
		return data, err
	}
	defer rowsUsers.Close()

	// Обработка результатов для пользователей
	var studentsMap = make(map[int]gv.Student)
	for rowsUsers.Next() {
		var Student gv.Student
		err := rowsUsers.Scan(&Student.Id, &Student.ChatID, &Student.Fio, &Student.Group)
		if err != nil {
			return data, err
		}
		studentsMap[Student.ChatID] = Student
	}

	// Проверка на ошибки при итерации по строкам
	if err = rowsUsers.Err(); err != nil {
		return data, err
	}

	// Если не найдено пользователей с данной группой, возвращаем пустой массив
	if len(studentsMap) == 0 {
		return data, nil
	}

	// Подготовка списка chatid для запроса достижений
	var userChatIDs []int
	for chatID := range studentsMap {
		userChatIDs = append(userChatIDs, chatID)
	}

	// Подготовка запроса для получения данных о достижениях
	queryAchievements := `SELECT chatid, achievements, filename FROM files WHERE chatid = ANY($1)`

	rowsAchievements, err := db.Query(queryAchievements, pq.Array(userChatIDs))
	if err != nil {
		return data, err
	}
	defer rowsAchievements.Close()

	// Обработка результатов для достижений
	for rowsAchievements.Next() {
		var d gv.Data
		var chatID int
		err := rowsAchievements.Scan(&chatID, &d.Achievement.Name, &d.Achievement.Filename)
		if err != nil {
			return data, err
		}
		// Дублирование данных о студенте при необходимости
		if Student, ok := studentsMap[chatID]; ok {
			d.Student = Student
			data = append(data, d)
		}
	}

	// Проверка на ошибки при итерации по строкам
	if err = rowsAchievements.Err(); err != nil {
		return data, err
	}

	// Сортировка слайса data по полю Fio
	sort.Sort(ByFio(data))

	return data, nil
}

// FetchData извлекает все данные пользователей и их достижения
func FetchData() ([]gv.Data, error) {
	var data []gv.Data

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return data, err
	}
	defer db.Close()

	// Выполнение запроса для получения данных о пользователях
	queryUsers := `SELECT id, chatid, fio, usergroup FROM users`
	rowsUsers, err := db.Query(queryUsers)
	if err != nil {
		return data, err
	}
	defer rowsUsers.Close()

	// Обработка результатов для пользователей
	var studentsMap = make(map[int]gv.Student)
	for rowsUsers.Next() {
		var Student gv.Student
		err := rowsUsers.Scan(&Student.Id, &Student.ChatID, &Student.Fio, &Student.Group)
		if err != nil {
			return data, err
		}
		studentsMap[Student.ChatID] = Student
	}

	// Проверка на ошибки при итерации по строкам
	if err = rowsUsers.Err(); err != nil {
		return data, err
	}

	// Если не найдено пользователей, возвращаем пустой массив
	if len(studentsMap) == 0 {
		return data, nil
	}

	// Подготовка списка chatid для запроса достижений
	var userChatIDs []int
	for chatID := range studentsMap {
		userChatIDs = append(userChatIDs, chatID)
	}

	// Подготовка запроса для получения данных о достижениях
	queryAchievements := `SELECT chatid, achievements, filename FROM files WHERE chatid = ANY($1)`

	rowsAchievements, err := db.Query(queryAchievements, pq.Array(userChatIDs))
	if err != nil {
		return data, err
	}
	defer rowsAchievements.Close()

	// Обработка результатов для достижений
	for rowsAchievements.Next() {
		var d gv.Data
		var chatID int
		err := rowsAchievements.Scan(&chatID, &d.Achievement.Name, &d.Achievement.Filename)
		if err != nil {
			return data, err
		}
		// Дублирование данных о студенте при необходимости
		if Student, ok := studentsMap[chatID]; ok {
			d.Student = Student
			data = append(data, d)
		}
	}

	// Проверка на ошибки при итерации по строкам
	if err = rowsAchievements.Err(); err != nil {
		return data, err
	}

	// Сортировка слайса data по полю Fio
	sort.Sort(ByFio(data))

	return data, nil
}

func SearchDataByGroup(group string) ([]gv.Student, error) {
	// Слайс для хранения результатов
	var students []gv.Student

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return students, err
	}
	defer db.Close()

	// Выполнение запроса
	query := `SELECT id, chatid, fio, usergroup FROM users WHERE usergroup = $1`

	rows, err := db.Query(query, group)
	if err != nil {
		return students, err
	}
	defer rows.Close()

	// Обработка результатов
	for rows.Next() {
		var st gv.Student
		err := rows.Scan(&st.Id, &st.ChatID, &st.Fio, &st.Group)
		if err != nil {
			return students, err
		}
		students = append(students, st)
	}

	// Проверка на ошибки при итерации по строкам
	err = rows.Err()
	if err != nil {
		return students, err
	}
	return students, nil
}

func SearchDataByFIO(substr string) ([]gv.Student, error) {
	// Слайс для хранения результатов
	var students []gv.Student

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return students, err
	}
	defer db.Close()

	searchPattern := "%" + substr + "%"

	// Выполнение запроса
	query := `SELECT id, chatid, fio, usergroup FROM users WHERE fio ILIKE $1`

	rows, err := db.Query(query, searchPattern)
	if err != nil {
		return students, err
	}
	defer rows.Close()

	// Обработка результатов
	for rows.Next() {
		var st gv.Student
		err := rows.Scan(&st.Id, &st.ChatID, &st.Fio, &st.Group)
		if err != nil {
			return students, err
		}
		students = append(students, st)
	}

	// Проверка на ошибки при итерации по строкам
	err = rows.Err()
	if err != nil {
		return students, err
	}
	return students, nil
}

// Проверяем, существует ли пользователь
func IsUserRegistered(chatID int64) (bool, error) {
	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return false, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM students WHERE chatid = $1)`, chatID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func IsUserRegisteredStudent(chatID int64) (bool, error) {
	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return false, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE chatid = $1)`, chatID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Регистрируем пользователя
func RegisterUser(chatID int64, fio, userGroup string) error {
	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	// Вставляем нового пользователя
	_, err = db.Exec(`INSERT INTO users(chatid, fio, usergroup) VALUES($1, $2, $3)`, chatID, fio, userGroup)
	if err != nil {
		return err
	}

	return nil
}

// Обновляем ФИО пользователя
func EditFIO(chatID int64) error {
	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	// Обновляем ФИО пользователя
	_, err = db.Exec(`UPDATE users SET fio = $1 WHERE chatid = $2`, gv.UserData[chatID]["FIO"], chatID)
	if err != nil {
		return err
	}

	return nil
}

func EditGroup(chatID int64) error {
	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	// Обновляем ФИО пользователя
	_, err = db.Exec(`UPDATE users SET usergroup = $1 WHERE chatid = $2`, gv.UserData[chatID]["group"], chatID)
	if err != nil {
		return err
	}

	return nil
}

// SaveFile сохраняет файл в базе данных
func SaveFile(bot *tgbotapi.BotAPI, update tgbotapi.Update, chatID int64, new bool) error {
	if update.Message.Document == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка: пришел не файл\nПожалуйста, отправьте файл с расширением pdf")
		bot.Send(msg)
		return errors.New("пришел не файл")
	}

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	fileID := update.Message.Document.FileID
	// Проверка, что файл имеет расширение .pdf
	if filepath.Ext(update.Message.Document.FileName) != ".pdf" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка: файл не соответствует формату pdf\nПожалуйста, отправьте файл с расширением pdf")
		bot.Send(msg)
		return errors.New("файл не имеет расширения .pdf")
	}

	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return err
	}

	url := file.Link(bot.Token)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if new {
		_, err = db.Exec(`
			INSERT INTO files (chatid, achievements, filename, data) 
			VALUES ($1, $2, $3, $4)
		`, chatID, gv.UserData[chatID]["achievements"], update.Message.Document.FileName, data)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec(`
			UPDATE files 
			SET data = $4 
			WHERE chatid = $1 AND achievements = $2 AND filename = $3
		`, chatID, gv.UserData[chatID]["achievements"], update.Message.Document.FileName, data)
		if err != nil {
			return err
		}
	}

	return nil
}

// SendFile извлекает файлы из базы данных и отправляет их пользователю через Telegram
func SendFile(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, chatID int, achievements string) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	var filename string
	var data []byte

	// Updated query to include username and achievements
	query := `SELECT filename, data FROM files WHERE chatid = $1 AND achievements = $2`
	row := db.QueryRow(query, chatID, achievements)

	err = row.Scan(&filename, &data)
	if err != nil {
		return err
	}

	fileBytes := tgbotapi.FileBytes{Name: filename, Bytes: data}
	_, err = bot.Send(tgbotapi.NewDocumentUpload(callback.Message.Chat.ID, fileBytes))
	return err
}

func CreateExcelFile(data []gv.Data, chatID int) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	f := excelize.NewFile()
	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Headers
	f.SetCellValue(sheetName, "A1", "№")
	f.SetCellValue(sheetName, "B1", "ChatID")
	f.SetCellValue(sheetName, "C1", "Группа")
	f.SetCellValue(sheetName, "D1", "ФИО")
	f.SetCellValue(sheetName, "E1", "Достижение")

	// Fill in data
	row := 2
	for i, d := range data {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), d.Student.ChatID)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), d.Student.Group)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), d.Student.Fio)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), d.Achievement.Name)
		row++
	}

	f.SetColWidth(sheetName, "B", "B", 10)
	f.SetColWidth(sheetName, "C", "C", 6)
	f.SetColWidth(sheetName, "D", "D", 25)
	f.SetColWidth(sheetName, "D", "D", 25)
	f.SetActiveSheet(index)
	filename := "students_achievements.xlsx"
	if err := f.SaveAs(filename); err != nil {
		return err
	}

	// Read file content
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE students SET DATA = $1 WHERE CHATID = $2 AND ISSTUDENT = $3`,
		fileData, chatID, false)
	if err != nil {
		return err
	}

	return nil
}

func SendExcelFile(bot *tgbotapi.BotAPI, chatID int64) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	var data []byte
	filename := "students_achievements.xlsx"
	// Updated query to include username and achievements
	query := `SELECT data FROM students WHERE chatid = $1`
	row := db.QueryRow(query, chatID)

	err = row.Scan(&data)
	if err != nil {
		return err
	}

	fileBytes := tgbotapi.FileBytes{Name: filename, Bytes: data}
	_, err = bot.Send(tgbotapi.NewDocumentUpload(chatID, fileBytes))
	if err != nil {
		return err
	}

	// Удаление файла из базы данных после успешной отправки
	_, err = db.Exec(`UPDATE students SET data = NULL WHERE chatid = $1`, chatID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, "✅ Файл успешно выгружен")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Вернуться обратно в меню", "button_download_menu"),
		))
	bot.Send(msg)

	return nil
}

func GetFIOAndGroup(chatID int64) (string, string, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	var fio, usergroup string
	// Updated query to include username and achievements
	query := `SELECT fio, usergroup FROM users WHERE chatid = $1`
	row := db.QueryRow(query, chatID)

	err = row.Scan(&fio, &usergroup)
	if err != nil {
		return "", "", err
	}

	return fio, usergroup, nil

}

func SearchAchievements(substr string, chatID int64) ([]gv.Achievement, error) {
	// Слайс для хранения результатов
	var achievements []gv.Achievement

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return achievements, err
	}
	defer db.Close()

	searchPattern := "%" + substr + "%"

	// Выполнение запроса
	query := `SELECT id, achievements, filename, chatid FROM files WHERE achievements ILIKE $1 AND chatid = $2`

	rows, err := db.Query(query, searchPattern, chatID)
	if err != nil {
		return achievements, err
	}
	defer rows.Close()

	// Обработка результатов
	for rows.Next() {
		var ach gv.Achievement
		err := rows.Scan(&ach.Id, &ach.Name, &ach.Filename, &ach.ChatID)
		if err != nil {
			return achievements, err
		}
		achievements = append(achievements, ach)
	}

	// Проверка на ошибки при итерации по строкам
	err = rows.Err()
	if err != nil {
		return achievements, err
	}
	return achievements, nil
}

func AllAchievements(chatID int64) ([]gv.Achievement, error) {
	// Слайс для хранения результатов
	var achievements []gv.Achievement

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return achievements, err
	}
	defer db.Close()

	// Выполнение запроса
	query := `SELECT id, achievements, filename, chatid FROM files WHERE chatid = $1`

	rows, err := db.Query(query, chatID)
	if err != nil {
		return achievements, err
	}
	defer rows.Close()

	// Обработка результатов
	for rows.Next() {
		var ach gv.Achievement
		err := rows.Scan(&ach.Id, &ach.Name, &ach.Filename, &ach.ChatID)
		if err != nil {
			return achievements, err
		}
		achievements = append(achievements, ach)
	}

	// Проверка на ошибки при итерации по строкам
	err = rows.Err()
	if err != nil {
		return achievements, err
	}
	return achievements, nil
}

func DeleteRecord(chatID int64, achievements string) error {

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM files WHERE CHATID = $1 AND ACHIEVEMENTS = $2;`

	// Выполнение SQL-запроса
	result, err := db.Exec(query, chatID, achievements)
	if err != nil {
		return err
	}

	// Проверка количества удалённых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("Успешное удаления %d записией где chatID=%d и achievements=%s", rowsAffected, chatID, achievements)

	return nil
}

func DeleteRecordsByChatIDFromFiles(chatID int64) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	// Проверка на наличие записей с данным chatID
	var exists bool
	checkQuery := `SELECT EXISTS (SELECT 1 FROM files WHERE CHATID = $1);`
	err = db.QueryRow(checkQuery, chatID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("таблица пуста")
	}

	query := `DELETE FROM files WHERE CHATID = $1;`

	// Выполнение SQL-запроса
	result, err := db.Exec(query, chatID)
	if err != nil {
		return err
	}

	// Проверка количества удалённых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("✅ Успешное удаление %d записей где chatID=%d", rowsAffected, chatID)

	return nil
}

// Обновление достижений в базе данных
func UpdateAchievementsInDatabase(chatID int64, name string) error {

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	newAchievements, ok := gv.UserData[chatID]["achievements"]
	if !ok {
		return fmt.Errorf("не удалось получить новые достижения из userData для чата %d", chatID)
	}

	_, err = db.Exec(`UPDATE files SET ACHIEVEMENTS = $1 WHERE CHATID = $2 AND ACHIEVEMENTS = $3`, newAchievements, chatID, name)
	if err != nil {
		return fmt.Errorf("произошла ошибка при обновлении достижений: %v", err)
	}

	return nil
}
