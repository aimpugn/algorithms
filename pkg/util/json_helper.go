package util

import (
	"reflect"
	"strconv"
	"strings"
)

/*func JsonDecode(data []byte, v interface{}) error {
	// Check for well-formedness.
	// Avoids filling out half a data structure
	// before discovering a JSON syntax error.
	var d decodeState
	err := checkValid(data, &d.scan)
	if err != nil {
		return err
	}

	d.init(data)

	return d.unmarshal(v)
}*/

/*
scanner.state 에 할당되는 상태 transition 함수 및 the method scanner.eof 메서드에 의해 반환되는 값들
호출자가 관심 있어할 스캔의 현재 상태 상세 정보를 보여준다
scanner.state에 대한 어떤 특정한 호출의 반환 값은 무시해도 괜찮다:
	만약 scanError 반환하면, 모든 순차적인 호출 역시 scanError 반환
It is okay to ignore the return value of any particular call to scanner.state:
	if one call returns scanError, every subsequent call will return scanError too.
*/
const (
	// Continue.
	scanContinue     = iota // 0, uninteresting byte
	scanBeginLiteral        // 1, end implied by next result != scanContinue
	scanBeginObject         // 2, begin object
	scanObjectKey           // 3, just finished object key (string)
	scanObjectValue         // 4, just finished non-last object value
	scanEndObject           // 5, end object (implies scanObjectValue if possible)
	scanBeginArray          // 6, begin array
	scanArrayValue          // 7, just finished array value
	scanEndArray            // 8, end array (implies scanArrayValue if possible)
	scanSkipSpace           // 9, space byte; can skip; known to be last "continue" result

	// Stop.
	scanEnd   // 10, top-level value ended *before* this byte; known to be first "stop" result
	scanError // 11, hit an error, scanner.err.
)

/*
parseState 스택에 쌓이는 값들
스캔된 복합적인 값들의 현재 상태 보여준다
파서가 중첩된 값 내부에 있다면, parseState는 중첩된 상태를 나타내며, 가장 바깥은 0이다
*/
const (
	/*
		iota 없으면 Missing value in const declaration 에러 발생
		0부터 증감. 0, 1, 2
	*/
	parseObjectKey   = iota // parsing object key (before colon)
	parseObjectValue        // parsing object value (after colon)
	parseArrayValue         // parsing array value
)

/*
스택오버플로 방지 위한 최대 중첩 뎁쓰
https://tools.ietf.org/html/rfc7159#section-9
*/
const maxNestingDepth = 10000

type scanner struct {
	/*
		step 함수는 다음 변화(transition) 실행하기 위해 호출되는 함수
		정수 상수 및 switch로 분기처리 되는 단일 함수 사용했지만, 64비트 Mac Mini에서 직접 호출하는 것이 10% 더 빠르고 읽기 더 좋다
	*/
	step func(*scanner, byte) int

	/*
		최상위 레벨의 값의 끝에 도달
	*/
	endTop bool

	/*
		배열의 값, 오브젝트의 키, 오브젝트의 값 등 무엇을 처리중에 있는지 확인하는 스택
	*/
	parseState []int

	/*
		배열의 수, maxNestingDepth 최대 중첩 깊이와 비교
	*/
	parseStateCnt int

	err error

	/*
		decoder.Decode에 의해 업데이트 되는, 소비되는 총 바이트
		scan.reset 시 일부러 0으로 설정하지 않는다
	*/
	bytes int64
}

/*
유효하게 JSON 인코딩 된 데이터인지 검증
할당(allocation)을 피하기 위해 checkValid 함수에 의해 사용되도록 scan 전달
*/
func checkValid(data []byte, scan *scanner) error {
	scan.reset()
	for _, c := range data {
		// 소비되는 총 바이트
		scan.bytes++
		if scan.step(scan, c) == scanError {
			return scan.err
		}
	}
	if scan.eof() == scanError {
		return scan.err
	}
	return nil
}

