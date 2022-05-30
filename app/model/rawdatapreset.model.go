package model

import (
	"github.com/google/uuid"
)

const TableNameRawdatapreset = "rawdatapreset"

// TODO: Delete if there are no sensors
type RawDataPreset struct {
	Base
	Name      string      `gorm:"column:name;not null;uniqueIndex:unique_rawdatapreset_name_in_thing" json:"name"`
	ThingId   uuid.UUID   `gorm:"type:uuid;column:thing_id;not null;uniqueIndex:unique_rawdatapreset_name_in_thing" json:"thingId"`
	Thing     Thing       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	SensorIds []uuid.UUID `gorm:"-" json:"sensorIds"`
}

func (*RawDataPreset) TableName() string {
	return TableNameRawdatapreset
}
