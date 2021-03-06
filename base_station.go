package ipdb

import (
	"encoding/json"
	"reflect"
	"time"
)

type BaseStationInfo struct {
	CountryName	string	`json:"country_name"`
	RegionName string 	`json:"region_name"`
	CityName string 	`json:"city_name"`
	OwnerDomain string 	`json:"owner_domain"`
	IspDomain string 	`json:"isp_domain"`
	BaseStation string 	`json:"base_station"`
}

func (info BaseStationInfo) ToJson() []byte {
	all, err := json.Marshal(info)
	if err == nil {
		return all
	}

	return nil
}

type BaseStation struct {
	reader *Reader
}

func NewBaseStation(name string) (*BaseStation, error) {

	r, e := New(name, &BaseStationInfo{})
	if e != nil {
		return nil, e
	}

	return &BaseStation{
		reader: r,
	}, nil
}

func (db *BaseStation) Find(addr, language string) ([]string, error) {
	return db.reader.find1(addr, language)
}

func (db *BaseStation) FindMap(addr, language string) (map[string]string, error) {

	data, err := db.reader.find1(addr, language)
	if err != nil {
		return nil, err
	}
	info := make(map[string]string, len(db.reader.meta.Fields))
	for k, v := range data {
		info[db.reader.meta.Fields[k]] = v
	}

	return info, nil
}

func (db *BaseStation) FindInfo(addr, language string) (*BaseStationInfo, error) {

	data, err := db.reader.FindMap(addr, language)
	if err != nil {
		return nil, err
	}

	info := &BaseStationInfo{}

	for k, v := range data {
		sv := reflect.ValueOf(info).Elem()
		sfv := sv.FieldByName(db.reader.refType[k])

		if !sfv.IsValid() {
			continue
		}
		if !sfv.CanSet() {
			continue
		}

		sft := sfv.Type()
		fv := reflect.ValueOf(v)
		if sft == fv.Type() {
			sfv.Set(fv)
		}
	}

	return info, nil
}

func (db *BaseStation) IsIPv4() bool {
	return db.reader.IsIPv4Support()
}

func (db *BaseStation) IsIPv6() bool {
	return db.reader.IsIPv6Support()
}

func (db *BaseStation) Languages() []string {
	return db.reader.Languages()
}

func (db *BaseStation) Fields() []string {
	return db.reader.meta.Fields
}

func (db *BaseStation) BuildTime() time.Time {
	return db.reader.Build()
}