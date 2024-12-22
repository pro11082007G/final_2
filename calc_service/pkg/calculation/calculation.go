package calculation

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	case '^':
		return 3
	}
	return 0
}

func applyOp(a, b float64, op rune) float64 {
	switch op {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		if b != 0 {
			return a / b
		}
		panic("division by zero")
	case '^':
		return math.Pow(a, b)
	}
	return 0
}

func Calc(expression string) (float64, error) {
	var values []float64
	var ops []rune

	// Убираем пробелы
	expression = strings.ReplaceAll(expression, " ", "")
	fmt.Println("Обработанное выражение:", expression)

	for i := 0; i < len(expression); i++ {
		ch := rune(expression[i])
		fmt.Printf("Текущий символ: %c\n", ch)

		// Если символ - цифра, парсим полное число
		if unicode.IsDigit(ch) {
			var numStr string
			for i < len(expression) && (unicode.IsDigit(rune(expression[i])) || expression[i] == '.') {
				numStr += string(expression[i])
				i++
			}
			i-- // Уменьшаем i, так как он будет увеличен в for
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, errors.New("invalid expression: insufficient values for operation")
			}
			values = append(values, num)
			fmt.Printf("Добавлено число: %f\n", num)
		} else if ch == '(' {
			ops = append(ops, ch)
			fmt.Println("Добавлены открывающие скобки")
		} else if ch == ')' {
			// Обрабатываем скобки
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				if len(values) < 2 {
					return 0, errors.New("invalid expression: insufficient values for operation")
				}
				val2 := values[len(values)-1]
				values = values[:len(values)-1]
				val1 := values[len(values)-1]
				values = values[:len(values)-1]
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				result := applyOp(val1, val2, op)
				values = append(values, result)
				fmt.Printf("Выполнена операция: %f %c %f = %f\n", val1, op, val2, result)
			}
			if len(ops) == 0 {
				return 0, errors.New("mismatched parentheses: too many closing parentheses")
			}
			ops = ops[:len(ops)-1] // Удаляем '('
		} else if precedence(ch) > 0 {
			// Обрабатываем операции
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(ch) {
				if len(values) < 2 {
					return 0, errors.New("invalid expression: insufficient values for operation")
				}
				val2 := values[len(values)-1]
				values = values[:len(values)-1]
				val1 := values[len(values)-1]
				values = values[:len(values)-1]
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				result := applyOp(val1, val2, op)
				values = append(values, result)
				fmt.Printf("Выполнена операция: %f %c %f = %f\n", val1, op, val2, result)
			}
			ops = append(ops, ch)
			fmt.Printf("Добавлен оператор: %c\n", ch)
		} else if !unicode.IsSpace(ch) {
			return 0, fmt.Errorf("invalid character: %c", ch)
		}
	}

	// Применяем оставшиеся операции
	for len(ops) > 0 {
		if len(values) < 2 {
			return 0, errors.New("invalid expression: insufficient values for operation")
		}
		val2 := values[len(values)-1]
		values = values[:len(values)-1]
		val1 := values[len(values)-1]
		values = values[:len(values)-1]
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]

		result := applyOp(val1, val2, op)
		values = append(values, result)
		fmt.Printf("Выполнена операция: %f %c %f = %f\n", val1, op, val2, result)
	}

	if len(values) == 1 {
		return values[0], nil
	}
	return 0, errors.New("invalid expression: unmatched operators and values")
}
