package domain

// CalculateBestStreak recebe uma lista de is_on_diet ordenada por data ASC
// e retorna o maior número de refeições consecutivas dentro da dieta.
func CalculateBestStreak(onDietList []bool) int {
	best, current := 0, 0
	for _, onDiet := range onDietList {
		if onDiet {
			current++
			if current > best {
				best = current
			}
		} else {
			current = 0
		}
	}
	return best
}
