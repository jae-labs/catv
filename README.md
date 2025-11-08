# CATV - Cards Against The Void

[![Build Status](https://github.com/luiz1361/catv/actions/workflows/go.yml/badge.svg)](https://github.com/luiz1361/catv/actions) [![GitHub tag](https://img.shields.io/github/tag/luiz1361/catv.svg)](https://github.com/luiz1361/catv/releases)

CATV (Cards Against The Void) is a minimal command-line utility for reviewing notes. 'Flashcards' can be written in markdown-like syntax.

**Goal:** Be super simple to use—just generate and review flashcards with minimal setup.

CATV uses Ollama's language models to create and review spaced repetition flashcards.

**Platform Note:** This project was made for macOS and tested on a MacBook M2 with 16GB of RAM, but it should work on any platform where Go and Ollama are supported.

## Quick Start

1. **Install Ollama Desktop:**
  - Download and install from [ollama.com](http://ollama.com/)

2. **Pull the llama3.1 model:**
  ```bash
  ollama pull llama3.1
  ```

3. **Clone and set up CATV:**
  ```bash
  git clone https://github.com/luiz1361/catv.git
  cd catv
  go mod tidy
  ```

4. **Generate flashcards from markdown files:**
  ```bash
  go run . generate --file <folder-with-markdowns>
  ```

5. **Review your flashcards:**
  ```bash
  go run .
  ```

That's it! No extra configuration needed. CATV will use the local Ollama API and store flashcards in a SQLite database.

## Features

- Generate flashcards from markdown files using AI
- Review flashcards with spaced repetition
- Colorful, user-friendly terminal interface

## Screenshots

Below are some screenshots of CATV (Cards Against The Void) in action:

![Screenshot 1](screenshots/1.png)

![Screenshot 2](screenshots/2.png)

![Screenshot 3](screenshots/3.png)

![Screenshot 4](screenshots/4.png)

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
