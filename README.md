# 🦙 ALPAKA NOTES v2.0 - Ultimate TUI Edition

**Najpiękniejsza** aplikacja do notatek w terminalu! Stworzona z ❤️ używając Charm Bracelet Bubble Tea.

```
   ▄▄▄       ██▓     ██▓███   ▄▄▄       ██ ▄█▀▄▄▄      
  ▒████▄    ▓██▒    ▓██░  ██▒▒████▄     ██▄█▒▒████▄    
  ▒██  ▀█▄  ▒██░    ▓██░ ██▓▒▒██  ▀█▄  ▓███▄░▒██  ▀█▄  
  ░██▄▄▄▄██ ▒██░    ▒██▄█▓▒ ▒░██▄▄▄▄██ ▓██ █▄░██▄▄▄▄██ 
   ▓█   ▓██▒░██████▒▒██▒ ░  ░ ▓█   ▓██▒▒██▒ █▄▓█   ▓██▒
```

## ✨ Cechy Premium

### 🎨 **Piękny Interfejs**
- Gradient kolorów i animacje
- ASCII art logo
- Animowany splash screen
- Responsywny layout
- Kolorowe tagi i ikony
- Migający kursor

### 🔐 **Bezpieczeństwo**
- Szyfrowanie AES (wkrótce)
- Hash hasła SHA-256
- Własny format `.alpaka`
- Brak przechowywania hasła
- Pokaż/ukryj hasło (Ctrl+H)

### 📝 **Funkcje Notatek**
- Nieograniczona liczba notatek
- Wiele linii w treści
- System tagów
- Liczniki znaków
- Automatyczne timestampy
- Wyszukiwanie w czasie rzeczywistym

### 📊 **Statystyki i Analiza**
- Liczba notatek, słów, tagów
- Średnia słów na notatkę
- Tag cloud
- Ostatnia aktywność
- Wykresy użycia

### 🎯 **3 Tryby Widoku**
1. **Lista** - Szczegółowy podgląd wszystkich notatek
2. **Siatka** - Kompaktowy widok 2 kolumny
3. **Podgląd** - Pełny widok pojedynczej notatki

### 🔄 **3 Tryby Sortowania**
- Po dacie (najnowsze pierwsze)
- Po tytule (alfabetycznie)
- Po tagach

## 🚀 Instalacja

### Wymagania
- Go 1.18+ (zalecane 1.21+)

### Szybki start

```bash
# Sklonuj lub utwórz katalog
mkdir alpaka-notes && cd alpaka-notes

# Skopiuj pliki: main.go, screens.go, notebook.go

# Inicjalizuj moduł
go mod init github.com/alpaka/notes

# Pobierz zależności
go get github.com/charmbracelet/bubbletea@v0.23.2
go get github.com/charmbracelet/lipgloss@v0.7.1

# Uruchom!
go run .
```

### Kompilacja

```bash
# Linux
go build -o alpaka

# Windows
GOOS=windows GOARCH=amd64 go build -o alpaka.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o alpaka-mac
```

## 🎮 Sterowanie

### Globalne
- **Ctrl+C** - Wyjście z aplikacji
- **Esc** - Powrót do menu głównego
- **↑/↓** lub **j/k** - Nawigacja (Vim keys!)

### Ekran logowania
- Wpisz hasło
- **Ctrl+H** - Pokaż/ukryj hasło
- **Enter** - Zaloguj

### Menu główne
- **↑/↓** lub **j/k** - Wybór opcji
- **Enter** - Potwierdź
- **q** - Wyjście

### Dodawanie notatki
- **Tab** - Następne pole
- **Shift+Tab** - Poprzednie pole
- **Enter** - Nowa linia (w treści)
- **Ctrl+S** - Zapisz notatkę
- **Esc** - Anuluj

### Przeglądanie notatek
- **↑/↓** lub **j/k** - Przewijanie
- **d** - Usuń notatkę
- **v** - Zmień widok (Lista/Siatka/Podgląd)
- **s** - Zmień sortowanie
- **Esc** - Powrót

### Wyszukiwanie
- Wpisz zapytanie
- Wyniki w czasie rzeczywistym
- **Esc** - Powrót

### Statystyki
- Przeglądaj dane
- **Esc** - Powrót

### Ustawienia
- **↑/↓** - Wybór opcji
- **Enter/Space** - Zmień ustawienie
- **Esc** - Powrót

## 📁 Struktura projektu

```
alpaka-notes/
├── main.go          # Główna aplikacja + style
├── screens.go       # Wszystkie ekrany (Login, Menu, etc.)
├── notebook.go      # Model danych + szyfrowanie
├── go.mod           # Zależności
├── go.sum           # Checksums
├── README.md        # Ta dokumentacja
└── notatki.alpaka   # Twoje zaszyfrowane notatki
```

## 🎨 Paleta kolorów

```
Pink:    #FF6B9D  - Primary
Purple:  #C792EA  - Secondary
Blue:    #82AAFF  - Accent
Cyan:    #89DDFF  - Info
Green:   #C3E88D  - Success
Yellow:  #FFCB6B  - Warning
Orange:  #F78C6C  - Highlight
Red:     #FF5370  - Danger
```

## 📦 Format pliku .alpaka

```
ALPAKA
VERSION:1.0
HASH:<sha256_password_hash>
---ENCRYPTED---
<encrypted_json_data>
```

## 🔒 Bezpieczeństwo

**Obecna implementacja:**
- XOR cipher (demonstracja)
- SHA-256 hash hasła
- JSON serialization

**Planowane ulepszenia:**
- [ ] AES-256-GCM encryption
- [ ] PBKDF2/Argon2 key derivation
- [ ] Salt generation
- [ ] Backup encryption

## 🎯 Roadmap

### v2.1
- [ ] Eksport do Markdown/PDF
- [ ] Import z innych formatów
- [ ] Kategorie/foldery
- [ ] Przypięte notatki
- [ ] Archiwum

### v2.2
- [ ] Załączniki do notatek
- [ ] Obrazy inline
- [ ] Markdown rendering
- [ ] Syntax highlighting
- [ ] Motywy kolorystyczne

### v3.0
- [ ] Synchronizacja (opcjonalna)
- [ ] Współdzielenie notatek
- [ ] Web interface
- [ ] Mobile app
- [ ] Wtyczki

## 🏆 Funkcje Premium

✅ **Animowany splash screen** z gradientami  
✅ **5 ekranów** (Login, Menu, Dodaj, Przeglądaj, Statystyki, Ustawienia, Szukaj)  
✅ **ASCII art** logo z gradientem  
✅ **Migający kursor** we wszystkich polach  
✅ **Kolorowe tagi** (5 kolorów rotacyjnie)  
✅ **3 tryby widoku** notatek  
✅ **3 tryby sortowania**  
✅ **Statystyki** z licznikami i wykresami  
✅ **Wyszukiwanie** w czasie rzeczywistym  
✅ **Liczniki znaków** w formularzach  
✅ **Progress bar** przy ładowaniu  
✅ **Statusy** (Success/Error/Warning/Info)  
✅ **Vim keybindings** (j/k)  
✅ **Responsywny** layout  

## 📚 Biblioteki

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI Framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style & Layout

## 🤝 Wkład

Chcesz pomóc? Świetnie!

1. Fork projektu
2. Stwórz branch (`git checkout -b feature/amazing`)
3. Commit (`git commit -m 'Add amazing feature'`)
4. Push (`git push origin feature/amazing`)
5. Pull Request

## 📄 Licencja

MIT License
