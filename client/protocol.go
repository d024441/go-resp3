// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strconv"

	"github.com/stfnmllr/go-resp3/client/internal/conv"
	"github.com/stfnmllr/go-resp3/client/internal/monitor"
)

const (
	protocolVersion = 3
)

// Redis reply constants.
const (
	ReplyOK     = "OK"
	ReplyQueued = "QUEUED"
)

const (
	blobStringType             = '$'
	simpleStringType           = '+'
	simpleErrorType            = '-'
	numberType                 = ':'
	nullType                   = '_'
	doubleType                 = ','
	booleanType                = '#'
	blobErrorType              = '!'
	verbatimStringType         = '='
	bigNumberType              = '('
	arrayType                  = '*'
	mapType                    = '%'
	setType                    = '~'
	attributeType              = '|'
	pushType                   = '>'
	streamedType               = '?'
	streamedStringToken        = ';'
	streamedDataTypeTerminator = '.'

	booleanTrue  = 't'
	booleanFalse = 'f'

	lf = '\r'
	nl = '\n'

	lineBreak = string(lf) + string(nl)
)

const maxInt = int64(^uint(0) >> 1)

// A InvalidValueError is raised by a redis commend
// provided with an invalid parameter value.
// - Name:  Parameter name.
// - Value: Invalid value.
type InvalidValueError struct {
	Name  string
	Value interface{}
}

func newInvalidValueError(name string, value interface{}) *InvalidValueError {
	return &InvalidValueError{Name: name, Value: value}
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("Invalid value %v for %s", e.Value, e.Name)
}

// An InvalidTypeError is raised by encoding
// an unsupported value type.
type InvalidTypeError struct {
	Value interface{}
}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("encode: unsupported type %[1]T value %[1]v", e.Value)
}

// An UnexpectedCharacterError is raised by decoding
// an unexpected character.
type UnexpectedCharacterError struct {
	ActChar byte
	ExpChar byte
}

func (e *UnexpectedCharacterError) Error() string {
	return fmt.Sprintf("decode: unsupported character %q expected %q"+string(e.ExpChar), e.ActChar, e.ExpChar)
}

// An InvalidNumberError is raised by decoding
// a number including unexpected characters.
type InvalidNumberError struct {
	Value string
}

func (e *InvalidNumberError) Error() string {
	return fmt.Sprintf("decode: invalid number %s", e.Value)
}

// An InvalidDoubleError is raised by decoding
// a double including unexpected characters.
type InvalidDoubleError struct {
	Value string
}

func (e *InvalidDoubleError) Error() string {
	return fmt.Sprintf("decode: invalid double %s", e.Value)
}

// An InvalidBigNumberError is raised by decoding
// a bignumber including unexpected characters.
type InvalidBigNumberError struct {
	Value string
}

func (e *InvalidBigNumberError) Error() string {
	return fmt.Sprintf("decode: invalid bignumber %s", e.Value)
}

// StringPtr returns a pointer to the string parameter.
func StringPtr(s string) *string {
	return &s
}

// Int64Ptr returns a pointer to the int64 parameter.
func Int64Ptr(i int64) *int64 {
	return &i
}

// check if objects implement interfaces
var _ Encoder = (*encode)(nil)
var _ Decoder = (*decode)(nil)

// Encoder is the interface that wraps the Encode method.
//
// Encode is not a general redis encoder but encodes redis commands only.
// Redis commands are send as ARRAY of BULK STRINGS.
// Encode is encoding all elements of v to BULK STRINGS.
// In case an element is a not suported type, Encode will panic.
type Encoder interface {
	Encode([]interface{}) error
	Flush() error
}

// NewEncoder returns an Encoder for redis commands.
func NewEncoder(w io.Writer) Encoder {
	return newEncode(w)
}

type encode struct {
	w   *bufio.Writer
	err error
}

