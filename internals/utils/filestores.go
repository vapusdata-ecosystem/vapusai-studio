package utils

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	os "os"
	"path"
	filepath "path/filepath"
	"strings"
	"sync"

	toml "github.com/pelletier/go-toml/v2"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
	"gopkg.in/yaml.v3"
)

// Function to write toml file
func WriteTomlFile(data interface{}, filename, path string) error {
	bytes, err := toml.Marshal(data)
	if err != nil {
		return err
	}

	file := filepath.Join(path, filename+DOT+DEFAULT_CONFIG_TYPE)
	dmlogger.CoreLogger.Info().Msgf("Writing to file: %v", file)
	err = os.WriteFile(file, bytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

// Function to read toml file
func ReadTomlFile(data interface{}, filename, path string) error {
	file := filepath.Join(path, filename+DOT+DEFAULT_CONFIG_TYPE)
	dmlogger.CoreLogger.Info().Msgf("Reading from file: %v", file)
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(bytes, data)
	if err != nil {
		return err
	}
	return nil
}

func GetConfFileType(fileName string) string {
	return strings.Replace(path.Ext(fileName), ".", "", -1)
}

func CreateFile(filename, filePath string, data any, base64encoded bool) error {
	if filePath == "" {
		curPath, err := os.Getwd()
		if err != nil {
			curPath = os.TempDir()
		}
		filePath = curPath
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if _, err := os.Create(filePath); err != nil {
			return err
		}
	}
	fType := GetConfFileType(filename)
	switch strings.ToLower(fType) {
	case "yaml":
		return WriteYAMLFile(data, filepath.Join(filePath, filename), base64encoded)
	case "json":
		return WriteJSONFile(data, filepath.Join(filePath, filename), base64encoded)
	case "toml":
		return WriteTOMLFile(data, filepath.Join(filePath, filename), base64encoded)
	default:
		return dmerrors.ErrInvalidArgs
	}
}

func ReadFile(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}
func GenericUnMarshaler(bytes []byte, result any, format string) error {
	var err error
	switch strings.ToLower(format) {
	case "yaml":
		err = yaml.Unmarshal(bytes, result)
		if err != nil {
			return err
		}
	case "json":
		err = json.Unmarshal(bytes, result)
		if err != nil {
			return err
		}
	case "toml":
		err = toml.Unmarshal(bytes, result)
		if err != nil {
			return err
		}
	case "csv":
		result, err = CSVBytesToArrayOfMap(bytes)
		if err != nil {
			return err
		}
		log.Println("Result ->>>>>>>>>>>>>>> ", result)
	default:
		return dmerrors.ErrInvalidArgs
	}
	return err
}

func DataFileLoader(bytes []byte, format string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	var err error
	switch strings.ToLower(format) {
	case "yaml":
		err = yaml.Unmarshal(bytes, &result)
		if err != nil {
			return nil, err
		}
	case "json":
		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return nil, err
		}
	case "toml":
		err = toml.Unmarshal(bytes, &result)
		if err != nil {
			return nil, err
		}
	case "csv":
		result, err = CSVBytesToArrayOfMap(bytes)
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, dmerrors.ErrInvalidArgs
	}
	return result, err
}

func GenericMarshaler(object any, format string) ([]byte, error) {
	format = strings.ToLower(format)
	switch strings.ToLower(format) {
	case "yaml":
		return yaml.Marshal(object)
	case "json":
		return json.Marshal(object)
	case "toml":
		return toml.Marshal(object)
	default:
		return nil, dmerrors.ErrInvalidArgs
	}
}

func MapArrayCSVMarshaler(data []map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer

	// Create a CSV writer
	writer := csv.NewWriter(&buf)

	// Write the header (keys from the first map)
	var header []string
	if len(data) > 0 {
		for key := range data[0] {
			header = append(header, key)
		}
		if err := writer.Write(header); err != nil {
			return nil, fmt.Errorf("failed to write header to CSV: %v", err)
		}
	}

	// Write the rows
	for _, row := range data {
		var record []string
		for _, key := range header {
			record = append(record, fmt.Sprintf("%v", row[key]))
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write row to CSV: %v", err)
		}
	}

	// Flush the CSV writer
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("failed to flush writer: %v", err)
	}

	return buf.Bytes(), nil
}

