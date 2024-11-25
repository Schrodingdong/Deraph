package parser

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// PYTHON CONSTRUCTION
// type ImportStatment struct{
//   stmt        string
//   fromPart    string
//   importPart  string
//   isRelative  bool
// }
//
// func ExtractImportsFromTree(root *FileNode) {
//   println(root.Path)
//   if root.IsDir {
//     for _, child := range root.ChildNodes {
//       ExtractImportsFromTree(child)
//     }
//   } else {
//     ExtractImportsFromFileNode(root)
//   }
// }
//
// func ExtractImportsFromFileNode(fileNode *FileNode) {
//   content := fileNode.Content
//
//   rg := regexp.MustCompile(`(?:from\s+([\w\.]+)\s+)*(?:import\s+([\w, ]+))`)
//   importStatements := rg.FindAllSubmatch(content, -1)
//   importStmtList := []ImportStatment{}
//   for _, matches := range importStatements {
//     importStmt := ImportStatment{
//       stmt      : strings.Trim(string(matches[0]), "\n\t "),
//       fromPart  : string(matches[1]),
//       importPart: string(matches[2]),
//     }
//     importStmt.isRelative = strings.HasPrefix(importStmt.fromPart, ".")
//     importStmtList = append(importStmtList, importStmt)
//   }
//
//   for _, imprt := range importStmtList {
//     println("stmt: ", imprt.stmt)
//     println("fromPart  : ", imprt.fromPart)
//     println("importPart: ")
//     analyzeImport(imprt)
//     println("isRelative: ", imprt.isRelative)
//     println("-------------------------------------")
//   }
// }
//
// func analyzeImport(imprt ImportStatment) {
//   importList := strings.Split(imprt.importPart, ",")
//   for _, imprt := range importList {
//     imprt = strings.Trim(imprt, " \n")
//     println("\t" + imprt)
//     // examples
//     // ... import module
//     // ... import module as mdl, module
//
//
//     // TODO:
//     // Search for file with same name as the importPart
//     // NOTE: to contextualize our searche, we will use the FromPart of the import
//     // if exist and is a file (.py)
//     // if exist and is a dir
//     // if not exist
//     //    check if it is a pythonObject (aka a variable, function, or class in the file)
//   }
// }

// Dirs containing py files (and/or __init__.py file)
type PyPackage struct {
  Name       string
  FileRef    *FileNode
  ModuleList []*PyModule
  SubPackageList []*PyPackage
}


// python files
type PyModule struct {
  Name      string
  FileRef   *FileNode
  ObjList   []*PyObject
}

// Variables, functions, classes ... defined in modules
type PyObject struct {
  Name      string
  ObjType   string
}

func BuildPackageTree(root *FileNode) *PyPackage{
  if !root.IsDir {
    return nil
  }
  packageNode := PyPackage{
    Name: root.Name,
    FileRef: root,
  }
  // Search for Modules
  packageNode.ModuleList = extractModulesFromPackage(&packageNode)
  // Search subpackages
  packageList := []*PyPackage{}
  for _, child := range root.ChildNodes {
    if !child.IsDir {
      continue
    }
    packageList = append(packageList, BuildPackageTree(child))
  }
  packageNode.SubPackageList = packageList
  return &packageNode
}

func extractModulesFromPackage(pyPackage *PyPackage) []*PyModule{
  packagePath := pyPackage.FileRef.Path
  moduleList := []*PyModule{}
  dirEntry, err := os.ReadDir(packagePath)
  if err != nil {
    panic(err)
  }
  for _, entry := range dirEntry {
    if entry.IsDir() {
      continue
    }
    modulePath := packagePath + "/" + entry.Name()
    moduleContent, err := os.ReadFile(modulePath)
    if err != nil {
      panic(err)
    }
    moduleFile := &FileNode{
      Name: entry.Name(),
      Path: modulePath,
      IsDir: false,
      ChildNodes: nil,
      Content: moduleContent,
    }
    pyModule := PyModule{
      Name: entry.Name(),
      FileRef: moduleFile,
    }
    pyModule.ObjList = extractPyObjectFromModule(&pyModule)
    moduleList = append(moduleList, &pyModule)
  }
  return moduleList
}

