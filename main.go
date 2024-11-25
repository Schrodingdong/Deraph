package main

import (
	"fmt"
	"os"

	"github.com/schrodi/deraph/grapher"
	"github.com/schrodi/deraph/parser"
)

func main(){
  args := os.Args
  if len(args) == 1 {
    fmt.Println("DERAPH - Python dependency grapher")
    fmt.Println("usage: ")
    fmt.Println("\tderaph <path/to/python/project>")
  } else if len(args) > 2 {
    fmt.Println("Wrong number of args")
    fmt.Println("usage: ")
    fmt.Println("\tderaph <path/to/python/project>")
  }
  // projectPath := "../Zeta-functions/src/docker_proxy"
  projectPath := args[1]
  
  // Construct the project dir
  root := parser.TraverseDir(
    projectPath,
    func(node *parser.FileNode) bool {
      return parser.GetFileExtension(node) == "py"
    },
  )
  fmt.Println("FILE TREE ============================")
  parser.PrintFileTree(root)
  fmt.Println()
  
  // Construct the proje-ct package/module tree
  rootPackage := parser.BuildPackageTree(root)
  fmt.Println("PACKAGE TREE ============================")
  parser.PrintPackageTree(rootPackage)
  fmt.Println()
  
  // Extract imports from 
  fmt.Println("IMPORT EXTRACTION ============================")
  FileToImportDepMap := parser.GenerateFileToImportDepMapForPackageTree(rootPackage)

  // Generate the graphviz file
  graphvizContent := grapher.GenerateGraphvizFromFileToImportDepMap(rootPackage, FileToImportDepMap)
  grapher.GenerateGraphvizFile("testuuu.gv", graphvizContent)
}

