package main

import (
    "fmt"
	"log"
    "embed"
    "os/exec"
    "strings"
    "regexp"

    "github.com/gotk3/gotk3/gdk"
    "github.com/gotk3/gotk3/gtk"
)

// We define every GTK related variable here because it gives us a pretty decent boost in boot up time.

// Window variables
var win 			*gtk.Window
var vbox 			*gtk.Box
var bottom 			*gtk.Box
var rest            *gtk.Box
var css             *gtk.CssProvider

// Buttons
var charSwitch		*gtk.Button
var charSwitchImage *gtk.Image
var attack			*gtk.Button
var attackImage     *gtk.Image
var inventory		*gtk.Button
var inventoryImage  *gtk.Image

// Input areas and their needed variables.
var	input			*gtk.Entry
var	textarea		*gtk.Label
var textareaBuffer  []string


// Other global variables

var inIRC           bool        = false

//go:embed gfx/*
var fs              embed.FS

var re              *regexp.Regexp

var err             error

func main() {
    winInit()
    winBuild()

    win.Add(vbox)
    win.SetDefaultSize(320, 200)
    win.ShowAll()
    re = regexp.MustCompile(`(␛|\033|\x1B|)(\[)([0-9]{1,2}|m)([A-z]|[0-9]{1,2})(;([0-9]{1,3}(m|))){0,2}`)
    winContext, _ := win.GetStyleContext()
    winContext.AddProvider(css, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
    gtk.Main()
}

func winInit() {
    // Initialize GTK without parsing any command line arguments.
    gtk.Init(nil)

    // Create a new toplevel window, set its title, and connect it to the
    // "destroy" signal to exit the GTK main loop when it is destroyed.
    win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    if err != nil {
        log.Fatal("Unable to create window:", err)
    }
    win.SetTitle("Terminal")

    win.Connect("destroy", func() {
        gtk.MainQuit()
    })
}

func winBuild() {
    // from here on we start to ignore GTK errors because the even the devs of the library
    // don't seem to do much with them in examples, probably because gtk just prints the error and
    // go usually can't handle it and it's a waste of time and memory to expect err to ever return
    // anything useful.

    // Firstly, load in the CSS file.
    css, _ = gtk.CssProviderNew()
    cssData, err := fs.ReadFile("gfx/style.css")
    if(err != nil) {panic(err)}
    css.LoadFromData(string(cssData))

    // Main VBox
    vbox, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

    // Bottommost section of the window
    bottom, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

    // Buttons

    // The errors are never useful anyways
    charSwitch, _ = gtk.ButtonNew()
    charSwitchFile, err := fs.ReadFile("gfx/charswitch.png")
    if(err != nil) {panic(err)}
    charSwitchPixbuf, _ := gdk.PixbufNewFromBytesOnly(charSwitchFile)
    charSwitchImage, _ = gtk.ImageNewFromPixbuf(charSwitchPixbuf)
    charSwitch.SetImage(charSwitchImage)

    attack, _ = gtk.ButtonNew()
    attackFile, err := fs.ReadFile("gfx/attack.png")
    if(err != nil) {panic(err)}
    attackPixbuf, _ := gdk.PixbufNewFromBytesOnly(attackFile)
    attackImage, _ = gtk.ImageNewFromPixbuf(attackPixbuf)
    attack.SetImage(attackImage)

    inventory, _ = gtk.ButtonNew()
    inventoryFile, err := fs.ReadFile("gfx/bag.png")
    if(err != nil) {panic(err)}
    inventoryPixbuf, _ := gdk.PixbufNewFromBytesOnly(inventoryFile)
    inventoryImage, _ = gtk.ImageNewFromPixbuf(inventoryPixbuf)
    inventory.SetImage(inventoryImage)

    // Input Textbox
    input, _ = gtk.EntryNew()
    input.Connect("activate", sendCommand)

    // The rest of the window
    rest, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    
    // Output Textarea
    
    textarea, _ = gtk.LabelNew("")
    textarea.SetXAlign(-1)
    textarea.SetLineWrap(true)
    textarea.SetMaxWidthChars(1)

    // Add all this shit to the window.
    
    bottom.PackStart(charSwitch, false, false, 0)
    bottom.PackStart(input, true, true, 0)
    bottom.PackStart(attack, false, false, 0)
    bottom.PackEnd(inventory, false, false, 0)
    
    rest.PackStart(textarea, false, false, 0)

    vbox.PackStart(rest, true, true, 0)
    vbox.PackEnd(bottom, false, false, 0)
}

func userSwitchDropdown() {
    fmt.Println("Would've opened user dropdown (todo)")
}

func bagOpen() {
    fmt.Println("Would've opened bag (todo)")
}

func attackDropdown() {
    fmt.Println("Would've opened attack dropdown (todo)")
}

func sendCommand() {
    text, _ := input.GetText()
    // are we in IRC mode?
    if(inIRC) {
        // send a command to the IRC server we're connected to (todo)
    } else {
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
        output, err := cmd.Output()
        var outputs []string
        if(err != nil) {
            outputs = strings.Split(err.Error(), "\n")
        } else {
            outputs = strings.Split(string(output), "\n")
        }
        for _, v := range outputs {
            textareaBuffer = append(textareaBuffer, v)
        }
        input.SetText("")
        updateTextarea()
    }
}

func updateTextarea() {
    var text string
    areaLength := len(textareaBuffer)
    width, height := win.GetSize()
    textareaBuffer_ := textareaBuffer
    if(areaLength > height/20) {
        textareaBuffer_ = textareaBuffer[areaLength-(height/20):areaLength]
    }
    for _, v := range textareaBuffer_ {
        toAppend := v
        // Strip any bash control characters out
        // We do this four times since some things come up something is removed.
        toAppend = re.ReplaceAllString(toAppend, "")
        if(len(toAppend) > width/6) {
            toAppend = v[0:width/6]
        }
        fmt.Println(len(toAppend))
        fmt.Println(width/6)
        text += toAppend+"\n"
    }
    for _, v := range text {
        fmt.Print(string(v)+"|")
    }
    textarea.SetText(text)
}