# Mapwrap - GO

This project was started as a way to learn [GO](http://golang.org). It wraps [mapserver.org's](http://mapserver.org)  mapserv cgi to make serving mapfiles a little easier,
and is based on the original [mapwrap](http://github.com/gina-alaska/mapwrap).  Thanks go to @spruceboy and @dayne for the original mapwrap.


## Configuring mapwrap-go

When mapwrap starts up, it looks first for the environment variable $MAPWRAP_CONFIG and then in the current directory for a file named 'mapwrap.json'.

Mapwrap-go assumes (hopefully) sane defaults, but the following options can be configured.

```json
"mapserv":  "/path/to/mapserv"  //Mapwrap-go will try to find this in your $PATH
"directory": "/path/to/maps"  //This is the working directory for the http server. Defaults to the current directory
"port": "8080"  //port mapwrap-go listens on,  default is 8080
"environment": ["VAR=value","VAR2=something"],  //Any additional environment variables that mapserv might need to run. Default []
"maps": [{...}]  //This is where your mapfiles are configured.
```

###  Configuring maps
Map configuration has only one required value: 'name'. It will try to determine what to do based on that name, but its behavior can be overridden if necessary.

Ex:  
```json
{
  "directory": "/wms/maps",
  "maps": [{
    "name": "bdl"
  }]
}
```

This would look in '/wms/maps' for a file named 'bdl.map', and bind http://server/bdl/ to serve that mapfile.

If you would like to change path that it binds to, you can use the 'path' option.
```json
{
  "directory": "/wms/maps",
  "maps": [{
    "name": "bdl",
    "path": "/wms/bdl"
  }]
}
```

Additionally, you can specify different mapfiles based on the SRS requested, and even alias them to a common name.
```json
  "maps": [{
    "name": "bdl",
    "projections": ["4326", "alaska"],
    "aliases": [{
      "alaska": ["3338", "102006"]
    }]
  }]
```

This would cause mapserv to look for a file named bdl_4326.map when it recives an SRS=EPSG:4326, and bdl_alaska.map when SRS=EPSG:3338 or EPSG:102006. All other codes would go to bdl.map


##TODO
- [ ] Logging doesn't conform to the common log format, yet. 
- [ ] Logfile is always in the current directory, this should be configurable.
- [ ] Check if necessary mapfiles exist on startup.
- [ ] Tests! Both test coverage and real-world. 
- [ ] Better documentation

