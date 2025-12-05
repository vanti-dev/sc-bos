package goproto

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// analyzeProtoFiles returns the required generators for each proto file in protoDir.
// The keys of the returned map are relative paths from protoDir for each proto file.
// The values are the combined Generator flags needed for that file.
// analyzeProtoFiles recursively walks protoDir to find all .proto files.
func analyzeProtoFiles(protoDir string) (map[string]Generator, error) {
	fileGenerators := make(map[string]Generator)

	err := filepath.Walk(protoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".proto") {
			return nil
		}

		relPath, err := filepath.Rel(protoDir, path)
		if err != nil {
			return fmt.Errorf("getting relative path: %w", err)
		}

		gen, err := determineGenerators(path)
		if err != nil {
			return fmt.Errorf("analyzing %s: %w", relPath, err)
		}

		fileGenerators[relPath] = gen
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileGenerators, nil
}

// determineGenerators analyzes a proto file to determine which generators it needs.
func determineGenerators(filePath string) (Generator, error) {
	// Get the proto directory (parent of the file)
	protoDir := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)

	fileDesc, err := parseProtoFile(protoDir, fileName)
	if err != nil {
		return 0, fmt.Errorf("parsing proto file: %w", err)
	}

	return determineGeneratorsFromDescriptor(fileDesc), nil
}

// determineGeneratorsFromDescriptor analyzes a file descriptor to determine which generators it needs.
// This is separated from determineGenerators to allow testing without file I/O.
func determineGeneratorsFromDescriptor(fileDesc *descriptorpb.FileDescriptorProto) Generator {
	var gen Generator

	if len(fileDesc.GetService()) == 0 {
		// No services, no special generators needed
		return gen
	}

	// Files with services get wrappers
	gen |= GenWrapper
	if isRoutedAPI(fileDesc) {
		gen |= GenRouter
	}
	return gen
}

// parseProtoFile uses protomod protoc to parse a proto file and return its descriptor.
func parseProtoFile(protoDir, fileName string) (*descriptorpb.FileDescriptorProto, error) {
	// Use protomod protoc to handle proto.mod dependencies
	protomodPath, err := exec.LookPath("protomod")
	if err != nil {
		return nil, fmt.Errorf("protomod not found in PATH: %w", err)
	}

	cmd := exec.Command(protomodPath, "protoc", "--",
		"-I", protoDir,
		"--descriptor_set_out=/dev/stdout",
		"--include_imports",
		fileName,
	)
	cmd.Dir = protoDir

	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("protomod protoc failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("running protomod protoc: %w", err)
	}

	// Parse the descriptor set
	var fds descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(output, &fds); err != nil {
		return nil, fmt.Errorf("unmarshaling descriptor set: %w", err)
	}

	// Find the file descriptor for our file
	for _, fd := range fds.GetFile() {
		if fd.GetName() == fileName {
			return fd, nil
		}
	}

	return nil, fmt.Errorf("file descriptor not found for %s", fileName)
}

// isRoutedAPI determines if a proto file defines a routed API.
// A routed API has services where ALL request messages have a 'name' string field.
func isRoutedAPI(fileDesc *descriptorpb.FileDescriptorProto) bool {
	services := fileDesc.GetService()
	if len(services) == 0 {
		return false
	}

	// Build a map of message types in this file
	messages := make(map[string]*descriptorpb.DescriptorProto)
	pkg := fileDesc.GetPackage()
	for _, msg := range fileDesc.GetMessageType() {
		fullName := msg.GetName()
		if pkg != "" {
			fullName = pkg + "." + fullName
		}
		messages[fullName] = msg
		// Also register without package prefix for local references
		messages[msg.GetName()] = msg
	}

	// Check all methods in all services
	hasAnyRequestMessages := false
	for _, service := range services {
		for _, method := range service.GetMethod() {
			inputType := method.GetInputType()
			inputType = strings.TrimPrefix(inputType, ".")
			// Try with and without package prefix
			simpleName := filepath.Base(strings.ReplaceAll(inputType, ".", "/"))

			msg, ok := messages[inputType]
			if !ok {
				msg, ok = messages[simpleName]
			}

			if !ok {
				// Request message is defined in another file, skip this check
				continue
			}

			hasAnyRequestMessages = true

			if !hasNameField(msg) {
				return false
			}
		}
	}

	// If we didn't find any request messages defined in this file, it's not routed
	if !hasAnyRequestMessages {
		return false
	}

	return true
}

// hasNameField checks if a message has a 'name' string field.
func hasNameField(msg *descriptorpb.DescriptorProto) bool {
	for _, field := range msg.GetField() {
		if field.GetName() == "name" && field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_STRING {
			return true
		}
	}
	return false
}
