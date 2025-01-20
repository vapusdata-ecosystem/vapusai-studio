package router

import (
	"encoding/json"
	"math/rand"
	"slices"
	"strings"
	"text/template"

	utils "github.com/vapusdata-oss/aistudio/core/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/encoding/protojson"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	"gopkg.in/yaml.v3"
)

func limitWords(text string, limit int) string {
	words := strings.Fields(text)
	if len(words) > limit {
		return strings.Join(words[:limit], " ") + "..."
	}
	return text
}

func limitletters(text string, limit int) string {
	if len(text) == 0 {
		return "--"
	}
	letters := strings.Split(text, "")
	if len(letters) > limit {
		return strings.Join(letters[:limit], "") + "..."
	}
	return text
}

func EpochConverter(epoch int64) string {
	if epoch == 0 {
		return "--"
	}
	return utils.GetFormattedTime(epoch, "2006-01-02")
}

func EpochConverterFull(epoch int64) string {
	if epoch == 0 {
		return "--"
	}
	return utils.GetFormattedTime(epoch, "2006-01-02 15:04")
}

func EpochConverterTextDate(epoch int64) string {
	if epoch == 0 {
		return "--"
	}
	return utils.GetFormattedTime(epoch, "1 January 2006")
}

func InSlice(value string, list ...string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func SliceContains(value string, list []string) bool {
	return slices.Contains(list, value)
}

func toJSON(v interface{}) string {
	a, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(a)
}

func protoToJSON(v any) string {
	l, ok := v.(protoreflect.ProtoMessage)
	if !ok {
		a, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(a)
	}
	a, err := protojson.MarshalOptions{
		Multiline:      true,
		UseEnumNumbers: false,
	}.Marshal(l)
	if err != nil {
		return ""
	}
	return string(a)
}

func stringCheck(s string) string {
	if s == "" {
		return "--"
	}
	return s
}

func escapeHTML(input string) string {
	return template.HTMLEscapeString(input)
}

func escapeJSON(input string) string {
	return template.JSEscapeString(input)
}

func addRand(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func randBool() bool {
	return rand.Intn(2) == 0
}

func marshalToYaml(v any) string {
	a, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}
	return string(a)
}

func strContains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func sliceLen[T any](slice []T, expectedLen int, condition string) bool {
	if condition == "==" {
		return len(slice) == expectedLen
	} else if condition == ">" {
		return len(slice) > expectedLen
	} else if condition == "<" {
		return len(slice) < expectedLen
	}
	return false
}

func slugToTitle(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	return cases.Title(language.English).String(s)
}

func enumoTitle(s any) string {
	d := s.(protoreflect.Enum)
	value := d.Descriptor().Values().ByNumber(d.Number())
	cc := string(value.Name())
	return cases.Title(language.English).String(cc)
}

func joinSlice(s []string, separator string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.Join(s, separator)
}
