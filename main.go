package main

import (
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
  parser.PrintFileTree(root)
}
