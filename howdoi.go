package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	requestURL   = "https://www.google.com/search?q=site:%s%%20%s"
	answerHeader = "%s \n----\n Answer from %s"
)

// Howdoi is the main struct for the howdoi command with its flags
type Howdoi struct {
	Position    int
	ShowAllText bool
	ColorOutput bool
	NumAnswers  int
	Question    string
}

// Init will create a new Howdoi object and set the flags to it
func Init() *Howdoi {
	h := &Howdoi{}

	flag.IntVar(&h.Position, "p", 1, "select answer in specified position")
	flag.BoolVar(&h.ShowAllText, "a", false, "display the full text of the answer")
	//flag.IntVar(&h.NumAnswers, "n", 1, "number of answers to return")

	return h
}

// Execute is the main function for the howdoi command
func (h *Howdoi) Execute() {
	flag.Parse()

	// smal check on position value
	if h.Position <= 0 {
		h.Position = 1
	}

	err := h.sanitizeQuestion(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	links, err := h.getLinks()
	if err != nil {
		log.Fatal(err)
	}

	answer, err := h.getAnswer(links[h.Position-1])
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
			log.Fatal("No answers found")
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

func (h *Howdoi) getAnswer(link string) (string, error) {
	link = fmt.Sprintf("%s?answertab=votes", link)
	req, err := goquery.NewDocument(link)
	if err != nil {
		log.Fatal(errors.New("Could not get answer. Pleasy try again later"))
	}

	answerDiv := req.Find(".answer")
	if len(answerDiv.Nodes) == 0 {
		return "", errors.New("No answers found")
	}

	answerDiv = answerDiv.First()

	if !h.ShowAllText {
		// grabbing <code> or <pre> content
		instructions := answerDiv.Find("pre")

		if len(instructions.Nodes) == 0 {
			instructions = answerDiv.Find("code")
		}
		return instructions.First().Text(), nil
	}
	output := fmt.Sprintf(answerHeader, answerDiv.Find(".post-text > *").Text(), link)
	return output, nil
}
