package main

import (
    "fmt"
    "gopkg.in/yaml.v2"
    "bufio"
    "os"
    "strings"
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

    if err = scanner.Err(); err !=nil {
       fmt.Println(err)
    }

    for _, y := range yamls {
        kind := getResourceKind(y)
        name := getResourceName(y)
        filename := name + "_" + kind + ".yaml"
        writeYamlFile(y, filename)
    }
    return
}

func getResourceKind(yamlString string) string{
    k8sManifest := make(map[string]interface{})
    yaml.Unmarshal([]byte(yamlString), &k8sManifest)
    return k8sManifest["kind"].(string)
}

func getResourceName(yamlString string) string {
    k8sManifest := make(map[string]interface{})
    yaml.Unmarshal([]byte(yamlString), &k8sManifest)
    return  k8sManifest["metadata"].(map[interface{}]interface{})["name"].(string)
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
