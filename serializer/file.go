package serializer

import (
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
)

func WriteToFile(msg proto.Message, filename string) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("cannot mashal proto message to binary: %v", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary to file: %v", err)
	}
	return nil
}

func ReadFromFile(fileName string, msg proto.Message) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("cannot read binry data from file: %v", err)
	}
	err = proto.Unmarshal(data, msg)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary data to proto message: %v", err)
	}
	return nil
}

func WriteToJSON(msg proto.Message, fileName string) error {
	data, err := ProtobufToJSON(msg)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to JSON: %v", err)
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write JSON to file: %v", err)
	}
	return nil
}
