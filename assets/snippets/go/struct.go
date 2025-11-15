package main

import "fmt"

type User struct {
    ID        int
    Name      string
    Email     string
    Age       int
}

func (u *User) Greet() string {
    return fmt.Sprintf("Hello, I'm %s", u.Name)
}

func main() {
    user := &User{
        ID:    1,
        Name:  "Alice",
        Email: "alice@example.com",
        Age:   30,
    }
    fmt.Println(user.Greet())
}
