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
  
  loadConfig()
}

func main() {
  flag.Parse()
  
  logFile, err := os.OpenFile("mapwrap.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
  if err != nil {
      log.Fatalf("Error opening file: %v", err)
  }
  defer logFile.Close()

  log.SetOutput(logFile)
  //Don't prefix the time in these logs
  log.SetFlags(0)
  
  for _, m := range GetConfig().Maps {
    path := fmt.Sprintf("%s", m.UrlPath())
    http.HandleFunc(path, m.serveMap)
  }


  err = http.ListenAndServe(":" + GetConfig().Port, nil)
  if err != nil {
    fmt.Printf("Unable to start: %v", err)    
  }

}