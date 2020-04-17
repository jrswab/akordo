package manuals

// Man is the man page for `--man man`
const Man string = `
Man(1) User Commands Man(1)

NAME
       man - Used to get information about a bot command.

SYNOPSIS
       --man [command]

DESCRIPTION
       Man returns to the requester the informait about a comannd and how
       that command is to be used. The message is sent as a DM to the user
       requesting the information.

EXAMPLE
       --man pong
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
