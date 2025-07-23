# Vault — Secure Password Manager

Encrypted CLI password manager written in Go.  
Stores account data in an encrypted file (`data.vault`) using AES-256.

## Key Features
- Create, search, and delete accounts
- Encrypt data with a secret key
- Store encrypted data in `data.vault`
- Read key from `.env` or environment variables
- Colorful CLI interface

## How to Run

### 1. Install dependencies
`go mod tidy`

### 2. Set up the encryption key
You must provide a 64-character hex-encoded key (32 bytes) for AES-256.

Option A: Use .env file (recommended)

`echo "KEY=your_32_char_hex_key_here" > .env`

You can generate a key through the terminal (Linux/macOS):
`openssl rand -hex 32`

Alternatively, you can generate an AES-256 key on special websites. But this is not a safe practice.

Example of adding a key:
`echo "KEY=bf19575b7a4bd0302f587c3a19d8bf3a9ed077fcb9706a6a1c775079a463f37f" > .env`

The old file is encrypted with the previous key and cannot be decrypted with a new one.
If you decide to change the key, you must first delete the `data.vault` file, which contains the encrypted data and is created after saving your account. 
Otherwise, the program will not run with the old encrypted file and the new key. 

Deleting `data.vault` will erase all saved accounts. 

Option B: Set environment variable

`KEY=your_32_char_hex_key_here go run main.go`

If you use this option, you will need to specify the AES-256. key each time you launch the app. If you decide to replace the key, remember to delete the `data.vault` file.

### 3. Run the app

`go run main.go`