package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/schrodi/deraph/grapher"
	"github.com/schrodi/deraph/parser"
)

var hasExtDep bool
var projectPath string 
var outputPath string 
var verbose bool 
func init() {
  cwd, _ := os.Getwd()
  defaultOuputPath := filepath.Join(cwd, "graphviz.gv")
  flag.BoolVar(&hasExtDep, "ext", false, "Add external dependencies to the output")
  flag.BoolVar(&verbose, "v", false, "Verbose output")
  flag.StringVar(&projectPath, "path", "", "Python project to analyze")
  flag.StringVar(&outputPath, "out", defaultOuputPath, "output path")
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
  if verbose {
    fmt.Println("PACKAGE TREE =====================================")
    parser.PrintPackageTree(rootPackage)
    fmt.Println()
  }

  FileToImportDepMap := parser.GenerateFileToImportDepMapForPackageTree(rootPackage)
  if verbose {
    fmt.Println("FILE TO IMPORT DEPENDENCY MAP ====================")
    parser.PrintFileToImportDepMap(FileToImportDepMap)
    fmt.Println()
  }
  
  fmt.Printf("Generating graphviz content to %q\n\n", outputPath)
  graphvizContent := grapher.GenerateGraphvizFromFileToImportDepMap(rootPackage, FileToImportDepMap, hasExtDep)
  grapher.GenerateGraphvizFile(outputPath, graphvizContent)
  absPath, err := filepath.Abs(outputPath)
  if err != nil { panic(err) }
  fmt.Printf("Generation successful! Graphviz content available in: \n%v\n", absPath)
}

func printUsage() {
    fmt.Println("usage: ")
    fmt.Println("deraph [-ext] [-v] --path </path/to/project> [--out <outFileName>]")
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

