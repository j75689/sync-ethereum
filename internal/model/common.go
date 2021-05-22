package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type GormBigInt big.Int

func (bi GormBigInt) Int64() int64 {
	bigI := big.Int(bi)
	return bigI.Int64()
}

func (bi GormBigInt) Bytes() []byte {
	bigI := big.Int(bi)
	return bigI.Bytes()
}

func (bi *GormBigInt) Scan(val interface{}) error {
	if val == nil {
		return nil
	}
	var data string
	switch v := val.(type) {
	case []byte:
		data = string(v)
	case string:
		data = v
	case int64:
		bigI := new(big.Int).SetInt64(v)
		*bi = GormBigInt(*bigI)
		return nil
	default:
		return fmt.Errorf("bigint: can't convert %s type to *big.Int", reflect.TypeOf(val).Kind())
	}

	bigI, ok := new(big.Int).SetString(data, 10)
	if !ok {
		return fmt.Errorf("bigint can't convert %s to *big.Int", data)
	}
	*bi = GormBigInt(*bigI)
	return nil
}

func (bi GormBigInt) Value() (driver.Value, error) {
	bigI := big.Int(bi)
	return bigI.String(), nil
}

func (bi GormBigInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(bi.Int64())
}

func (bi GormBigInt) UnmarshalJSON(data []byte) error {
	var i64 int64
	err := json.Unmarshal(data, &i64)
	if err != nil {
		return err
	}
	bigI := big.NewInt(i64)
	bi = GormBigInt(*bigI)
	return nil
}

type BlockData []byte

func (d BlockData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("failed convert value to []byte", value))
	}
	d = BlockData(b)
	return nil
}

func (d BlockData) Value() (driver.Value, error) {
	return []byte(d), nil
}

func (d BlockData) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(common.BytesToHash(d).Hex())), nil
}

func (d BlockData) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	d = common.FromHex(s)
	return nil
}

type Pagination struct {
	Page    int64 `json:"page,omitempty"`
	PerPage int64 `json:"per_page,omitempty"`
}

func (p Pagination) LimitAndOffset(db *gorm.DB) *gorm.DB {
	if p.PerPage != 0 || p.Offset() != 0 {
		db = db.Limit(int(p.PerPage)).Offset(int(p.Offset()))
	}
	return db
}

func (p Pagination) Offset() int64 {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PerPage
}

type SortOrder string

func (field SortOrder) String() string {
	return string(field)
}

const (
	SortASC  SortOrder = "ASC"
	SortDESC SortOrder = "DESC"
)

type SortField struct {
	Field string    `json:"sort_field,omitempty"`
	Order SortOrder `json:"sort_order,omitempty"`
}

type Sorting []SortField

func (s Sorting) Sort(db *gorm.DB) *gorm.DB {
	sortfield := []string{}
	for _, sort := range s {
		if len(sort.Field) != 0 && len(sort.Order) != 0 {
			sortfield = append(sortfield, fmt.Sprintf("%s %s", sort.Field, sort.Order))
		}
	}

	if len(sortfield) > 0 {
		db = db.Order(strings.Join(sortfield, ","))
	}

	return db
}
