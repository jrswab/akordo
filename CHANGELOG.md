# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Unit tests
- The bot prefix to be a constant in controller/controller.go

### Changed
- ReceiveMessage to CheckMessage and moved switch statement to its own method

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