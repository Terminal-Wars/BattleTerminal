package main

import (
	"fmt"
	"log"

    "github.com/gotk3/gotk3/gtk"
)

// We define every variable here because it gives us a pretty decent boost in boot up time.

var win 			*gtk.Window
var vbox 			*gtk.Box
var bottom 			*gtk.Box
var charswitch		*gtk.Button
var attack			*gtk.Button
var inventory		*gtk.Button
var	input			*gtk.Entry
var rest			*gtk.Box
var	textarea		*gtk.Label

var	err				error

func main() {
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

    // Main VBox
    vbox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if(err != nil) {fmt.Println(err)}

    // Bottommost section of the window
    bottom, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if(err != nil) {fmt.Println(err)}

    // Buttons

    charswitch, err = gtk.ButtonNewWithLabel("C")
    if(err != nil) {fmt.Println(err)}

    attack, err = gtk.ButtonNewWithLabel("A")
    if(err != nil) {fmt.Println(err)}

    inventory, err = gtk.ButtonNewWithLabel("Is")
    if(err != nil) {fmt.Println(err)}

    // Input Textbox
    input, err = gtk.EntryNew()
    if(err != nil) {fmt.Println(err)}
    // The rest of the window
    
    rest, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if(err != nil) {fmt.Println(err)}
    
    // Output Textarea
    
    textarea, err = gtk.LabelNew("testtesttesttesttesttesttesttesttesttest\ntesttesttesttesttesttesttesttesttesttest\ntesttesttesttesttesttesttesttesttesttest\ntesttesttesttesttesttesttesttesttesttest\n")
    if(err != nil) {fmt.Println(err)}
    textarea.SetXAlign(-1)
    textarea.SetLineWrap(true)
    textarea.SetMaxWidthChars(1)

    // Add all this shit to the window.
    
	bottom.PackStart(charswitch, false, false, 0)
	bottom.PackStart(input, true, true, 0)
	bottom.PackStart(attack, false, false, 0)
	bottom.PackEnd(inventory, false, false, 0)
	
	rest.PackStart(textarea, false, false, 0)

	vbox.PackStart(rest, true, true, 0)
	vbox.PackEnd(bottom, false, false, 0)
    
    win.Add(vbox)

    // Set the default window size.
    win.SetDefaultSize(320, 200)

    // Recursively show all widgets contained in this window.
    win.ShowAll()

    // Begin executing the GTK main loop.  This blocks until
    // gtk.MainQuit() is run.
    gtk.Main()
}