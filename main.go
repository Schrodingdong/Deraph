package main

import (
	"fmt"

	"github.com/schrodi/deraph/parser"
)

func main(){
  projectPath := "../Zeta-functions/src/docker_proxy"

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

  // Construct the project package/module tree
  rootPackage := parser.BuildPackageTree(root)
  fmt.Println("PACKAGE TREE ============================")
  parser.PrintPackageTree(rootPackage)
  fmt.Println()
  
  // Extract imports from 
  fmt.Println("IMPORT EXTRACTION ============================")
  parser.GenerateFileToImportDepMapForPackageTree(rootPackage)
}

