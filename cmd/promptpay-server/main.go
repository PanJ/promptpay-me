package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/skip2/go-qrcode"

	promptpay "github.com/PanJ/promptpay-me"
)

type promptpayHandler struct {
	recipientID string
}

func (h promptpayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var amount float64
	var err error

	if r.URL.Path != "/" {
		// get amount from path
		p := strings.Trim(r.URL.Path, "/")
		amount, err = strconv.ParseFloat(p, 32)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusInternalServerError)
			return
		}
	}

	payload := promptpay.Generate(h.recipientID, float32(amount))

	png, err := qrcode.Encode(payload, qrcode.High, 512)
	if err != nil {
		http.Error(w, "Unable to generate QR Code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func main() {
	recipientID := os.Getenv("PROMPTPAY_RECIPIENT_ID")
	if len(recipientID) == 0 {
		log.Fatal("Recipient ID required")
	}

	http.Handle("/", &promptpayHandler{recipientID})
	log.Println("start server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("server error;", err)
	}
}