/*
scanner 사용할 수 있도록 준비
scanner.step 호출 전에 반드시 호출되어야 한다
*/
func (s *scanner) reset() {
	s.step = stateBeginValue
	s.parseState = s.parseState[0:0]
	s.err = nil
	s.endTop = false
}

/*
scanner 에게 입력의 끝에 도달했음을 알려준다
scanner.step 같이 [ scanError | scanEnd ] 스캔의 상태를 반환
*/
func (s *scanner) eof() int {
	if s.err != nil {
		return scanError
	}
	if s.endTop {
		return scanEnd
	}
	s.step(s, ' ')
	if s.endTop {
		return scanEnd
	}
	if s.err == nil {
		s.err = &SyntaxError{"unexpected end of JSON input", s.bytes}
	}
	return scanError
}

/*
JSON 구문 오류 설명
*/
type SyntaxError struct {
	msg    string // 에러 설명
	Offset int64  // Offset 바이트 읽은 후 에러 발생
}

func (e *SyntaxError) Error() string {
	return e.msg
}

/*
입력의 시작 부분의 상태
*/
func stateBeginValue(s *scanner, c byte) int {
	if isSpace(c) {
		return scanSkipSpace
	}
	switch c {
	case '{':
		s.step = stateBeginStringOrEmpty
		return s.pushParseState(c, parseObjectKey, scanBeginObject)
	case '[':
		s.step = stateBeginValueOrEmpty
		return s.pushParseState(c, parseArrayValue, scanBeginArray)
	case '"':
		s.step = stateInString
		return scanBeginLiteral
	case '-':
		s.step = stateNeg
		return scanBeginLiteral
	case '0': // beginning of 0.123
		s.step = state0
		return scanBeginLiteral
	case 't': // beginning of true
		s.step = stateT
		return scanBeginLiteral
	case 'f': // beginning of false
		s.step = stateF
		return scanBeginLiteral
	case 'n': // beginning of null
		s.step = stateN
		return scanBeginLiteral
	}
	if '1' <= c && c <= '9' { // beginning of 1234.5
		s.step = state1
		return scanBeginLiteral
	}
	return s.error(c, "looking for beginning of value")
}

func isSpace(c byte) bool {
	return c <= ' ' && (c == ' ' || c == '\t' || c == '\r' || c == '\n')
}

/*
parse 상태 p를 parseState 스택에 쌓는다
maxNestingDepth 초과하면 error 상태가 반환되며, 그렇지 않으면 successState 반환
*/
func (s *scanner) pushParseState(c byte, newParseState int, successState int) int {
	s.parseState = append(s.parseState, newParseState)
	if len(s.parseState) <= maxNestingDepth {
		return successState
	}
	return s.error(c, "exceeded max depth")
}

/*
스택에서 parse 상태값 제거하고 그에 따라 s.step 업데이트
*/
func (s *scanner) popParseState() {
	n := len(s.parseState) - 1
	s.parseState = s.parseState[0:n]
	if n == 0 {
		s.step = stateEndTop
		s.endTop = true
	} else {
		s.step = stateEndValue
	}
}

/*
오브젝트의 시작인 `{` 읽은 후의 상태
*/
func stateBeginStringOrEmpty(s *scanner, c byte) int {
	if isSpace(c) {
		return scanSkipSpace
	}
	if c == '}' {
		n := len(s.parseState)
		s.parseState[n-1] = parseObjectValue
		return stateEndValue(s, c)
	}
	return stateBeginString(s, c)
}

