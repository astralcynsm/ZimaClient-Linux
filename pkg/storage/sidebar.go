package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func AddToSidebar(label string, localPath string) error {
	homeDir := getRealUserHome()
	absPath, _ := filepath.Abs(localPath)
	uri := "file://" + absPath

	// GTK
	gtkPath := filepath.Join(homeDir, ".config/gtk-3.0/bookmarks")
	_ = addGTKBookmark(gtkPath, uri, label)

	// Qt
	kdePath := filepath.Join(homeDir, ".local/share/user-places.xbel")
	if _, err := os.Stat(kdePath); err == nil {
		_ = addKDEBookmark(kdePath, uri, label)
	}

	return nil
}

func getRealUserHome() string {
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		return filepath.Join("/home", sudoUser)
	}
	home, _ := os.UserHomeDir()
	return home
}

func addGTKBookmark(path, uri, label string) error {
	b, _ := os.ReadFile(path)
	lines := strings.Split(string(b), "\n")
	newLine := fmt.Sprintf("%s %s", uri, label)

	// remove duplicates
	var newLines []string
	newLines = append(newLines, newLine) // on top
	for _, line := range lines {
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, uri) {
			newLines = append(newLines, line)
		}
	}
	return os.WriteFile(path, []byte(strings.Join(newLines, "\n")+"\n"), 0644)
}

func addKDEBookmark(path, uri, label string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	sContent := string(content)

	// skip if exists
	if strings.Contains(sContent, uri) {
		return nil
	}

	bookmarkXML := fmt.Sprintf(` <bookmark href="%s">
  <title>%s</title>
  <info>
   <metadata owner="http://freedesktop.org">
    <bookmark:icon name="network-server"/>
   </metadata>
   <metadata owner="http://www.kde.org">
    <ID>rnctl_%d</ID>
    <isSystemItem>false</isSystemItem>
   </metadata>
  </info>
 </bookmark>`, uri, label, time.Now().Unix())

	// find the exact insert location
	xbelIndex := strings.Index(sContent, "<xbel")
	if xbelIndex == -1 {
		return fmt.Errorf("invalid xbel")
	}

	closingBracket := strings.Index(sContent[xbelIndex:], ">")
	if closingBracket == -1 {
		return fmt.Errorf("invalid xbel")
	}

	insertPos := xbelIndex + closingBracket + 1
	newContent := sContent[:insertPos] + "\n" + bookmarkXML + sContent[insertPos:]

	return os.WriteFile(path, []byte(newContent), 0644)
}

func RemoveFromSidebar(localPath string) error {
	homeDir := getRealUserHome()
	absPath, _ := filepath.Abs(localPath)
	uri := "file://" + absPath

	gtkPath := filepath.Join(homeDir, ".config/gtk-3.0/bookmarks")
	_ = removeGTKBookmark(gtkPath, uri)

	// rmv qt bookmark
	kdePath := filepath.Join(homeDir, ".local/share/user-places.xbel")
	if _, err := os.Stat(kdePath); err == nil {
		_ = removeKDEBookmark(kdePath, uri)
	}

	return nil
}

func removeGTKBookmark(path, uri string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(b), "\n")
	var newLines []string

	for _, line := range lines {
		// keep those non-targets
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, uri) {
			newLines = append(newLines, line)
		}
	}

	// rew
	return os.WriteFile(path, []byte(strings.Join(newLines, "\n")+"\n"), 0644)
}

func removeKDEBookmark(path, uri string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	sContent := string(b)

	// locate where it starts
	startTag := fmt.Sprintf(`<bookmark href="%s">`, uri)
	startIndex := strings.Index(sContent, startTag)
	if startIndex == -1 {
		return nil
	}
	endTag := "</bookmark>"
	endOffset := strings.Index(sContent[startIndex:], endTag)
	if endOffset == -1 {
		return fmt.Errorf("XML 结构异常，找不到闭合标签")
	}

	// cal end tag loc
	endIndex := startIndex + endOffset + len(endTag)

	// 搞半天还要自己拼
	newContent := sContent[:startIndex] + sContent[endIndex:]

	return os.WriteFile(path, []byte(newContent), 0644)
}
