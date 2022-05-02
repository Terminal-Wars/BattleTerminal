package src

import (
    //"fmt"
    "os/exec"
    "strings"
)

var inIRC           bool            = false

func sendCommand() {
    text, _ := input.GetText()
    input.SetText("")
    // is our message prefixed with /?
    if(string(text[0]) == "/") {
        text := strings.Replace(text,"/","",2)
        // Are we in an IRC server?
        if(inIRC) {
            client.Write(text) // if yes, send it as an irc command.
        } else { // Otherwise, execute it as a system command.
            // execute a system command
            commands := strings.Split(text, " ")
            // exec command wants different strings for each argument but doesn't support arrays.
            // so we make do with that, only supporting five arguments (and even that's overkill).
            var cmd *exec.Cmd
            switch(len(commands)) {
                case 2: cmd = exec.Command(commands[0], commands[1])
                case 3: cmd = exec.Command(commands[0], commands[1], commands[2])
                case 4: cmd = exec.Command(commands[0], commands[1], commands[2], commands[3])
                case 5: cmd = exec.Command(commands[0], commands[1], commands[2], commands[3], commands[4])
                default: cmd = exec.Command(commands[0])
            }
            switch(commands[0]) {
                // built in IRC shit
                case "user", "nick":
                    if(len(commands) == 1) {
                        sendToTextarea("Usage: /"+commands[0]+" {desiredName}")
                    } else {
                        userID = commands[1]
                        sendToTextarea("You are now "+userID)
                    }
                case "join":
                    if(userID == "") {
                        sendToTextarea("Please use /user first to set your nickname")
                    } else {
                        inIRC = true
                        go func() {
                            err := join()
                            sendToTextarea(err)
                        }()
                    }
                // regular commands
                default:
                    output, err := cmd.Output()
                    var outputs []string
                    if(err != nil) {
                        outputs = strings.Split(err.Error(), "\n")
                    } else {
                        outputs = strings.Split(string(output), "\n")
                    }
                    for _, v := range outputs {
                        sendToTextarea(v)
                    }
            }
        }
    } else { // If not...
        // If we're in an IRC server, send it as an IRC message.
        if(inIRC) { 
            client.Write("PRIVMSG #testing "+text)
        } 
        // Regardless, put it in the terminal.
        sendToTextarea(text)
    }
}