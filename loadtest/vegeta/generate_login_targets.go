package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Create("login.targets")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	const totalUsers = 100 // целевая 100_000

	for i := 0; i < totalUsers; i++ {
		username := fmt.Sprintf("user_%06d", i)
		body := fmt.Sprintf(`{"username":"%s","password":"Password123"}`, username)

		// Каждый target: метод + URL, заголовок, пустая строка, тело, пустая строка
		target := fmt.Sprintf(
			"POST http://localhost:8080/api/auth\n"+
				"Content-Type: application/json\n\n"+
				"%s\n\n", body)

		if _, err := f.WriteString(target); err != nil {
			panic(err)
		}
	}
}
