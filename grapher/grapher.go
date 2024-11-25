package grapher

import (
	"bytes"
	"context"
	"os"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/schrodi/deraph/parser"
)

// Node Styling
var PACKAGE_COLOR string = "turquoise"
var PACKAGE_SHAPE string = "tab"
var PACKAGE_STYLE string = "filled"
var MODULE_STYLE string = "rounded"
var MODULE_SHAPE string = "squre"

// Edge Styling
var EDGE_LEN float64 = 2

func initializeGraph() (context.Context, *graphviz.Graphviz, *graphviz.Graph, error){
  ctx := context.Background()
  g, err := graphviz.New(ctx)
  if err != nil { return nil, nil, nil, err }
  graph, err := g.Graph()
  graph.SetLayout("dot")
  if err != nil { return nil, nil, nil, err }
  return ctx, g, graph, nil
}

func GenerateGraphvizFile(graphvizFilePath string, content []byte) {
  // Generate the graphviz file
  file, err := os.Create(graphvizFilePath)
  if err != nil { panic(err) }
  file.Write(content)
  file.Close()
}

func GenerateGraphNodes(graph *graphviz.Graph, pckg *parser.PyPackage) *graphviz.Graph{
  pckgName := pckg.Name
  pckgNode, err := graph.CreateNodeByName(pckgName)
  if err != nil { panic(err) }
  pckgNode.SetFillColor(PACKAGE_COLOR)
  pckgNode.SetStyle(cgraph.NodeStyle(PACKAGE_STYLE))
  pckgNode.SetShape(cgraph.Shape(PACKAGE_SHAPE))

  for _, module := range pckg.ModuleList {
    modName := module.Name
    modNode, err := graph.CreateNodeByName(modName)
    if err != nil { panic(err) }
    modNode.SetStyle(cgraph.NodeStyle(MODULE_STYLE))
    modNode.SetShape(cgraph.Shape(MODULE_SHAPE))

    edgeName := pckgName+"."+modName
    e, err := graph.CreateEdgeByName(edgeName, pckgNode, modNode)
    if err != nil { panic(err) }
    e.SetLen(EDGE_LEN)
  }
  for _, subPckg := range pckg.SubPackageList{
    subPckgGraph := GenerateGraphNodes(graph, subPckg)
    subPckgGraphRootNode, err := subPckgGraph.NodeByName(subPckg.Name)
    if err != nil { panic(err) }
    edgeName := pckg.Name+"."+subPckg.Name
    e, err := graph.CreateEdgeByName(edgeName, pckgNode, subPckgGraphRootNode)
    if err != nil { panic(err) }
    e.SetLen(EDGE_LEN)
  }
  return graph
}


func GenerateGraphvizFromFileToImportDepMap(rootPackage *parser.PyPackage, ftiMap map[*parser.FileNode][]*parser.Dependency) []byte {
  ctx, g, graph, err := initializeGraph()
  if err != nil { panic(err) }
  defer func() {
    if err := graph.Close(); err != nil { panic(err) }
    g.Close()
  }()

  // Traverse the package tree
  graph = GenerateGraphNodes(graph, rootPackage)

  // Render to buffer
  var buf bytes.Buffer
  renderErr := g.Render(ctx, graph, "dot", &buf)
  if renderErr != nil { panic(renderErr) }
  return buf.Bytes()
}