func CSVBytesToArrayOfMap(data []byte) ([]map[string]interface{}, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading headers: %v", err)
	}
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}
	var result []map[string]interface{}
	for _, row := range rows {
		if len(row) != len(headers) {
			return nil, fmt.Errorf("row has a different number of columns than headers")
		}

		// Create a map for the current row
		rowMap := make(map[string]interface{})
		for i, value := range row {
			rowMap[headers[i]] = value
		}
		result = append(result, rowMap)
	}
	return result, nil
}

func CreateTarFile(tarFile string, files2Add []string, fileDest string) error {
	tarFileHandle, err := os.Create(tarFile)
	if err != nil {
		return err
	}
	defer tarFileHandle.Close()

	tw := tar.NewWriter(tarFileHandle)
	defer tw.Close()

	for _, fl := range files2Add {
		file, err := os.Open(fl)
		if err != nil {
			return err
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		hdr := &tar.Header{
			Name: filepath.Join(fileDest, filepath.Base(file.Name())),
			Mode: 0644,
			Size: fileInfo.Size(),
		}
		log.Println("File to be added to tar - ", filepath.Join(fileDest, filepath.Base(file.Name())))
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}

		_, err = io.Copy(tw, file)
		if err != nil {
			return err
		}
	}

	return err
}

func WriteYAMLFile(data any, fileName string, base64encoded bool) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	if base64encoded {
		bytes = []byte(base64.StdEncoding.EncodeToString(bytes))
	}
	log.Println("Writing to file - ", fileName, " data - ", string(bytes))
	err = os.WriteFile(fileName, bytes, 0644)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrWriteYAMLFile, err)
	}
	return nil
}

func WriteJSONFile(data any, fileName string, base64encoded bool) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	if base64encoded {
		bytes = []byte(base64.StdEncoding.EncodeToString(bytes))
	}
	err = os.WriteFile(fileName, bytes, 0644)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrWriteYAMLFile, err)
	}
	return nil
}

func WriteTOMLFile(data any, fileName string, base64encoded bool) error {
	bytes, err := toml.Marshal(data)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	if base64encoded {
		bytes = []byte(base64.StdEncoding.EncodeToString(bytes))
	}
	err = os.WriteFile(fileName, bytes, 0644)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrWriteYAMLFile, err)
	}
	return nil
}

var FileContentUtils = map[string]func([]byte) bool{
	"pdf":  isPDF,
	"json": isJSON,
	"yaml": isYAML,
	"csv":  isCSV,
	"jpg":  isJPG,
	"jpeg": isJPEG,
	"png":  isPNG,
	"webp": isWebP,
	"svg":  isSVG,
	"pptx": isPPTX,
	"xlsx": isXLSX,
	"docx": isDOCX,
	// "txt":  isTXT,
}

func DetectFileTypeFromContent(content []byte) string {
	var wg sync.WaitGroup
	var foundChan = make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for n, f := range FileContentUtils {
		wg.Add(1)
		go func(ctx context.Context, f func([]byte) bool, fileType string) {
			defer wg.Done()
			if f(content) {
				log.Println("DetectFileTypeFromContent: found:+++++++++++++++++++++++++ ", fileType)
				select {
				case foundChan <- fileType:
					cancel()
				default:
					return
				}
			}
		}(ctx, f, n)
	}
	go func() {
		wg.Wait()
		close(foundChan)
	}()
	found := <-foundChan
	log.Println("DetectFileTypeFromContent: found:================== ", found)
	if found == "" {
		return "txt"
	}
	return found
}

func isPDF(content []byte) bool {
	return bytes.HasPrefix(content, []byte("%PDF-"))
}

// isJSON attempts to unmarshal the content as JSON.
func isJSON(content []byte) bool {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	var js interface{}
	return decoder.Decode(&js) == nil
}

// isYAML attempts to unmarshal the content as YAML.
func isYAML(content []byte) bool {
	var ys interface{}
	return yaml.Unmarshal(content, &ys) == nil
}

// isCSV uses heuristics to determine if the content is CSV.
// It checks for multiple lines with consistent number of fields.
func isCSV(content []byte) bool {
	reader := csv.NewReader(bytes.NewReader(content))
	reader.FieldsPerRecord = -1 // Disable field count check

	var fieldCount int
	lineCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false
		}
		// Skip empty lines
		if len(record) == 0 {
			continue
		}
		if lineCount == 0 {
			fieldCount = len(record)
		} else {
			if len(record) != fieldCount {
				return false
			}
		}
		lineCount++
	}
	// Assume CSV if there are at least two lines with consistent fields
	return lineCount >= 2
}

