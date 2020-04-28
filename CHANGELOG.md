# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Clear command to clear messages of the bot or a user for the past 100 messages.
- Deletion of command message after the bot returns the request
- The ability to add banned words.
- `bannedWords.json` to `.gitignore`
- Removal of user from the guild when a blacklisted word is used.
- Rule agreement command to give users a role to chat after reading the rules.
- A command to see an XP leader board

### Changed
- Controller.go `msgType` to be "chan" by default.
- `checkSyntax()` to `checkMessage()`
- The controller reply method to be unexported
- Saved data file location to the `data` directory
- `.gitignore` to ignore the new data directory
- AkSession interface to be located in `plugins.go`
- `loadSavedData()` to it's own package: `load.SavedData()`
- Bot check in `xp.go` to check the message bot bool instead of checking an ID
- `xp/commands_test.go` to use `reflect.DeepEqual()`

### Fixed
- Bug where xp was looking at the bot token instead of the bot ID environment variable.

## v0.8.0
### Added
- The roles package.
- Comments for context in xp/commands.go
- Bot owner check for xp add auto role command
- selfAssignRoles.json to .gitignore
- The embed massage type to contoller/controller.go

### Changed
- This file.
- File name variable for xp auto rank json
- load data statements into their own function in main.go

## v0.7.0
### Added
- XP system
- xp.json to gitignore
- .idea directory to gitingore (GoLand editor files)
- envars file to .gitignore
- .autoRanks.json to .gitignore
- XP reward to crypto game
- XP load, save, and autosave to main.go
- Unit tests for the XP package
- XP command catch to switch clause in controller.go
- Auto promotion of roles based on user XP

### Changed
- akordo in gitignore to akordo* to exclude binaries with version tags.
- con to sd as variable for `NewSessionData()` in main.go
- `checkLastAsk` to be exported for use by xp.go
- All every usage of `checkLastAsk` to `CheckLastAsk`
- Changed hard coded file name in main.go to the exported variable from xp.go

### Fixed
- gif and rule34 commands when they return only one image and caused a panic

## v0.6.0
### Added
- "Crypto Game" where the user has to decrypt a message to win

### Changed
- The default bot prefix to `=`

### Fixed
- `man` to send as a direct message.

## v0.5.0
### Added
- Unit tests
- Gif to manpages.go
- Execute function to controller/controller.go to call the correct plugin.
- A reply function to controller/controller.go to handle sending data back to the user.

### Changed
- The bot prefix to be specified in the sessionData struct
- ReceiveMessage to CheckMessage and moved switch statement to its own method
- Command switch statement to receive a string and error
- Command switch to send message upon success instead of each plugin.
- Plugins to return data instead of executing their own message send
- `Rule34` takes an AkSession (interface) instead of `*dg.Session`
- The manual method to return one value instead of two.

## v0.4.0
### Added
- Sending a random gif on user request `--gif <word>`
- The option to pass the bot token as an environment variable.
- Plugins.go to hold general plugin functions.
- Cool-down to meme, pong, and rule34 commands

### Changed
- main.go to create a struct to hold the data needed to check user request frequency.
- Logging to show the username of the requester

## Fixed
- `invalid argument to Intn` when rule34 returns an empty list.

### v0.3.0
### Added
- `--man` to return a UNIX style manual page for given commands.
- This changelog
- `CONTRIBUTING.md`

## v0.2.0
### Added
- Plugins directory for bot command logic.
- Controller to contain the logic used for parsing user request.
- Rule34 plugin (only works in a NSFW channel)
- MemeGen plugin to create custom memes in chat.
- `.gitignore` to ignore compiled binary

## v0.1.0
### Added
- Base for running the bot as a binary
- The ability to pass the bot token by cli flag
- Added Pong command