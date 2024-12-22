package application

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/pro11082007G/calc_service/pkg/calculation"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

// Функция запуска приложения
func (a *Application) Run() error {
	for {
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
			continue
		}
		text = strings.TrimSpace(text)
		if text == "exit" {
			log.Println("application was successfully closed")
			return nil
		}

		result, err := calculation.Calc(text)
		if err != nil {
			log.Printf("Calculation failed for expression %s: %v\n", text, err)
		} else {
			log.Printf("%s = %f\n", text, result)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Проверка на валидность выражения
	if !isValidExpression(request.Expression) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Response{Error: "Expression is not valid"})
		return
	}

	result, err := calculation.Calc(request.Expression)
	response := Response{}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: "Internal server error"})
		return
	}

	response.Result = result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Функция для проверки валидности выражения
func isValidExpression(expr string) bool {
	// Регулярное выражение для проверки корректности выражения
	validExpr := regexp.MustCompile(`^[0-9+\-*/().\s]+$`)
	return validExpr.MatchString(expr)
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	log.Printf("Starting server on port %s...\n", a.config.Addr)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
