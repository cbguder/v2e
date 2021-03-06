package evernote

import (
	"encoding/xml"
	"os"
	"strings"
	"time"

	"github.com/cbguder/unc/models"
)

const (
	timeLayout         = "20060102T150405Z07:00"
	exportDoctype      = `<!DOCTYPE en-export SYSTEM "http://xml.evernote.com/pub/evernote-export3.dtd">` + "\n"
	noteDoctype        = `<!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd">` + "\n"
	noStandaloneHeader = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>` + "\n"
)

type Exporter struct {
}

func (e Exporter) Export(outputPath string, notes []models.Note) error {
	enexNotes := make([]enexNote, len(notes))

	for i, note := range notes {
		eNote, err := convertNote(note)
		if err != nil {
			return err
		}
		enexNotes[i] = eNote
	}

	export := enexExport{
		Notes: enexNotes,
	}

	marshaledNote, err := xml.MarshalIndent(export, "", "  ")
	if err != nil {
		return err
	}

	f, _ := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	f.Write([]byte(xml.Header))
	f.Write([]byte(exportDoctype))
	f.Write(marshaledNote)
	f.Write([]byte("\n"))

	f.Close()

	return nil
}

func convertNote(note models.Note) (enexNote, error) {
	innerNote := convertBody(note.Body)

	ein, err := xml.Marshal(innerNote)
	if err != nil {
		return enexNote{}, err
	}

	content := enexContent{
		Body: noStandaloneHeader + noteDoctype + string(ein),
	}

	return enexNote{
		Title:    note.Title,
		Content:  content,
		Tags:     note.Tags,
		Created:  formatTime(note.Created),
		Modified: formatTime(note.Modified),
	}, nil
}

func convertBody(body string) enexInnerNote {
	lines := strings.Split(body, "\n")

	nodes := make([]enexNode, len(lines))

	for i, line := range lines {
		nodes[i] = enexNode{
			XMLName: xml.Name{Local: "div"},
		}

		if line == "" {
			nodes[i].InnerXml = "<br/>"
		} else {
			nodes[i].CharData = line
		}
	}

	return enexInnerNote{
		Children: nodes,
	}
}

func formatTime(t time.Time) string {
	return t.In(time.UTC).Format(timeLayout)
}
