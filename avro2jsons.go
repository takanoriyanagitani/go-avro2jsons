package avro2jsons

type AvroRow map[string]any

type AvroBlobSizeMax int

const (
	AvroBlobSizeMaxDefault AvroBlobSizeMax = 1048576
)

type AvroConfig struct {
	blobSizeMax AvroBlobSizeMax
}

func (c AvroConfig) WithBlobSizeMax(m AvroBlobSizeMax) AvroConfig {
	c.blobSizeMax = m
	return c
}

func (c AvroConfig) BlobSizeMax() AvroBlobSizeMax { return c.blobSizeMax }

var AvroConfigDefault AvroConfig = AvroConfig{}.
	WithBlobSizeMax(AvroBlobSizeMaxDefault)
