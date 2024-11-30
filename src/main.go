package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/schrodi/deraph/grapher"
	"github.com/schrodi/deraph/parser"
)

var verbose bool 
var hasExtDep bool
var toGraphviz bool
var projectPath string 
var outputDir string 
func init() {
  defaultOuputDir, _ := os.Getwd()
  flag.BoolVar(&hasExtDep, "ext", false, "Add external dependencies to the output")
  flag.BoolVar(&verbose, "v", false, "Verbose output")
  flag.BoolVar(&toGraphviz, "toGraphviz", false, "Generates the graphviz file only")
  flag.StringVar(&projectPath, "path", "", "Python project to analyze")
  flag.StringVar(&outputDir, "outDir", defaultOuputDir, "output dir")
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

  // Generate the graphviz content
  graphvizContent := grapher.GenerateGraphvizFromFileToImportDepMap(rootPackage, FileToImportDepMap, hasExtDep)
  // Save it
  var graphvizFileName = "graphviz.gv"
  var graphvizOutputPath string
  if toGraphviz {
    graphvizOutputPath = path.Join(outputDir, graphvizFileName)
  } else {
    graphvizOutputPath = path.Join("/tmp", graphvizFileName)
  }
  _, err := grapher.GenerateGraphvizFile(graphvizOutputPath, graphvizContent)
  if err != nil { panic(err) }
  absPath, err := filepath.Abs(graphvizOutputPath)
  defer func() { 
    // delete temp file if not toGraphviz output
    if !toGraphviz {
      os.Remove(absPath)
    }
  }()
  if err != nil { panic(err) }
  if toGraphviz {
    fmt.Printf("Generation successful! Graphviz content available in: \n%v\n", absPath)
    return
  }

  // Generate the image
  imageFileName := "project_dep_graph.png"
  imageOutputPath := path.Join(outputDir, imageFileName)
  cmd := exec.Command("dot", "-Tpng", absPath, "-o", imageOutputPath)
  if err := cmd.Run(); err != nil { panic(err) }
  fmt.Printf("Generation successful! Image available in: \n%v\n", imageOutputPath)
  return
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

