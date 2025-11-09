package geoip

import (
	"net"

	_ "embed"

	"github.com/oschwald/geoip2-golang"
)

type Country struct {
	ISOCode string
	EU      bool
}

//go:embed country.mmdb
var mmdb []byte

type GeoIP struct{ db *geoip2.Reader }

func New() (*GeoIP, error) {
	db, err := geoip2.FromBytes(mmdb)
	if err != nil {
		return nil, err
	}
	return &GeoIP{db: db}, nil
}

func (s *GeoIP) GetCountry(ip net.IP) (Country, error) {
	record, err := s.db.Country(ip)
	if err != nil {
		return Country{}, err
	}
	return Country{
		ISOCode: record.Country.IsoCode,
		EU:      record.Country.IsInEuropeanUnion,
	}, nil
}
