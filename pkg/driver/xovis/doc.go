// Package xovis contains a driver for camera-based occupancy sensors produced by Xovis.
// This package is intended to work with sensors running firmware version 5.x.
//
// Xovis exposes data organised under Logics.
// Zone-based (occupancy) logics are exposed by this driver using the OccupancySensor trait.
// Line-based (person count in/out) logics are exposed by this driver using the EnterLeaveSensor trait.
//
// A Xovis sensor can operate in single-sensor or multi-sensor modes. The data and operation available are similar
// between these two modes, but in the API they are accessed using separate endpoints.
//
// For one-shot data access, the driver uses the Live Data API endpoints to fetch the most up-to-date info on the
// appropriate logic. This involves us performing an HTTP(S) request to the sensor's built-in server.
// For streaming data (e.g. PullXxx Smart Core APIs) data flows in the other direction using a webhook -
// we add a route to the area controller's HTTP server for receiving the data. A Data Push agent must be manually
// configured on the sensor as follows:
//
// Data Push Agent:
//   - Type: Logics push
//   - Scheduler Type: Immediate
//   - Resolution: As desired, 1 minute makes most sense
//   - Format: JSON
//   - Time Format: RFC 3339
package xovis
