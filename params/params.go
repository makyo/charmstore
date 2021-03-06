// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// The params package holds types that are a part of the charm store's external
// contract - they will be marshalled (or unmarshalled) as JSON
// and delivered through the HTTP API.
package params

import (
	"encoding/json"
	"time"

	"github.com/juju/utils/debugstatus"
	"gopkg.in/juju/charm.v5-unstable"
)

const (
	// ContentHashHeader specifies the header attribute
	// that will hold the content hash for archive GET responses.
	ContentHashHeader = "Content-Sha384"

	// EntityIdHeader specifies the header attribute that will hold the
	// id of the entity for archive GET responses.
	EntityIdHeader = "Entity-Id"
)

// Special user/group names.
const (
	Everyone = "everyone"
	Admin    = "admin"
)

// MetaAnyResponse holds the result of a meta/any request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetaany
type MetaAnyResponse struct {
	Id   *charm.Reference
	Meta map[string]interface{} `json:",omitempty"`
}

// ArchiveUploadResponse holds the result of a post or a put to /id/archive.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#post-idarchive
type ArchiveUploadResponse struct {
	Id *charm.Reference
	PromulgatedId *charm.Reference		`json:",omitempty"`
}

// ExpandedId holds a charm or bundle fully qualified id.
// A slice of ExpandedId is used as response for
// id/expand-id GET requests.
type ExpandedId struct {
	Id string
}

// ArchiveSizeResponse holds the result of an
// id/meta/archive-size GET request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetaarchive-size
type ArchiveSizeResponse struct {
	Size int64
}

// HashResponse holds the result of id/meta/hash and id/meta/hash256 GET
// requests.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetahash
// and https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetahash256
type HashResponse struct {
	Sum string
}

// ManifestFile holds information about a charm or bundle file.
// A slice of ManifestFile is used as response for
// id/meta/manifest GET requests.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetamanifest
type ManifestFile struct {
	Name string
	Size int64
}

// ArchiveUploadTimeResponse holds the result of an id/meta/archive-upload-time
// GET request. See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetaarchive-upload-time
type ArchiveUploadTimeResponse struct {
	UploadTime time.Time
}

// RelatedResponse holds the result of an id/meta/charm-related GET request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetacharm-related
type RelatedResponse struct {
	// Requires holds an entry for each interface provided by
	// the charm, containing all charms that require that interface.
	Requires map[string][]MetaAnyResponse `json:",omitempty"`

	// Provides holds an entry for each interface required by the
	// the charm, containing all charms that provide that interface.
	Provides map[string][]MetaAnyResponse `json:",omitempty"`
}

// RevisionInfoResponse holds the result of an id/meta/revision-info GET
// request. See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetarevision-info
type RevisionInfoResponse struct {
	Revisions []*charm.Reference
}

// BundleCount holds the result of an id/meta/bundle-unit-count
// or bundle-machine-count GET request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetabundle-unit-count
// and https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetabundle-machine-count
type BundleCount struct {
	Count int
}

// TagsResponse holds the result of an id/meta/tags GET request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetatags
type TagsResponse struct {
	Tags []string
}

// Published holds the result of a changes/published GET request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-changespublished
type Published struct {
	Id          *charm.Reference
	PublishTime time.Time
}

// DebugStatus holds the result of the status checks.
// This is defined for backward compatibility: new clients should use
// debugstatus.CheckResult directly.
type DebugStatus debugstatus.CheckResult

// SearchResult holds a single result from a search operation.
type SearchResult struct {
	Id *charm.Reference
	// Meta holds at most one entry for each meta value
	// specified in the include flags, holding the
	// data that would be returned by reading /meta/meta?id=id.
	// Metadata not relevant to a particular result will not
	// be included.
	Meta map[string]interface{} `json:",omitempty"`
}

// SearchResponse holds the response from a search operation.
type SearchResponse struct {
	SearchTime time.Duration
	Total      int
	Results    []SearchResult
}

// IdUserResponse holds the result of an id/meta/id-user GET request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetaid-user
type IdUserResponse struct {
	User string
}

// IdSeriesResponse holds the result of an id/meta/id-series GET request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetaid-series
type IdSeriesResponse struct {
	Series string
}

// IdNameResponse holds the result of an id/meta/id-name GET request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetaid-name
type IdNameResponse struct {
	Name string
}

// IdRevisionResponse holds the result of an id/meta/id-revision GET request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetaid-revision
type IdRevisionResponse struct {
	Revision int
}

// IdResponse holds the result of an id/meta/id GET request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetaid
type IdResponse struct {
	Id       *charm.Reference
	User     string `json:",omitempty"`
	Series   string `json:",omitempty"`
	Name     string
	Revision int
}

// PermResponse holds the result of an id/meta/perm GET request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetaperm
type PermResponse struct {
	Read  []string
	Write []string
}

// PermRequest holds the request of an id/meta/perm PUT request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#put-idmetaperm
type PermRequest struct {
	Read  []string
	Write []string
}

// PromulgatedResponse holds the result of an id/meta/promulgated GET request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-idmetapromulgated
type PromulgatedResponse struct {
	Promulgated bool
}

// PromulgateRequest holds the request of an id/promulgate PUT request.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#put-idpromulgate
type PromulgateRequest struct {
	Promulgated bool
}

const (
	// BzrDigestKey is the extra-info key used to store the Bazaar digest
	BzrDigestKey = "bzr-digest"

	// LegacyDownloadStats is the extra-info key used to store the legacy
	// download counts, and to retrieve them when
	// charmstore.LegacyDownloadCountsEnabled is set to true.
	// TODO (frankban): remove this constant when removing the legacy counts
	// logic.
	LegacyDownloadStats = "legacy-download-stats"
)

// Log holds the representation of a log message.
// This is used by clients to store log events in the charm store.
type Log struct {
	// Data holds the log message as a JSON-encoded value.
	Data *json.RawMessage

	// Level holds the log level as a string.
	Level LogLevel

	// Type holds the log type as a string.
	Type LogType

	// URLs holds a slice of entity URLs associated with the log message.
	URLs []*charm.Reference `json:",omitempty"`
}

// LogResponse represents a single log message and is used in the responses
// to /log GET requests.
// See https://github.com/juju/charmstore/blob/v4/docs/API.md#get-log
type LogResponse struct {
	// Data holds the log message as a JSON-encoded value.
	Data json.RawMessage

	// Level holds the log level as a string.
	Level LogLevel

	// Type holds the log type as a string.
	Type LogType

	// URLs holds a slice of entity URLs associated with the log message.
	URLs []*charm.Reference `json:",omitempty"`

	// Time holds the time of the log.
	Time time.Time
}

// LogLevel defines log levels (e.g. "info" or "error") to be used in log
// requests and responses.
type LogLevel string

const (
	InfoLevel    LogLevel = "info"
	WarningLevel LogLevel = "warning"
	ErrorLevel   LogLevel = "error"
)

// LogType defines log types (e.g. "ingestion") to be used in log requests and
// responses.
type LogType string

const (
	IngestionType        LogType = "ingestion"
	LegacyStatisticsType LogType = "legacyStatistics"

	IngestionStart    = "ingestion started"
	IngestionComplete = "ingestion completed"

	LegacyStatisticsImportStart    = "legacy statistics import started"
	LegacyStatisticsImportComplete = "legacy statistics import completed"
)
