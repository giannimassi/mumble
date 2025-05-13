# mumble

> **⚠️ WARNING: This project was built for my own entertainment and learning. I am not a security expert, and I strongly discourage using this tool for real-world applications or storing sensitive information. Use at your own risk!**

**Mumble** is a minimal, dependency-free CLI tool that generates secure, pronounceable passwords — fast to type, easy to remember, and hard to guess.

No special characters you need a PhD to type. Just strong, speakable passwords with optional symbols, clipboard support, and entropy validation.

---

## ✨ Features

- 🔐 Secure random password generation (`crypto/rand`)
- 👄 Pronounceable and typable on any US keyboard
- 🔢 Customizable length via `--length`
- 🔣 Optional symbols via `--symbols`
- 📋 Clipboard support via `--copy` (macOS/Linux)
- 🧠 Entropy estimation and enforcement (min 60 bits)
- ❌ Excludes ambiguous characters (`1`, `l`, `I`, `0`, `O`)

---

## 📦 Installation

Install directly using Go:

```bash
go install github.com/giannimassi/mumble@latest
```

> Requires Go 1.20+

## 🚀 Usage

```bash
mumble [--length=N] [--symbols=true|false] [--copy]
```

### Examples

```bash
mumble
# e.g. "futiboruda9o"
# Estimated entropy: 67.84 bits

mumble --length=16
# e.g. "dopemirukinozuva"

mumble --copy
# Copies password to clipboard (macOS/Linux)
```

## 🔐 How Secure Is It?

Passwords are generated using Go's crypto/rand for true randomness. Each generated password must pass:
 • Structure validation (must contain a consonant, vowel, digit, and optional symbol)
 • Minimum entropy check (>= 60 bits)

Entropy is estimated using:

```
Entropy = log2(effective charset size) × password length
```

This ensures all generated passwords are both usable and strong.

## 📄 License

MIT. Use it, fork it, improve it.
