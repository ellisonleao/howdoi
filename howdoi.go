package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	flag "github.com/ogier/pflag"
)

const (
	version       = "0.1.0"
	requestURL    = "https://www.google.com/search?q=site:%s%%20%s"
	answerFooter  = "%s \n----\n Answer from %s"
	answerHeader  = "--- Answer %d --\n %s"
	noAnswer      = "< no answer given >"
	noAnswerFound = "Sorry, couldn't find any help with that topic"
	help          = `
	usage: howdoi [-h|--help] [-p|--pos POS] [-a|--all] [-l|--link] [-n|--num-answers NUM_ANSWERS] [-v|--version]
              [QUERY [QUERY ...]]

instant coding answers via the command line

positional arguments:
  QUERY                 the question to answer

optional arguments:
  -h, --help            show this help message and exit
  -p POS, --pos POS     select answer in specified position (default: 1)
  -a, --all             display the full text of the answer
  -l, --link            display only the answer link
  -n NUM_ANSWERS, --num-answers NUM_ANSWERS
                        number of answers to return
  -v, --version         displays the current version of howdoi`
)

// Howdoi is the main struct for the howdoi command with its flags
type Howdoi struct {

	// Position argument grabs a specific answer position in a list of answers.
	// Negative values and numbers greater than the number of answers will be handled as well
	Position int

	// ShowAllText displays the full answer instead of just the code part
	ShowAllText bool

	// ShowLinkOnly displays only the answer link
	ShowLinkOnly bool

	// NumAnswers will show a number of answers between 1 and total answers
	NumAnswers int

	// Question records the current input question
	Question string

	// ShowHelp will output the help
	ShowHelp bool

	// ShowVersion will output the current version
	ShowVersion bool
}

// Init will create a new Howdoi object and set the flags to it
func Init() *Howdoi {
	h := &Howdoi{}

	flag.IntVar(&h.Position, "pos", 1, "select answer in specified position")
	flag.BoolVar(&h.ShowAllText, "all", false, "display the full text of the answer")
	flag.BoolVar(&h.ShowLinkOnly, "link", false, "display only the answer link")
	flag.IntVar(&h.NumAnswers, "num-answers", 1, "number of answers to return")
	flag.BoolVar(&h.ShowHelp, "help", false, "show this help message and exit")
	flag.BoolVar(&h.ShowVersion, "version", false, "show current version")

	return h
}

// Execute is the main function for the howdoi command
func (h *Howdoi) Execute() {
	flag.Parse()

	if h.ShowHelp {
		fmt.Println(help)
		os.Exit(0)
	}

	if h.ShowVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// smal check on position value
	if h.Position <= 0 {
		h.Position = 1
	}

	err := h.sanitizeQuestion(flag.Args())
	if err != nil {
		fmt.Println(help)
		os.Exit(1)
	}

	links, err := h.getLinks()
	if err != nil {
		log.Fatal(err)
	}

	answer, err := h.getAnswer(links)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(answer)
}

// sanitizeQuestion parses the args input and set the urlencoded question on the howdoi object
func (h *Howdoi) sanitizeQuestion(args []string) error {
	h.Question = strings.Join(args, " ")
	h.Question = strings.TrimSpace(h.Question)

	if len(h.Question) == 0 {
		return errors.New("Input not valid")
	}

	h.Question = strings.Replace(h.Question, "?", "", -1)
	h.Question = url.QueryEscape(h.Question)
	return nil
}

// getLinks will grab the link for the answer pages
func (h *Howdoi) getLinks() ([]string, error) {
	req := fmt.Sprintf(requestURL, "stackoverflow.com", h.Question)
	doc, err := goquery.NewDocument(req)
	if err != nil {
		return nil, err
	}

	result := doc.Find(".l")
	if len(result.Nodes) == 0 {
		result = doc.Find(".r a")
		if len(result.Nodes) == 0 {
			fmt.Println(noAnswerFound)
			os.Exit(0)
		}
	}
	links := []string{}

	result.Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		parsed, err := url.Parse(link)
		query := parsed.Query()
		link = query["q"][0]

		if err != nil {
			//TODO!
			fmt.Println("ERROR on Link", link)
		} else if strings.Contains(link, "question") {
			// adding only the questions links
			links = append(links, link)
		}
	})
	return links, nil
}

func (h *Howdoi) getAnswer(links []string) (string, error) {
	var output string

	// validating position
	if h.Position > len(links) {
		h.Position = 1
	}

	// do not show answer header if there is only one answer to return
	if h.NumAnswers == 1 {
		if h.ShowLinkOnly {
			return links[h.Position-1], nil
		}
		return getAnswerText(links[h.Position-1]), nil
	}

	links = links[0:h.NumAnswers]
	for i, link := range links {
		answer := getAnswerText(link)
		output += fmt.Sprintf(answerHeader, i+1, answer)
	}

	return output, nil
}

func getAnswerText(link string) string {
	link = fmt.Sprintf("%s?answertab=votes", link)
	req, err := goquery.NewDocument(link)
	if err != nil {
		log.Fatal(errors.New("Could not get answer. Pleasy try again later"))
	}

	answerDiv := req.Find(".answer")
	if len(answerDiv.Nodes) == 0 {
		return noAnswer
	}

	answerDiv = answerDiv.First()

	if !h.ShowAllText {
		// grabbing <code> or <pre> content
		instructions := answerDiv.Find("pre")

		if len(instructions.Nodes) == 0 {
			instructions = answerDiv.Find("code")
		}
		return instructions.First().Text()
	}
	return fmt.Sprintf(answerFooter, answerDiv.Find(".post-text > *").Text(), link)
}