func extractPyObjectFromModule(pyModule *PyModule) []*PyObject{
  content := pyModule.FileRef.Content
  moduleObjects := []*PyObject{}

  rgFunc := regexp.MustCompile(`(?:def\s+)([a-zA-Z_][a-zA-Z0-9_]*)`)
  rgClass := regexp.MustCompile(`(?:class\s+)([a-zA-Z_][a-zA-Z0-9_]*)`)
  funcMatches  := rgFunc.FindAllStringSubmatch(string(content), -1)
  classMatches := rgClass.FindAllStringSubmatch(string(content), -1)
  for _, matchEntry := range funcMatches {
    funcName := matchEntry[1]
    funcObject := PyObject {
      Name: funcName,
      ObjType: "FUNCTION",
    }
    moduleObjects = append(moduleObjects, &funcObject)
  }
  for _, matchEntry := range classMatches {
    className := matchEntry[1]
    classObject := PyObject {
      Name: className,
      ObjType: "CLASS",
    }
    moduleObjects = append(moduleObjects, &classObject)
  }
  return moduleObjects
}








type ImportStatment struct{
  stmt        string
  fromPart    string
  importPart  string
  isRelative  bool
}

type Dependency struct {
  Module *PyModule
  IsExternal bool
}

// Recursively extract imports for the modules in the package tree.
func GenerateFileToImportDepMapForPackageTree(packageRoot *PyPackage) map[*FileNode][]*Dependency{
  return generateFileToImportDepMapForPackageTree(packageRoot, packageRoot)
}


func generateFileToImportDepMapForPackageTree(currPackage *PyPackage, packageRoot *PyPackage) map[*FileNode][]*Dependency{
  // mapping: fileWeAnalyze --> []moduleFile
  fileToImportDepMap := make(map[*FileNode][]*Dependency)
  for _, module := range currPackage.ModuleList {
    importList := extractImportsFromModule(module)
    moduleDepList := []*Dependency{}
    for _, imprtStmt := range importList {
      imprtList := strings.Split(imprtStmt.importPart, ",")
      for _,imprt := range imprtList{
        // Clear import from "as .."
        var moduleName string
        moduleName = strings.Trim(imprt, " \n")
        moduleName = strings.Split(moduleName, " ")[0]
        moduleRef, isExt := checkImportNameIsModuleInPackageTree(moduleName, packageRoot) 
        moduleDep := Dependency{
          Module: moduleRef,
          IsExternal: isExt,
        }
        moduleDepList = append(moduleDepList, &moduleDep)
      }
    }
    fileToImportDepMap[module.FileRef] = moduleDepList
  }
  
  // For subpackages
  for _, subPackage := range currPackage.SubPackageList {
    for k, v := range generateFileToImportDepMapForPackageTree(subPackage, packageRoot) {
      fileToImportDepMap[k] = v
    }
  }
  return fileToImportDepMap
}

func checkImportNameIsModuleInPackageTree(importName string, pyPackage *PyPackage) (*PyModule, bool) {
  moduleRef, isExt := checkImportNameIsModuleInPackage(importName, pyPackage)
  if isExt {
    for _, subpackage := range pyPackage.SubPackageList {
      moduleRef, isExt := checkImportNameIsModuleInPackageTree(importName, subpackage)
      if !isExt {
        return moduleRef, isExt
      }
    }
  }
  return moduleRef, isExt
}

func checkImportNameIsModuleInPackage(importName string, pyPackage *PyPackage) (*PyModule, bool) {
  for _, module := range pyPackage.ModuleList {
    // check modules
    if importName+".py" == module.Name {
      return module, false
    }
    // Check module Objects
    for _, pyObj := range module.ObjList {
      if importName == pyObj.Name {
        return module, false
      }
    }
  }
  extModule := PyModule{
    Name: "EXT-" + importName,
    FileRef: nil,
    ObjList: nil,
  }
  return &extModule, true
}

func extractImportsFromModule(module *PyModule) []*ImportStatment{
  content := module.FileRef.Content
  rg := regexp.MustCompile(`(?:from\s+([\w\.]+)\s+)*(?:import\s+([\w, ]+))`)
  importStatements := rg.FindAllSubmatch(content, -1)
  importStmtList := []*ImportStatment{}
  for _, matches := range importStatements {
    importStmt := ImportStatment{
      stmt      : strings.Trim(string(matches[0]), "\n\t "),
      fromPart  : string(matches[1]),
      importPart: string(matches[2]),
    }
    importStmt.isRelative = strings.HasPrefix(importStmt.fromPart, ".")
    importStmtList = append(importStmtList, &importStmt)
  }
  return importStmtList
}



























// FILE TREE CONSTRUCTION ====================================================
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
      if dirNode != nil {
        root.ChildNodes = append(root.ChildNodes, dirNode)
      }
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
  // If a dir doesn't have (fitlered) files, return nil rootNode
  if root.IsDir && len(root.ChildNodes) == 0 {
    return nil
  }
  return &root
}
