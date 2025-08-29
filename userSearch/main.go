package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

var auditFile = "audit.data"

func main() {
	insertFile := "users.data"

	if _, err := os.Stat(insertFile); err == nil {
		fmt.Println("User Data File already exists.")
		//Retrieve(insertFile)
		Insert(insertFile)
	} else if os.IsNotExist(err) {
		info, err := os.Create(insertFile)
		if err != nil {
			fmt.Println("User Data File created : ", info)
		}
		Insert(insertFile)
		//Retrieve(insertFile)
	}
}

func Insert(insertFile string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter name (or type 'exit' to quit): ")
		name1, _ := reader.ReadString('\n')
		name := strings.TrimSpace(name1)
		if strings.ToLower(name) == "exit" {
			fmt.Println("Exiting input loop.")
			break
		}

		fmt.Print("Enter your email: ")
		emailRaw, _ := reader.ReadString('\n')
		email := strings.TrimSpace(emailRaw)

		cm1 := NewUser(uint(rand.IntN(100)), name, email, "Active")
		bytes, err := json.Marshal(cm1)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			continue
		}

		bytes = append(bytes, '\n')
		WritetoAudit(cm1,"Insert")
		f, err := os.OpenFile(insertFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer f.Close()

		if _, err := f.Write(bytes); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Println("Data appended successfully")
	}
}

func Retrieve(filePath string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter email to search (or type 'exit' to quit): ")
		email1, _ := reader.ReadString('\n')
		email := strings.TrimSpace(email1)
		if strings.ToLower(email) == "exit" {
			fmt.Println("Exiting input loop.")
			break
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		found := false

		for scanner.Scan() {
			line := scanner.Text()
			var user User
			if err := json.Unmarshal([]byte(line), &user); err != nil {
				fmt.Println("Error parsing JSON:", err)
				continue
			}

			if strings.EqualFold(user.Email, email) {
				fmt.Println("User found:")
				fmt.Printf("Id: %d\nName: %s\nEmail: %s\nStatus: %s\n", user.Id, user.Name, user.Email, user.Status)
				//WritetoAudit(scanner.Bytes())
				found = true
				break
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}

		if !found {
			fmt.Println("No user found with that email.")
		}
	}
}

func WritetoAudit(user *User,operation string) {
    f, err := os.OpenFile(auditFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error writing to log file:", err)
        return
    }
    defer f.Close()

    // Add timestamp to log entry
    timestamped := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), user.Name,user.Email,operation)
    f.WriteString(timestamped)
}


func NewUser(id uint, name string, email string, status string) *User {
	return &User{id, name, email, status}
}

type User struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}
