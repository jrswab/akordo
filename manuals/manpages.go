package manuals

// Gif is the man page for `<prefix>man gif`
const gif string = `
Gif(1) User Commands Gif(1)

NAME
       gif - Selects a random gif to display based off the user's input.

SYNOPSIS
       <prefix>gif [tag]

DESCRIPTION
       Gif returns to the requester a single giy based off the input tag
       provided. The gif is sent to the same channel that the user executed
       the command within.

EXAMPLE
       <prefix>gif laugh
`

// Man is the man page for `<prefix>man man`
const man string = `
Man(1) User Commands Man(1)

NAME
       man - Used to get information about a bot command.

SYNOPSIS
       <prefix>man [command]

DESCRIPTION
       Man returns to the requester the information about a command and how
       that command is to be used. The message is sent as a DM to the user
       requesting the information.

EXAMPLE
       <prefix>man pong
`

// Meme is the man page for `<prefix>man man`
const meme string = `
MEME(1) User Commands MEME(1)

NAME
       meme - Create a meme on the fly.

SYNOPSIS
       <prefix>meme [meme name] [top text] [bottom text]

DESCRIPTION
       Meme sends back a new meme created with the text the user provides
       by calling the memegen.link api. Text with spaces must be formatted
       as this_is_top or this-is-top.

COMMAND SYNOPSIS
       This is just a brief synopsis of <prefix>meme commands to serve as a
       reminder to those who already know the memegen.link api; for more
       information please refer to https://memegen.link/api/

       <prefix>meme list
              Returns the template url on memegen.link with all available images.

       <prefix>meme spongebob why_are_you here
              Returns the url to the newly created mocking spongebob meme with
              the top text "why are you" and the bottom text "here".

       <prefix>meme spongebob why_are_you
              Returns the url to the newly created mocking spongebob meme with
              the top text "why are you" and no bottom text.
`

// Ping is the man page for `<prefix>man rule34`
const ping string = `
PING(1) User Commands PING(1)

NAME
       ping - Used to check the bots responce.

SYNOPSIS
       <prefix>ping

DESCRIPTION
       Ping will return "pong" when the bot is online.

EXAMPLE
        <prefix>ping
`

// Rule34 is the man page for `<prefix>man rule34`
const rule34 string = `
RULE34(1) User Commands RULE34(1)

NAME
       rule34 - NSFW command to grab a random image from rule34.xxx

SYNOPSIS
       <prefix>rule34 [tag]

DESCRIPTION
       Rule34 states that "if it exists there is porn for it." This command
       will only trigger in a channel marked as NFSW.

EXAMPLE
        <prefix>rule34 boobs
`
