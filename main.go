package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/docker/libcompose/config"
	//"github.com/docker/libcompose/utils"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

func convertMapKeysToStrings(existingMap map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})

	for k, v := range existingMap {
		newMap[k] = convertKeysToStrings(v)
	}

	return newMap
}

func convertServiceMapKeysToStrings(serviceMap config.RawServiceMap) config.RawServiceMap {
        newServiceMap := make(config.RawServiceMap)

        for k, v := range serviceMap {
                newServiceMap[k] = convertServiceKeysToStrings(v)
        }

        return newServiceMap
}

func convertServiceKeysToStrings(service config.RawService) config.RawService {
        newService := make(config.RawService)

        for k, v := range service {
                newService[k] = convertKeysToStrings(v)
        }

        return newService
}

func convertKeysToStrings(item interface{}) interface{} {
	switch typedDatas := item.(type) {

	case map[interface{}]interface{}:
		newMap := make(map[string]interface{})

		for key, value := range typedDatas {
			stringKey := key.(string)
			newMap[stringKey] = convertKeysToStrings(value)
		}
		return newMap

	case []interface{}:
		// newArray := make([]interface{}, 0) will cause golint to complain
		var newArray []interface{}
		newArray = make([]interface{}, 0)

		for _, value := range typedDatas {
			newArray = append(newArray, convertKeysToStrings(value))
		}
		return newArray

	default:
		return item
	}
}

func main() {
	schema, err := ioutil.ReadFile("schema.json")
	if err != nil {
		log.Fatal(err)
	}
	schemaLoader := gojsonschema.NewStringLoader(string(schema))

	index, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}

	compose, err := ioutil.ReadFile("compose.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(index))
	})

	http.HandleFunc("/compose", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(compose))
	})

	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprint(w, err)
			return
		}
		cloudConfig := r.FormValue("cc")

		var data map[string]interface{}
		if err := yaml.Unmarshal([]byte(cloudConfig), &data); err != nil {
			fmt.Fprint(w, err)
			return
		}
		data = convertServiceKeysToStrings(data)

		loader := gojsonschema.NewGoLoader(data)
		result, err := gojsonschema.Validate(schemaLoader, loader)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		if result.Valid() {
			fmt.Fprint(w, "Valid!")
		} else {
			for _, desc := range result.Errors() {
				fmt.Fprintf(w, "%s<br>", desc)
			}
		}
	})

	http.HandleFunc("/compose-validate", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprint(w, err)
			return
		}
		cloudConfig := r.FormValue("cc")

		var data config.RawServiceMap
		if err := yaml.Unmarshal([]byte(cloudConfig), &data); err != nil {
			fmt.Fprint(w, err)
			return
		}
		data = convertServiceMapKeysToStrings(data)

                if err := config.Validate(data); err == nil {
			fmt.Fprint(w, "Valid!")
                } else {
			fmt.Fprint(w, err)
                }
	})

	log.Fatal(http.ListenAndServe(":9000", nil))
}
