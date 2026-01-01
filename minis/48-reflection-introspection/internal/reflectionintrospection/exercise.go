//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package reflectionintrospection
// TODO: implement GetTypeName.
func GetTypeName(v interface{}) string { panic("TODO: implement") }
// TODO: implement GetKind.
func GetKind(v interface{}) string { panic("TODO: implement") }
// TODO: implement CountFields.
func CountFields(v interface{}) int { panic("TODO: implement") }
// TODO: implement GetJSONTag.
func GetJSONTag(v interface{}, fieldName string) string { panic("TODO: implement") }
// TODO: implement GetAllTags.
func GetAllTags(v interface{}, fieldName string) string { panic("TODO: implement") }
// TODO: implement GetFieldValue.
func GetFieldValue(val interface{}, fieldName string) interface{} { panic("TODO: implement") }
// TODO: implement GetFieldValues.
func GetFieldValues(val interface{}) map[string]interface{} { panic("TODO: implement") }
// TODO: implement SetFieldValue.
func SetFieldValue(val interface{}, fieldName string, newValue interface{}) error {
	panic("TODO: implement")
}
// TODO: implement CallMethod.
func CallMethod(obj interface{}, methodName string, args ...interface{}) []interface{} {
	panic("TODO: implement")
}
// TODO: implement HasMethod.
func HasMethod(obj interface{}, methodName string) bool { panic("TODO: implement") }
// TODO: implement SameType.
func SameType(a, b interface{}) bool { panic("TODO: implement") }
// TODO: implement IsPointer.
func IsPointer(v interface{}) bool { panic("TODO: implement") }
// TODO: implement SliceLength.
func SliceLength(slice interface{}) int { panic("TODO: implement") }
// TODO: implement MapKeys.
func MapKeys(m interface{}) []interface{} { panic("TODO: implement") }
// TODO: implement NewInstance.
func NewInstance(v interface{}) interface{} { panic("TODO: implement") }
// TODO: implement GetFieldNames.
func GetFieldNames(v interface{}) []string { panic("TODO: implement") }
// TODO: implement DeepCopy.
func DeepCopy(v interface{}) interface{} { panic("TODO: implement") }
