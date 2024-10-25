package inject

import (
	"fmt"
	"github.com/xgodev/boost/extra/graph"
	"os"
)

func ExportInjectGraphToGraphviz(g *graph.Graph[Component], filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Cabeçalho do arquivo DOT
	_, err = file.WriteString(`digraph InjectGraph {
        graph [splines=true, overlap=false];
        node [fontname="Arial"];
        edge [color="#606060"];
    `)
	if err != nil {
		return err
	}

	// Iterar sobre os vértices do grafo
	for _, vertex := range g.Vertices {
		component := vertex.Value
		entry := component.Entry

		// Nome completo da função com o pacote
		fullFunctionName := fmt.Sprintf("%s.%s", entry.Package, entry.Func.Name)

		// Adiciona o nó da função ao DOT
		_, err = file.WriteString(fmt.Sprintf(
			"\t\"%s\" [label=\"%s\", shape=box, style=filled, fillcolor=\"#A0A0FF\"];\n",
			fullFunctionName, entry.Func.Name))
		if err != nil {
			return err
		}

		// Adicionar as arestas (dependências) para as injeções
		injectName := component.An.Name

		// Verificar se o vértice correspondente à injeção existe
		if paramVertex, exists := g.Vertices[injectName]; exists {
			paramComponent := paramVertex.Value
			paramFunction := fmt.Sprintf("%s.%s", paramComponent.Entry.Package, paramComponent.Entry.Func.Name)

			// Adicionar a aresta entre a função que fornece o parâmetro e a função que consome
			_, err = file.WriteString(fmt.Sprintf(
				"\t\"%s\" -> \"%s\";\n",
				paramFunction, fullFunctionName))
			if err != nil {
				return err
			}
		} else {
			// Aqui você pode colocar um log ou tratamento de erro, caso não encontre o parâmetro.
			fmt.Printf("Parâmetro não encontrado: %s\n", injectName)
		}
	}

	// Fechamento do arquivo DOT
	_, err = file.WriteString("}\n")
	if err != nil {
		return err
	}

	return nil
}
