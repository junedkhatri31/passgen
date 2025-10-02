package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
)

// Character sets excluding similar characters (0, O, I, l, 1)
const (
	uppercase = "ABCDEFGHJKLMNPQRSTUVWXYZ"
	lowercase = "abcdefghijkmnpqrstuvwxyz"
	numbers   = "23456789"
	special   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

func printUsage(programName string) {
	fmt.Printf("Usage: %s [OPTIONS]\n", programName)
	fmt.Println("Options:")
	fmt.Println("  -l LENGTH    Password length (default: 12)")
	fmt.Println("  -s           Include special characters")
	fmt.Println("  -c COUNT     Number of passwords to generate (default: 1)")
	fmt.Println("  -h           Show this help message")
	fmt.Println("\nExamples:")
	fmt.Printf("  %s                    # Generate 12-character password\n", programName)
	fmt.Printf("  %s -l 16 -s           # Generate 16-character password with special chars\n", programName)
	fmt.Printf("  %s -l 10 -c 5         # Generate 5 passwords of 10 characters each\n", programName)
}

func getRandomChar(charset string) (byte, error) {
	max := big.NewInt(int64(len(charset)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return charset[n.Int64()], nil
}

func shuffleString(str []byte) error {
	length := len(str)
	for i := length - 1; i > 0; i-- {
		max := big.NewInt(int64(i + 1))
		jBig, err := rand.Int(rand.Reader, max)
		if err != nil {
			return err
		}
		j := jBig.Int64()
		str[i], str[j] = str[j], str[i]
	}
	return nil
}

func generatePassword(length int, includeSpecial bool) (string, error) {
	// Validate minimum length
	minLength := 3
	if includeSpecial {
		minLength = 4
	}
	if length < minLength {
		return "", fmt.Errorf("password length must be at least %d", minLength)
	}

	password := make([]byte, length)
	pos := 0

	// Ensure at least one character from each required set
	var err error
	password[pos], err = getRandomChar(uppercase)
	if err != nil {
		return "", err
	}
	pos++

	password[pos], err = getRandomChar(lowercase)
	if err != nil {
		return "", err
	}
	pos++

	password[pos], err = getRandomChar(numbers)
	if err != nil {
		return "", err
	}
	pos++

	if includeSpecial && length >= 4 {
		password[pos], err = getRandomChar(special)
		if err != nil {
			return "", err
		}
		pos++
	}

	// Fill remaining positions randomly
	for i := pos; i < length; i++ {
		var charsetChoice int
		if includeSpecial {
			max := big.NewInt(4)
			n, err := rand.Int(rand.Reader, max)
			if err != nil {
				return "", err
			}
			charsetChoice = int(n.Int64())
		} else {
			max := big.NewInt(3)
			n, err := rand.Int(rand.Reader, max)
			if err != nil {
				return "", err
			}
			charsetChoice = int(n.Int64())
		}

		switch charsetChoice {
		case 0:
			password[i], err = getRandomChar(uppercase)
		case 1:
			password[i], err = getRandomChar(lowercase)
		case 2:
			password[i], err = getRandomChar(numbers)
		case 3:
			password[i], err = getRandomChar(special)
		}

		if err != nil {
			return "", err
		}
	}

	// Shuffle the password to randomize character positions
	if err := shuffleString(password); err != nil {
		return "", err
	}

	return string(password), nil
}

func main() {
	length := flag.Int("l", 12, "Password length")
	includeSpecial := flag.Bool("s", false, "Include special characters")
	count := flag.Int("c", 1, "Number of passwords to generate")
	help := flag.Bool("h", false, "Show help message")

	flag.Parse()

	if *help {
		printUsage(os.Args[0])
		os.Exit(0)
	}

	// Validate input
	if *length < 3 {
		fmt.Fprintln(os.Stderr, "Error: Password length must be at least 3")
		os.Exit(1)
	}
	if *length > 128 {
		fmt.Fprintln(os.Stderr, "Error: Password length cannot exceed 128")
		os.Exit(1)
	}
	if *count < 1 {
		fmt.Fprintln(os.Stderr, "Error: Count must be at least 1")
		os.Exit(1)
	}
	if *count > 100 {
		fmt.Fprintln(os.Stderr, "Error: Count cannot exceed 100")
		os.Exit(1)
	}
	if *includeSpecial && *length < 4 {
		fmt.Fprintln(os.Stderr, "Error: Password length must be at least 4 when using special characters")
		os.Exit(1)
	}

	// Generate passwords
	plural := ""
	if *count > 1 {
		plural = "s"
	}

	fmt.Printf("Generated password%s:\n", plural)
	fmt.Printf("Length: %d characters\n", *length)
	fmt.Print("Character sets: Uppercase, Lowercase, Numbers")
	if *includeSpecial {
		fmt.Print(", Special characters")
	}
	fmt.Println()
	fmt.Println("Excluded similar characters: 0, O, I, l, 1")
	fmt.Println()

	for i := 0; i < *count; i++ {
		password, err := generatePassword(*length, *includeSpecial)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating password: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%d: %s\n", i+1, password)
	}
}