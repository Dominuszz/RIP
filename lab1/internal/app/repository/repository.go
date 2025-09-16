package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type ComplexClass struct {
	ID          int
	IMG         string
	Complexity  string
	Degree      string
	Description string
}

func (r *Repository) GetComplexClasses() ([]ComplexClass, error) {
	complexityClasses := []ComplexClass{
		{
			ID:          1,
			IMG:         "linear.jpg",
			Complexity:  "n",
			Degree:      "Линейная",
			Description: "время прямо пропорционально размеру данных",
		},
		{
			ID:          2,
			IMG:         "log.jpg",
			Complexity:  "Log(n)",
			Degree:      "Логарифмическая",
			Description: "время растет медленно, каждый шаг уменьшает задачу вдвое.",
		},
		{
			ID:          3,
			IMG:         "square.jpg",
			Complexity:  "n^2",
			Degree:      "Квадратичная",
			Description: "время растет пропорционально квадрату размера данных (вложенные циклы).",
		},
		{
			ID:          4,
			IMG:         "const.jpg",
			Complexity:  "1",
			Degree:      "Константная",
			Description: "выполняется за фиксированное время, независимо от размера данных.",
		},
		{
			ID:          5,
			IMG:         "loglin.jpg",
			Complexity:  "n log(n)",
			Degree:      "Линейно-логарифмическая",
			Description: "сочетание линейного и логарифмического роста, часто в эффективных алгоритмах сортировки.",
		},
		{
			ID:          6,
			IMG:         "exponential.jpg",
			Complexity:  "2^n",
			Degree:      "Экспоненциальная",
			Description: "Очень быстрый рост времени. Характерен для задач полного перебора.",
		},
		{
			ID:          7,
			IMG:         "fact.jpg",
			Complexity:  "n!",
			Degree:      "Факториальная",
			Description: "время растет факториально, самый быстрый рост, практически нерешаемо для больших n.",
		},
	}

	if len(complexityClasses) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return complexityClasses, nil
}
func (r *Repository) GetComplexClass(id int) (ComplexClass, error) {
	complexityClasses, err := r.GetComplexClasses()
	if err != nil {
		return ComplexClass{}, err
	}

	for _, сomplexClass := range complexityClasses {
		if сomplexClass.ID == id {
			return сomplexClass, nil
		}
	}
	return ComplexClass{}, fmt.Errorf("заказ не найден")
}
func (r *Repository) GetComplexClassByDegree(title string) ([]ComplexClass, error) {
	complexityClasses, err := r.GetComplexClasses()
	if err != nil {
		return []ComplexClass{}, err
	}

	var result []ComplexClass
	for _, сomplexClass := range complexityClasses {
		if strings.Contains(strings.ToLower(сomplexClass.Degree), strings.ToLower(title)) {
			result = append(result, сomplexClass)
		}
	}

	return result, nil
}
func (r *Repository) GetCart() ([]ComplexClass, error) {
	complexityClasses, err := r.GetComplexClasses()
	if err != nil {
		return []ComplexClass{}, err
	}

	var result []ComplexClass
	for _, compclass := range complexityClasses {
		if compclass.ID == 1 || compclass.ID == 2 {
			result = append(result, compclass)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}

	return result, nil
}
