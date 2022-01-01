package cyclonedx

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

var bomLinkRegex = regexp.MustCompile(`^urn:cdx:(?P<serialNumber>[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})\/(?P<version>[0-9]+)(?:#(?P<bomRef>[0-9a-zA-Z\-._~%!$&'()*+,;=:@\/?]+))?$`)

// BOMLink TODO
type BOMLink struct {
	SerialNumber uuid.UUID // Serial number of the linked BOM
	Version      int       // Version of the linked BOM
	Reference    string    // Reference of the linked element
}

// NewBOMLink TODO
func NewBOMLink(bom *BOM, elem interface{}) (*BOMLink, error) {
	if bom == nil {
		return nil, fmt.Errorf("bom is nil")
	}
	if bom.SerialNumber == "" {
		return nil, fmt.Errorf("missing serial number")
	}
	if bom.Version < 1 {
		return nil, fmt.Errorf("versions below 1 are not allowed")
	}

	serial, err := uuid.Parse(bom.SerialNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid serial number: %w", err)
	}

	if elem == nil {
		return &BOMLink{
			SerialNumber: serial,
			Version:      bom.Version,
		}, nil
	}

	var bomRef string
	switch e := elem.(type) {
	case Component:
		bomRef = e.BOMRef
	case *Component:
		bomRef = e.BOMRef
	case Service:
		bomRef = e.BOMRef
	case *Service:
		bomRef = e.BOMRef
	default:
		return nil, fmt.Errorf("element of type %T is not referenceable", e)
	}
	if bomRef == "" {
		return nil, fmt.Errorf("element has no bom reference")
	}

	return &BOMLink{
		SerialNumber: serial,
		Version:      bom.Version,
		Reference:    bomRef,
	}, nil
}

// String TODO
func (b BOMLink) String() string {
	if b.Reference == "" {
		return fmt.Sprintf("urn:cdx:%s/%d", b.SerialNumber, b.Version)
	}

	return fmt.Sprintf("urn:cdx:%s/%d#%s", b.SerialNumber, b.Version, url.QueryEscape(b.Reference))
}

// IsBOMLink TODO
func IsBOMLink(s string) bool {
	return bomLinkRegex.MatchString(s)
}

// ParseBOMLink TODO
func ParseBOMLink(s string) (*BOMLink, error) {
	matches := bomLinkRegex.FindStringSubmatch(s)
	if len(matches) < 3 || len(matches) > 4 {
		return nil, fmt.Errorf("")
	}

	serial, err := uuid.Parse(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid serial number: %w", err)
	}
	version, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid version: %w", err)
	}

	if len(matches) == 4 {
		bomRef, err := url.QueryUnescape(matches[3])
		if err != nil {
			return nil, fmt.Errorf("invalid reference: %w", err)
		}

		return &BOMLink{
			SerialNumber: serial,
			Version:      version,
			Reference:    bomRef,
		}, nil
	}

	return &BOMLink{
		SerialNumber: serial,
		Version:      version,
	}, nil
}
