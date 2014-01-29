package main

import(
  "fmt"
  "strings"
  "net"
  "net/http"
  "net/url"
  "net/http/cgi"
  "log"
  "time"
)
var exceptions = []string{"blank", "image", "xml"}

type Map struct {
  Name string
  Projections []string
  Aliases map[string][]string
  Path string 

}

func (m Map) Mapfile(proj string) string {
  srs := strings.Split(proj, ":")
  srid := srs[len(srs)-1:len(srs)][0]
  
  for projection, aliases := range m.Aliases {
    for _, alias := range aliases {
      if alias == srid {
        srid = projection
      }
    }
  }
  for _, p := range m.Projections {
    if p == srid {
      return fmt.Sprintf("%s_%s.map", m.Name, srid)
    }
  }  
  return fmt.Sprintf("%s.map", m.Name)
}

func (m Map) UrlPath() string {
  p := m.Path
  
  if p == "" {
    p = m.Name
  }
  if !strings.HasPrefix(p, "/") {
    p = "/" + p
  }
  if !strings.HasSuffix(p, "/") {
    p = p + "/"
  }

  return p
}

func (m Map) serveMap(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()
  if err != nil {
    //TODO: Check that this is the correct status code to useg
    log.Printf(buildCommonLogFormat(r, time.Now(), 404, 0))
    return
  }
  
  normalizeKeys(r.Form, strings.ToUpper)
  
  if r.Form.Get("REQUEST") == "" {
    r.Form.Set("REQUEST", "GetCapabilities")
  }

  if r.Form.Get("SERVICE") == "" {
    r.Form.Set("SERVICE", "WMS")
  }
  //Don't let the user set the mapfile.
  r.Form.Del("MAP")
  r.Form.Set("MAP", m.Mapfile(r.Form.Get("SRS")))
  
  //ESRI software sends an invalid value?
  //Force it to xml unless it's a valid value
  if invalidException(r.Form.Get("EXCEPTIONS")) {
    r.Form.Set("EXCEPTIONS", "xml")
  }

  fmt.Printf("%v", config)

  queryString := "QUERY_STRING=" + r.Form.Encode()
  // env := append(config.Environment, queryString)
  handler := cgi.Handler{
    Path: config.Mapserv,
    Dir: config.Directory,
    Env: []string{queryString},
  }

  handler.ServeHTTP(w, r)

  //Should be able to say w.Header().Get("Status"), w.Header().Get("Length(?)")
  log.Println(buildCommonLogFormat(r, time.Now(), 200, 0))
}

func buildCommonLogFormat(r *http.Request, ts time.Time, status, size int) string {
  username := "-"
  host, _, err := net.SplitHostPort(r.RemoteAddr)

  if err != nil {
          host = r.RemoteAddr
  }
  return fmt.Sprintf("%s - %s [%s] \"%s %s %s\" %d %d",
                  host,
                  username,
                  ts.Format("02/Jan/2006:15:04:05 -0700"),
                  r.Method,
                  r.URL.RequestURI(),
                  r.Proto,
                  status,
                  size,
          )
}

func normalizeKeys(v url.Values, normalFunc func(string) string) {
  for param, values := range v {
    normalizedParam := normalFunc(param)
    v.Del(param)
    // Mapserv doesn't take multiple values per param
    //  Save a little time and only set the first one
    v.Set(normalizedParam, values[0])
  }    
}

func invalidException(value string) bool {
  for _, e := range exceptions {
    if strings.ToLower(value) == e {
      return false
    }
  }
  return true
}