package main

import (
	"log"
	"syscall/js"
	"time"

	"github.com/enescakir/emoji"
)

func main() {
	js.Global().Set(
		"emojifyMyText",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			doc := js.Global().Get("document")
			textarea := doc.Call("getElementById", "my_text_area")
			if textarea.IsUndefined() {
				log.Println("textarea is undefined")
				return nil
			}
			value := textarea.Get("value")
			if value.IsUndefined() {
				log.Println("innerText is undefined")
				return nil
			}
			valueStr := value.String()
			log.Printf("Got text: %s", valueStr)
			emojiStr := emoji.Sprint(valueStr)
			textarea.Set("value", emojiStr)
			return nil
		}),
	)
	log.Printf("WASM is running")
	for {
		time.Sleep(time.Second) // we need to keep running
	}
}