/*
`{}`, `true`, `["x"` 같은 문자들을 읽은 후와 같이, 값을 완료한 후의 상태
*/
func stateEndValue(s *scanner, c byte) int {
	n := len(s.parseState)

	if n == 0 {
		// Completed top-level before the current byte.
		s.step = stateEndTop
		s.endTop = true
		return stateEndTop(s, c)
	}

	if isSpace(c) { /* 공백인 경우 다음 `c byte`로 이동 */
		s.step = stateEndValue
		return scanSkipSpace
	}

	ps := s.parseState[n-1] /* 가장 최근 parse stack의 상태 값 */
	switch ps {
	case parseObjectKey:
		if c == ':' { /* end of parsing object key */
			s.parseState[n-1] = parseObjectValue
			s.step = stateBeginValue
			return scanObjectKey
		}
		return s.error(c, "after object key")
	case parseObjectValue:
		if c == ',' { /* end of parsing object value */
			s.parseState[n-1] = parseObjectKey
			s.step = stateBeginString
			return scanObjectValue
		}
		if c == '}' { /* end of parsing object itself */
			s.popParseState() /* pop state from state stack */
			return scanEndObject
		}
		return s.error(c, "after object key:value pair")
	case parseArrayValue:
		if c == ',' { /* end of parsing array value */
			s.step = stateBeginValue
			return scanArrayValue
		}
		if c == ']' { /* end of parsing array itself */
			s.popParseState() /* pop state from state stack */
			return scanEndArray
		}
		return s.error(c, "after array element")
	}
	return s.error(c, "")
}

/*
`{"key":value,` 읽은 후의 상태
*/
func stateBeginString(s *scanner, c byte) int {
	if isSpace(c) {
		return scanSkipSpace
	}
	if c == '"' {
		s.step = stateInString
		return scanBeginLiteral
	}
	return s.error(c, "looking for beginning of object key string")
}

/*
배열의 시작인 `[` 읽은 후의 상태
*/
func stateBeginValueOrEmpty(s *scanner, c byte) int {
	if isSpace(c) {
		return scanSkipSpace
	}
	if c == ']' {
		return stateEndValue(s, c)
	}
	return stateBeginValue(s, c)
}

/*
`"` 읽은 후의 상태
*/
func stateInString(s *scanner, c byte) int {
	if c == '"' {
		s.step = stateEndValue
		return scanContinue
	}
	if c == '\\' {
		s.step = stateInStringEsc
		return scanContinue
	}
	if c < 0x20 {
		return s.error(c, "in string literal")
	}
	return scanContinue
}

/*
`{}` 또는 `[1, 2, 3]` 같은 값을 읽은 후, 최상위 레벨 값 완료 후의 상태
이때는 오직 공백 문자만 나타나야 한다
*/
func stateEndTop(s *scanner, c byte) int {
	if !isSpace(c) {
		// Complain about non-space byte on next call.
		s.error(c, "after top-level value")
	}
	return scanEnd
}

/*
================================================================================================================
======================================================== 쌍따옴표 문자열에서 Escape 된 백슬래시 읽은 경우 시작
================================================================================================================

쌍따옴표된 문자열에서 `"\` 읽은 후
*/
func stateInStringEsc(s *scanner, c byte) int {
	switch c {
	case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
		s.step = stateInString
		return scanContinue
	case 'u':
		s.step = stateInStringEscU
		return scanContinue
	}
	return s.error(c, "in string escape code")
}

/*
유니코드0
쌍따옴표된 문자열에서 `"\u` 읽은 후
*/
func stateInStringEscU(s *scanner, c byte) int {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU1
		return scanContinue
	}
	// numbers
	return s.error(c, "in \\u hexadecimal character escape")
}

/*
유니코드1
쌍따옴표된 문자열에서 `"\u1` 읽은 후
*/
func stateInStringEscU1(s *scanner, c byte) int {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU12
		return scanContinue
	}
	// numbers
	return s.error(c, "in \\u hexadecimal character escape")
}

/*
유니코드2
쌍따옴표된 문자열에서 `"\u12` 읽은 후
*/
func stateInStringEscU12(s *scanner, c byte) int {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU123
		return scanContinue
	}
	// numbers
	return s.error(c, "in \\u hexadecimal character escape")
}

