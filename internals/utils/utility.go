package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	os "os"
	"reflect"
	"strconv"
	"strings"

	guuid "github.com/google/uuid"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"google.golang.org/protobuf/types/known/structpb"
)

func AStructToAMap(data interface{}) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}

	err = json.Unmarshal(bytes, &result)

	if err != nil {
		return nil, dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	return result, nil
}

func StructToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}

	err = json.Unmarshal(bytes, &result)

	if err != nil {
		return nil, dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	return result, nil
}

func AStructToAString(data any) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Error converting struct to string: ", err)
		return EMPTYSTR
	}

	return string(bytes)
}

func StrToStruct(data any, resp interface{}) error {
	log.Println("data to be marshelled - ", data)
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Error converting interface to string: ", err)
		return err
	}

	err = json.Unmarshal(bytes, resp)

	if err != nil {
		return dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	return nil
}

func Int32Ptr(x int32) *int32 { return &x }

func Int64Ptr(x int64) *int64 { return &x }

func IntPtr(x int) *int { return &x }

func Float32Ptr(x float32) *float32 { return &x }

func Float64Ptr(x float64) *float64 { return &x }

func Bool2Ptr(x bool) *bool { return &x }

func Str2Ptr(x string) *string { return &x }

func GetObjectTypeName(obj any) string {
	return reflect.TypeOf(obj).Name()
}

func TrailingSlash(path string, verify, add bool) (bool, string) {
	if verify && strings.HasSuffix(path, "/") {
		return true, path
	}
	if add && !strings.HasSuffix(path, "/") {
		return true, path + "/"
	}
	return false, path
}

func CreateSecretName(opts ...string) string {
	return strings.Join(opts, "::")
}

func IsInt(s string) (bool, int) {
	val, err := strconv.Atoi(s)
	return err == nil, val
}

func IsFloat(s string) (bool, float64) {
	val, err := strconv.ParseFloat(s, 64)
	return err == nil, val
}

func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func SetCtxValue(ctx context.Context, key ContextKeys, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetCtxValue(ctx context.Context, key ContextKeys) any {
	return ctx.Value(key)
}

func ListStructPbToSlice(data []*structpb.Struct) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, v := range data {
		result = append(result, StructPbToMap(v))
	}
	return result
}

func StructPbToMap(obj *structpb.Struct) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range obj.Fields {
		switch v := value.GetKind().(type) {
		case *structpb.Value_NumberValue:
			result[key] = v.NumberValue
		case *structpb.Value_StringValue:
			result[key] = v.StringValue
		case *structpb.Value_BoolValue:
			result[key] = v.BoolValue
			// case *structpb.Value_StructValue:
			// 	result[key] = StructPbToMap(v.StructValue, data)
			// case *structpb.Value_ListValue:
			// 	// Handle lists (arrays)
			// 	var list []interface{}
			// 	for _, item := range v.ListValue.Values {
			// 		// Recursively handle items in the list
			// 		itemValue, err := StructPbToMap(&structpb.Struct{Fields: map[string]*structpb.Value{"item": item}})
			// 		if err != nil {
			// 			return nil, err
			// 		}
			// 		list = append(list, itemValue["item"])
			// 	}
			// 	result[key] = list
		}
	}
	return result
}

func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

func DeepCopy(src, dest any) (any, error) {
	bytes, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func GetUUID() string {
	return guuid.New().String()
}

func IsNil[T any](val T) bool {
	// Use reflection to check if the value is nil
	return reflect.ValueOf(val).Kind() == reflect.Ptr && reflect.ValueOf(val).IsNil()
}
