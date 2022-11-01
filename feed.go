package prismprojectcalendar

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/emersion/go-ical"
)

var DefaultCalendarID string = sha256Hex("github.com/c032/prismproject-calendar")

const defaultDuration = 1 * time.Hour

var (
	youtubeStreamURL  *url.URL
	youtubeChannelURL *url.URL
)

const (
	youtubeStreamURLStr  = "https://youtu.be/"
	youtubeChannelURLStr = "https://www.youtube.com/channel/"
)

func init() {
	var err error

	youtubeStreamURL, err = url.Parse(youtubeStreamURLStr)
	if err != nil {
		panic(err)
	}

	youtubeChannelURL, err = url.Parse(youtubeChannelURLStr)
	if err != nil {
		panic(err)
	}
}

func sha256Hex(value string) string {
	sum := sha256.Sum256([]byte(value))

	return fmt.Sprintf("%x", sum)
}

type Feed struct {
	Live     []FeedItem `json:"live"`
	Upcoming []FeedItem `json:"upcoming"`
}

func (f *Feed) Calendar(calendarID string, calendarName string) (*ical.Calendar, error) {
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, calendarID)
	cal.Props.SetText(ical.PropName, calendarName)
	cal.Props.SetText("X-WR-CALNAME", calendarName)

	propPublishedTTL := ical.NewProp("X-PUBLISHED-TTL")
	propPublishedTTL.SetDuration(8 * time.Hour)
	cal.Props.Set(propPublishedTTL)

	for i, item := range f.Upcoming {
		var (
			err   error
			event *ical.Event
		)

		event, err = item.Event()
		if err != nil {
			return nil, fmt.Errorf("could not parse event at index %d: %w", i, err)
		}

		cal.Children = append(cal.Children, event.Component)
	}

	return cal, nil
}

func (f *Feed) CalendarString(calendarID string, calendarName string) (string, error) {
	cal, err := f.Calendar(calendarID, calendarName)
	if err != nil {
		return "", fmt.Errorf("could not generate calendar: %w", err)
	}

	var buf bytes.Buffer
	enc := ical.NewEncoder(&buf)
	err = enc.Encode(cal)
	if err != nil {
		return "", fmt.Errorf("could not encode calendar: %w", err)
	}

	calStr := buf.String()

	return calStr, nil
}

type FeedItem struct {
	ID           string `json:"id"`
	ChannelID    string `json:"channelId"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	Viewers      int64  `json:"viewers"`
	RawPublished int64  `json:"published"`
	RawScheduled int64  `json:"scheduled"`
	//RawStart     *int64 `json:"start"`
}

func (fi *FeedItem) URL() (*url.URL, error) {
	var (
		err         error
		feedItemURL *url.URL
	)

	var path = "/" + url.PathEscape(fi.ID)

	feedItemURL, err = youtubeStreamURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("could not create url: %w", err)
	}

	return feedItemURL, nil
}

func (fi *FeedItem) ChannelURL() (*url.URL, error) {
	var (
		err        error
		channelURL *url.URL
	)

	var path = "/channel/" + url.PathEscape(fi.ChannelID)

	channelURL, err = youtubeChannelURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("could not create url: %w", err)
	}

	return channelURL, nil
}

func (fi *FeedItem) Published() (time.Time, error) {
	t := time.Unix(fi.RawPublished, 0)

	return t, nil
}

func (fi *FeedItem) Scheduled() (time.Time, error) {
	t := time.Unix(fi.RawScheduled, 0)

	return t, nil
}

func (fi *FeedItem) Event() (*ical.Event, error) {
	var (
		err error

		published time.Time
		scheduled time.Time
		eventURL  *url.URL
	)

	published, err = fi.Published()
	if err != nil {
		return nil, fmt.Errorf("could not read published time: %w", err)
	}

	scheduled, err = fi.Scheduled()
	if err != nil {
		return nil, fmt.Errorf("could not read scheduled time: %w", err)
	}

	eventURL, err = fi.URL()
	if err != nil {
		return nil, fmt.Errorf("could not read event URL: %w", err)
	}

	event := ical.NewEvent()
	event.Props.SetText(ical.PropUID, fi.ID)
	event.Props.SetDateTime(ical.PropDateTimeStamp, published)
	event.Props.SetText(ical.PropSummary, fi.Title)
	event.Props.SetDateTime(ical.PropDateTimeStart, scheduled)
	event.Props.SetURI(ical.PropURL, eventURL)

	propDuration := ical.NewProp(ical.PropDuration)
	propDuration.SetDuration(defaultDuration)
	event.Props.Set(propDuration)

	return event, nil
}

func ParseFeed(r io.Reader) (*Feed, error) {
	d := json.NewDecoder(r)
	var feed *Feed

	err := d.Decode(&feed)
	if err != nil {
		return nil, fmt.Errorf("could not decode json: %w", err)
	}

	return feed, nil
}
