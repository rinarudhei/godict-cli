# godict-cli

`godict-cli` is a lightweight English dictionary command-line application written in Go.  
It fetches word definitions, phonetics, examples, and more from the free [Dictionary API](https://dictionaryapi.dev/).  

With `godict-cli`, you can quickly look up words directly in your terminal, with both **simple** and **verbose** output modes.

---

## Features

- 🔍 Look up English words from the command line  
- 📖 Simple mode: quick definition and example  
- 📚 Verbose mode: full dictionary-style output (phonetics, origin, synonyms, antonyms)  
- ⚡ Fast and minimal dependencies
- 🔈 Pronunciation audio 
- 🌐 Uses [Free Dictionary API](https://dictionaryapi.dev/)  

---

## Installation

Clone the repo and build:

```bash
git clone https://github.com/rinarudhei/godict-cli.git
cd godict-cli
go build -o godict-cli
```

Or Install directly with:

```bash
go install github.com/rinarudhei/godict-cli@latest
```

--- 

## Usage 

```bash
godict-cli [flags] <word>
```

## Flags

```bash
--api-timeout duration   Set API timeout in duration second (default 5ns)
--api-url string         Free Dictionary API (default "https://api.dictionaryapi.dev/api/v2/entries/en/")
-h, --help               Show help
-s, --sound              Play pronunciation audio
-v, --verbose            Verbose response (full dictionary output)
```

---

## Examples

### Simple response

```bash
godict-cli cool
```

```bash
Word:      cool
Phonetic:  /kuːl/
Meaning:   A moderate or refreshing state of cold; moderate temperature of the air between hot and cold; coolness.
Example:   in the cool of the morning
```


### Verbose response

```bash
godict-cli -v hello
```

### With audio pronunciation

```bash
godict-cli -s hello      # Play audio with simple output
godict-cli -v -s hello   # Play audio with verbose output
```

```bash
Word:     hello
Phonetic: həˈləʊ
Phonetics:
    - həˈləʊ (audio: https://ssl.gstatic.com/dictionary/static/sounds/20200429/hello--_gb_1.mp3)
    - hɛˈləʊ
Origin:   early 19th century: variant of earlier hollo ; related to holla.
Meanings:
    Part of Speech: exclamation
    Definition:     used as a greeting or to begin a phone conversation.
    Example:        hello there, Katie!
```

--- 

## Contribute

Contributions are welcome! 🎉

1. Fork the repository

2. Create a new feature branch (git checkout -b feature/awesome-feature)

3. Commit your changes (git commit -m "Add awesome feature")

4. Push to the branch (git push origin feature/awesome-feature)

5. Open a Pull Request

Please make sure your code is formatted (go fmt ./...) and passes basic linting before submitting.

--- 

## License

MIT License @ 2025 [rinarudhei]
