package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v3"
)

type pathUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path

		if destination, found := pathsToUrls[path]; found {
			http.Redirect(w, r, destination, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrls []pathUrl

	err := yaml.Unmarshal(yamlBytes, &pathUrls)
	if err != nil {
		return nil, err
	}

	pathsToUrls := make(map[string]string)

	for _, path_to_url := range pathUrls {
		pathsToUrls[path_to_url.Path] = path_to_url.Url
	}

	return MapHandler(pathsToUrls, fallback), nil
}
