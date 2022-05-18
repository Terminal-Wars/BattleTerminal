package src

import (
	"fmt"
    "log"

    "github.com/jezek/xgb/xproto"
    "github.com/Terminal-Wars/TermUI"
)

var textareaBuffer          []string

const (
	// Labels
	TextareaID  			= 0
	// Textboxes
	InputID  				= 0
	// Buttons
	CharSwitchButtonID 		= 0
	AttackButtonID 			= 1
	BagSwitchButtonID 		= 2

	width 					= 320
	height 					= 240
)

var (
	win 					TermUI.Window

	textarea 				*TermUI.UIEvent
	input 					*TermUI.UIEvent

	err 					error
)

func InitWin() {
    win, err = TermUI.NewWindow(width,height,
        []uint32{
            0xffffffff,
            xproto.EventMaskKeyPress |
            xproto.EventMaskKeyRelease,
        })
    if(err != nil) {log.Fatalln(err)}
}

func FillWin() {
    // Main textarea
    textarea = win.Label("",
        TextareaID,                         // ui event id
        uint16(win.PercentOfWidth(98)),     // Width
        uint16(win.PercentOfHeight(82)),    // Height
        win.PercentOfWidth(1),              // X
        win.PercentOfHeight(1),             // Y
        1,                                  // states for labels can be set to have a background and scrollbar.
    )

    // Character Switch Button
    go win.Button("Chr",
        CharSwitchButtonID,
        uint16(win.PercentOfWidth(9)),
        uint16(win.PercentOfHeight(14)),
        win.PercentOfWidth(1),
        win.PercentOfHeight(85),
    )

    // Input
    input = win.Textbox("",
        InputID,
        uint16(win.PercentOfWidth(69)),
        uint16(win.PercentOfHeight(14)),
        win.PercentOfWidth(11),
        win.PercentOfHeight(85),
    )

    // Attack Button
    go win.Button("Atk",
        AttackButtonID,
        uint16(win.PercentOfWidth(9)),
        uint16(win.PercentOfHeight(14)),
        win.PercentOfWidth(81),
        win.PercentOfHeight(85),
    )

    // Bag Button
    go win.Button("Bag",
        BagSwitchButtonID,
        uint16(win.PercentOfWidth(9)),
        uint16(win.PercentOfHeight(14)),
        win.PercentOfWidth(90),
        win.PercentOfHeight(85),
    )
}

func WinLoop() {
	defer leave() // Leave the IRC server if we're in one
    // UI Events
    go func() {
        for {
            ev := win.WaitForUIEvent()
            switch ev.(type) {
                case TermUI.UITextboxSubmitEvent:
                    SendCommand(input.Name)
                    input.Name = ""
                    win.DrawUITextbox(InputID)
                case TermUI.UIPressEvent:
                    //
            }
        }
    }()
    // Checking for X events
    for {
        ev, xerr := win.Conn.WaitForEvent()
        if xerr != nil {fmt.Printf("Error: %s\n", xerr)}

        // (in some WMs this happens when you close the program)
        if ev == nil && xerr == nil {
        	leave() // redundant but the above defer doesn't get called when this happens
        	return
        }

        win.DefaultListeners(ev);
    }
}

func sendToTextarea(str string) {
	if(textarea == nil) {return}
	textarea.Name += str+"\n"
	win.DrawUILabel(TextareaID)
	win.DrawUILabel(TextareaID)
}