package manuals

// Man is the man page for `--man man`
const Man string = `
Man(1) User Commands Man(1)

NAME
       man - Used to get information about a bot command.

SYNOPSIS
       --man [command]

DESCRIPTION
       Man returns to the requester the informait about a command and how
       that command is to be used. The message is sent as a DM to the user
       requesting the information.

EXAMPLE
       --man pong
`

// Meme is the man page for `--man man`
const Meme string = `
MEME(1) User Commands MEME(1)

NAME
       meme - Create a meme on the fly.

SYNOPSIS
       --meme [meme name] [top text] [bottom text]

DESCRIPTION
       Meme sends back a new meme created with the text the user provides
       by calling the memegen.link api. Text with spaces must be formatted
       as this_is_top or this-is-top.

COMMAND SYNOPSIS
       This is just a brief synopsis of --meme commands to serve as a
       reminder to those who already know the memegen.link api; for more
       information please refer to https://memegen.link/api/

       --meme list
              Returns the template url on memegen.link with all available images.
       
       --meme spongebob why_are_you here
              Returns the url to the newly created mocking spongebob meme with
              the top text "why are you" and the bottom text "here".
       
       --meme spongebob why_are_you
              Returns the url to the newly created mocking spongebob meme with
              the top text "why are you" and no bottom text.
`

// Ping is the man page for `--man rule34`
const Ping string = `
PING(1) User Commands PING(1)

NAME
       ping - Used to check the bots responce.

SYNOPSIS
       --ping

DESCRIPTION
       Ping will return "pong" when the bot is online.

EXAMPLE
        --ping
`

// Rule34 is the man page for `--man rule34`
const Rule34 string = `
RULE34(1) User Commands RULE34(1)

NAME
       rule34 - NSFW command to grab a random image from rule34.xxx

SYNOPSIS
       --rule34 [tag]

DESCRIPTION
       Rule34 states that "if it exists there is porn for it." This command
       will only trigger in a channel marked as NFSW.

EXAMPLE
        --rule34 boobs
`
