// Package pubcache implements a cache of selected publications from a remote Smart Core publication server.
// This can be used to ensure that publications (such as configuration data) remain available even when communication
// with the publication server is unavailable, and to allow multiple subsystems to share a publication without
// retrieving it multiple times.
//
// The server is the authoritative source of publication versions.
// The cache stores only the most recent version of each cached publication.
package pubcache
