package database

import "gorm.io/gorm"

type Keys struct {
	Guid       string `gorm:"primaryKey" json:"guid"`
	BusinessId string `json:"business_id"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Timestamp  uint64 `json:"timestamp"`
}

type KeysView interface {
	QueryKeysByBusinessId(string, uint64, uint64) ([]Keys, error)
}

type KeysDB interface {
	KeysView

	StoreKeys([]Keys, uint64) error
}

type addressDB struct {
	gorm *gorm.DB
}

func NewKeysDB(db *gorm.DB) KeysDB {
	return &addressDB{gorm: db}
}

func (db *addressDB) StoreKeys(keyList []Keys, keyLength uint64) error {
	result := db.gorm.CreateInBatches(&keyList, int(keyLength))
	return result.Error
}

func (db *addressDB) QueryKeysByBusinessId(businessId string, page, pageSize uint64) ([]Keys, error) {
	panic("implement me")
}
