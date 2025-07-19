#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>
#include <sys/time.h>

// Character sets excluding similar characters (0, O, I, l, 1)
const char UPPERCASE[] = "ABCDEFGHJKLMNPQRSTUVWXYZ";
const char LOWERCASE[] = "abcdefghijkmnpqrstuvwxyz";
const char NUMBERS[] = "23456789";
const char SPECIAL[] = "!@#$%^&*()_+-=[]{}|;:,.<>?";

void print_usage(const char* program_name) {
    printf("Usage: %s [OPTIONS]\n", program_name);
    printf("Options:\n");
    printf("  -l LENGTH    Password length (default: 12)\n");
    printf("  -s           Include special characters\n");
    printf("  -c COUNT     Number of passwords to generate (default: 1)\n");
    printf("  -h           Show this help message\n");
    printf("\nExamples:\n");
    printf("  %s                    # Generate 12-character password\n", program_name);
    printf("  %s -l 16 -s           # Generate 16-character password with special chars\n", program_name);
    printf("  %s -l 10 -c 5         # Generate 5 passwords of 10 characters each\n", program_name);
}

char get_random_char(const char* charset) {
    int len = strlen(charset);
    return charset[rand() % len];
}

void shuffle_string(char* str, int length) {
    for (int i = length - 1; i > 0; i--) {
        int j = rand() % (i + 1);
        char temp = str[i];
        str[i] = str[j];
        str[j] = temp;
    }
}

void generate_password(char* password, int length, int include_special) {
    int pos = 0;
    
    // Ensure at least one character from each required set
    password[pos++] = get_random_char(UPPERCASE);
    password[pos++] = get_random_char(LOWERCASE);
    password[pos++] = get_random_char(NUMBERS);
    
    if (include_special && length >= 4) {
        password[pos++] = get_random_char(SPECIAL);
    }
    
    // Fill remaining positions randomly
    for (int i = pos; i < length; i++) {
        int charset_choice;
        
        if (include_special) {
            charset_choice = rand() % 4;
        } else {
            charset_choice = rand() % 3;
        }
        
        switch (charset_choice) {
            case 0:
                password[i] = get_random_char(UPPERCASE);
                break;
            case 1:
                password[i] = get_random_char(LOWERCASE);
                break;
            case 2:
                password[i] = get_random_char(NUMBERS);
                break;
            case 3:
                password[i] = get_random_char(SPECIAL);
                break;
        }
    }
    
    // Shuffle the password to randomize character positions
    shuffle_string(password, length);
    
    // Null terminate the string
    password[length] = '\0';
}

int main(int argc, char* argv[]) {
    int length = 12;
    int include_special = 0;
    int count = 1;
    int opt;
    
    // Parse command line arguments
    while ((opt = getopt(argc, argv, "l:sc:h")) != -1) {
        switch (opt) {
            case 'l':
                length = atoi(optarg);
                if (length < 3) {
                    fprintf(stderr, "Error: Password length must be at least 3\n");
                    return 1;
                }
                if (length > 128) {
                    fprintf(stderr, "Error: Password length cannot exceed 128\n");
                    return 1;
                }
                break;
            case 's':
                include_special = 1;
                break;
            case 'c':
                count = atoi(optarg);
                if (count < 1) {
                    fprintf(stderr, "Error: Count must be at least 1\n");
                    return 1;
                }
                if (count > 100) {
                    fprintf(stderr, "Error: Count cannot exceed 100\n");
                    return 1;
                }
                break;
            case 'h':
                print_usage(argv[0]);
                return 0;
            default:
                print_usage(argv[0]);
                return 1;
        }
    }
    
    // Check minimum length when special characters are required
    if (include_special && length < 4) {
        fprintf(stderr, "Error: Password length must be at least 4 when using special characters\n");
        return 1;
    }
    
    // Initialize random seed with high resolution time
    struct timeval tv;
    gettimeofday(&tv, NULL);
    srand(tv.tv_sec ^ tv.tv_usec);
    
    // Generate passwords
    char password[129]; // Max length + null terminator
    
    printf("Generated password%s:\n", count > 1 ? "s" : "");
    printf("Length: %d characters\n", length);
    printf("Character sets: Uppercase, Lowercase, Numbers");
    if (include_special) {
        printf(", Special characters");
    }
    printf("\n");
    printf("Excluded similar characters: 0, O, I, l, 1\n\n");
    
    for (int i = 0; i < count; i++) {
        generate_password(password, length, include_special);
        printf("%d: %s\n", i + 1, password);
    }
    
    return 0;
}