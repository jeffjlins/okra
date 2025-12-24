package domain

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gookit/validate"
)

// TODO: Smallest Uom is the pivot
//
//	Use the smallest uom in each group as the pivot so we don't need an explicit pivot field stored in the data.
//	Convert all NI data to smallest uom of the group so it is all aligned on a pivot.
//	Make all conversions go through the smallest uom of the group.
//	The pivot field can still be present but can be computed after the fact for a set transform and set validation.
//
// TODO: Need to add a validation that either short name or full name is present.
//
//	It should also align with name default type.
//	Also check that either both singular and plural are there or that they aren't.
//
// TODO: Need to think through group transitions more.
//
//	You might say 1 tsp, 2 tsp, 1 tbsp, 4 tsp, 5 tsp, 2 tbsp, 2.25 tbsp, etc.
//	You might have multiple lists of pivots like high pivots and low pivots so that if the lower uom has a higher matching snap then it might override the higher uoms snap and use the lower uom.
//	You might have overlapping group ranges and then do this negotiation.
//
// TODO: Just change the names in the csv to match this, remove differentiation field and pivot field
// TODO: Match names must be greater than 0 length strings
// TODO: Set validation to ensure no conflicting match names
type PreciseFloat32 float32

func (f PreciseFloat32) MarshalJSON() ([]byte, error) {
	val := float64(f)
	rounded := math.Round(val*1e6) / 1e6
	s := strconv.FormatFloat(rounded, 'f', 6, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return json.Marshal(s)
}

func (f *PreciseFloat32) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		var n float64
		if err2 := json.Unmarshal(data, &n); err2 != nil {
			return fmt.Errorf("cannot unmarshal %s into PreciseFloat32", string(data))
		}
		*f = PreciseFloat32(n)
		return nil
	}
	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}
	*f = PreciseFloat32(val)
	return nil
}

func (f PreciseFloat32) Float32() float32 {
	return float32(f)
}

// TODO: Set validation to ensure unique Label
type BaseUom struct {
	Label string `json:"label" validate:"required"`

	Enabled     bool             `json:"enabled" validate:"required"` // does not exist in csv
	MeasureType UomMeasureType   `json:"measure_type" validate:"required"`
	Group       *string          `json:"group,omitempty" validate:"-"` // Even when filling this in, it should be a reference because there will be one group instance per group
	GroupMin    *PreciseFloat32  `json:"group_min,omitempty" validate:"-"`
	GroupMax    *PreciseFloat32  `json:"group_max,omitempty" validate:"-"`
	SnapAmount  []PreciseFloat32 `json:"snap_amount" validate:"required"`    // This is to ensure it doesn't do values in between these. (e.g. [0.25, 0.001], [1], etc)
	SnapSelect  *PreciseFloat32  `json:"snap_select,omitempty" validate:"-"` // The snap to total ratio must be at least this much to use the snap, otherwise it checks the next highest snap. (e.g. 0.1, etc)

	MatchNamesRecipe    []string `json:"match_names_recipe" validate:"-"`     // was "recipe_match_names"
	MatchNamesFoodLabel []string `json:"match_names_food_label" validate:"-"` // was "food_label_match_names"

	PrintedNameDefaultType   UomPrintedNameType `json:"default_name_type" validate:"required"`
	PrintedNameShortSingular *string            `json:"short_name_singular,omitempty" validate:"-"`
	PrintedNameShortPlural   *string            `json:"short_name_plural,omitempty" validate:"-"`
	PrintedNameFullSingular  *string            `json:"full_name_singular,omitempty" validate:"-"`
	PrintedNameFullPlural    *string            `json:"full_name_plural,omitempty" validate:"-"`

	AdditionalInfo *UomAdditionalInfo `json:"info,omitempty" validate:"-"`
}

type Uom struct {
	BaseUom
	Id string `json:"id" validate:"required"`
}

type UomAdditionalInfo struct {
	Systems   []string `json:"systems" validate:"-"`
	NameGroup *string  `json:"name_group,omitempty" validate:"-"`
}

type UomMeasureType = string

const (
	VOL    UomMeasureType = "volume"
	WEIGHT UomMeasureType = "weight"
	ITEM   UomMeasureType = "item"
	PKG    UomMeasureType = "package"
)

type UomPrintedNameType = string

const (
	SHORT UomPrintedNameType = "short"
	FULL  UomPrintedNameType = "full"
)

type UomOption func(*BaseUom)

func NewUom(
	label string,
	measureType UomMeasureType,
	snapAmount []PreciseFloat32,
	printedNameDefaultType UomPrintedNameType,
	opts ...UomOption,
) (*BaseUom, error) {
	uom := &BaseUom{
		Label:                  label,
		Enabled:                true,
		MeasureType:            measureType,
		SnapAmount:             snapAmount,
		PrintedNameDefaultType: printedNameDefaultType,
		MatchNamesRecipe:       []string{},
		MatchNamesFoodLabel:    []string{},
	}
	// Apply optional configurations
	for _, opt := range opts {
		opt(uom)
	}
	if err := uom.Validate(); err != nil {
		return nil, err
	}
	return uom, nil
}

func Create(base *BaseUom) (*Uom, error) {
	if err := base.Validate(); err != nil {
		return nil, err
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	uom := &Uom{
		BaseUom: *base,
		Id:      id.String(),
	}
	if err := uom.Validate(); err != nil {
		return nil, err
	}
	return uom, nil
}

func WithGroup(group string, groupMin, groupMax *PreciseFloat32) UomOption {
	return func(u *BaseUom) {
		u.Group = &group
		u.GroupMin = groupMin
		u.GroupMax = groupMax
	}
}

func WithSnapSelect(snapSelect PreciseFloat32) UomOption {
	return func(u *BaseUom) {
		u.SnapSelect = &snapSelect
	}
}

func WithMatchNamesRecipe(names []string) UomOption {
	return func(u *BaseUom) {
		if names != nil {
			u.MatchNamesRecipe = names
		}
	}
}

func WithMatchNamesFoodLabel(names []string) UomOption {
	return func(u *BaseUom) {
		if names != nil {
			u.MatchNamesFoodLabel = names
		}
	}
}

func WithPrintedNames(shortSingular, shortPlural, fullSingular, fullPlural *string) UomOption {
	return func(u *BaseUom) {
		u.PrintedNameShortSingular = shortSingular
		u.PrintedNameShortPlural = shortPlural
		u.PrintedNameFullSingular = fullSingular
		u.PrintedNameFullPlural = fullPlural
	}
}

func WithAdditionalInfo(info *UomAdditionalInfo) UomOption {
	return func(u *BaseUom) {
		u.AdditionalInfo = info
	}
}

func WithEnabled(enabled bool) UomOption {
	return func(u *BaseUom) {
		u.Enabled = enabled
	}
}

func (u *BaseUom) Validate() error {
	v := validate.Struct(u)
	if !v.Validate() {
		return v.Errors
	}
	return nil
}

func (u *Uom) Validate() error {
	v := validate.Struct(u)
	if !v.Validate() {
		return v.Errors
	}
	return nil
}

func (u *Uom) String() string {
	return print(reflect.ValueOf(u))
}
