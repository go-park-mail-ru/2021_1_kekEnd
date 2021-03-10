package localstorage

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"strconv"
	"sync"
)

type MovieLocalStorage struct {
	movies    map[string]*models.Movie
	currentID uint64
	mutex     sync.Mutex
}

func NewMovieLocalStorage() *MovieLocalStorage {
	// dummy data for testing
	movies := map[string]*models.Movie{
		"1": {
			ID:             "1",
			Title:          "Чужой",
			Description:    "Группа космонавтов высаживается на неизвестной планете и знакомится с ксеноморфом. Шедевр Ридли Скотта",
			Voiceover:      []string{"Русский", "Английский"},
			Subtitles:      []string{"Русские"},
			Quality:        "HD",
			ProductionYear: 1979,
			Country:        []string{"Великобритания", "США"},
			Genre:          []string{"Хоррор", "Драма"},
			Slogan:         "slogan",
			Director:       "Ридли Скотт",
			Scriptwriter:   "---",
			Producer:       "---",
			Operator:       "---",
			Composer:       "---",
			Artist:         "---",
			Montage:        "---",
			Budget:         "---",
			Duration:       "---",
			Actors:         []string{"Сигурни Уивер", "Иэн Холм"},
		},

		"2": {
			ID:             "2",
			Title:          "Назад в будущее",
			Description:    "Безумный ученый и 17-летний оболтус тестируют машину времени и наводят шороху в 1950-х. Классика кинофантастики",
			Voiceover:      []string{"Русский"},
			Subtitles:      []string{"Нет"},
			Quality:        "HD",
			ProductionYear: 1985,
			Country:        []string{"США"},
			Genre:          []string{"Фантастика", "Комедия"},
			Slogan:         "«Семнадцатилетний Марти МакФлай пришел вчера домой пораньше. На 30 лет раньше»",
			Director:       "Роберт Земекис",
			Scriptwriter:   "---",
			Producer:       "---",
			Operator:       "---",
			Composer:       "---",
			Artist:         "---",
			Montage:        "---",
			Budget:         "---",
			Duration:       "---",
			Actors:         []string{"Майкл Дж. Фокс", "Кристофер Ллойд"},
		},

		"3": {
			ID:             "3",
			Title:          "Титаник",
			Description:    "Гибель легендарного лайнера на фоне запретной любви. Великий фильм-катастрофа — в отреставрированной версии",
			Voiceover:      []string{"Русский", "Английский"},
			Subtitles:      []string{"Русские"},
			Quality:        "HD",
			ProductionYear: 1997,
			Country:        []string{"США", "Мексика"},
			Genre:          []string{"Мелодрама", "История"},
			Slogan:         "«Ничто на Земле не сможет разлучить их»",
			Director:       "Джеймс Кэмерон",
			Scriptwriter:   "---",
			Producer:       "---",
			Operator:       "---",
			Composer:       "---",
			Artist:         "---",
			Montage:        "---",
			Budget:         "---",
			Duration:       "---",
			Actors:         []string{"Леонардо ДиКаприо", "Кейт Уинслет"},
		},

		"4": {
			ID:             "4",
			Title:          "Зеленая книга",
			Description:    "Путешествие итальянца-вышибалы и чернокожего пианиста — комедия-лауреат премии «Оскар» за лучший фильм",
			Voiceover:      []string{"Русский", "Английский"},
			Subtitles:      []string{"Русские", "Английские"},
			Quality:        "HD",
			ProductionYear: 2018,
			Country:        []string{"США", "Китай"},
			Genre:          []string{"Комедия", "Драма"},
			Slogan:         "«Inspired by a True Friendship»",
			Director:       "Питер Фаррелли",
			Scriptwriter:   "---",
			Producer:       "---",
			Operator:       "---",
			Composer:       "---",
			Artist:         "---",
			Montage:        "---",
			Budget:         "---",
			Duration:       "---",
			Actors:         []string{"Вигго Мортенсен", "Махершала Али"},
		},

		"5": {
			ID:             "5",
			Title:          "Властелин колец: Братство кольца",
			Description:    "Фродо Бэггинс отправляется спасать Средиземье. Первая часть культовой фэнтези-трилогии Питера Джексона",
			Voiceover:      []string{"Русский", "Английский"},
			Subtitles:      []string{"Русские"},
			Quality:        "HD",
			ProductionYear: 2001,
			Country:        []string{"Новая Зеландия", "США"},
			Genre:          []string{"Фэнтези", "Приключения"},
			Slogan:         "«Power can be held in the smallest of things...»",
			Director:       "Питер Джексон",
			Scriptwriter:   "---",
			Producer:       "---",
			Operator:       "---",
			Composer:       "---",
			Artist:         "---",
			Montage:        "---",
			Budget:         "---",
			Duration:       "---",
			Actors:         []string{"Элайджа Вуд", "Иэн Маккеллен"},
		},

		"6": {
			ID:             "6",
			Title:          "Побег из Шоушенка",
			Description:    "Выдающаяся драма о силе таланта, важности дружбы, стремлении к свободе и Рите Хэйворт",
			Voiceover:      []string{"Русский", "Английский"},
			Subtitles:      []string{"Русские", "Английские"},
			Quality:        "HD",
			ProductionYear: 1994,
			Country:        []string{"США"},
			Genre:          []string{"Драма"},
			Slogan:         "«Страх - это кандалы. Надежда - это свобода»",
			Director:       "Фрэнк Дарабонт",
			Scriptwriter:   "---",
			Producer:       "---",
			Operator:       "---",
			Composer:       "---",
			Artist:         "---",
			Montage:        "---",
			Budget:         "---",
			Duration:       "---",
			Actors:         []string{"Тим Роббинс", "Морган Фриман"},
		},
	}

	return &MovieLocalStorage{
		movies:    movies,
		currentID: 2,
	}
}

func (storage *MovieLocalStorage) CreateMovie(movie *models.Movie) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	movie.ID = strconv.FormatUint(storage.currentID, 10)
	storage.movies[movie.ID] = movie
	storage.currentID++

	return nil
}

func (storage *MovieLocalStorage) GetMovieByID(id string) (*models.Movie, error) {
	movie, exists := storage.movies[id]
	if !exists {
		return nil, errors.New("movie not found")
	}
	return movie, nil
}