/*
유니코드3
쌍따옴표된 문자열에서 `"\u123` 읽은 후
*/
func stateInStringEscU123(s *scanner, c byte) int {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInString
		return scanContinue
	}
	// numbers
	return s.error(c, "in \\u hexadecimal character escape")
}

/*
================================================================================================================
======================================================== 쌍따옴표 문자열에서 Escape 된 백슬래시 읽은 경우 끝
================================================================================================================
*/

/*
================================================================================================================
======================================================== 숫자 시작
================================================================================================================

숫자에서 `-` 읽은 후의 상태
*/
func stateNeg(s *scanner, c byte) int {
	if c == '0' {
		s.step = state0
		return scanContinue
	}
	if '1' <= c && c <= '9' {
		s.step = state1
		return scanContinue
	}
	return s.error(c, "in numeric literal")
}

/*
숫자중 `0` 읽은 후의 상태
*/
func state0(s *scanner, c byte) int {
	if c == '.' {
		s.step = stateDot
		return scanContinue
	}
	if c == 'e' || c == 'E' {
		s.step = stateE
		return scanContinue
	}
	return stateEndValue(s, c)
}

/*
`1.`처럼 숫자에서 정수와 소수점(decimal point)을 읽은 후의 상태
*/
func stateDot(s *scanner, c byte) int {
	if '0' <= c && c <= '9' {
		s.step = stateDot0
		return scanContinue
	}
	return s.error(c, "after decimal point in numeric literal")
}

/*
숫자중 `0` 읽은 후의 상태
`1` 또는 `100`처럼 숫자중 `0` 아닌 숫자를 읽은 후의 상태
*/
func state1(s *scanner, c byte) int {
	if '0' <= c && c <= '9' {
		s.step = state1
		return scanContinue
	}
	return state0(s, c)
}

/*
`3.14`처럼 정수, 소수점 그리고 소수점 다음 숫자를 읽은 후의 상태
*/
func stateDot0(s *scanner, c byte) int {
	if '0' <= c && c <= '9' {
		return scanContinue
	}
	if c == 'e' || c == 'E' {
		s.step = stateE
		return scanContinue
	}
	return stateEndValue(s, c)
}

/*
`314e` 또는 `0.314e`를 읽은 후처럼, 가수와 e를 읽은 후의 상태
*/
func stateE(s *scanner, c byte) int {
	if c == '+' || c == '-' {
		s.step = stateESign
		return scanContinue
	}
	return stateESign(s, c)
}

/*
`314e-` 또는 `0.314e+`처럼 가수, e, 그리고 부호(sign)을 읽은 후의 상태
*/
func stateESign(s *scanner, c byte) int {
	if '0' <= c && c <= '9' {
		s.step = stateE0
		return scanContinue
	}
	return s.error(c, "in exponent of numeric literal")
}

/*
`314e-2` 또는 `0.314e+1` 또는 `3.14e0`를 읽은 후처럼, 가수, e, 부호(옵션), 최소 한 자리의 지수를 읽은 후의 상태
*/
func stateE0(s *scanner, c byte) int {
	if '0' <= c && c <= '9' {
		return scanContinue
	}
	return stateEndValue(s, c)
}

/*
================================================================================================================
======================================================== 숫자 끝
================================================================================================================
*/

/*
================================================================================================================
======================================================== true 시작
================================================================================================================

t를 읽은 후의 상태
*/
func stateT(s *scanner, c byte) int {
	if c == 'r' {
		s.step = stateTr
		return scanContinue
	}
	return s.error(c, "in literal true (expecting 'r')")
}

/*
tr 읽은 후의 상태
*/
func stateTr(s *scanner, c byte) int {
	if c == 'u' {
		s.step = stateTru
		return scanContinue
	}
	return s.error(c, "in literal true (expecting 'u')")
}

/*
tru 읽은 후의 상태
*/
func stateTru(s *scanner, c byte) int {
	if c == 'e' {
		s.step = stateEndValue
		return scanContinue
	}
	return s.error(c, "in literal true (expecting 'e')")
}

