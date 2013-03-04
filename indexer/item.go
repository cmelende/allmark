// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package model defines the basic
	data structures of the docs engine.
*/
package indexer

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/andreaskoch/docs/util"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	UnknownItemType      = "unknown"
	DocumentItemType     = "document"
	PresentationItemType = "presentation"
	CollectionItemType   = "collection"
	MessageItemType      = "message"
	ImageGalleryItemType = "imagegallery"
	LocationItemType     = "location"
	CommentItemType      = "comment"
	TagItemType          = "tag"
	RepositoryItemType   = "repository"
)

type Item struct {
	Path         string
	RenderedPath string
	Files        []File
	ChildItems   []Item
	Blocks       []Block
	Type         string
}

// Create a new repository item
func NewItem(path string, files []File, childItems []Item) (item Item, err error) {

	itemType := getItemType(path)

	if itemType == UnknownItemType {
		err = errors.New(fmt.Sprintf("The item %q does not match any of the known item types.", path))
	}

	return Item{
		Path:         path,
		RenderedPath: getRenderedItemPath(path),
		Files:        files,
		ChildItems:   childItems,
		Type:         itemType,
	}, err
}

func (item Item) GetFilename() string {
	return filepath.Base(item.Path)
}

func (item Item) GetHash() string {
	itemBytes, readFileErr := ioutil.ReadFile(item.Path)
	if readFileErr != nil {
		return ""
	}

	sha1 := sha1.New()
	sha1.Write(itemBytes)

	return fmt.Sprintf("%x", string(sha1.Sum(nil)[0:6]))
}

func (item Item) Walk(walkFunc func(item Item)) {

	walkFunc(item)

	// add all children
	for _, child := range item.ChildItems {
		child.Walk(walkFunc)
	}
}

func (item Item) IsRendered() bool {
	return util.FileExists(item.RenderedPath)
}

func (item Item) GetAbsolutePath() string {
	return item.RenderedPath
}

func (item Item) GetRelativePath(basePath string) string {

	fullItemPath := item.RenderedPath
	relativePath := strings.Replace(fullItemPath, basePath, "", 1)
	relativePath = "/" + strings.TrimLeft(relativePath, "/")

	return relativePath
}

func (item Item) GetBlockValue(name string) string {
	if item.Blocks == nil || len(item.Blocks) == 0 {
		return ""
	}

	for _, element := range item.Blocks {

		if strings.ToLower(element.Name) == strings.ToLower(name) {
			return element.Value
		}

	}

	return ""
}

func (item *Item) AddBlock(name string, value string) {

	block, err := NewBlock(name, value)
	if err != nil {
		panic("Cannot add a block without a name")
	}

	if item.Blocks == nil {
		item.Blocks = make([]Block, 1, 1)
		item.Blocks[0] = block
		return
	}

	item.Blocks = append(item.Blocks, block)

}

// Get the item type from the given item path
func getItemType(itemPath string) string {
	filename := filepath.Base(itemPath)
	return getItemTypeFromFilename(filename)
}

// Get the filepath of the rendered repository item
func getRenderedItemPath(itemPath string) string {
	itemDirectory := filepath.Dir(itemPath)
	renderedFilePath := filepath.Join(itemDirectory, "index.html")
	return renderedFilePath
}

func getItemTypeFromFilename(filename string) string {

	lowercaseFilename := strings.ToLower(filename)

	switch lowercaseFilename {
	case "document.md", "readme.md":
		return DocumentItemType

	case "presentation.md":
		return PresentationItemType

	case "collection.md":
		return CollectionItemType

	case "message.md":
		return MessageItemType

	case "imagegallery.md":
		return ImageGalleryItemType

	case "location.md":
		return LocationItemType

	case "comment.md":
		return CommentItemType

	case "tag.md":
		return TagItemType

	case "repository.md":
		return RepositoryItemType
	}

	return UnknownItemType
}
