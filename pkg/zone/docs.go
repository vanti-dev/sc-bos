// Package zone represents the concept of a physical space and what you can do with that space.
// For example a Meeting Room can be represented by a zone and you might want to turn all the lights off in that zone.
//
// Zones themselves are organised into types, the most general of which is the `area`. Typically a zone is made up of
// a collection of features, each providing some function in the zone, for example lighting control.
//
// # Example
//
// The below configuration can be used to setup an area zone (type area, implemented by the `area` package) named "Room1".
// The area is built to support the lighting feature (implemented by the `feature/lighting` package)
// which looks for the "lights" and "lightGroups" properties. When the zone is started the area configures and starts
// each of the features, announcing them with the controller node as needed.
//
//	{
//	  "name": "Room1",
//	  "type": "area",
//	  "lights": ["lights/01", "lights/03"],
//	  "lightGroups": {
//	    "speaker": ["lights/01"],
//	    "audience": ["lights/02", "lights/03"]
//	  }
//	}
package zone
