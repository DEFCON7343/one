package goca

import "errors"

// MarketPlace represents an OpenNebula MarketPlace
type MarketPlace struct {
	XMLResource
	ID   uint
	Name string
}

// MarketPlacePool represents an OpenNebula MarketPlacePool
type MarketPlacePool struct {
	XMLResource
}

// NewMarketPlacePool returns a marketplace pool. A connection to OpenNebula is
// performed.
func NewMarketPlacePool(args ...int) (*MarketPlacePool, error) {
	var who, start, end int

	switch len(args) {
	case 0:
		who = PoolWhoMine
		start = -1
		end = -1
	case 1:
		who = args[0]
		start = -1
		end = -1
	case 3:
		who = args[0]
		start = args[1]
		end = args[2]
	default:
		return nil, errors.New("Wrong number of arguments")
	}

    response, err := client.Call("one.marketpool.info", who, start, end)
	if err != nil {
		return nil, err
	}

	marketpool := &MarketPlacePool{XMLResource{body: response.Body()}}

	return marketpool, err
}

// NewMarketPlace finds a marketplace object by ID. No connection to OpenNebula.
func NewMarketPlace(id uint) *MarketPlace {
	return &MarketPlace{ID: id}
}

// NewMarketPlaceFromName finds a marketplace object by name. It connects to
// OpenNebula to retrieve the pool, but doesn't perform the Info() call to
// retrieve the attributes of the marketplace.
func NewMarketPlaceFromName(name string) (*MarketPlace, error) {
	marketPool, err := NewMarketPlacePool()
	if err != nil {
		return nil, err
	}

	id, err := marketPool.GetIDFromName(name, "/MARKETPLACE_POOL/MARKETPLACE")
	if err != nil {
		return nil, err
	}

	return NewMarketPlace(id), nil
}

// CreateMarketPlace allocates a new marketplace. It returns the new marketplace ID.
// * tpl: template of the marketplace
func CreateMarketPlace(tpl string) (uint, error) {
	response, err := client.Call("one.market.allocate", tpl)
	if err != nil {
		return 0, err
	}

	return uint(response.BodyInt()), nil
}

// Delete deletes the given marketplace from the pool.
func (market *MarketPlace) Delete() error {
	_, err := client.Call("one.market.delete", market.ID)
	return err
}

// Update replaces the marketplace template contents.
// * tpl: The new template contents. Syntax can be the usual attribute=value or XML.
// * appendTemplate: Update type: 0: Replace the whole template. 1: Merge new template with the existing one.
func (market *MarketPlace) Update(tpl string, appendTemplate int) error {
	_, err := client.Call("one.market.update", market.ID, tpl, appendTemplate)
	return err
}

// Chmod changes the permission bits of a marketplace
// * uu: USER USE bit. If set to -1, it will not change.
// * um: USER MANAGE bit. If set to -1, it will not change.
// * ua: USER ADMIN bit. If set to -1, it will not change.
// * gu: GROUP USE bit. If set to -1, it will not change.
// * gm: GROUP MANAGE bit. If set to -1, it will not change.
// * ga: GROUP ADMIN bit. If set to -1, it will not change.
// * ou: OTHER USE bit. If set to -1, it will not change.
// * om: OTHER MANAGE bit. If set to -1, it will not change.
// * oa: OTHER ADMIN bit. If set to -1, it will not change.
func (market *MarketPlace) Chmod(uu, um, ua, gu, gm, ga, ou, om, oa int) error {
	_, err := client.Call("one.market.chmod", market.ID, uu, um, ua, gu, gm, ga, ou, om, oa)
	return err
}

// Chown changes the ownership of a marketplace.
// * userID: The User ID of the new owner. If set to -1, it will not change.
// * groupID: The Group ID of the new group. If set to -1, it will not change.
func (market *MarketPlace) Chown(userID, groupID int) error {
	_, err := client.Call("one.market.chown", market.ID, userID, groupID)
	return err
}

// Rename renames a marketplace.
// * newName: The new name.
func (market *MarketPlace) Rename(newName string) error {
	_, err := client.Call("one.market.rename", market.ID, newName)
	return err
}

// Info retrieves information for the marketplace.
func (market *MarketPlace) Info() error {
	response, err := client.Call("one.market.info", market.ID)
	if err != nil {
		return err
	}
	market.body = response.Body()
	return nil
}
