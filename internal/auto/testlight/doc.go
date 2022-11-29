// Package testlight implements automatic test data retrieval from DALI Part 202 Emergency lighting fixtures.
// Test results are stored locally persistently, and cleared from the light as soon as they are read.
//
// The automation polls the lights in cycles. The interval between cycles is set in configuration.
// Every cycle, each light will be polled sequentially. The minimum interval between polling lights can be set in
// configuration to limit activity on the DALI bus.
//
// The database stores events noticed during the polling procedure:
//   - Status events: the set of failures reported by the light at the time it is polled. To reduce storage requirements,
//     status events are only saved if the status is different to last time.
//   - Function test pass event: when a light reports a function test has passed, this is recorded in the database.
//   - Duration test pass event: when a light reports a duration test has passed, the duration achieved is recorded.
//
// For test pass events, the data is automatically cleared from the light once saved in the database, to prevent
// duplicated events. Because test failures can't be cleared from the lights in this way, they are treated as part
// of the status, rather than a discrete event.
// The latest status value received from each light is also stored in the database, to allow easy "is this light OK?"
// queries, and to determine when the status has changed.
package testlight
