# Passgen

A secure password generator written in Go that creates cryptographically strong passwords.

## Features

- Cryptographically secure random password generation
- Customizable password length
- Optional special characters
- Generate multiple passwords at once
- Excludes similar-looking characters (0, O, I, l, 1) to avoid confusion
- Guarantees at least one character from each selected character set

## Installation

```sh
sudo curl -SL https://github.com/junedkhatri31/passgen/releases/download/v2.0.0/passgen -o /usr/local/bin/passgen
```
```sh
sudo chmod +x /usr/local/bin/passgen
```

### Build from Source

```bash
git clone https://github.com/junedkhatri31/passgen.git
cd passgen
go build -o passgen
```

## Usage

```bash
passgen [OPTIONS]
```

### Options

- `-l LENGTH` - Password length (default: 12)
- `-s` - Include special characters
- `-c COUNT` - Number of passwords to generate (default: 1)
- `-h` - Show help message

### Examples

Generate a 12-character password (default):
```bash
passgen
```

Generate a 16-character password with special characters:
```bash
passgen -l 16 -s
```

Generate 5 passwords of 10 characters each:
```bash
passgen -l 10 -c 5
```

## Character Sets

- **Uppercase**: A-Z (excluding O, I)
- **Lowercase**: a-z (excluding l)
- **Numbers**: 2-9 (excluding 0, 1)
- **Special** (optional): `!@#$%^&*()_+-=[]{}|;:,.<>?`

## Testing

Run the test suite:
```bash
go test -v
```

## License

This project is open source and available under the MIT License.
