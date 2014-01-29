package main

import(
  "encoding/json"
  "os"
  "os/exec"
  "path/filepath"
  "log"
  "bytes"
  "io"
)

type Config struct {
  Environment []string
  Mapserv string
  Port string 
  Directory string
  Maps []Map
}

const defaultConfig = `
{
  "environment": [],
  "port": "8080"
}
`

var config *Config

func GetConfig() *Config {
  return config
}


func loadConfig() {
  var temp Config
  
  if err := decodeConfig(bytes.NewBufferString(defaultConfig), &temp); err != nil {
    log.Fatal(err)
  }
  
  //Look for a configuration file in the following order:
  // Environment:  MAPWRAP_CONFIG
  // Current Directory: mapwrap.json
  configFile := os.Getenv("MAPWRAP_CONFIG")
  if configFile == "" {
    cwd, err := os.Getwd()
    if err != nil {
      log.Fatal(err)
    }
    configFile = filepath.Join(cwd, "mapwrap.json")
  }

  f, err := os.Open(configFile)
  defer f.Close()
  
  if err != nil {
    log.Printf("Error opening configuration file: %s\n", configFile)
    log.Fatal(err)
  }
  
  if err := decodeConfig(f, &temp); err != nil {
    log.Printf("Error loading configuration file: %s\n", configFile)
    log.Fatal(err)
  }
  //Set the working directory if it's not already set
  if temp.Directory == "" {
    temp.Directory, err = os.Getwd()
    if err != nil {
      log.Fatal(err)
    } 
  }
  //Make sure the directory exists
  _, err = os.Stat(temp.Directory)
  if err != nil {
    log.Fatal(err)
  }
  
  if temp.Mapserv == "" {
    out, err := exec.Command("which", "mapserv").Output()
  
    if err != nil {
      log.Fatal("Error attempting to find mapserv: ", err)
    } 
    temp.Mapserv = string(out)
  }
  _, err = exec.Command(temp.Mapserv).Output()
  if err != nil {
    log.Fatal("Error attempting to run mapserv: ", err)
  }
  
  config = &temp

}

// Decodes configuration in JSON format from the given io.Reader into
// the config object pointed to.
func decodeConfig(r io.Reader, c *Config) error {
  decoder := json.NewDecoder(r)
  return decoder.Decode(c)
}
