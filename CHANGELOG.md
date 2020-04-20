# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed
- `man` to send as a direct message.

### Security

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