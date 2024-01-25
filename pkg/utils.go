package pkg

import (
	"reflect"
	"strings"

	"github.com/mss-boot-io/mss-boot/pkg/enum"
	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mgm "github.com/kamva/mgm/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/schema"
)

const (
	// TrafficKey traffic key
	TrafficKey = "X-Request-ID"
	// LoggerKey logger key
	LoggerKey = "_go-admin-logger-request"
)

// CompareHashAndPassword compare hash and password
func CompareHashAndPassword(hash string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

// GenerateMsgIDFromContext 生成msgID
func GenerateMsgIDFromContext(c *gin.Context) string {
	requestID := c.GetHeader(TrafficKey)
	if requestID == "" {
		requestID = uuid.New().String()
		c.Header(TrafficKey, requestID)
	}
	return requestID
}

// ModelDeepCopy model deep copy
func ModelDeepCopy(m mgm.Model) mgm.Model {
	return reflect.New(reflect.TypeOf(m).Elem()).Interface().(mgm.Model)
}

// TablerDeepCopy model deep copy
func TablerDeepCopy(m schema.Tabler) schema.Tabler {
	return reflect.New(reflect.TypeOf(m).Elem()).Interface().(schema.Tabler)
}

// DeepCopy deep copy
func DeepCopy(d any) any {
	return reflect.New(reflect.TypeOf(d).Elem()).Interface()
}

// BuildMap build map
func BuildMap(keys []string, value string, dataType enum.DataType) map[string]any {
	data := make(map[string]any)
	if len(keys) > 1 {
		data[keys[0]] = BuildMap(keys[1:], value, dataType)
	} else {
		var v any
		switch dataType {
		case enum.DataTypeInt:
			v, _ = cast.ToIntE(value)
		case enum.DataTypeFloat:
			v, _ = cast.ToFloat64E(value)
		case enum.DataTypeBool:
			v, _ = cast.ToBoolE(value)
		default:
			v = value
		}
		return map[string]any{keys[0]: v}
	}
	return data
}

// MergeMapsDepth deep merge multi map
func MergeMapsDepth(ms ...map[string]any) map[string]any {
	data := make(map[string]any)
	for i := range ms {
		data = MergeMapDepth(data, ms[i])
	}
	return data
}

// MergeMapDepth deep merge map
func MergeMapDepth(m1, m2 map[string]any) map[string]any {
	for k := range m2 {
		if v, ok := m1[k]; ok {
			if m, ok := v.(map[string]any); ok {
				m1[k] = MergeMapDepth(m, m2[k].(map[string]any))
			} else {
				m1[k] = m2[k]
			}
		} else {
			m1[k] = m2[k]
		}
	}
	return m1
}

// MergeMap merge map
func MergeMap(m1, m2 map[string]any) map[string]any {
	for k := range m2 {
		m1[k] = m2[k]
	}
	return m1
}

// SupportMultiTenant support multi tenant
func SupportMultiTenant(data any) bool {
	return supportColumn(data, "tenantID", "tenant_id")
}

func SupportCreator(data any) bool {
	return supportColumn(data, "creatorID", "creator_id")
}

func supportColumn(data any, fields ...string) bool {
	typeOf := reflect.TypeOf(data)
	valueOf := reflect.ValueOf(data)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}

	var exist bool
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		if field.Type.Kind() == reflect.Struct {
			exist = SupportMultiTenant(valueOf.Field(i).Interface())
		}
		if field.Type.Kind() == reflect.Ptr {
			continue
		}
		for j := range fields {
			exist = exist || strings.ToLower(field.Name) == strings.ToLower(fields[j])
			if exist {
				break
			}
		}
		if exist {
			break
		}
	}
	return exist
}

func SetValue(data any, key string, value any) {
	typeOf := reflect.TypeOf(data)
	valueOf := reflect.ValueOf(data)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	key = strings.ToLower(key)
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		if field.Type.Kind() == reflect.Struct {
			SetValue(valueOf.Field(i).Interface(), key, value)
			continue
		}
		if field.Type.Kind() == reflect.Ptr {
			continue
		}
		if strings.ToLower(field.Name) == key {
			valueOf.Field(i).Set(reflect.ValueOf(value))
		}
	}
}
