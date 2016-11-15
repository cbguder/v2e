package evernote_test

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/cbguder/v2e/evernote"
	"github.com/cbguder/v2e/helpers"
	"github.com/cbguder/v2e/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Evernote Exporter", func() {
	var (
		exporter evernote.Exporter
	)

	BeforeEach(func() {
		os.Setenv("TZ", "America/Los_Angeles")

		exporter = evernote.NewExporter()
	})

	AfterEach(func() {
		os.Unsetenv("TZ")
	})

	It("exports notes", func() {
		notes := []models.Note{
			{
				Title:    "My Special Title 1",
				Body:     "My Special Body 1",
				Tags:     []string{"Work", "Life"},
				Created:  time.Date(2016, time.April, 20, 16, 20, 0, 0, time.Local),
				Modified: time.Date(2017, time.January, 1, 8, 0, 0, 0, time.Local),
			},
			{
				Title:    "My Special Title 2",
				Body:     "My Special Body 2",
				Tags:     []string{"Harder", "Better", "Faster", "Stronger"},
				Created:  time.Date(2012, time.July, 4, 9, 5, 0, 0, time.Local),
				Modified: time.Date(2014, time.February, 1, 15, 30, 0, 0, time.Local),
			},
		}

		outputFile := helpers.CreateTempFile("enex")
		defer helpers.DiscardTempFile(outputFile)

		err := exporter.Export(outputFile.Name(), notes)
		Expect(err).NotTo(HaveOccurred())

		outputBytes, err := ioutil.ReadAll(outputFile)
		Expect(err).NotTo(HaveOccurred())

		Expect(outputBytes).To(ContainSubstring("<title>My Special Title 1</title>"))
		Expect(outputBytes).To(ContainSubstring(`<en-note><div>My Special Body 1</div></en-note>`))
		Expect(outputBytes).To(ContainSubstring("<tag>Work</tag>"))
		Expect(outputBytes).To(ContainSubstring("<tag>Life</tag>"))
		Expect(outputBytes).To(ContainSubstring("<created>20160420T232000Z</created>"))
		Expect(outputBytes).To(ContainSubstring("<updated>20170101T160000Z</updated>"))

		Expect(outputBytes).To(ContainSubstring("<title>My Special Title 2</title>"))
		Expect(outputBytes).To(ContainSubstring(`<en-note><div>My Special Body 2</div></en-note>`))
		Expect(outputBytes).To(ContainSubstring("<tag>Harder</tag>"))
		Expect(outputBytes).To(ContainSubstring("<tag>Better</tag>"))
		Expect(outputBytes).To(ContainSubstring("<tag>Faster</tag>"))
		Expect(outputBytes).To(ContainSubstring("<tag>Stronger</tag>"))
		Expect(outputBytes).To(ContainSubstring("<created>20120704T160500Z</created>"))
		Expect(outputBytes).To(ContainSubstring("<updated>20140201T233000Z</updated>"))
	})
})
