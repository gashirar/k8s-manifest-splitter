package main

import (
    "fmt"
    "gopkg.in/yaml.v2"
    "bufio"
    "os"
    "strings"
    "errors"
)

func main() {
    fp, err := os.Open(os.Args[1])
    if err != nil {
       fmt.Println(err)
       return
    }
    defer fp.Close()

    scanner := bufio.NewScanner(fp)

    yamls := []string{}
    line := ""
    for scanner.Scan() {
        s := scanner.Text()
        if strings.HasPrefix(s, "---") {
            yamls = append(yamls, line)
            line = ""
        } else {
            line = line + s + "\n"
        }
    }
    yamls = append(yamls, line)

    if err = scanner.Err(); err !=nil {
       fmt.Println(err)
    }

    for _, y := range yamls {
        if kind, err := getResourceKind(y); err == nil {
            if name, err := getResourceName(y); err == nil {
                filename := name + "_" + kind + ".yaml"
                writeYamlFile(y, filename)        
            }    
        }
    }
    return
}

func getResourceKind(yamlString string) (string, error) {
    k8sManifest := make(map[string]interface{})
    yaml.Unmarshal([]byte(yamlString), &k8sManifest)
    if v, ok := k8sManifest["kind"]; ok {
        return v.(string), nil
    }
    return "", errors.New("Error: Invalid YAML String.")
}

func getResourceName(yamlString string) (string, error) {
    k8sManifest := make(map[string]interface{})
    yaml.Unmarshal([]byte(yamlString), &k8sManifest)
    if _, ok := k8sManifest["metadata"]; ok {
        if v, ok := k8sManifest["metadata"].(map[interface{}]interface{})["name"]; ok {
            return v.(string), nil
        }
    }
    return "", errors.New("Error: Invalid YAML String.")
}

func writeYamlFile(yamlString string, filename string) error {

    fp, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer fp.Close()

    fp.Write([]byte(yamlString))
    fmt.Println("[CREATE] : " + filename)

    return nil
}
