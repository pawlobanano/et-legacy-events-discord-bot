## ET: Legacy Events Discord bot
A helper tool for tournaments management.

### Screenshot examples
`!cup edition 1` | `!cup teams` | Google Sheets UI
:-:|:-:|:-:
![!cup edition 1 command example](/assets/cup-edition-1-command-example.png)|![!cup teams command example](/assets/cup-teams-command-example.png)|![Google Sheets example](/assets/google-sheets-example.png)

### Bot keywords
```
!cup edition <number> | e <num>
!cup help | h
!cup teams | t
(TODO) !cup team <letter>
!cupbot status | s
```

### Run bot
```sh
make run 
```

### API Credentials
- .env
- key.json

### Notes 
1. [Discord API](https://discord.com/developers/docs/intro)
2. [Discord Markdown formatting - examples](https://support.discord.com/hc/en-us/articles/210298617-Markdown-Text-101-Chat-Formatting-Bold-Italic-Underline-)
3. [Google Sheets API - setup example](https://thriveread.com/golang-google-sheets-and-spreadsheet-api)
4. [Go formatting - Google standard](https://google.github.io/styleguide/go/decisions)
5. [Go routing - research, examples, benchmarks](https://benhoyt.com/writings/go-routing)
6. [Go ServerMux and Handlers](https://www.alexedwards.net/blog/an-introduction-to-handlers-and-servemuxes-in-go)
