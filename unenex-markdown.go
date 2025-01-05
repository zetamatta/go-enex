package enex

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/godown"
)

func ToMarkdowns(rootDir, enexName string, source []byte, styleSheet string, wDebug, wLog io.Writer) error {
	exports, err := ParseMulti(source, wDebug)
	if err != nil {
		return err
	}
	err = makeDir(rootDir, enexName, wLog)
	if err != nil {
		return err
	}
	rootDir = filepath.Join(rootDir, enexName)

	index, err := os.Create(filepath.Join(rootDir, "README.md"))
	if err != nil {
		return err
	}
	defer index.Close()

	if enexName != "" {
		fmt.Fprintf(index, "# %s\n\n", enexName)
	}

	for _, note := range exports {
		safeName := ToSafe.Replace(note.Title)

		fmt.Fprintf(index, "* [%s](%s)\n",
			note.Title,
			url.PathEscape(safeName+".md"),
		)

		html, imgSrc := note.HtmlAndDir()

		var markdown strings.Builder
		godown.Convert(&markdown, strings.NewReader(html), nil)
		fname := filepath.Join(rootDir, safeName+".md")
		fd, err := os.Create(fname)
		if err != nil {
			return err
		}
		ShrinkMarkdown(strings.NewReader(markdown.String()), fd)
		fd.Close()
		fmt.Fprintln(wLog, "Create File:", fname)

		if err := extractAttachment(rootDir, imgSrc.Images, wLog); err != nil {
			return err
		}
	}
	return nil
}

func FilesToMarkdowns(rootDir string, enexFiles []string, wDebug, wLog io.Writer) error {
	wReadme, err := os.Create(filepath.Join(rootDir, "README.md"))
	if err != nil {
		return err
	}
	defer wReadme.Close()

	for _, enexFileName := range enexFiles {
		data, err := os.ReadFile(enexFileName)
		if err != nil {
			return err
		}
		enexName := getEnexBaseName(enexFileName)
		if err := ToMarkdowns(rootDir, enexName, data, "", wDebug, wLog); err != nil {
			return err
		}
		fmt.Fprintf(wReadme, "- [%s](%s/README.md)\n",
			enexName, url.PathEscape(enexName))
	}
	return nil
}
