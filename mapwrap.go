package main

import( 
  "net/http"
  "log"
  "fmt"
  "os"
  "flag"
)

var configFile string
var mapDirectory string
var config *Config

func init() {
	const (
    defaultConfig = "mapwrap.json"
		configUsage   = "the path to the mapwrap file"
    defaultMapDir = "maps"
    mapDirUsage   = "the path to map files"
	)
	flag.StringVar(&configFile, "config", defaultConfig, configUsage)
	flag.StringVar(&configFile, "c", defaultConfig, configUsage+" (shorthand)")
  flag.StringVar(&mapDirectory, "maps", defaultMapDir, mapDirUsage)
  flag.StringVar(&mapDirectory, "m", defaultMapDir, mapDirUsage)
  
}

func main() {
  flag.Parse()
  // Load the configuration files
  config, err := loadConfig()
  
  logFile, err := os.OpenFile("mapwrap.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
  if err != nil {
      log.Fatalf("Error opening file: %v", err)
  }
  defer logFile.Close()

  log.SetOutput(logFile)
  
  for _, m := range config.Maps {
    path := fmt.Sprintf("%s", m.UrlPath())
    log.Printf("Attaching %s to %v\n", m.Name, path)
    http.HandleFunc(path, m.serveMap)
  }

  err = http.ListenAndServe(":" + config.Port, nil)
  if err != nil {
    fmt.Sprintf("Unable to start: %v", err)    
  }

}