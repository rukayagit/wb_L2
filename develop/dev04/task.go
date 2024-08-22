package main

import (
	"fmt"
	"sort"
	"strings"
)

// FindAnagramGroups находит все множества анаграмм в переданном массиве слов.
// Возвращает мапу, где ключ — первое слово из множества, значение — срез слов из этого множества.
func FindAnagramGroups(words []string) map[string][]string {
	anagramMap := make(map[string][]string)

	// Приведение всех слов к нижнему регистру и удаление дубликатов
	wordSet := make(map[string]struct{})
	for _, word := range words {
		word = strings.ToLower(word)
		if _, exists := wordSet[word]; !exists {
			wordSet[word] = struct{}{}
		}
	}

	// Преобразование уникальных слов в срез
	var uniqueWords []string
	for word := range wordSet {
		uniqueWords = append(uniqueWords, word)
	}

	// Функция для сортировки букв в слове
	sortString := func(s string) string {
		sortedRunes := []rune(s)
		sort.Slice(sortedRunes, func(i, j int) bool {
			return sortedRunes[i] < sortedRunes[j]
		})
		return string(sortedRunes)
	}

	// Группировка слов по их анаграммам
	for _, word := range uniqueWords {
		sortedWord := sortString(word)
		anagramMap[sortedWord] = append(anagramMap[sortedWord], word)
	}

	// Формирование результата, исключая множества из одного элемента и дубли
	result := make(map[string][]string)
	for _, group := range anagramMap {
		if len(group) > 1 {
			sort.Strings(group)
			// Удаление дубликатов в итоговом срезе
			uniqueGroup := []string{}
			wordSet := make(map[string]struct{})
			for _, word := range group {
				if _, exists := wordSet[word]; !exists {
					wordSet[word] = struct{}{}
					uniqueGroup = append(uniqueGroup, word)
				}
			}
			if len(uniqueGroup) > 1 {
				result[uniqueGroup[0]] = uniqueGroup
			}
		}
	}

	return result
}

func main() {
	// Пример входных данных
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "слиток", "листок", "пятак"}

	// Вызов функции FindAnagramGroups
	anagramGroups := FindAnagramGroups(words)

	// Вывод результата
	for key, group := range anagramGroups {
		fmt.Printf("%s: %v\n", key, group)
	}
}
