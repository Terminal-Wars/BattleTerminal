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

    // Firstly, load in the CSS file.
    css, err = gtk.CssProviderNew()
    if(err != nil) {panic(err)}
    cssData, err := fs.ReadFile("gfx/style.css")
    if(err != nil) {panic(err)}
    css.LoadFromData(string(cssData))

    // Main VBox
    vbox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if(err != nil) {panic(err)}

    // Bottommost section of the window
    bottom, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if(err != nil) {panic(err)}

    // Buttons

    // The errors are never useful anyways
    charSwitch, err = gtk.ButtonNew()
    if(err != nil) {panic(err)}
    charSwitchFile, err := fs.ReadFile("gfx/charswitch.png")
    if(err != nil) {panic(err)}
    charSwitchPixbuf, err := gdk.PixbufNewFromBytesOnly(charSwitchFile)
    if(err != nil) {panic(err)}
    charSwitchImage, err = gtk.ImageNewFromPixbuf(charSwitchPixbuf)
    if(err != nil) {panic(err)}
    charSwitch.SetImage(charSwitchImage)

    attack, err = gtk.ButtonNew()
    if(err != nil) {panic(err)}
    attackFile, err := fs.ReadFile("gfx/attack.png")
    if(err != nil) {panic(err)}
    attackPixbuf, err := gdk.PixbufNewFromBytesOnly(attackFile)
    if(err != nil) {panic(err)}
    attackImage, err = gtk.ImageNewFromPixbuf(attackPixbuf)
    if(err != nil) {panic(err)}
    attack.SetImage(attackImage)

    inventory, err = gtk.ButtonNew()
    if(err != nil) {panic(err)}
    inventoryFile, err := fs.ReadFile("gfx/bag.png")
    if(err != nil) {panic(err)}
    inventoryPixbuf, err := gdk.PixbufNewFromBytesOnly(inventoryFile)
    if(err != nil) {panic(err)}
    inventoryImage, err = gtk.ImageNewFromPixbuf(inventoryPixbuf)
    if(err != nil) {panic(err)}
    inventory.SetImage(inventoryImage)

    // Input Textbox
    input, err = gtk.EntryNew()
    if(err != nil) {panic(err)}
    input.Connect("activate", sendCommand)


    // The rest of the window
    rest, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if(err != nil) {panic(err)}
    
    // Output Textarea
    textareaWrapper, err = gtk.ScrolledWindowNew(nil, nil)
    if(err != nil) {panic(err)}
    textareaWrapper.SetMaxContentWidth(4)
    textareaWrapper.SetMaxContentHeight(4)
    textareaWrapper.SetPlacement(gtk.CORNER_BOTTOM_RIGHT)
    textareaWrapper.SetOverlayScrolling(true)
    textareaAdjustment = textareaWrapper.GetVAdjustment()

    textarea, err = gtk.LabelNew("")
    if(err != nil) {panic(err)}
    textarea.SetXAlign(-1)
    textarea.SetLineWrap(true)
    textarea.SetMaxWidthChars(1)

    textareastyle, err := textarea.GetStyleContext()
    if(err != nil) {panic(err)}
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
    defer func(){
        if r := recover(); r != nil {
            fmt.Println("\n\n\n\n\n\nRecovered in GTK", r, "\n\n\n\n\n\n")
        }
    }()
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
