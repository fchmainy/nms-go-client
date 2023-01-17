package nmsclient

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Compile-time interface checks
var _ NMSObjectManager = new(ObjectManager)

type NMSObjectManager interface {
	GetDNSView(name string) (*DNSView, error)
	CreateEnvironment(name string, description string, tags []string, eas EA) (*Network, error)
	GetEnvironment(name string, uuid string, eas EA) (*Network, error)
	UpdateEnvironment(name string, description string, tags []string, eas EA) (*Network, error)
	DeleteEnvironment(name string, uuid string, eas EA) (*Network, error)
}

type ObjectManager struct {
	connector NMSConnector
	cmpType   string
	tenantID  string
}

func NewObjectManager(connector IBConnector, cmpType string, tenantID string) IBObjectManager {
	objMgr := &ObjectManager{}

	objMgr.connector = connector
	objMgr.cmpType = cmpType
	objMgr.tenantID = tenantID

	return objMgr
}

// CreateMultiObject unmarshals the result into slice of maps
func (objMgr *ObjectManager) CreateMultiObject(req *MultiRequest) ([]map[string]interface{}, error) {

	conn := objMgr.connector.(*Connector)
	queryParams := NewQueryParams(false, nil)
	res, err := conn.makeRequest(CREATE, req, "", queryParams)

	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(res, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetUpgradeStatus returns the grid upgrade information
func (objMgr *ObjectManager) GetUpgradeStatus(statusType string) ([]UpgradeStatus, error) {
	var res []UpgradeStatus

	if statusType == "" {
		// TODO option may vary according to the WAPI version, need to
		// throw relevant  error.
		msg := fmt.Sprintf("Status type can not be nil")
		return res, errors.New(msg)
	}
	upgradestatus := NewUpgradeStatus(UpgradeStatus{})

	sf := map[string]string{
		"type": statusType,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(upgradestatus, "", queryParams, &res)

	return res, err
}

// GetAllMembers returns all members information
func (objMgr *ObjectManager) GetAllMembers() ([]Member, error) {
	var res []Member

	memberObj := NewMember(Member{})
	err := objMgr.connector.GetObject(
		memberObj, "", NewQueryParams(false, nil), &res)
	return res, err
}

// GetCapacityReport returns all capacity for members
func (objMgr *ObjectManager) GetCapacityReport(name string) ([]CapacityReport, error) {
	var res []CapacityReport

	capacityReport := NewCapcityReport(CapacityReport{})

	sf := map[string]string{
		"name": name,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(capacityReport, "", queryParams, &res)
	return res, err
}

// GetLicense returns the license details for member
func (objMgr *ObjectManager) GetLicense() ([]License, error) {
	var res []License

	licenseObj := NewLicense(License{})
	err := objMgr.connector.GetObject(
		licenseObj, "", NewQueryParams(false, nil), &res)
	return res, err
}

// GetLicense returns the license details for grid
func (objMgr *ObjectManager) GetGridLicense() ([]License, error) {
	var res []License

	licenseObj := NewGridLicense(License{})
	err := objMgr.connector.GetObject(
		licenseObj, "", NewQueryParams(false, nil), &res)
	return res, err
}

// GetGridInfo returns the details for grid
func (objMgr *ObjectManager) GetGridInfo() ([]Grid, error) {
	var res []Grid

	gridObj := NewGrid(Grid{})
	err := objMgr.connector.GetObject(
		gridObj, "", NewQueryParams(false, nil), &res)
	return res, err
}

// CreateZoneAuth creates zones and subs by passing fqdn
func (objMgr *ObjectManager) CreateZoneAuth(
	fqdn string,
	eas EA) (*ZoneAuth, error) {

	zoneAuth := NewZoneAuth(ZoneAuth{
		Fqdn: fqdn,
		Ea:   eas})

	ref, err := objMgr.connector.CreateObject(zoneAuth)
	zoneAuth.Ref = ref
	return zoneAuth, err
}

// Retreive a authortative zone by ref
func (objMgr *ObjectManager) GetZoneAuthByRef(ref string) (*ZoneAuth, error) {
	res := NewZoneAuth(ZoneAuth{})

	if ref == "" {
		return nil, fmt.Errorf("empty reference to an object is not allowed")
	}

	err := objMgr.connector.GetObject(
		res, ref, NewQueryParams(false, nil), res)
	return res, err
}

// DeleteZoneAuth deletes an auth zone
func (objMgr *ObjectManager) DeleteZoneAuth(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}

// GetZoneAuth returns the authoritatives zones
func (objMgr *ObjectManager) GetZoneAuth() ([]ZoneAuth, error) {
	var res []ZoneAuth

	zoneAuth := NewZoneAuth(ZoneAuth{})
	err := objMgr.connector.GetObject(
		zoneAuth, "", NewQueryParams(false, nil), &res)

	return res, err
}

// GetZoneDelegated returns the delegated zone
func (objMgr *ObjectManager) GetZoneDelegated(fqdn string) (*ZoneDelegated, error) {
	if len(fqdn) == 0 {
		return nil, nil
	}
	var res []ZoneDelegated

	zoneDelegated := NewZoneDelegated(ZoneDelegated{})

	sf := map[string]string{
		"fqdn": fqdn,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(zoneDelegated, "", queryParams, &res)

	if err != nil || res == nil || len(res) == 0 {
		return nil, err
	}

	return &res[0], nil
}

// CreateZoneDelegated creates delegated zone
func (objMgr *ObjectManager) CreateZoneDelegated(fqdn string, delegate_to []NameServer) (*ZoneDelegated, error) {
	zoneDelegated := NewZoneDelegated(ZoneDelegated{
		Fqdn:       fqdn,
		DelegateTo: delegate_to})

	ref, err := objMgr.connector.CreateObject(zoneDelegated)
	zoneDelegated.Ref = ref

	return zoneDelegated, err
}

// UpdateZoneDelegated updates delegated zone
func (objMgr *ObjectManager) UpdateZoneDelegated(ref string, delegate_to []NameServer) (*ZoneDelegated, error) {
	zoneDelegated := NewZoneDelegated(ZoneDelegated{
		Ref:        ref,
		DelegateTo: delegate_to})

	refResp, err := objMgr.connector.UpdateObject(zoneDelegated, ref)
	zoneDelegated.Ref = refResp
	return zoneDelegated, err
}

// DeleteZoneDelegated deletes delegated zone
func (objMgr *ObjectManager) DeleteZoneDelegated(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}
