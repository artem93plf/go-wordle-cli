// Package game реализует игровую логику игры «5 букв».
package game

// DefaultDictionary возвращает новый слайс встроенных слов для каждой игры.
func DefaultDictionary() []string {
	return []string{
		"арбуз",
		"банан",
		"берег",
		"билет",
		"вагон",
		"ведро",
		"ветер",
		"вилка",
		"вишня",
		"голос",
		"город",
		"доска",
		"дятел",
		"жизнь",
		"забор",
		"земля",
		"кабан",
		"камин",
		"киоск",
		"книга",
		"кошка",
		"лампа",
		"лодка",
		"масло",
		"метла",
		"мороз",
		"океан",
		"парус",
		"паста",
		"песня",
		"пирог",
		"пламя",
		"полка",
		"почта",
		"радио",
		"робот",
		"ручка",
		"сапог",
		"сахар",
		"свеча",
		"слово",
		"сосна",
		"стена",
		"трава",
		"туман",
		"улица",
		"факел",
		"чашка",
		"школа",
		"ягода",
	}
}

// DictionaryByOptions возвращает словарь для заданных параметров.
func DictionaryByOptions(category Category, difficulty Difficulty) []string {
	switch category {
	case CategoryAnimal:
		switch difficulty {
		case DifficultyEasy:
			return []string{
				"кошка",
				"кабан",
				"зебра",
				"панда",
				"выдра",
				"хомяк",
				"белка",
				"бизон",
			}
		case DifficultyHard:
			return []string{
				"дятел",
				"коала",
				"гиена",
				"лемур",
				"сайга",
				"варан",
				"динго",
			}
		}

	case CategoryCity:
		switch difficulty {
		case DifficultyEasy:
			return []string{
				"томск",
				"париж",
				"токио",
				"пермь",
				"курск",
				"минск",
				"псков",
				"киров",
			}
		case DifficultyHard:
			return []string{
				"киото",
				"дакка",
				"аккра",
				"дакар",
				"лхаса",
				"салех",
				"табук",
			}
		}

	case CategoryFood:
		switch difficulty {
		case DifficultyEasy:
			return []string{
				"банан",
				"сахар",
				"пицца",
				"арбуз",
				"булка",
				"лимон",
				"сырок",
				"салат",
			}
		case DifficultyHard:
			return []string{
				"паста",
				"пирог",
				"рамен",
				"кимчи",
				"карри",
				"такос",
				"цукат",
				"кешью",
			}
		}
	}

	return nil
}
