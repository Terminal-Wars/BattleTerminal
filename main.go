package main

import (
    "time"

    "github.com/Terminal-Wars/BattleTerminal/src"
)

func main() {
    src.Start = time.Now()
    src.WinInit()
    src.WinBuild()
    src.WinLoop()
}