// isJPG checks if the content starts with the JPG magic number.
func isJPG(content []byte) bool {
	return bytes.HasPrefix(content, []byte{0xFF, 0xD8, 0xFF})
}

// isJPEG checks if the content starts with the JPEG magic number.
func isJPEG(content []byte) bool {
	return bytes.HasPrefix(content, []byte{0xFF, 0xD8, 0xFF})
}

// isPNG checks if the content starts with the PNG magic number.
func isPNG(content []byte) bool {
	pngSignature := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	return bytes.HasPrefix(content, pngSignature)
}

// isWebP checks if the content has the WebP magic number.
// WebP files start with "RIFF" followed by file size and "WEBP".
func isWebP(content []byte) bool {
	if len(content) < 12 {
		return false
	}
	return bytes.Equal(content[0:4], []byte("RIFF")) && bytes.Equal(content[8:12], []byte("WEBP"))
}

// isSVG checks if the content starts with an SVG tag.
func isSVG(content []byte) bool {
	contentStr := strings.TrimSpace(string(content))
	return strings.HasPrefix(contentStr, "<svg")
}

// isPPTX checks if the content is a PPTX file by inspecting ZIP entries.
func isPPTX(content []byte) bool {
	return isOfficeFormat(content, "ppt/")
}

// isXLSX checks if the content is an XLSX file by inspecting ZIP entries.
func isXLSX(content []byte) bool {
	return isOfficeFormat(content, "xl/")
}

// isDOCX checks if the content is a DOCX file by inspecting ZIP entries.
func isDOCX(content []byte) bool {
	return isOfficeFormat(content, "word/")
}

// isOfficeFormat checks if a ZIP archive contains specific directories indicating Office formats.
func isOfficeFormat(content []byte, requiredPath string) bool {
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return false
	}
	for _, file := range reader.File {
		if strings.HasPrefix(file.Name, requiredPath) {
			return true
		}
	}
	return false
}

// isTXT checks if the content is plain text.
// Simple heuristic: attempt to convert bytes to string and check for printable characters.
func isTXT(content []byte) bool {
	// Define a threshold for printable characters
	printable := 0
	total := 0
	for _, b := range content {
		if b >= 32 && b <= 126 || b == 10 || b == 13 || b == 9 {
			printable++
		}
		total++
	}
	if total == 0 {
		return false
	}
	return float64(printable)/float64(total) > 0.95 // 95% printable
}

// TODO: Add more file type detection functions
func ConvertCsvTo(fileType, expectedformat string, contentBytes []byte) ([]byte, error) {
	var fBytes []byte
	var err error
	if fileType == expectedformat {
		return contentBytes, nil
	}
	result := []map[string]interface{}{}
	result, err = CSVBytesToArrayOfMap(contentBytes)
	if err != nil {
		return nil, err
	}
	log.Println("ConvertFileBytes: result: ", result)
	switch strings.ToLower(expectedformat) {
	case "json":
		fBytes, err = GenericMarshaler(result, strings.ToUpper(expectedformat))
	case "yaml":
		fBytes, err = GenericMarshaler(result, strings.ToUpper(expectedformat))
	case "csv":
		fBytes, err = MapArrayCSVMarshaler(result)
	default:
		fBytes, err = MapArrayCSVMarshaler(result)
	}
	if err != nil {
		return nil, err
	}
	return fBytes, nil
}

func ConvertFile(fileType, expectedformat string, contentBytes []byte) ([]byte, error) {
	var fBytes []byte
	var err error
	if fileType == expectedformat {
		return contentBytes, nil
	}
	result := []map[string]interface{}{}
	err = GenericUnMarshaler(contentBytes, &result, fileType)
	if err != nil {
		return nil, err
	}
	log.Println("ConvertFileBytes: result: ", result)
	switch strings.ToLower(expectedformat) {
	case "json":
		fBytes, err = GenericMarshaler(result, strings.ToUpper(expectedformat))
	case "yaml":
		fBytes, err = GenericMarshaler(result, strings.ToUpper(expectedformat))
	case "csv":
		fBytes, err = MapArrayCSVMarshaler(result)
	default:
		fBytes, err = MapArrayCSVMarshaler(result)
	}
	if err != nil {
		return nil, err
	}
	return fBytes, nil
}
