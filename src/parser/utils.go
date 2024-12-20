package parser

import (
	"fmt"
	"regexp"
	"strings"
)

func PrintFileToImportDepMap(m map[*FileNode][]*Dependency) {
  for k, v := range m {
    if len(v) == 0 {
      continue
    }
    fmt.Println(k.Name, ":")
    for _, dep := range v {
      fmt.Printf("%30q: IS_EXT=%v\n", dep.Module.Name, dep.IsExternal)
    }
  }
}


func PrintPackageTree(node *PyPackage) {
  printPackageTreeWithDepth(node, 0)
}

func printPackageTreeWithDepth(node *PyPackage, depth int8) {
  prefix := strings.Repeat("  ", int(depth))
  fmt.Println(prefix + "[PACKAGE] ", node.Name)
  for _, module := range node.ModuleList {
    fmt.Println(prefix + "  " + "[MODULE] ", module.Name, ":")
    for _, obj := range module.ObjList {
      fmt.Println(prefix + "    " + "> " + obj.ObjType + " ", obj.Name)
    }
  }
  for _,subpackage := range node.SubPackageList {
    printPackageTreeWithDepth(subpackage, depth + 1)
  }
}

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
