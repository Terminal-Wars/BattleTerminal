package src


import (
	"log"
    "embed"
    "regexp"
    "fmt"
    "time"

    "github.com/gotk3/gotk3/gdk"
    "github.com/gotk3/gotk3/gtk"
)

var Start time.Time
var end time.Time

// We define every GTK related variable here because it gives us a pretty decent boost in boot up time.

// Window variables
var win 			        *gtk.Window
var vbox 			        *gtk.Box
var bottom 			        *gtk.Box
var rest                    *gtk.Box
var css                     *gtk.CssProvider

// Buttons
var charSwitch		        *gtk.Button
var charSwitchImage         *gtk.Image
var attack			        *gtk.Button
var attackImage             *gtk.Image
var inventory		        *gtk.Button
var inventoryImage          *gtk.Image

// Input areas and their needed variables.
var	input			        *gtk.Entry
var	textarea		        *gtk.Label
var textareaWrapper         *gtk.ScrolledWindow
var textareaBuffer          []string
var textareaAdjustment      *gtk.Adjustment

// Other global variables

var offset          int16          = 0
var re              *regexp.Regexp

//go:embed gfx/*
var fs              embed.FS

var err             error


func WinInit() {
    // Initialize GTK without parsing any command line arguments.
    gtk.Init(nil)

    // Create a new toplevel window, set its title, and connect it to the
    // "destroy" signal to exit the GTK main loop when it is destroyed.
    win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    if err != nil {
        log.Fatal("Unable to create window:", err)
    }
    win.SetTitle("Terminal")
    win.AddEvents(int(gdk.EVENT_CONFIGURE) | 2097152)
    win.Connect("destroy", func() {
        gtk.MainQuit()
    })
}

func WinBuild() {
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
    textareaWrapper, _ = gtk.ScrolledWindowNew(nil, nil)
    textareaWrapper.SetMaxContentWidth(4)
    textareaWrapper.SetMaxContentHeight(4)
    textareaWrapper.SetPlacement(gtk.CORNER_BOTTOM_RIGHT)
    textareaWrapper.SetOverlayScrolling(true)
    textareaAdjustment = textareaWrapper.GetVAdjustment()

    textarea, _ = gtk.LabelNew("")
    textarea.SetXAlign(-1)
    textarea.SetLineWrap(true)
    textarea.SetMaxWidthChars(1)

    textareastyle, _ := textarea.GetStyleContext()
    textareastyle.AddClass("textarea")
    textareastyle.AddProvider(css, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

    textareaWrapper.Add(textarea)

    // Add all this shit to the window.
    
    bottom.PackStart(charSwitch, false, false, 0)
    bottom.PackStart(input, true, true, 0)
    bottom.PackStart(attack, false, false, 0)
    bottom.PackEnd(inventory, false, false, 0)
    
    rest.PackEnd(textareaWrapper, true, true, 0)

    vbox.PackStart(rest, true, true, 0)
    vbox.PackEnd(bottom, false, false, 0)
}

func WinLoop() {
	win.Add(vbox)
    win.SetDefaultSize(320, 200)
    win.ShowAll()
    re = regexp.MustCompile(`(␛|\033|\x1B|)(\[)([0-9]{1,2}|m)([A-z]|[0-9]{1,2})(;([0-9]{1,3}(m|))){0,2}`)
    winContext, err := win.GetStyleContext()
    if(err == nil) {winContext.AddProvider(css, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)}
    end = time.Now()
    fmt.Println(end.Sub(Start))
    bitch()
}

// function to try and catch when the gtk.Main process crashes and doesn't let it.
func bitch() {
    defer bitch()
    gtk.Main()
}

func sendToTextarea(str string) {
	textareaBuffer = append(textareaBuffer, str)
	updateTextarea()
}

func updateTextarea() {
    _, height := win.GetSize()
    var text string
    areaLength := len(textareaBuffer)
    areaMax := (height/20)-1
    for _, v := range textareaBuffer {
        // Strip any bash control characters out
        toAppend := re.ReplaceAllString(v, "")
        text += toAppend+"\n"
    }
    // If we're within the limit for lines on screen...
    if(areaLength < areaMax) {
        // Forcefully push our text up based on how much is left
        // (because GTK is a little silly)
        for i := 0; i < areaMax-areaLength; i++ {
            text += "\n"
        }
    }
    textarea.SetText(text)
    textareaAdjustment.SetValue(textareaAdjustment.GetUpper()+10)
}
