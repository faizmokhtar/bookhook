package hooks

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type CRC struct {
	ResponseToken string `json:"response_token"`
}

func TwitterWebhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		log.Println("info: challenge-response check")
		token := r.URL.Query()["crc_token"]
		if len(token) < 1 {
			log.Println("error: no crc_token given")
			return
		}

		log.Println("info: crc token ", token)

		h := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
		h.Write([]byte(token[0]))
		encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))

		response := CRC{
			ResponseToken: "sha256=" + encoded,
		}

		responseJSON, _ := json.Marshal(response)
		log.Println("response: ", response)
		fmt.Fprintf(w, string(responseJSON))
	case "POST":
		log.Println("listening to twitter account activity")
		body, _ := ioutil.ReadAll(r.Body)

		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, body, "", "\t")
		if error != nil {
			log.Println("JSON parse error: ", error)
			return
		}

		log.Println(string(prettyJSON.Bytes()))
	default:
		fmt.Fprintln(w, "go away!")
	}
}
