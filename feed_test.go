package prismprojectcalendar_test

import (
	"io/ioutil"
	"os"
	"testing"

	ppc "github.com/c032/prismproject-calendar"
)

const (
	jsonFilePath = "testdata/youtube-1667312263.json"
	icalFilePath = "testdata/youtube-1667312263.ical"
)

func testFeed(t *testing.T) *ppc.Feed {
	var (
		err error
		f   *os.File
	)

	f, err = os.Open(jsonFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var feed *ppc.Feed

	feed, err = ppc.ParseFeed(f)
	if err != nil {
		t.Fatal(err)
	}

	if feed == nil {
		t.Fatal("ppc.ParseFeed(f) = nil; want non-nil")
	}

	return feed
}

func TestParseFeed(t *testing.T) {
	var (
		err error
		f   *os.File
	)

	f, err = os.Open(jsonFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var feed *ppc.Feed

	feed, err = ppc.ParseFeed(f)
	if err != nil {
		t.Fatal(err)
	}

	if feed == nil {
		t.Fatal("ppc.ParseFeed(f) = nil; want non-nil")
	}

	// TODO

	t.Skip("Incomplete.")
}

func TestFeed_Calendar(t *testing.T) {
	// TODO

	t.Skip("Not implemented.")
}

func TestFeed_CalendarString(t *testing.T) {
	const expectedFilePath = icalFilePath

	var (
		err error

		rawExpectedIcal []byte
		expectedIcalStr string

		feed *ppc.Feed
	)

	rawExpectedIcal, err = ioutil.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatal(err)
	}
	expectedIcalStr = string(rawExpectedIcal)

	feed = testFeed(t)

	calendarStr, err := feed.CalendarString(ppc.DefaultCalendarID, "PRISM Project")
	if err != nil {
		panic(err)
	}

	if calendarStr != expectedIcalStr {
		t.Errorf("want feed.CalendarString(ppc.DefaultCalendarID) to return content of %s", expectedFilePath)
	}
}

func TestFeedItem_Published(t *testing.T) {
	// TODO

	t.Skip("Not implemented.")
}

func TestFeedItem_Scheduled(t *testing.T) {
	// TODO

	t.Skip("Not implemented.")
}

func TestFeedItem_URL(t *testing.T) {
	// TODO

	t.Skip("Not implemented.")
}

func TestFeedItem_ChannelURL(t *testing.T) {
	// TODO

	t.Skip("Not implemented.")
}

func TestFeedItem_Event(t *testing.T) {
	// TODO

	t.Skip("Not implemented.")
}
