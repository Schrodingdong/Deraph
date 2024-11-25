# DERAPH
Simple CLI to generate graphviz graphs of your python project's dependencies.

## Quickstart
- Run the following commands:
    - Download the binary for your OS/ARCH
    - Unzip it
    - Give it execution permissions
    - Move it to a dir included in $PATH
```bash
wget "https://github.com/Schrodingdong/Deraph/releases/download/1.0.0/deraph-(OS)-(ARCH).zip"
unzip deraph-OS-ARCH.zip
chmod 777 deraph-linux-amd64
sudo mv deraph-linux-amd64 /usr/local/bin/deraph  # Or any dir, just make the dir path is included in the $PATH variable
```

- That's it ! to run your program:
```bash
deraph
# Usage template:
# deraph [-ext] [-v] --path <projectDirPath> [--out <outputFilePath>]
```

- To generate the image graph, you can use the [dot](https://graphviz.org/download/) command
```bash
dot -Tpng path/to/outFileName -o imageName.png
```

## From source code
- Clone the project
- Build the program
```bash
go build .
```

- To use the cli:
```
./deraph [-ext] [-v] --path <projectDirPath> [--out <outputFilePath>]
```

## Example
```bash
deraph --path ./example/python_project # output: $(your_cwd)/graphviz.gv
deraph --ext --path ./example/python_project # Includes external dependencies
deraph --ext -v --path ./example/python_project # verbose output
deraph --ext -v --path ./example/python_project --out superdupercoolgraph.gv
```

- Here is the image output

![Project dependency graph](images/pyproject_graph.png "Project dependency graph")

- And here with external dependencies shown

![Project dependency graph with external dependencies shown](images/ext_dep_pyproject_graph.png "Project dependency graph with external dependencies shown")
