// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package v4

import (
	"encoding/base64"
	"net/http"
	"strings"

	"gopkg.in/errgo.v1"
	"gopkg.in/macaroon-bakery.v0/bakery"
	"gopkg.in/macaroon-bakery.v0/bakery/checkers"
	"gopkg.in/macaroon-bakery.v0/httpbakery"
	"gopkg.in/macaroon.v1"

	"gopkg.in/juju/charmstore.v4/internal/router"
	"gopkg.in/juju/charmstore.v4/params"
)

const (
	basicRealm        = "CharmStore4"
	promulgatorsGroup = "promulgators"
)

// authorize checks that the current user is authorized based on the provided
// ACL. If an authenticated user is required, authorize tries to retrieve the
// current user in the following ways:
// - by checking that the request's headers HTTP basic auth credentials match
//   the superuser credentials stored in the API handler;
// - by checking that there is a valid macaroon in the request's cookies.
// A params.ErrUnauthorized error is returned if superuser credentials fail;
// otherwise a macaroon is minted and a httpbakery discharge-required
// error is returned holding the macaroon.
func (h *Handler) authorize(req *http.Request, acl []string) error {
	logger.Infof(
		"authorize, bakery %p, auth location %q, acl %q, path: %q, method: %q",
		h.store.Bakery,
		h.config.IdentityLocation,
		acl,
		req.URL.Path,
		req.Method)

	// No need to authenticate if the ACL is open to everyone.
	for _, name := range acl {
		if name == params.Everyone {
			return nil
		}
	}

	auth, verr := h.checkRequest(req)
	if verr == nil {
		logger.Infof("authenticated with auth: %#v", auth)
		if err := h.checkACLMembership(auth, acl); err != nil {
			return errgo.WithCausef(err, params.ErrUnauthorized, "")
		}
		return nil
	}
	if _, ok := errgo.Cause(verr).(*bakery.VerificationError); !ok {
		return errgo.Mask(verr, errgo.Is(params.ErrUnauthorized))
	}

	// Macaroon verification failed: mint a new macaroon.
	m, err := h.newMacaroon()
	if err != nil {
		return errgo.Notef(err, "cannot mint macaroon")
	}
	// Request that this macaroon be supplied for all requests
	// to the whole handler.
	// TODO use a relative URL here: router.RelativeURLPath(req.RequestURI, "/")
	cookiePath := "/"
	return httpbakery.NewDischargeRequiredError(m, cookiePath, verr)
}

// checkRequest checks for any authorization tokens in the request and returns any
// found as an authorization. If no suitable credentials are found, or an error occurs,
// then a zero valued authorization is returned.
func (h *Handler) checkRequest(req *http.Request) (authorization, error) {
	user, passwd, err := parseCredentials(req)
	if err == nil {
		if user != h.config.AuthUsername || passwd != h.config.AuthPassword {
			return authorization{}, errgo.WithCausef(nil, params.ErrUnauthorized, "invalid user name or password")
		}
		return authorization{Admin: true}, nil
	}
	if errgo.Cause(err) != errNoCreds || h.store.Bakery == nil || h.config.IdentityLocation == "" {
		return authorization{}, errgo.WithCausef(err, params.ErrUnauthorized, "authentication failed")
	}
	attrMap, err := httpbakery.CheckRequest(h.store.Bakery, req, nil, checkers.New())
	if err != nil {
		return authorization{}, errgo.Mask(err, errgo.Any)
	}
	return authorization{
		Admin:    false,
		Username: attrMap[usernameAttr],
		Groups:   strings.Fields(attrMap[groupsAttr]),
	}, nil
}

func (h *Handler) authorizeEntity(id *router.ResolvedURL, req *http.Request) error {
	baseEntity, err := h.store.FindBaseEntity(&id.URL, "acls")
	if err != nil {
		if errgo.Cause(err) == params.ErrNotFound {
			return errgo.WithCausef(nil, params.ErrNotFound, "entity %q not found", id)
		}
		return errgo.Notef(err, "cannot retrieve entity %q for authorization", id)
	}
	return h.authorizeWithPerms(req, baseEntity.ACLs.Read, baseEntity.ACLs.Write)
}

func (h *Handler) authorizeWithPerms(req *http.Request, read, write []string) error {
	var acl []string
	switch req.Method {
	case "DELETE", "PATCH", "POST", "PUT":
		acl = write
	default:
		acl = read
	}
	return h.authorize(req, acl)
}

const (
	usernameAttr = "username"
	groupsAttr   = "groups"
)

// authorization conatains authorization information extracted from an HTTP request.
// The zero value for a authorization contains no privileges.
type authorization struct {
	Admin    bool
	Username string
	Groups   []string
}

func (h *Handler) checkACLMembership(auth authorization, acl []string) error {
	if auth.Admin {
		return nil
	}
	if auth.Username == "" {
		return errgo.New("no username declared")
	}
	members := map[string]bool{
		params.Everyone: true,
		auth.Username:   true,
	}
	for _, name := range auth.Groups {
		members[name] = true
	}
	for _, name := range acl {
		if members[name] {
			return nil
		}
	}
	return errgo.Newf("access denied for user %q", auth.Username)
}

func (h *Handler) newMacaroon() (*macaroon.Macaroon, error) {
	// TODO generate different caveats depending on the requested operation
	// and whether there's a charm id or not.
	// Mint an appropriate macaroon and send it back to the client.
	return h.store.Bakery.NewMacaroon("", nil, []checkers.Caveat{checkers.NeedDeclaredCaveat(checkers.Caveat{
		Location:  h.config.IdentityLocation,
		Condition: "is-authenticated-user",
	}, usernameAttr, groupsAttr)})
}

var errNoCreds = errgo.New("missing HTTP auth header")

// parseCredentials parses the given request and returns the HTTP basic auth
// credentials included in its header.
func parseCredentials(req *http.Request) (username, password string, err error) {
	auth := req.Header.Get("Authorization")
	if auth == "" {
		return "", "", errNoCreds
	}
	parts := strings.Fields(auth)
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", errgo.New("invalid HTTP auth header")
	}
	// Challenge is a base64-encoded "tag:pass" string.
	// See RFC 2617, Section 2.
	challenge, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", errgo.New("invalid HTTP auth encoding")
	}
	tokens := strings.SplitN(string(challenge), ":", 2)
	if len(tokens) != 2 {
		return "", "", errgo.New("invalid HTTP auth contents")
	}
	return tokens[0], tokens[1], nil
}
