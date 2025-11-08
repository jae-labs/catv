
<p align="center">
<img src="logo/logo-xs.png" alt="CATV Logo" width="120"/>
</p>

<p align="center">
  <a href="https://github.com/luiz1361/catv/actions"><img src="https://github.com/luiz1361/catv/actions/workflows/go.yml/badge.svg" alt="Build Status"></a>
  <a href="https://github.com/luiz1361/catv/releases"><img src="https://img.shields.io/github/tag/luiz1361/catv.svg" alt="GitHub tag"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  <a href="https://buymeacoffee.com/luiz1361"><img src="https://img.shields.io/badge/Buy%20Me%20A%20Coffee-donate-orange.svg?logo=buymeacoffee" alt="Buy Me A Coffee"></a>
</p>

CATV (Cards Against The Void) is a fast, minimal command-line tool for turning your notes into interactive flashcards and reviewing them with spaced repetition. Simply point CATV at your folder of markdown notes, and it uses Ollama's local AI models (LLMs) to automatically generate flashcards and quiz you in a colorful terminal interface. CATV is written in Go for minimal dependencies—all you need is the Ollama client app from [ollama.com](https://ollama.com) and the compiled binary from this repository. No cloud, no subscriptions, no hidden costs.

**Why CATV?**
- Effortlessly convert your markdown notes into flashcards using local AI (LLM)
- Review and reinforce knowledge with spaced repetition
- Enjoy a distraction-free, user-friendly terminal experience
- 100% private: your notes and flashcards never leave your device
- Secure and offline: no internet required, no data sent to third parties
- Free and open-source: no hidden costs or paywalls

Designed for simplicity, privacy, and security, working entirely offline and storing your cards locally in SQLite. Originally built and tested for macOS, it runs anywhere Ollama is supported.

## Quick Start

1. **Install Ollama Desktop:**
  - Download and install from [ollama.com](http://ollama.com/)

2. **Pull the llama3.1 model:**
  ```bash
  ollama pull llama3.1
  ```

3. **Download and set up:**
  ```bash
  curl -L https://github.com/luiz1361/catv/releases/latest/download/catv-darwin-arm64 -o catv && chmod +x catv
  ```

4. **Generate flashcards from your markdown notes:**
  ```bash
  ./catv generate --path <path-to-your-markdown-files>
  ```

5. **Review your flashcards:**
  ```bash
  ./catv
  ```

That's it! No extra configuration needed. It will use the local Ollama API and store flashcards in a SQLite database.



## Features

| Feature                        | Description                                         |
|--------------------------------|-----------------------------------------------------|
| AI Flashcard Generation        | Create flashcards from markdown using Ollama AI      |
| Spaced Repetition Review       | Review cards with spaced repetition algorithm        |
| Terminal User Interface        | Colorful, user-friendly TUI for reviewing cards      |
| SQLite Storage                 | Flashcards stored locally in SQLite database         |
| No Extra Configuration         | Works out-of-the-box with minimal setup              |

## FAQ
<details>
<summary>What are the system requirements?</summary>
Tested on a MacBook M2 with 16GB of RAM using the llama3.1 model. Performance and compatibility may vary on other systems.
</details>

<details>
<summary>What platforms are supported?</summary>
Any platform with Go and Ollama (tested on macOS)
</details>

<details>
<summary>Do I need an internet connection?</summary>
No, Ollama runs locally.
</details>

<details>
<summary>Where are flashcards stored?</summary>
In a local SQLite database.
</details>

<details>
<summary>Can I use my own markdown files?</summary>
Yes, just point CATV to your folder of markdown notes.
</details>

<details>
<summary>How do I update the Ollama model?</summary>
Use `ollama pull <model>` to update or change models.
</details>

<details>
<summary>Is CATV open source?</summary>
Yes, licensed under MIT.
</details>

<details>
<summary>How do I contribute?</summary>
Open a pull request or issue on GitHub.
</details>

## Screenshots

Below are some screenshots of CATV in action:

![Screenshot 1](screenshots/1.png)

![Screenshot 2](screenshots/2.png)

![Screenshot 3](screenshots/3.png)

![Screenshot 4](screenshots/4.png)

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
