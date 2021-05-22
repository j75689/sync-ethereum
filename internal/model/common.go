package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type GormBigInt big.Int

func (bi GormBigInt) Scan(value interface{}) error {
	n, ok := value.(int64)
	if !ok {
		return errors.New(fmt.Sprint("failed convert value to int64", value))
	}

	bi = GormBigInt(*big.NewInt(n))
	return nil
}

func (bi GormBigInt) Value() (driver.Value, error) {
	bigI := big.Int(bi)
	return bigI.Int64(), nil
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
