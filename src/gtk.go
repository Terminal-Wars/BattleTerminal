package src


import (
	"log"
    "embed"
    "regexp"
    "image/color"
    "os"

    "gioui.org/app"
    "gioui.org/font/gofont"
    "gioui.org/io/system"
    "gioui.org/layout"
    "gioui.org/unit"
    "gioui.org/op"
    "gioui.org/widget"
    "github.com/Terminal-Wars/ClassicDesign"
)

var (
    // gioui widgets
    charSwitchButton  = new(widget.Clickable)
    attackButton      = new(widget.Clickable)
    inventoryButton   = new(widget.Clickable)
    inputBox          = widget.Editor {
        SingleLine: true,
        Submit: true,
    }
    textarea          = widget.Label {
        MaxLines: 12,
    }
    list              = &widget.List{
        List: layout.List{
            Axis: layout.Vertical,
        },
    }
    textareaBuffer      []string

    // textarea stuff

    offset              int16           = 0
    re                  *regexp.Regexp  = regexp.MustCompile(`(␛|\033|\x1B|)(\[)([0-9]{1,2}|m)([A-z]|[0-9]{1,2})(;([0-9]{1,3}(m|))){0,2}`)
    
    // embedded files
    //go:embed gfx/*
    fs                  embed.FS

    // misc.
    width               float32           = 320
    height              float32           = 240
    err                 error
)

type (
    D = layout.Dimensions
    C = layout.Context
)

func WinInit() {
    go func() {
        w := app.NewWindow(
            app.Title("BattleTerm"),
            app.Size(unit.Dp(width), unit.Dp(height)),
        )
        err := run(w)
        if(err != nil) {
            log.Fatal(err)
        }
        os.Exit(0)
    }()
    app.Main()
}

func run(w *app.Window) error {
    // operations from the UI
    var ops op.Ops

    // todo: new design style
    th := material.NewTheme(gofont.Collection())

    // listen for events in the window.
    for e := range w.Events() {

      // detect what type of event
      switch e := e.(type) {

      // this is sent when the application should re-render.
      case system.FrameEvent:
        gtx := layout.NewContext(&ops, e)
        inputEvent(gtx)
        fill(gtx, th)
        e.Frame(gtx.Ops)
      }
    }
    /*
    charSwitchFile, err := fs.ReadFile("gfx/charswitch.png")
    if(err != nil) {panic(err)}

    attackFile, err := fs.ReadFile("gfx/attack.png")
    if(err != nil) {panic(err)}

    inventoryFile, err := fs.ReadFile("gfx/bag.png")
    if(err != nil) {panic(err)}
    */
    return nil
}

func inputEvent(gtx layout.Context) {
    for _, e := range inputBox.Events() {
        if _, ok := e.(widget.SubmitEvent); ok {
            sendCommand()
        }
    }
}

func fill(gtx layout.Context, th *material.Theme) layout.Dimensions {
    in := layout.UniformInset(unit.Dp(4))
    //sp := layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4), Top: unit.Dp(4), Bottom: unit.Dp(0),}
    // layout our widgets
    widgets := []layout.Widget{
        // Textarea that shows the output.
        func(gtx C) D {
            return layout.Flex{Alignment: layout.Start, Spacing: layout.SpaceEvenly, Axis: layout.Vertical}.Layout(gtx,
                layout.Rigid(func(gtx C) D {
                    size(gtx,320,150)
                    border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(0), Width: unit.Px(1)}
                    return border.Layout(gtx, func(gtx C) D {
                        return in.Layout(gtx, material.Label(th, unit.Dp(12), getTextarea()).Layout)
                    })
                }),
            )
        },
        // Bottom menu
        func(gtx C) D {
            return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceEvenly, Axis: layout.Horizontal}.Layout(gtx,
                // character switcher
                layout.Flexed(10,func(gtx C) D  {
                    return material.Button(th, charSwitchButton, "C").Layout(gtx)
                }),
                // input textbox
                layout.Flexed(60,func(gtx C) D  {
                    border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(0), Width: unit.Px(1)}
                    return in.Layout(gtx, func(gtx C) D {
                        return border.Layout(gtx, func(gtx C) D {
                            return in.Layout(gtx, material.Editor(th, &inputBox, "  ").Layout)
                        })
                    })
                }),
                // attack button
                layout.Flexed(10,func(gtx C) D  {
                    return material.Button(th, attackButton, "A").Layout(gtx)
                }),
                // inventory button
                layout.Flexed(10,func(gtx C) D  {
                    return material.Button(th, inventoryButton, "I").Layout(gtx)
                }),
            )
        },
    }
    return material.List(th, list).Layout(gtx, len(widgets),
        func(gtx C, i int) D {
            return layout.UniformInset(unit.Dp(4)).Layout(gtx, widgets[i])
        },
    )
}

func size(gtx layout.Context, width float32, height float32) {
    gtx.Constraints.Min.X = gtx.Px(unit.Dp(0))
    gtx.Constraints.Max.X = gtx.Px(unit.Dp(width))
    gtx.Constraints.Min.Y = gtx.Px(unit.Dp(0))
    gtx.Constraints.Max.Y = gtx.Px(unit.Dp(height))
}

func sendToTextarea(str string) {
	textareaBuffer = append(textareaBuffer, str)
}

func getTextarea() (string) {
    var text string
    // create a clone of the textarea buffer since we may need to change what it points to in a bit
    textareaBuffer_ := textareaBuffer
    areaLength := len(textareaBuffer_)
    areaMax := (int(height)/20)-1
    // If we're outside the limits of the screen...
    if(areaLength > areaMax) {
        // change what we're looking at.
        textareaBuffer_ = textareaBuffer[areaLength-areaMax:areaLength]
    }
    for _, v := range textareaBuffer_ {
        // Strip any bash control characters out
        toAppend := re.ReplaceAllString(v, "")
        text += toAppend+"\n"
    }
    // If we're within the limit for lines on screen...
    if(areaLength < areaMax) {
        // Forcefully push our textarea towards the edge of the screen
        // (this used to be because GTK is silly and sadly gioui is silly too)
        for i := 0; i < areaMax-areaLength; i++ {
            text += "\n"
        }
    }
    return text
}
