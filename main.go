package main

import (
	"strconv"
	"github.com/skip2/go-qrcode"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"os"
)

var id string

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	png, err := qrcode.Encode(GeneratePayload(id, 0), qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Unable to generate QR Code", 500)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func WithAmount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	amount, err := strconv.ParseFloat(ps.ByName("amount"), 32)
	if err != nil {
		http.Error(w, "Invalid input", 500)
		return
	}
	png, err := qrcode.Encode(GeneratePayload(id, float32(amount)), qrcode.High, 512)
	if err != nil {
		http.Error(w, "Unable to generate QR Code", 500)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func main() {
	id = os.Getenv("PROMPTPAY_RECIPIENT_ID")
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/:amount", WithAmount)
	log.Fatal(http.ListenAndServe(":8080", router))
	log.Println("Server started")
}