func newEncode(w io.Writer) *encode {
	//return &encode{w: bufio.NewWriterSize(w, 4096)}
	//return &encode{w: bufio.NewWriterSize(w, 8192)}
	//return &encode{w: bufio.NewWriterSize(w, 16384)}
	return &encode{w: bufio.NewWriterSize(w, 32768)}
}

func (e *encode) Encode(values []interface{}) error {
	e.w.WriteByte(arrayType)
	e.w.WriteString(strconv.Itoa(len(values)))
	e.w.WriteString(lineBreak)
	for _, v := range values {
		err := e.encode(v)
		if err != nil {
			e.err = err
			return err
		}
	}
	return nil
}

func (e *encode) Flush() error {
	if e.err != nil {
		return e.err
	}
	//println(e.w.Buffered())
	return e.w.Flush()
}

func (e *encode) encode(v interface{}) error {
	switch v := v.(type) {
	case string:
		e.encodeString(v)
	case []byte:
		e.encodeString(string(v))
	case bool:
		e.encodeString(strconv.FormatBool(v))
	case int:
		e.encodeString(strconv.FormatInt(int64(v), 10))
	case uint:
		e.encodeString(strconv.FormatUint(uint64(v), 10))
	case float32:
		e.encodeString(strconv.FormatFloat(float64(v), 'g', -1, 64))
	case float64:
		e.encodeString(strconv.FormatFloat(v, 'g', -1, 64))
	case Zfloat64:
		e.encodeString(v.String())
	default:
		return e.encodeValue(reflect.ValueOf(v))
	}
	return nil
}

