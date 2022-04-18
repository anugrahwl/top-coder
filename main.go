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

func GetInput(signalTerminate <-chan string, inputCh chan<- string) {
	inputReader := bufio.NewReader(os.Stdin)
	inputs := []string{}

	for {
		select {
		case <-signalTerminate:
			inputCh <- inputs[len(inputs)-1]
		default:
			str, _ := inputReader.ReadString('\n')
			str = strings.Trim(str, "\n")
			inputs = append(inputs, str)
		}
	}
}

func main() {
	inputController := make(chan string)
	inputCh := make(chan string)

	go GetInput(inputController, inputCh)
	animation.AnimateSentenceForward(START_SCREEN, inputController)
	<-inputCh
	animation.AnimateSentenceBackward(START_SCREEN)

	done := false
	for !done {
		go GetInput(inputController, inputCh)
		animation.AnimateSentenceForward("From [xx] : ", inputController)
		source := <-inputCh
		animation.AnimateSentenceBackward("From [xx] : " + source)

		go GetInput(inputController, inputCh)
		animation.AnimateSentenceForward("To [xx] : ", inputController)
		target := <-inputCh
		animation.AnimateSentenceBackward("To [xx] : " + target)

		go GetInput(inputController, inputCh)
		animation.AnimateSentenceForward("Sentence : ", inputController)
		sentence := <-inputCh
		animation.AnimateSentenceBackward("Sentence : " + sentence)

		trans, err := Translate(source, target, sentence)
		if err != nil {
			trans = "SOMETHING WENT WRONG. CHECK YOUR INTERNET CONNECTION."
		}
		prompt := "\nWanna translate again [y/n]? "

		go GetInput(inputController, inputCh)
		animation.AnimateSentenceForward(trans+prompt, inputController)
		input := <-inputCh
		animation.AnimateSentenceBackward(trans + prompt + input)

		promptAgain := true
		for promptAgain {
			if input == "y" || input == "Y" {
				promptAgain = false
			} else if input == "n" || input == "N" {
				animation.AnimateSentenceForwardWithNoInput("Bye Bye")
				time.Sleep(time.Millisecond * 500)
				animation.AnimateSentenceBackward("Bye Bye")

				done = true
				promptAgain = false
			} else {
				go GetInput(inputController, inputCh)
				msg := "What are you trying to say bruh? "
				animation.AnimateSentenceForward(msg+prompt, inputController)
				input = <-inputCh
				animation.AnimateSentenceBackward(msg + prompt + input)
			}
		}
	}
}
