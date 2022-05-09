package src

import (
    "fmt"
    "net"
    "os/exec"
    "strings"
)

var inIRC           bool            = false

func sendCommand() {
    text, _ := input.GetText()
    input.SetText("")
    if(len(text) <= 0) {return}
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
            // !!! THIS IS DISCOURAGED! MOST COMMANDS SHOULD BE PROGRAMS.
            // THIS SHOULD ONLY BE USED FOR THE BUILT IN IRC SHIT AND TESTING.
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
                        go func() {
                            err := join()
                            if(err == nil) {
                                inIRC = true
                            } else {
                                sendToTextarea(err.Error())
                            }
                        }()
                    }
                // net testing shit
                case "pinglocal":
                    conn, err := net.Dial("tcp", ":48889")
                    if(err != nil) {
                        fmt.Println(err)
                        return
                    }
                    fmt.Fprintf(conn, "Ping!")
                    conn.Close()
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
                        // some text editors fuck with unicode ig so we have to compare with this
                        clear := string([]byte{27,91,72,27,91,50,74,27,91,51,74})
                        clear_x := string([]byte{27,91,72,27,91,50,74})
                        switch(v) {
                            case clear:
                                textareaBuffer = textareaBuffer[:0]
                                updateTextarea()
                            case clear_x:
                                _, height := win.GetSize()
                                areaMax := (height/20)
                                for i := 0; i < areaMax; i++ {
                                    sendToTextarea("")
                                }
                            default: sendToTextarea(v)
                        }
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