func (e *encode) encodeValue(v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		e.encodeString(v.String())
	case reflect.Bool:
		e.encodeString(strconv.FormatBool(v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		e.encodeString(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		e.encodeString(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		e.encodeString(strconv.FormatFloat(v.Float(), 'g', -1, 64))
	case reflect.Interface:
		e.encodeValue(v.Elem())
	case reflect.Ptr:
		e.encodeValue(v.Elem())
	default:
		return &InvalidTypeError{v}
	}
	return nil
}

func (e *encode) encodeString(s string) {
	e.w.WriteByte(blobStringType)
	e.w.WriteString(strconv.Itoa(len(s)))
	e.w.WriteString(lineBreak)
	e.w.WriteString(s)
	e.w.WriteString(lineBreak)
}

type decodeReader struct {
	r   *bufio.Reader
	buf []byte
}

func newDecodeReader(r io.Reader) *decodeReader {
	return &decodeReader{
		r:   bufio.NewReaderSize(r, 32768),
		buf: make([]byte, 0, 128),
	}
}

func (r *decodeReader) readType() (byte, error) {
	return r.r.ReadByte()
}

func (r *decodeReader) peek() (byte, error) {
	b, err := r.r.Peek(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (r *decodeReader) discardByte(db byte) error {
	b, err := r.r.ReadByte()
	if err != nil {
		return err
	}
	if b != db {
		return &UnexpectedCharacterError{ActChar: b, ExpChar: db}
	}
	return nil
}

func (r *decodeReader) readLineBreak() error {
	if err := r.discardByte(lf); err != nil {
		return err
	}
	if err := r.discardByte(nl); err != nil {
		return err
	}
	return nil
}

func (r *decodeReader) readByte() (byte, error) {
	b, err := r.readBlob(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (r *decodeReader) resize(size64 int64) {
	if size64 > maxInt {
		panic("maximum integer size exceeded")
	}
	size := int(size64)
	if cap(r.buf) < size {
		r.buf = make([]byte, size)
	}
	r.buf = r.buf[:size]
}

func (r *decodeReader) readBlob(size int64) ([]byte, error) {
	r.resize(size + 2)
	if _, err := io.ReadFull(r.r, r.buf); err != nil {
		return nil, err
	}
	if r.buf[size] != lf {
		return nil, &UnexpectedCharacterError{ActChar: r.buf[size], ExpChar: lf}
	}
	if r.buf[size+1] != nl {
		return nil, &UnexpectedCharacterError{ActChar: r.buf[size+1], ExpChar: nl}
	}
	return r.buf[:size], nil
}

// replace bufio.Reader readBytes, readString as they do allocate
func (r *decodeReader) readDelim(delim byte) ([]byte, error) {
	r.buf = r.buf[:0]

	for {
		frag, err := r.r.ReadSlice(delim)
		if err != nil && err != bufio.ErrBufferFull { // unexpected error
			return nil, err
		}
		r.buf = append(r.buf, frag...)
		if err == nil { // got final fragment
			break
		}
	}
	return r.buf, nil
}

func (r *decodeReader) readBytes() ([]byte, error) {
	b, err := r.readDelim(lf)
	if err != nil {
		return nil, err
	}
	if err := r.discardByte(nl); err != nil {
		return nil, err
	}
	return b[:len(b)-1], nil
}

func (r *decodeReader) readFixedSize() (int64, error) {
	b, err := r.readBytes()
	if err != nil {
		return 0, err
	}
	i, err := conv.ParseInt(b)
	if err != nil {
		return 0, &InvalidNumberError{Value: string(b)}
	}
	return i, nil
}

func (r *decodeReader) readSize() (int64, error) {
	b, err := r.readBytes()
	if err != nil {
		return 0, err
	}
	if len(b) == 1 && b[0] == streamedType {
		return -1, nil
	}
	size, err := conv.ParseInt(b)
	if err != nil {
		return 0, &InvalidNumberError{Value: string(b)}
	}
	return size, nil
}

// Decoder is the interface that wraps the Decode method.
type Decoder interface {
	Decode() (interface{}, error)
}

// NewDecoder returns a Decoder for Redis results.
func NewDecoder(r io.Reader) Decoder {
	return &decode{
		r:   newDecodeReader(r),
		buf: make([]byte, 0, 128),
	}
}

type decode struct {
	r   *decodeReader
	buf []byte
}

func (d *decode) Decode() (interface{}, error) {
	t, err := d.r.peek()
	if err != nil {
		return nil, err
	}

	switch t {

	case pushType, simpleErrorType, blobErrorType:
		return d.decodeTopLevel()
	case simpleStringType, blobStringType:
		return d.decodeMonitorNotification() // can be dropped after monitor notification is send as notification
	default:
		return d.decode()
	}
}

func (d *decode) decodeMonitorNotification() (interface{}, error) {
	t, err := d.r.readType()
	if err != nil {
		return nil, err
	}
	switch t {
	case simpleStringType:
		return d.decodeSimpleStringOrMonitor()
	case blobStringType:
		return d.decodeBlobStringOrMonitor()
	default:
		panic("wring data type")
	}
}

func (d *decode) decodeBlobStringOrMonitor() (interface{}, error) {
	b, err := d.decodeBlobStringType()
	if err != nil {
		return nil, err
	}
	if n, ok := monitor.Parse(b); ok {
		return n, nil
	}
	return _string(b), nil
}

func (d *decode) decodeSimpleStringOrMonitor() (interface{}, error) {
	b, err := d.r.readBytes()
	if err != nil {
		return nil, err
	}
	if n, ok := monitor.Parse(b); ok {
		return n, nil
	}
	return _string(b), nil
}

func (d *decode) decodeTopLevel() (interface{}, error) {
	t, err := d.r.readType()
	if err != nil {
		return nil, err
	}
	switch t {
	case pushType:
		return d.decodePushSlice()
	case blobErrorType:
		return d.decodeBlobError()
	case simpleErrorType:
		return d.decodeSimpleError()
	default:
		panic("wring data type")
	}
}

func (d *decode) decode() (RedisValue, error) {
	t, err := d.r.peek()
	if err != nil {
		return nil, err
	}

	if t == attributeType {
		return d.decodeAttr()
	}
	return d.decodeType()
}

func (d *decode) decodeAttr() (RedisValue, error) {
	_, err := d.r.readType()
	if err != nil {
		return nil, err
	}

	attr, err := d.decodeMap()
	if err != nil {
		return nil, err
	}

	v, err := d.decodeType()
	if err != nil {
		return nil, err
	}

	return attrRedisValue{RedisValue: v.(RedisValue), attr: attr}, nil
}

func (d *decode) decodeType() (RedisValue, error) {

	t, err := d.r.readType()
	if err != nil {
		return nil, err
	}

	switch t {
	case nullType:
		err := d.r.readLineBreak()
		return _Null, err
	case blobStringType:
		return d.decodeBlobString()
	case verbatimStringType:
		return d.decodeVerbatimString()
	case simpleStringType:
		return d.decodeSimpleString()
	case numberType:
		return d.decodeNumber()
	case doubleType:
		return d.decodeDouble()
	case bigNumberType:
		return d.decodeBigNumber()
	case booleanType:
		return d.decodeBoolean()
	case arrayType:
		return d.decodeSlice()
	case mapType:
		return d.decodeMap()
	case setType:
		return d.decodeSet()
	default:
		panic("unsupported redis type " + string(t))
	}
}

//
func (d *decode) decodeBlobStringType() ([]byte, error) {
	size, err := d.r.readSize()
	if err != nil {
		return nil, err
	}

	var b []byte

	if size == -1 {
		b, err = d.decodeStreamedString()
	} else {
		b, err = d.r.readBlob(size)
	}
	return b, err
}

// Error
func (d *decode) decodeSimpleError() (error, error) {
	b, err := d.r.readBytes()
	if err != nil {
		return nil, err
	}
	return newRedisError(string(b)), nil
}

func (d *decode) decodeBlobError() (error, error) {
	b, err := d.decodeBlobStringType()
	if err != nil {
		return nil, err
	}
	return newRedisError(string(b)), err
}

// String
func (d *decode) decodeBlobString() (RedisValue, error) {
	b, err := d.decodeBlobStringType()
	if err != nil {
		return nil, err
	}
	return _string(b), nil
}

func (d *decode) decodeVerbatimString() (RedisValue, error) {
	b, err := d.decodeBlobStringType()
	if err != nil {
		return nil, err
	}
	return _verbatimString(string(b)), err
}

func (d *decode) decodeStreamedString() ([]byte, error) {
	d.buf = d.buf[:0]

	for {
		t, err := d.r.readType()
		if err != nil {
			return nil, err
		}
		if t != streamedStringToken {
			return nil, &UnexpectedCharacterError{ActChar: t, ExpChar: streamedStringToken}
		}
		l, err := d.r.readFixedSize()
		if err != nil {
			return nil, err
		}
		if l == 0 {
			break
		}

		part, err := d.r.readBlob(l)
		if err != nil {
			return nil, err
		}

		d.buf = append(d.buf, part...)
	}
	return d.buf, nil
}

func (d *decode) decodeSimpleString() (RedisValue, error) {
	b, err := d.r.readBytes()
	if err != nil {
		return nil, err
	}
	return _string(b), nil
}

// Number
func (d *decode) decodeNumber() (RedisValue, error) {
	i, err := d.r.readFixedSize()
	return _number(i), err
}

// Double
func (d *decode) decodeDouble() (RedisValue, error) {
	b, err := d.r.readBytes()
	if err != nil {
		return nil, err
	}
	f, err := strconv.ParseFloat(string(b), 64) // supports inf, -inf
	if err != nil {
		return nil, &InvalidDoubleError{Value: string(b)}
	}
	return _double(f), nil
}

// Bignumber
func (d *decode) decodeBigNumber() (RedisValue, error) {
	b, err := d.r.readBytes()
	if err != nil {
		return nil, err
	}
	i, ok := new(big.Int).SetString(string(b), 10)
	if !ok {
		return nil, &InvalidBigNumberError{Value: string(b)}
	}
	return (*_bignumber)(i), nil
}

// Boolean
func (d *decode) decodeBoolean() (RedisValue, error) {
	b, err := d.r.readByte()
	if err != nil {
		return nil, err
	}
	return _boolean(b == booleanTrue), nil
}

// Push Notification
func (d *decode) decodePushSlice() (interface{}, error) {
	slice, err := d.decodeSlice()
	if err != nil {
		return nil, err
	}
	return newNotification(slice)
}

// Slice
func (d *decode) decodeSlice() (_slice, error) {
	var v _slice

	size, err := d.r.readSize()
	if err != nil {
		return v, err
	}

	if size == -1 {
		v, err = d.decodeStreamedSlice()
	} else {
		v, err = d.decodeFixedSlice(size)
	}
	return v, err
}

func (d *decode) decodeFixedSlice(size int64) (_slice, error) {
	s := make(_slice, size)
	for i := int64(0); i < size; i++ {
		val, err := d.decode()
		if err != nil {
			return s, err
		}
		s[i] = val
	}
	return s, nil
}

func (d *decode) decodeStreamedSlice() (_slice, error) {
	s := make(_slice, 0)
	for {
		t, err := d.r.peek()
		if err != nil {
			return s, err
		}
		if t == streamedDataTypeTerminator {
			break
		}
		val, err := d.decode()
		if err != nil {
			return s, err
		}
		s = append(s, val)
	}
	if _, err := d.r.readByte(); err != nil {
		return s, err
	}
	return s, nil
}

// Map
func (d *decode) decodeMap() (_map, error) {
	var v _map

	size, err := d.r.readSize()
	if err != nil {
		return v, err
	}

	if size == -1 {
		v, err = d.decodeStreamedMap()
	} else {
		v, err = d.decodeFixedMap(size)
	}
	return v, err
}

func (d *decode) decodeFixedMap(size int64) (_map, error) {
	m := make(_map, size)
	for i := int64(0); i < size; i++ {
		key, err := d.decode()
		if err != nil {
			return m, err
		}
		val, err := d.decode()
		if err != nil {
			return m, err
		}
		m[i] = MapItem{key, val}
	}
	return m, nil
}

func (d *decode) decodeStreamedMap() (_map, error) {
	m := make(_map, 0)
	for {
		t, err := d.r.peek()
		if err != nil {
			return m, err
		}
		if t == streamedDataTypeTerminator {
			break
		}
		key, err := d.decode()
		if err != nil {
			return m, err
		}
		val, err := d.decode()
		if err != nil {
			return m, err
		}
		m = append(m, MapItem{key, val})
	}
	if _, err := d.r.readByte(); err != nil {
		return m, err
	}
	return m, nil
}

// Set
func (d *decode) decodeSet() (_set, error) {
	var v _set

	size, err := d.r.readSize()
	if err != nil {
		return v, err
	}

	if size == -1 {
		v, err = d.decodeStreamedSet()
	} else {
		v, err = d.decodeFixedSet(size)
	}
	return v, err
}

func (d *decode) decodeFixedSet(size int64) (_set, error) {
	s := make(_set, size)
	for i := int64(0); i < size; i++ {
		val, err := d.decode()
		if err != nil {
			return s, err
		}
		s[i] = val
	}
	return s, nil
}

func (d *decode) decodeStreamedSet() (_set, error) {
	s := make(_set, 0)
	for {
		t, err := d.r.peek()
		if err != nil {
			return s, err
		}
		if t == streamedDataTypeTerminator {
			break
		}
		val, err := d.decode()
		if err != nil {
			return s, err
		}
		s = append(s, val)
	}
	if _, err := d.r.readByte(); err != nil {
		return s, err
	}
	return s, nil
}
