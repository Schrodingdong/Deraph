package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/schrodi/deraph/grapher"
	"github.com/schrodi/deraph/parser"
)

var hasExtDep bool
var projectPath string 
var outputFilename string 
func init() {
  flag.BoolVar(&hasExtDep, "ext", false, "Add external dependencies to the output")
  flag.StringVar(&projectPath, "path", "", "Python project to analyze")
  flag.StringVar(&outputFilename, "out", "graphviz.gv", "output filename")
  flag.Parse()
}

func main(){
  if len(os.Args) == 1 {
    printBanner()
    printUsage()
    os.Exit(0)
  }
  if !doesPathExist(projectPath) {
    fmt.Printf("Project path %q doesn't exist or is invalid\n", projectPath)
    printUsage()
    os.Exit(1)
  }

  // Construct the project dir
  root := parser.TraverseDir(
    projectPath,
    func(node *parser.FileNode) bool {
      return parser.GetFileExtension(node) == "py"
    },
  )

  rootPackage := parser.BuildPackageTree(root)
  FileToImportDepMap := parser.GenerateFileToImportDepMapForPackageTree(rootPackage)
  graphvizContent := grapher.GenerateGraphvizFromFileToImportDepMap(rootPackage, FileToImportDepMap, hasExtDep)
  grapher.GenerateGraphvizFile(outputFilename, graphvizContent)
}

func printUsage() {
    fmt.Println("usage: ")
    fmt.Println("deraph [-ext] --path </path/to/project> --out <outFileName>")
    fmt.Println()
    fmt.Println("FLAGS:")
    flag.PrintDefaults()
}

func printBanner() {
    fmt.Print(`
   ___                    __ __
  / _ \___ _______ ____  / // /
 / // / -_) __/ _ \/ _ \/ _  / 
/____/\__/_/  \_,_/ .__/_//_/  
                 /_/           
Simple CLI to generate graphviz graphs of your python project's dependencies.

`)
}

func doesPathExist(path string) bool{
  info, err := os.Stat(path)
  if err != nil { return false }
  return info.IsDir()
}

