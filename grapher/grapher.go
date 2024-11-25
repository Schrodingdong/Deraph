package grapher

import (
	"bytes"
	"context"
	"os"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/schrodi/deraph/parser"
)

// Graph style
var GRAPH_LAYOUT string = "dot"

// Node Styling
var PACKAGE_COLOR string = "#a8dadc"
var PACKAGE_SHAPE string = "tab"
var PACKAGE_STYLE string = "filled"
var MODULE_STYLE string = "rounded"
var MODULE_SHAPE string = "rect"

// Edge Styling
var EDGE_LEN float64 = 3
var DEP_EDGE_COLOR string = "#e63946"
var DEP_EDGE_STYLE string = "dashed"

func initializeGraph() (context.Context, *graphviz.Graphviz, *graphviz.Graph, error){
  ctx := context.Background()
  g, err := graphviz.New(ctx)
  if err != nil { return nil, nil, nil, err }
  graph, err := g.Graph()
  graph.SetLayout("dot")
  if err != nil { return nil, nil, nil, err }
  return ctx, g, graph, nil
}

// Initialize a cluster for external dependencies to link to
func initializeExtDepGraph(graph *graphviz.Graph) (*graphviz.Graph, error) {
  extDep, err := graph.CreateSubGraphByName("cluster_ext_dep")
  if err != nil { return nil, err }
  extDep.SetBackgroundColor("#eaeaea")
  extDep.SetClusterRank("TB")
  return extDep, nil
}

func GenerateGraphvizFile(graphvizFilePath string, content []byte) {
  // Generate the graphviz file
  file, err := os.Create(graphvizFilePath)
  if err != nil { panic(err) }
  file.Write(content)
  file.Close()
}

func generateGraphNodes(graph *graphviz.Graph, pckg *parser.PyPackage) *graphviz.Graph{
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
    subPckgGraph := generateGraphNodes(graph, subPckg)
    // TODO: for __init__.py files, will link to ssame  ones.
    subPckgGraphRootNode, err := subPckgGraph.NodeByName(subPckg.Name)
    if err != nil { panic(err) }
    edgeName := pckg.Name+"."+subPckg.Name
    e, err := graph.CreateEdgeByName(edgeName, pckgNode, subPckgGraphRootNode)
    if err != nil { panic(err) }
    e.SetLen(EDGE_LEN)
  }
  return graph
}

func addDependencyEdges(graph *graphviz.Graph, externalDependencies *graphviz.Graph, ftiMap map[*parser.FileNode][]*parser.Dependency, hasExtDep bool) *graphviz.Graph{
  for k, v := range ftiMap {
    fromNode, err := graph.NodeByName(k.Name)
    if err != nil { panic(err) }
    for _, dep := range v {
      depName := dep.Module.Name
      var depNode *cgraph.Node
      if strings.Contains(depName, "EXT") { 
        if hasExtDep {
          depNode, err = externalDependencies.CreateNodeByName(depName)
          if err != nil { panic(err) }
        }
      } else {
        depNode, err = graph.NodeByName(depName)
        if err != nil { panic(err) }
      }
      if depNode != nil {
        depEdgeName := k.Name+"->"+depName
        e, err := graph.CreateEdgeByName(depEdgeName, fromNode, depNode)
        if err != nil { panic(err) }
        e.SetStyle(cgraph.EdgeStyle(DEP_EDGE_STYLE))
        e.SetColor(DEP_EDGE_COLOR)
      }
    }
  }
  return graph
}

func GenerateGraphvizFromFileToImportDepMap(rootPackage *parser.PyPackage, ftiMap map[*parser.FileNode][]*parser.Dependency, hasExtDep bool) []byte {
  ctx, g, graph, err := initializeGraph()
  if err != nil { panic(err) }
  extDep, err := initializeExtDepGraph(graph)
  if err != nil { panic(err) }
  defer func() {
    if err := graph.Close(); err != nil { panic(err) }
    g.Close()
  }()

  // Traverse the package tree
  graph = generateGraphNodes(graph, rootPackage)
  graph = addDependencyEdges(graph, extDep, ftiMap, hasExtDep)

  // Render to buffer
  var buf bytes.Buffer
  renderErr := g.Render(ctx, graph, "dot", &buf)
  if renderErr != nil { panic(renderErr) }
  return buf.Bytes()
}