/*
================================================================================================================
======================================================== true 끝
================================================================================================================
*/

/*
================================================================================================================
======================================================== false 시작
================================================================================================================

`f` 읽은 후의 상태
*/
func stateF(s *scanner, c byte) int {
	if c == 'a' {
		s.step = stateFa
		return scanContinue
	}
	return s.error(c, "in literal false (expecting 'a')")
}

/*
`fa` 읽은 후의 상태
*/
func stateFa(s *scanner, c byte) int {
	if c == 'l' {
		s.step = stateFal
		return scanContinue
	}
	return s.error(c, "in literal false (expecting 'l')")
}

/*
`fal` 읽은 후의 상태
*/
func stateFal(s *scanner, c byte) int {
	if c == 's' {
		s.step = stateFals
		return scanContinue
	}
	return s.error(c, "in literal false (expecting 's')")
}

/*
`fals` 읽은 후의 상태
*/
func stateFals(s *scanner, c byte) int {
	if c == 'e' {
		s.step = stateEndValue
		return scanContinue
	}
	return s.error(c, "in literal false (expecting 'e')")
}

/*
================================================================================================================
======================================================== false 시작
================================================================================================================
*/

/*
================================================================================================================
======================================================== null 시작
================================================================================================================

`n` 읽은 후의 상태
*/
func stateN(s *scanner, c byte) int {
	if c == 'u' {
		s.step = stateNu
		return scanContinue
	}
	return s.error(c, "in literal null (expecting 'u')")
}

/*
`nu` 읽은 후의 상태
*/
func stateNu(s *scanner, c byte) int {
	if c == 'l' {
		s.step = stateNul
		return scanContinue
	}
	return s.error(c, "in literal null (expecting 'l')")
}

/*
`nul` 읽은 후의 상태
*/
func stateNul(s *scanner, c byte) int {
	if c == 'l' {
		s.step = stateEndValue
		return scanContinue
	}
	return s.error(c, "in literal null (expecting 'l')")
}

/*
================================================================================================================
======================================================== null 시작
================================================================================================================

`[1}` 또는 `5.1.2` 같은 문법 오류 발생 후의 상태
*/
func stateError(s *scanner, c byte) int {
	return scanError
}

/*
error 기록하고 error 상태로 전환
*/
func (s *scanner) error(c byte, context string) int {
	s.step = stateError
	s.err = &SyntaxError{"invalid character " + quoteChar(c) + " " + context, s.bytes}
	return scanError
}

/*
`c byte`에 quote('', "") 처리
*/
func quoteChar(c byte) string {
	// special cases - different from quoted strings
	if c == '\'' {
		return `'\''`
	}
	if c == '"' {
		return `'"'`
	}

	// use quoted string with different quotation marks
	s := strconv.Quote(string(c))
	return "'" + s[1:len(s)-1] + "'"
}

/*
================================================================================================================
======================================================== 구조체 선언
================================================================================================================

JSON 값 디코딩 중 상태값 보여주는 구조체
*/
type decodeState struct {
	data         []byte
	off          int /* 다음으로 읽을 데이터의 offset */
	opcode       int /* 마지막 read 결과로, 다음 operation code */
	scan         scanner
	errorContext struct { /* 타입 에러에 대한 context 제공 */
		Struct     reflect.Type
		FieldStack []string
	}
	savedError            error
	useNumber             bool
	disallowUnknownFields bool
}

/*
InvalidUnmarshalError 구조체는 Unmarshal 전달된 무효한 인자를 설명. 인자는 반드시 nil 아닌 포인터여야 한다
*/
type InvalidUnmarshalError struct {
	Type reflect.Type
}

