package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	TAG          = "validate"
	AND          = "|"
	ValidatorSep = ":"
	InSep        = ","
)

type ProgrammingError struct {
	Msg string
}

func (err ProgrammingError) Error() string {
	return err.Msg
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	result := strings.Builder{}

	for _, val := range v {
		s := fmt.Sprintf("field %s has error %s\n", val.Field, val.Err.Error())
		result.WriteString(s)
	}

	return result.String()
}

type Validator struct {
	valueType        reflect.Type
	value            reflect.Value
	validationErrors []ValidationError
	fields           []Field
}

type Field struct {
	kind             reflect.Kind
	field            reflect.StructField
	index            int
	valueType        reflect.Type
	value            reflect.Value
	tag              string
	validationErrors []ValidationError
	validationFn     []ValidationFn
}

func (fl *Field) IsInt() bool {
	return fl.kind == reflect.Int
}

func (fl *Field) IsString() bool {
	return fl.kind == reflect.String
}

func (fl *Field) IsSlice() bool {
	return fl.kind == reflect.Slice
}

func (fl *Field) AddValidationError(err ValidationError) {
	fl.validationErrors = append(fl.validationErrors, err)
}

func (fl *Field) ValidateSlice() error {
	var progErr error

	sliceLen := fl.value.Len()
	slice := fl.value.Slice(0, sliceLen)

LOOP:
	for i := 0; i < sliceLen; i++ {
		sliceValue := slice.Index(i)
		for _, valFn := range fl.validationFn {
			ok, err := valFn(sliceValue.String())
			if err != nil {
				progErr = ProgrammingError{Msg: getInvalidValidatorErrorMsg()}
				break LOOP
			}

			if !ok {
				fl.AddValidationError(
					ValidationError{
						Field: fmt.Sprintf("%s[%d]", fl.valueType.Name(), i),
						Err:   errors.New(getValidationErrorMsg()),
					})
			}
		}
	}

	if progErr != nil {
		return progErr
	}

	return nil
}

func (fl *Field) ValidatePrimitive() error {
	var progErr error

	for _, valFn := range fl.validationFn {
		var ok bool
		var err error

		if fl.IsInt() {
			ok, err = valFn(fl.value.Int())
		}

		if fl.IsString() {
			ok, err = valFn(fl.value.String())
		}

		if err != nil {
			progErr = ProgrammingError{Msg: getInvalidValidatorErrorMsg()}
			break
		}

		if !ok {
			fl.AddValidationError(
				ValidationError{
					Field: fl.valueType.Name(),
					Err:   errors.New(getValidationErrorMsg()),
				})
		}
	}
	if progErr != nil {
		return progErr
	}

	return nil
}

func (fl *Field) Validate() ([]ValidationError, error) {
	if fl.IsSlice() {
		err := fl.ValidateSlice()
		if err != nil {
			return nil, err
		}
	}

	if fl.IsInt() || fl.IsString() {
		err := fl.ValidatePrimitive()
		if err != nil {
			return nil, err
		}
	}

	return fl.validationErrors, nil
}

type (
	ValidationFn        func(value interface{}) (bool, error)
	ValidationFnCreator func(rawCond string) (ValidationFn, error)
	CheckTypeFn         func(field reflect.StructField) bool
)

func (vl *Validator) IsStruct() bool {
	valueType := vl.valueType.Kind()

	return valueType == reflect.Struct
}

func (vl *Validator) AddField(field Field) {
	vl.fields = append(vl.fields, field)
}

func (vl *Validator) Validate() (ValidationErrors, error) {
	progErr := false

	for _, field := range vl.fields {
		valErrors, err := field.Validate()
		if err != nil {
			progErr = true
			break
		}

		if valErrors != nil {
			vl.validationErrors = append(vl.validationErrors, valErrors...)
		}
	}

	if progErr {
		return nil, ProgrammingError{Msg: getInvalidValidatorErrorMsg()}
	}

	if len(vl.validationErrors) == 0 {
		return nil, nil
	}

	return vl.validationErrors, nil
}

func getTypeErrorMsg() string {
	return "value is not structure"
}

func getInvalidValidatorErrorMsg() string {
	return "invalid validator"
}

func getValidationErrorMsg() string {
	return "value is not valid"
}

var StringValidationFnMap = map[string]ValidationFnCreator{
	"len":    LenStringValidation,
	"regexp": RegexpStringValidation,
	"in":     inStringValidation,
}

var NumValidationFnMap = map[string]ValidationFnCreator{
	"min": MinNumValidation,
	"max": MaxNumValidation,
	"in":  inNumValidation,
}

func LenStringValidation(rawCond string) (ValidationFn, error) {
	cond, err := strconv.Atoi(rawCond)
	if err != nil {
		return nil, ProgrammingError{Msg: getInvalidValidatorErrorMsg()}
	}

	fn := func(value interface{}) (bool, error) {
		val, ok := value.(string)

		if !ok {
			return false, ProgrammingError{Msg: getTypeErrorMsg()}
		}

		if len(val) != cond {
			return false, nil
		}
		return true, nil
	}

	return fn, nil
}

func RegexpStringValidation(rawCond string) (ValidationFn, error) {
	cond, err := regexp.Compile(rawCond)
	if err != nil {
		return nil, ProgrammingError{Msg: getInvalidValidatorErrorMsg()}
	}

	fn := func(value interface{}) (bool, error) {
		val, ok := value.(string)

		if !ok {
			return false, ProgrammingError{Msg: getTypeErrorMsg()}
		}

		valid := cond.MatchString(val)

		return valid, nil
	}

	return fn, nil
}

