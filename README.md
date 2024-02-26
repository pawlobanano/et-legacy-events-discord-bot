## Tournament Discord Bot
A helper tool for tournaments management. With use of the Google Sheets as a database. 
Inspired by community Draft Cup started on a Discord channel for the game of Wolfenstein: Enemy Territory Legacy.

### Screenshot examples
`!cup edition 1` | `!cup team e` | `!cup teams` | Google Sheets UI
:-:|:-:|:-:|:-:
![!cup edition 1 command example](/assets/cup-edition-1-command-example.png)|![!cup team e command example](/assets/cup-team-e-command-example.png)|![!cup teams command example](/assets/cup-teams-command-example.png)|![Google Sheets example](/assets/google-sheets-example.png)

### Bot commands
```
!cup edition <number> | e <num>
!cup help | h
!cup playthrough | !cup p
!cup play quarterfinals | !cup play q | !cup p q
!cup team <letter> | t <let>
!cup teams | ts
!cupbot status | s
```
_*All commands require an admin role and privately messaging the bot. List of admins is read from one of the Google Sheet cells._

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
7. [Responsive Tournament Bracket - Codepen](https://codepen.io/jimmyhayek/pen/yJZdEB)
