package structure

import (
	"log"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

//===============================================================================
// Copy with UUID

// Custom type converter from string to uuid.UUID
func stringToUUID(src string) (uuid.UUID, error) {
	return uuid.Parse(src)
}

func stringToPUUID(src string) (*uuid.UUID, error) {
	if src == "" {
		return nil, nil
	}
	parsedUUID, err := uuid.Parse(src)
	if err != nil {
		return nil, err
	}
	return &parsedUUID, nil
}

// Here initially it was the other way around.
func CopyWithUUID(src, dest interface{}) {
	// Use copier with custom transformers
	err := copier.CopyWithOption(dest, src, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    false,
		Converters: []copier.TypeConverter{
			{
				SrcType: string(""),
				DstType: uuid.UUID{},
				Fn: func(src interface{}) (interface{}, error) {
					return stringToUUID(src.(string))
				},
			},
			{
				SrcType: string(""),
				DstType: (*uuid.UUID)(nil),
				Fn: func(src interface{}) (interface{}, error) {
					return stringToPUUID(src.(string))
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to copy: %v\n", err)
	}
}

//==============================================================================

// Here initially it was the other way around.
func Copy(s, ts interface{}) {
	err := copier.Copy(ts, s)
	if err != nil {
		panic(err)
	}
}
