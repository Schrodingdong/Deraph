package parser

import (
	"fmt"
	"os"
	"regexp"
)


type FileNode struct {
  Name string
  Path string
  IsDir bool 
  Content []byte
  ChildNodes []*FileNode
}


/* 
Traverse the directory, listing all its subdirs and files 

fileFilter: function to determine which files to keep
*/
func TraverseDir(dir string, fileFilter func(node *FileNode) bool) *FileNode{
  entryList, err := os.ReadDir(dir)
  if err != nil {
    fmt.Println(err)
    return nil
  }
  
  // get name from dir 
  rg := regexp.MustCompile(`[.*/]*(\w+)$`)
  match := rg.FindStringSubmatch(dir)
  dirName := match[1]

  root := FileNode {
    Name: dirName,
    Path: dir,
    IsDir: true,
    Content: nil,
    ChildNodes: []*FileNode{},
  }
  
  for _, entry := range entryList {
    if entry.IsDir() {
      entryDir := dir + "/" + entry.Name()
      dirNode := TraverseDir(entryDir, fileFilter)
      root.ChildNodes = append(root.ChildNodes, dirNode)
    } else {
      filePath := dir + "/" + entry.Name()
      fileContent, err := os.ReadFile(filePath)
      if err != nil {
        fmt.Println(err)
        continue
      }
      fileLeaf := FileNode{
        Name: entry.Name(),
        Path: filePath,
        IsDir: false,
        Content: fileContent,
        ChildNodes: nil,
      }
      if fileFilter(&fileLeaf) {
        root.ChildNodes = append(root.ChildNodes, &fileLeaf)
      }
    }
  }

  return &root
}
