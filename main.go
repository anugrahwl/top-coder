package main

import (
	"errors"
	"os"
	"time"

	"encoding/json"

	"net/http"
	"net/url"

	"bufio"
	"io/ioutil"

	"strconv"
	"strings"

	"github.com/anugrahwl/top-coder/animation"
)

const URL_API = "https://translate.google.com/translate_a/single?client=at&dt=t&dt=ld&dt=qca&dt=rm&dt=bd&dj=1&ie=UTF-8&oe=UTF-8&inputm=2&otf=2&iid=1dd3b944-fa62-4b55-b330-74909a99969e"

const (
	START_SCREEN = "Welcome to our language translator\nPress any key to continue..."
)

type Output struct {
	Sentences []Sentence `json:"sentences"`
}

type Sentence struct {
	Trans string `json:"trans"`
}

func Translate(origin string, target string, str string) (string, error) {
	// membuat objek value
	data := url.Values{}

	data.Set("sl", origin)
	data.Set("tl", target)
	data.Set("q", str)

	// membuat objek request
	request, err := http.NewRequest("POST", URL_API, strings.NewReader(data.Encode()))
	if err != nil {
		return "", errors.New("gagal membuat request object")
	}

	// mengeset header pada object request
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	// membuat object client dan menggunakannya untuk melakukan
	// request berdasarkan objek request yang dibuat sebelumnya
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return "", errors.New("gagal melakukan request ke API")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("gagal membaca response dari API")
	}

	result := Output{}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", errors.New("gagal mengkonversi response body")
	}

	return result.Sentences[0].Trans, nil
}

func GetInput() string {
	inputReader := bufio.NewReader(os.Stdin)

	str, _ := inputReader.ReadString('\n')
	str = strings.Trim(str, "\n")

	return str
}

func main() {
	animation.AnimateSentenceForward(START_SCREEN)
	GetInput()
	animation.AnimateSentenceBackward(START_SCREEN)

	done := false
	for !done {
		animation.AnimateSentenceForward("Source : ")
		source := GetInput()
		animation.AnimateSentenceBackward("Source : " + source)

		animation.AnimateSentenceForward("Target : ")
		target := GetInput()
		animation.AnimateSentenceBackward("Target : " + target)

		animation.AnimateSentenceForward("Sentence : ")
		sentence := GetInput()
		animation.AnimateSentenceBackward("Sentence : " + sentence)

		trans, _ := Translate(source, target, sentence)
		prompt := "\nWanna translate again [y/n]? "

		animation.AnimateSentenceForward(trans + prompt)
		input := GetInput()
		animation.AnimateSentenceBackward(trans + prompt + input)

		promptAgain := true
		for promptAgain {
			if input == "y" || input == "Y" {
				promptAgain = false
			} else if input == "n" || input == "N" {
				animation.AnimateSentenceForward("Bye Bye")
				time.Sleep(time.Millisecond * 500)
				animation.AnimateSentenceBackward("Bye Bye")
				done = true
				promptAgain = false
			} else {
				msg := "What are you trying to say bruh? "
				animation.AnimateSentenceForward(msg + prompt)
				input = GetInput()
				animation.AnimateSentenceBackward(msg + prompt + input)
			}
		}
	}
}
