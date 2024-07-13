package globalVar

// Achievement представляет структуру для строки в таблице files
type Achievement struct {
	Id       int
	Name     string
	Filename string
	ChatID   int64
}

type Student struct {
	Id     int
	ChatID int
	Fio    string
	Group  string
}

type Data struct {
	Student     Student
	Achievement Achievement
}

// Карта для хранения состояний пользователей
var UserStates = make(map[int64]string)

// Карта для хранения данных пользователей
var UserData = make(map[int64]map[string]string)

var UserAchievementIndex = make(map[int64]int)

var UserStudentIndex = make(map[int64]int)

var UserMessageChatID = make(map[int64]int)

// Добавляем эту карту для хранения результатов поиска
var UserSearchResults = make(map[int64][]Achievement)

var UserSearchStudents = make(map[int64][]Student)

var UserSearchResult = make(map[int64][]Data)
