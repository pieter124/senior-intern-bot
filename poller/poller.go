package poller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	domain "senior_intern_bot/domain"
	"time"
)

type Poller interface {
	Poll() ([]domain.Posting, error)
}

type BoardPoller struct {
	greenhouseCompanies []string
	ashbyCompanies      []string
	client              *http.Client
}

func getGreenhouseURL(company string) string {
	return fmt.Sprintf("https://boards-api.greenhouse.io/v1/boards/%s/jobs", company)
}

func getAshbyURL(company string) string {
	return fmt.Sprintf("https://api.ashbyhq.com/posting-api/job-board/%s", company)
}

func New() *BoardPoller {
	return &BoardPoller{
		greenhouseCompanies: GreenhouseCompanies,
		ashbyCompanies:      AshbyCompanies,
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: &http.Transport{MaxConnsPerHost: 32},
		},
	}
}

type pollResult struct {
	postings []domain.Posting
	err      error
}

func (poller *BoardPoller) Poll() ([]domain.Posting, error) {
	total := len(poller.greenhouseCompanies) + len(poller.ashbyCompanies)

	// Buffered to total so no send ever blocks.
	results := make(chan pollResult, total)
	for _, company := range poller.greenhouseCompanies {
		go poller.pollGreenhouseJob(company, results)
	}
	for _, company := range poller.ashbyCompanies {
		go poller.pollAshbyJob(company, results)
	}

	var (
		polled []domain.Posting
		errs   []error
	)
	for range total {
		r := <-results
		if r.err != nil {
			errs = append(errs, r.err)
			continue
		}
		polled = append(polled, r.postings...)
	}

	return polled, errors.Join(errs...)
}

func (poller *BoardPoller) pollGreenhouseJob(company string, results chan<- pollResult) {
	postings, err := poller.pollGreenhouse(company)
	if err != nil {
		err = fmt.Errorf("greenhouse %s: %w", company, err)
	}
	results <- pollResult{postings, err}
}

func (poller *BoardPoller) pollAshbyJob(company string, results chan<- pollResult) {
	postings, err := poller.pollAshby(company)
	if err != nil {
		err = fmt.Errorf("ashby %s: %w", company, err)
	}
	results <- pollResult{postings, err}
}

func (poller *BoardPoller) fetch(url string, into any) error {
	resp, err := poller.client.Get(url)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(into); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
}

func (poller *BoardPoller) pollAshby(company string) ([]domain.Posting, error) {
	var body struct {
		Jobs []ashbyJob `json:"jobs"`
	}
	if err := poller.fetch(getAshbyURL(company), &body); err != nil {
		return nil, err
	}

	postings := make([]domain.Posting, 0, len(body.Jobs))
	for _, j := range body.Jobs {
		postings = append(postings, j.toPosting(company))
	}
	return postings, nil
}

func (poller *BoardPoller) pollGreenhouse(company string) ([]domain.Posting, error) {
	var body struct {
		Jobs []greenhouseJob `json:"jobs"`
	}
	if err := poller.fetch(getGreenhouseURL(company), &body); err != nil {
		return nil, err
	}

	postings := make([]domain.Posting, 0, len(body.Jobs))
	for _, j := range body.Jobs {
		postings = append(postings, j.toPosting(company))
	}
	return postings, nil
}
