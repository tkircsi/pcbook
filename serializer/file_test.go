package serializer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tkircsi/pcbook/pb"
	"github.com/tkircsi/pcbook/samples"
	"github.com/tkircsi/pcbook/serializer"
	"google.golang.org/protobuf/proto"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop1 := samples.NewLaptop()
	err := serializer.WriteToFile(laptop1, binaryFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err = serializer.ReadFromFile(binaryFile, laptop2)
	require.NoError(t, err)
	require.True(t, proto.Equal(laptop1, laptop2))

	err = serializer.WriteToJSON(laptop1, jsonFile)
	require.NoError(t, err)
}