/*
` InvalidUnmarshalError ` 구조체만 선언할 경우 아래와 같은 에러 발생
> Cannot use '&InvalidUnmarshalError{reflect.TypeOf(v)}' (type *InvalidUnmarshalError)
> as type error Type does not implement 'error' as some methods are missing: Error() string
왜? ` unmarshal(v interface{}) error{ ... } ` 함수의 정의 보면 `error`를 반환하는데, ` InvalidUnmarshalError ` 구조체만으로는 반환되는 error 없기 때문

*/
func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}

/*
` UnmarshalTypeError ` 특정 Go 타입의 값에 적절하지 않은 JSON 값을 나타낸다
*/
type UnmarshalTypeError struct {
	Value  string       // description of JSON value - "bool", "array", "number -5"
	Type   reflect.Type // type of Go value it could not be assigned to
	Offset int64        // error occurred after reading Offset bytes
	Struct string       // name of the struct type containing the field
	Field  string       // the full path from root node to the field
}

func (e *UnmarshalTypeError) Error() string {
	if e.Struct != "" || e.Field != "" {
		return "json: cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
	}
	return "json: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

const phasePanicMsg = "JSON decoder out of sync - data changing underfoot?"

// Unmarshaler is the interface implemented by types that can unmarshal a JSON description of themselves.
// The input can be assumed to be a valid encoding of
// a JSON value. UnmarshalJSON must copy the JSON data
// if it wishes to retain the data after returning.
//
// By convention, to approximate the behavior of Unmarshal itself,
// Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
/*
` Unmarshaler `는
*/
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.off = 0
	d.savedError = nil
	d.errorContext.Struct = nil

	// Reuse the allocated space for the FieldStack slice.
	d.errorContext.FieldStack = d.errorContext.FieldStack[:0]
	return d
}

func (d *decodeState) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() { /* 포인터가 아니거나, ` nil `인 경우*/
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	d.scan.reset()
	d.scanWhile(scanSkipSpace) /* scanSkipSpace == 10 */
	/*
		Unmarshaler interface 테스트가 반드시 최상위 수준의 값에(at the top level of the value) 적용되어야 하기 때문에,
		`rv.Elem` 아닌 `rv`를 디코드
	*/
	// err := d.value(rv)  /* JSON 값을 `v reflect.Value`로 디코딩한다 */
	// if err != nil {
	// 	return d.addErrorContext(err)
	// }

	return d.savedError
}

/*
op와 같지 않은 scan code 받을 때까지 decodeState.data[decodeState.off:] 바이트를 처리
*/
func (d *decodeState) scanWhile(op int) {
	s := &d.scan
	data := d.data
	i := d.off

	for i < len(data) {
		newOp := s.step(s, data[i]) /* 바이트 배열의 값을 다음 스텝으로 전달 */
		i++
		if newOp != op {
			d.opcode = newOp
			d.off = i
			return
		}
	}

	d.off = len(data) + 1 /* len + 1로 처리 된 EOF 표시 */
	d.opcode = d.scan.eof()
}

/* `addErrorContext`는 `decodeState.errorContext`의 정보로 새로운 향상된(enhanced) 에러를 반환 */
func (d *decodeState) addErrorContext(err error) error {
	if d.errorContext.Struct != nil || len(d.errorContext.FieldStack) > 0 {
		switch err := err.(type) {
		case *UnmarshalTypeError:
			err.Struct = d.errorContext.Struct.Name()
			err.Field = strings.Join(d.errorContext.FieldStack, ".")
			return err
		}
	}
	return err
}