func inStringValidation(rawCond string) (ValidationFn, error) {
	parsedCond := strings.Split(rawCond, InSep)

	fn := func(value interface{}) (bool, error) {
		val, ok := value.(string)

		if !ok {
			return false, ProgrammingError{Msg: getTypeErrorMsg()}
		}

		valid := false

		for _, cond := range parsedCond {
			if cond == val {
				valid = true
				break
			}
		}

		return valid, nil
	}

	return fn, nil
}

func MinNumValidation(rawCond string) (ValidationFn, error) {
	cond, err := strconv.Atoi(rawCond)
	if err != nil {
		return nil, ProgrammingError{Msg: getInvalidValidatorErrorMsg()}
	}

	fn := func(value interface{}) (bool, error) {
		val, ok := value.(int64)

		if !ok {
			return false, ProgrammingError{Msg: getTypeErrorMsg()}
		}

		if int(val) < cond {
			return false, nil
		}
		return true, nil
	}

	return fn, nil
}

func MaxNumValidation(rawCond string) (ValidationFn, error) {
	cond, err := strconv.Atoi(rawCond)
	if err != nil {
		return nil, ProgrammingError{Msg: getInvalidValidatorErrorMsg()}
	}

	fn := func(value interface{}) (bool, error) {
		val, ok := value.(int64)

		if !ok {
			return false, ProgrammingError{Msg: getTypeErrorMsg()}
		}

		if int(val) > cond {
			return false, nil
		}
		return true, nil
	}

	return fn, nil
}

func inNumValidation(rawCond string) (ValidationFn, error) {
	progErr := false
	parsedCondStr := strings.Split(rawCond, InSep)
	parsedCond := make([]int, 0)

	for _, cond := range parsedCondStr {
		cnd, err := strconv.Atoi(cond)
		if err != nil {
			progErr = true
			break
		}

		parsedCond = append(parsedCond, cnd)
	}

	if progErr {
		return nil, ProgrammingError{Msg: getInvalidValidatorErrorMsg()}
	}

	fn := func(value interface{}) (bool, error) {
		val, ok := value.(int64)

		if !ok {
			return false, ProgrammingError{Msg: getTypeErrorMsg()}
		}
		valid := false

		for _, cond := range parsedCond {
			if cond == int(val) {
				valid = true
				break
			}
		}

		return valid, nil
	}

	return fn, nil
}

func isInt(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.Int
}

func isString(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.String
}

func isSliceInt(field reflect.StructField) bool {
	return reflect.SliceOf(reflect.TypeOf(1)) == field.Type
}

func isSliceString(field reflect.StructField) bool {
	return reflect.SliceOf(reflect.TypeOf("s")) == field.Type
}

func isValidValueType(field reflect.StructField) bool {
	valid := false
	for _, fn := range []CheckTypeFn{isInt, isSliceInt, isString, isSliceString} {
		if fn(field) {
			valid = true
			break
		}
	}

	return valid
}

func createValidator(v interface{}) (*Validator, error) {
	validator := &Validator{
		valueType:        reflect.TypeOf(v),
		value:            reflect.ValueOf(v),
		validationErrors: make([]ValidationError, 0),
	}

	if !validator.IsStruct() {
		return nil, ProgrammingError{Msg: getTypeErrorMsg()}
	}

	validator.fields = make([]Field, 0, validator.valueType.NumField())

	progErr := false

	for i := 0; i < validator.valueType.NumField(); i++ {
		field := validator.valueType.Field(i)

		if !field.IsExported() || !isValidValueType(field) {
			continue
		}

		tag, ok := field.Tag.Lookup(TAG)

		if !ok || len(tag) == 0 {
			continue
		}

		validationFn := make([]ValidationFn, 0)

		rules := strings.Split(tag, AND)

		for _, rule := range rules {
			parsedRule := strings.Split(rule, ValidatorSep)

			if len(parsedRule) < 2 {
				progErr = true
				break
			}

			validationMap := StringValidationFnMap

			if isInt(field) || isSliceInt(field) {
				validationMap = NumValidationFnMap
			}

			if fn, ok := validationMap[parsedRule[0]]; ok {
				valFn, err := fn(parsedRule[1])
				if err != nil {
					progErr = true
					break
				}

				validationFn = append(validationFn, valFn)
			}
		}

		var kind reflect.Kind

		switch {
		case isInt(field):
			kind = reflect.Int
		case isString(field):
			kind = reflect.String
		default:
			kind = reflect.Slice
		}

		validator.AddField(Field{
			kind:             kind,
			field:            field,
			index:            i,
			valueType:        field.Type,
			value:            validator.value.Field(i),
			tag:              tag,
			validationErrors: make([]ValidationError, 0),
			validationFn:     validationFn,
		})
	}

	if progErr {
		return nil, ProgrammingError{getInvalidValidatorErrorMsg()}
	}

	return validator, nil
}

func Validate(v interface{}) error {
	if v == nil {
		return ProgrammingError{Msg: getTypeErrorMsg()}
	}

	validator, err := createValidator(v)
	if err != nil {
		return err
	}

	valErr, pErr := validator.Validate()

	if pErr != nil {
		return pErr
	}

	return valErr
}
