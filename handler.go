package urlshort

import (
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// MapHandler retornará um http.HandlerFunc (que também
// implementa http.Handler) que tentará mapear qualquer
// caminhos (chaves no mapa) para sua URL correspondente (valores
// para onde cada chave no mapa aponta, em formato de string).
// Se o caminho não for fornecido no mapa, então o fallback
// http.Handler será chamado em seu lugar.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Println(pathsToUrls[path])
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler irá analisar o YAML fornecido e então retornar
// um http.HandlerFunc (que também implementa http.Handler)
// que tentará mapear quaisquer caminhos para seus correspondentes
// URL. Se o caminho não for fornecido no YAML, o
// o substituto http.Handler será chamado em seu lugar.
//
// Espera-se que YAML esteja no formato:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// Os únicos erros que podem ser retornados estão todos relacionados a ter
// dados YAML inválidos.
//
// Consulte MapHandler para criar um http.HandlerFunc semelhante via
// um mapeamento de caminhos para urls.
func YAMLHandler(ymlbytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(ymlbytes)
	if err != nil {
		return nil, err
	}
	urls := buildMap(parsedYaml)
	return MapHandler(urls, fallback), nil
}

func parseYAML(in []byte) (dst []map[string]string, err error) {
	err = yaml.Unmarshal(in, &dst)
	return dst, err
}

func buildMap(mapIn []map[string]string) map[string]string {
	var mapPathsUrl = make(map[string]string)

	for i := range mapIn {
		mapAux := mapIn[i]
		mapPathsUrl[mapAux["path"]] = mapAux["url"]
	}

	return mapPathsUrl
}

func ReadYaml(pathyaml string) (yml []byte, err error) {
	yml, err = os.ReadFile(pathyaml)
	if err != nil {
		return nil, err
	}

	return
}