/*
`value` 함수는 `decodeState.data[decodeState.off-1:]`의 JSON 값을 소비하하여 `v reflect.Value`로 디코딩하고, 다음 바이트를 읽어 나간다
만약 v 값이 무효(invalid)하다면, `JSON value`는 거부된다
`JSON value`의 첫 바이트는 이미 읽은 상태이다
*/
// func (d *decodeState) value(v reflect.Value) error {
// 	switch d.opcode {
// 	default:
// 		panic(phasePanicMsg)
//
// 	case scanBeginArray:  /* 배열 파싱 */
// 		if v.IsValid() {  /* ` v reflect.Value `가 값을 나타내는지 확인. ` zero Value `면 false 반환*/
// 			if err := d.array(v); err != nil {
// 				return err
// 			}
// 		} else {
// 			d.skip()
// 		}
// 		d.scanNext()
//
// 	case scanBeginObject:  /* 오브젝트 파싱 */
// 		if v.IsValid() {
// 			if err := d.object(v); err != nil {
// 				return err
// 			}
// 		} else {
// 			d.skip()
// 		}
// 		d.scanNext()
//
// 	case scanBeginLiteral:  /* 문자열 파싱 */
// 		// All bytes inside literal return scanContinue op code.
// 		start := d.readIndex()
// 		d.rescanLiteral()
//
// 		if v.IsValid() {
// 			if err := d.literalStore(d.data[start:d.readIndex()], v, false); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }
//
// /*
// ` array ` 함수는 ` decodeState.data[decodeState.off-1:] `의 배열을 소비하여 ` v reflect.Value `로 디코딩한다
// 배열의 첫 바이트 `[`는 이미 읽힌 상태이다
// 1. ` scanWhile ` 함수에서 ` decodeState.data[i] `의 바이트 읽어서 배열/오브젝트/문자열/음수/숫자/true/false/null 여부인지 판단하여
// 2. ` stateBeginValue `에서 이상 없으면 해당하는 ` opcode(operation code) ` 반환
// 3. 배열인 경우 ` value > scanBeginArray > array `로 이동
// */
// func (d *decodeState) array(v reflect.Value) error {
// 	// Check for unmarshaler.
// 	u, ut, pv := indirect(v, false)
// 	if u != nil {
// 		start := d.readIndex()
// 		d.skip()
// 		return u.UnmarshalJSON(d.data[start:d.off])
// 	}
// 	if ut != nil {
// 		d.saveError(&UnmarshalTypeError{Value: "array", Type: v.Type(), Offset: int64(d.off)})
// 		d.skip()
// 		return nil
// 	}
// 	v = pv
//
// 	// Check type of target.
// 	switch v.Kind() {
// 	case reflect.Interface:
// 		if v.NumMethod() == 0 {
// 			// Decoding into nil interface? Switch to non-reflect code.
// 			ai := d.arrayInterface()
// 			v.Set(reflect.ValueOf(ai))
// 			return nil
// 		}
// 		// Otherwise it's invalid.
// 		fallthrough
// 	default:
// 		d.saveError(&UnmarshalTypeError{Value: "array", Type: v.Type(), Offset: int64(d.off)})
// 		d.skip()
// 		return nil
// 	case reflect.Array, reflect.Slice:
// 		break
// 	}
//
// 	i := 0
// 	for {
// 		// Look ahead for ] - can only happen on first iteration.
// 		d.scanWhile(scanSkipSpace)
// 		if d.opcode == scanEndArray {
// 			break
// 		}
//
// 		// Get element of array, growing if necessary.
// 		if v.Kind() == reflect.Slice {
// 			// Grow slice if necessary
// 			if i >= v.Cap() {
// 				newcap := v.Cap() + v.Cap()/2
// 				if newcap < 4 {
// 					newcap = 4
// 				}
// 				newv := reflect.MakeSlice(v.Type(), v.Len(), newcap)
// 				reflect.Copy(newv, v)
// 				v.Set(newv)
// 			}
// 			if i >= v.Len() {
// 				v.SetLen(i + 1)
// 			}
// 		}
//
// 		if i < v.Len() {
// 			// Decode into element.
// 			if err := d.value(v.Index(i)); err != nil {
// 				return err
// 			}
// 		} else {
// 			// Ran out of fixed array: skip.
// 			if err := d.value(reflect.Value{}); err != nil {
// 				return err
// 			}
// 		}
// 		i++
//
// 		// Next token must be , or ].
// 		if d.opcode == scanSkipSpace {
// 			d.scanWhile(scanSkipSpace)
// 		}
// 		if d.opcode == scanEndArray {
// 			break
// 		}
// 		if d.opcode != scanArrayValue {
// 			panic(phasePanicMsg)
// 		}
// 	}
//
// 	if i < v.Len() {
// 		if v.Kind() == reflect.Array {
// 			// Array. Zero the rest.
// 			z := reflect.Zero(v.Type().Elem())
// 			for ; i < v.Len(); i++ {
// 				v.Index(i).Set(z)
// 			}
// 		} else {
// 			v.SetLen(i)
// 		}
// 	}
// 	if i == 0 && v.Kind() == reflect.Slice {
// 		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
// 	}
// 	return nil
// }
//
// /*
// ` indirect `는 non-pointer에 도달할 때까지, 필요에 따라 포인터를 할당하는 `v reflect.Value`를 따라간다
// ` Unmarshaler `를 만나면, ` indirect `는 멈추고 ` Unmarshaler `를 반환
// ` decodingNull `이 ` true `면, ` indirect `는 첫번째 설정 가능한 포인터에서 중지되므로 nil로 설정할 수 있다
// */
// func indirect(v reflect.Value, decodingNull bool) (Unmarshaler, encoding.TextUnmarshaler, reflect.Value) {
// 	// Issue #24153 indicates that it is generally not a guaranteed property
// 	// that you may round-trip a reflect.Value by calling Value.Addr().Elem()
// 	// and expect the value to still be settable for values derived from
// 	// unexported embedded struct fields.
// 	//
// 	// The logic below effectively does this when it first addresses the value
// 	// (to satisfy possible pointer methods) and continues to dereference
// 	// subsequent pointers as necessary.
// 	//
// 	// After the first round-trip, we set v back to the original value to
// 	// preserve the original RW flags contained in reflect.Value.
// 	v0 := v
// 	haveAddr := false
//
// 	/*
// 	` v reflect.Value `가 이름 있는 타입이고 메모리 주소를 지정할 수 있다면, 해당 주소에서 시작
// 	` reflect.Value.CanAddr `: ` reflect.Value.Addr() `로 메모리 주소 얻을 수 있는지 여부 확인
// 	*/
// 	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
// 		haveAddr = true
// 		v = v.Addr()  /* ` v reflect.Value `의 주소를 나타내는 포인터 값 반환 */
// 	}
// 	for {
// 		/*
// 		오직 ` v reflect.Value `의 타입이 ` Interface `이고 ` v reflect.Value `가 ` nil `이 아닌 경우, 인터페이스에서 값을 불러온다
// 		*/
// 		if v.Kind() == reflect.Interface && !v.IsNil() {
// 			e := v.Elem()  /* ` Interface v `가 포함하는 값 또는 ` Pointer v `가 가리키는 값 */
// 			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
// 				haveAddr = false
// 				v = e
// 				continue
// 			}
// 		}
//
// 		if v.Kind() != reflect.Ptr {
// 			break
// 		}
//
// 		if decodingNull && v.CanSet() {
// 			break
// 		}
//
// 		// Prevent infinite loop if v is an interface pointing to its own address:
// 		//     var v interface{}
// 		//     v = &v
// 		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
// 			v = v.Elem()
// 			break
// 		}
// 		if v.IsNil() {
// 			v.Set(reflect.New(v.Type().Elem()))
// 		}
// 		if v.Type().NumMethod() > 0 && v.CanInterface() {
// 			if u, ok := v.Interface().(Unmarshaler); ok {
// 				return u, nil, reflect.Value{}
// 			}
// 			if !decodingNull {
// 				if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
// 					return nil, u, reflect.Value{}
// 				}
// 			}
// 		}
//
// 		if haveAddr {
// 			v = v0 // restore original value after round-trip Value.Addr().Elem()
// 			haveAddr = false
// 		} else {
// 			v = v.Elem()
// 		}
// 	}
// 	return nil, nil, v
// }
