## Web Programming lab #5: CLI browser (basically)

### Features:

- Raw TCP sockets
- HTML page parser
- Web search
- Various content: can read html, json, images 

### Architecture:

- Written in Go
- Uses Cobra as cli builder

### Commands:

```bash
go2web -u "github.com" --no-cache --max-redirects 3 # searches github, no cache, with max 3 redirects
```

```bash
go2web -s "hot dogs" -e "mojeek" -d # searches for hotdogs with mojeek in dynamic mode (can move though results) 
```

### Demo <3

