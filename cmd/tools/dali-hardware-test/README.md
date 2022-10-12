# cmd/tools/dali-hardware-test

A program that periodically flashes all luminaires on one or more DALI buses, in order to test the the DALI bridge
PLC project is configured correctly, and the lighting panels are wired right etc.

Reads a config file, in JSON format, from the file specified by the `-config` flag.

```json
{
  "ads": {
    "netID": "1.2.3.4.1.1",
    "port": 851
  },
  
  "prefixes": [
    "GVL_Bridge.bus_T1_1"
  ]
}
```