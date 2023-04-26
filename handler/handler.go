package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math/big"
	"net/http"
	"os"

	cryptoRand "crypto/rand"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url" json:"url"`
}

func YAMLHandler(ymlbytes []byte) (urls map[string]string, err error) {
	parsedYaml, err := parseYAML(ymlbytes)
	if err != nil {
		return nil, err
	}
	urls = buildMap(parsedYaml)
	return urls, nil
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

func CadastrarUrl(w http.ResponseWriter, r *http.Request) {
	var url pathUrl
	json.NewDecoder(r.Body).Decode(&url)
	data, err := ReadYaml("teste.yaml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(data)

	randString := RandomString(5, 1, 3)
	url.Path = randString

	buff.WriteString(fmt.Sprintf("\n- path: %s\n  url: %s", url.Path, url.URL))

	os.WriteFile("teste.yaml", buff.Bytes(), fs.ModeAppend)

	w.Write([]byte(fmt.Sprintf("URL encurtada: localhost:8080/%s", url.Path)))
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]

	yaml, err := ReadYaml("teste.yaml")
	if err != nil {
		log.Fatal(err)
	}

	pathsToUrls, err := YAMLHandler(yaml)
	if err != nil {
		panic(err)
	}

	if dest, ok := pathsToUrls[path]; ok {
		http.Redirect(w, r, dest, http.StatusFound)
		return
	} else {
		http.Error(w, "URL não encontrada", http.StatusBadRequest)
	}
}

func RandomString(size, uppercase, number int) string {
	var (
		randString               string
		isNumber, isUppercase, c int
	)

	for i := 0; i < size; i++ {
		isNumber = randInt(100)
		if isNumber < number {
			c = randInt(10) + 48
		} else {
			isUppercase = randInt(100)
			c = randInt(26) + 97
			if isUppercase < uppercase {
				c -= 32
			}
		}
		randString += string(rune(c))
	}

	return randString
}

// RandInt return a non-crypt safe pseudo-random integer
func RandInt(ceil int) int {
	return randInt(ceil)
}

func randInt(top int) int {
	v, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(top)))
	return int(v.Int64())
}
