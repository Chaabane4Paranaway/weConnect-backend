// mailer.go - stub
package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SendMockOTP(email string, otp string) {
	fmt.Printf("📩 [MOCK] Sending OTP to %s: %s\n", email, otp)
}
