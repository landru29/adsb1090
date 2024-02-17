// Package aircraftdb manages the aircraft database.
package aircraftdb

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/landru29/adsb1090/internal/model"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/html"
)

const meanDatumSize = 90

type dbPath string

// Empty implements the model.Uniquer interface.
func (p dbPath) Empty() bool {
	return p.String() == ""
}

// Canonical implements the model.Uniquer interface.
func (p dbPath) Canonical() string {
	return path.Base(string(p))
}

// String implements the Stringer interface.
func (p dbPath) String() string {
	return string(p)
}

const (
	csvFieldsCount = 27
)

// Entry is a database entry.
type Entry struct {
	Addr             model.ICAOAddr `json:"a"`
	Registration     string         `json:"b,omitempty"`
	ManufacturerName string         `json:"c,omitempty"`
	Model            string         `json:"d,omitempty"`
	Operator         string         `json:"e,omitempty"`
	Owner            string         `json:"f,omitempty"`
	Built            *time.Time     `json:"g,omitempty"`
	ADSB             bool           `json:"h,omitempty"`
}

// String implements the Stringer interface.
func (e Entry) String() string {
	built := "?"

	if e.Built != nil {
		built = e.Built.Format(time.DateOnly)
	}

	return fmt.Sprintf(
		`Addr:             %s
Registration:     %s
ManufacturerName: %s
Model:            %s
Operator:         %s
Owner:            %s
Build:            %s
ASDB:             %s`,
		e.Addr,
		e.Registration,
		e.ManufacturerName,
		e.Model,
		e.Operator,
		e.Owner,
		built,
		map[bool]string{true: "yes", false: "no"}[e.ADSB],
	)
}

// Database is the aircraft database.
type Database map[model.ICAOAddr]Entry

// DownloadLatest downloads the latest file.
func DownloadLatest(baseURL string, display io.Writer) (Database, error) {
	dbReader, length, name, err := downloadCSV(baseURL)
	if err != nil {
		return nil, err
	}

	if display == nil {
		display = io.Discard
	}

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(dbReader)

	bar := progressbar.NewOptions64(length,
		progressbar.OptionSetWriter(display),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15), //nolint: gomnd
		progressbar.OptionSetDescription(name),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	reader := csv.NewReader(io.TeeReader(dbReader, bar))
	reader.Comma = ','

	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	database := Database{}

	// 00 icao24
	// 01 registration
	// 02 manufacturericao
	// 03 manufacturername
	// 04 model
	// 05 typecode
	// 06 serialnumber
	// 07 linenumber
	// 08 icaoaircrafttype
	// 09 operator
	// 10 operatorcallsign
	// 11 operatoricao
	// 12 operatoriata
	// 13 owner
	// 14 testreg
	// 15 registered
	// 16 reguntil
	// 17 status
	// 18 built
	// 19 firstflightdate
	// 20 seatconfiguration
	// 21 engines
	// 22 modes
	// 23 adsb
	// 24 acars
	// 25 notes
	// 26 categoryDescription
	for _, row := range data[1:] {
		if len(row) < csvFieldsCount {
			continue
		}

		if row[0] != "" {
			addr, err := model.ParseICAOAddr(row[0])
			if err != nil {
				continue
			}

			entry := Entry{
				Addr:             addr,
				Registration:     row[1],
				ManufacturerName: row[3],
				Model:            row[4],
				Operator:         row[10],
				Owner:            row[13],
				ADSB:             strings.ToLower(row[23]) == "true",
			}

			if row[19] != "" {
				built, err := time.Parse("2006-01-02", row[19])
				if err == nil {
					entry.Built = &built
				}
			}

			database[addr] = entry
		}
	}

	return database, nil
}

func downloadCSV(baseURL string) (io.ReadCloser, int64, string, error) {
	_, err := url.Parse(baseURL)
	if err != nil {
		return nil, 0, "", err
	}

	resp, err := http.Get(baseURL) //nolint: gosec,noctx
	if err != nil {
		return nil, 0, "", err
	}

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(resp.Body)

	urls := &model.UniqueList[dbPath]{}

	tokenizer := html.NewTokenizer(resp.Body)

	for {
		token := tokenizer.Next()

		if token == html.ErrorToken && errors.Is(tokenizer.Err(), io.EOF) {
			break
		}

		switch token { //nolint: exhaustive
		case html.ErrorToken:
			return nil, 0, "", fmt.Errorf("cannot parse HTML: %w", tokenizer.Err())
		case html.StartTagToken:
			name, withAttr := tokenizer.TagName()
			if string(name) == "a" && withAttr {
				urls.Add(href(baseURL, tokenizer))
			}
		}
	}

	if urls.Len() == 0 {
		return nil, 0, "", fmt.Errorf("no file found")
	}

	sort.Sort(urls)

	downloadResp, err := http.Get(urls.First().String()) //nolint: noctx
	if err != nil {
		return nil, 0, "", err
	}

	return downloadResp.Body, downloadResp.ContentLength, urls.First().String(), nil
}

func href(baseURL string, tokenizer *html.Tokenizer) dbPath {
	for {
		attrKey, attrValue, moreAttr := tokenizer.TagAttr()

		if string(attrKey) == "href" && len(attrValue) > 0 {
			link := string(attrValue)

			if strings.HasPrefix(path.Base(link), "aircraftDatabase-") {
				if attrValue[0] == '.' {
					linkURL, _ := url.Parse(baseURL)

					linkURL.Path = path.Join(linkURL.Path, link)

					link = linkURL.String()
				}

				return dbPath(link)
			}
		}

		if !moreAttr {
			break
		}
	}

	return ""
}

// String implements the Stringer interface.
func (d Database) String() string {
	return fmt.Sprintf("Aircraft database(%d)", len(d))
}

// Save saves the database in a bolt DB.
func (d Database) Save(filename string, display io.Writer) error {
	if display == nil {
		display = io.Discard
	}

	file, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return err
	}

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(file)

	writer := gzip.NewWriter(file)

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(writer)

	bar := progressbar.NewOptions(meanDatumSize*len(d),
		progressbar.OptionSetWriter(display),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(15), //nolint: gomnd
		progressbar.OptionSetDescription(fmt.Sprintf("saving to %s", filename)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	jsonEncoder := json.NewEncoder(io.MultiWriter(writer, bar))

	if err := jsonEncoder.Encode(d); err != nil {
		return err
	}

	_ = bar.Set(meanDatumSize * len(d))

	return nil
}

// Load loads aircraft database.
func (d *Database) Load(filename string, display io.Writer) error {
	if display == nil {
		display = io.Discard
	}

	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return err
	}

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(file)

	reader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(reader)

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions64(stat.Size(),
		progressbar.OptionSetWriter(display),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15), //nolint: gomnd
		progressbar.OptionSetDescription(fmt.Sprintf("loading from %s", filename)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	return json.NewDecoder(io.TeeReader(reader, bar)).Decode(d)
}
