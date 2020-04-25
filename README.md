# Akordo
A Discord chat bot written by the members of the Technonauts chat group.

[Join the Technonauts!](https://discord.gg/A2uuCUr)

## Using the Bot:
1. [Setup your Discord application](https://discordapp.com/developers/applications/)
3. Invite your bot to your server
4. Download the source code `git pull git.sr.ht/~jrswab/akordo`
5. Build the binary: `go build`
6. Run the binary: `./akordo -t <your bot token here>`
  - You may omit `-t` when the setting `BOT_TOKEN` environment variable

### Currently required environment variables:
- `BOT_OWNER`
- `BOT_ID`
- `GIPHY_KEY` (only if using the `gif` command)
