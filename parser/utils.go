package parser

import (
	"fmt"
	"regexp"
)

func PrintFileTree(node *FileNode) {
  fmt.Println(node.Path)
  for _, childNode := range node.ChildNodes {
    PrintFileTree(childNode)
  }
}

func GetFileExtension(node *FileNode) string {
  rg := regexp.MustCompile(`.*\.(\w+)$`)
  match := rg.FindStringSubmatch(node.Name)
  if len(match) < 2 {
    panic("Couldn't extract file extension for : " + node.Name)
  }
  return match[1]
}
