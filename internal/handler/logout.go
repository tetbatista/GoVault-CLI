package handler

import (
	"fmt"

	"github.com/tetbatista/govault/internal/infrastructure/auth"
)

func Logout() {
	if err := auth.DestroySession(); err != nil {
		fmt.Println("❌ No active session to terminate.")
		return
	}
	fmt.Println("✅ Session terminated.")
}