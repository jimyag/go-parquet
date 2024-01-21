# go-parquet

[![godoc for jimyag/go-parquet](https://godoc.org/github.com/nathany/looper?status.svg)](http://godoc.org/github.com/jimyag/go-parquet)

[English](./README.md)

go-parquet 是一个纯go实现的读写parquet格式文件 。

基于 [parquet-go](https://github.com/xitongsys/parquet-go) 并做了一些改进增强。

```sh
* 支持读写带有嵌套结构的Parquet文件
* 简单易用
* 高性能
```

## 安装

```sh
go get -u   github.com/jimyag/go-parquet
```

## 例子

`example/` 目录包含多个示例。

The `local_flat.go` example creates some data and writes it out to the `example/output/flat.parquet` file.

`local_flat.go`示例创建了一些数据，并将其写入`example/output/flat.parquet`文件。

```sh
cd example
go run local_flat.go
```

`local_flat.go`代码展示了如何轻松地将 Go 程序中的结构体写入到 Parquet 文件。

## Type

|INT64|int64|
|INT96([deprecated](https://github.com/apache/parquet-format/blob/master/CHANGES.md))|string|
|FLOAT|float32|
|DOUBLE|float64|
|BYTE_ARRAY|string|
|FIXED_LEN_BYTE_ARRAY|string|

### Logical Type

Parquet 原始的类型比较简单，为了能够和其他类型兼容，并且还能保持主类型的最小化。
Parquet 还设计了 Logical Types。例如 String 储存为 BYTE_ARRAY，但是需要是 UTF8 编码。
这样就需要一个 Logical Type 来做注解或者说扩充说明。

|逻辑类型|原始类型|Go Type|
|-|-|-|
|UTF8|BYTE_ARRAY|string|
|INT_8|INT32|int32|
|INT_16|INT32|int32|
|INT_32|INT32|int32|
|INT_64|INT64|int64|
|UINT_8|INT32|int32|
|UINT_16|INT32|int32|
|UINT_32|INT32|int32|
|UINT_64|INT64|int64|
|DATE|INT32|int32|
|TIME_MILLIS|INT32|int32|
|TIME_MICROS|INT64|int64|
|TIMESTAMP_MILLIS|INT64|int64|
|TIMESTAMP_MICROS|INT64|int64|
|INTERVAL|FIXED_LEN_BYTE_ARRAY|string|
|DECIMAL|INT32,INT64,FIXED_LEN_BYTE_ARRAY,BYTE_ARRAY|int32,int64,string,string|
|LIST|-|slice||
|MAP|-|map||

### 注意

- go-parquet 支持类型别名，如 `type MyString string`。但基本类型必须遵循表格说明。

- 一些类型转换函数: [converter.go](https://github.com/jimyag/go-parquet/blob/main/types/converter.go)

## 编码

#### PLAIN

所有类型

#### PLAIN_DICTIONARY/RLE_DICTIONARY

所有类型

#### DELTA_BINARY_PACKED

INT32, INT64, INT_8, INT_16, INT_32, INT_64, UINT_8, UINT_16, UINT_32, UINT_64, TIME_MILLIS, TIME_MICROS, TIMESTAMP_MILLIS, TIMESTAMP_MICROS

#### DELTA_BYTE_ARRAY

BYTE_ARRAY, UTF8

#### DELTA_LENGTH_BYTE_ARRAY

BYTE_ARRAY, UTF8

### Tips

- 有些平台不支持所有编码。如果不确定，请使用 PLAIN 和 PLAIN_DICTIONARY。
- 如果字段有许多不同的值，请不要使用 PLAIN_DICTIONARY 编码。因为它会将所有不同的值记录在一个映射中，这会占用大量内存。实际上，它使用 32 位整数来存储索引。如果您的主键数量大于 32 位，则不能使用该编码。
- 大的数组值可能会在页面统计中作为最小值和最大值重复出现，从而大大增加文件大小。如果统计信息对此类字段无用，可在字段标签中添加 `omitstats=true`，从而从写入的文件中省略统计信息。

## 重复类型

在 Parquet 文件中，字段可以是 REQUIRED（必须的）、OPTIONAL（可选的）或 REPEATED（重复的），这些都是重复类型的一部分。

|重复类型|例子 |说明|
|-|-|-|
|REQUIRED|```V1 int32 `parquet:"name=v1, type=INT32"` ```|No extra description|
|OPTIONAL|```V1 *int32 `parquet:"name=v1, type=INT32"` ```|Declare as pointer|
|REPEATED|```V1 []int32 `parquet:"name=v1, type=INT32, repetitiontype=REPEATED"` ```|添加'repetitiontype=REPEATED' 标签|

### 注意

- List 与 REPEATED 变量的区别在于标签中的 "重复类型"。虽然这两个变量在 go 中都存储为 slice，但在 parquet 中是不同的。您可以在 [此处](https://github.com/apache/parquet-format/blob/master/LogicalTypes.md) 中找到 parquet 中 List 的详细信息。我建议使用 List。
- 对于 LIST 和 MAP，一些存在的 parquet 文件使用一些非标准格式（见 [此处](https://github.com/apache/parquet-format/blob/master/LogicalTypes.md)）。对于标准格式，go-parquet 将把它们转换为 go slice 和 go map。对于非标准格式，go-parquet 将把它们转换为相应的结构体。

## 例子

```go
 Bool              bool    `parquet:"name=bool, type=BOOLEAN"`
 Int32             int32   `parquet:"name=int32, type=INT32"`
 Int64             int64   `parquet:"name=int64, type=INT64"`
 Int96             string  `parquet:"name=int96, type=INT96"`
 Float             float32 `parquet:"name=float, type=FLOAT"`
 Double            float64 `parquet:"name=double, type=DOUBLE"`
 ByteArray         string  `parquet:"name=bytearray, type=BYTE_ARRAY"`
 FixedLenByteArray string  `parquet:"name=FixedLenByteArray, type=FIXED_LEN_BYTE_ARRAY, length=10"`

 Utf8             string `parquet:"name=utf8, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
 Int_8            int32   `parquet:"name=int_8, type=INT32, convertedtype=INT32, convertedtype=INT_8"`
 Int_16           int32  `parquet:"name=int_16, type=INT32, convertedtype=INT_16"`
 Int_32           int32  `parquet:"name=int_32, type=INT32, convertedtype=INT_32"`
 Int_64           int64  `parquet:"name=int_64, type=INT64, convertedtype=INT_64"`
 Uint_8           int32  `parquet:"name=uint_8, type=INT32, convertedtype=UINT_8"`
 Uint_16          int32 `parquet:"name=uint_16, type=INT32, convertedtype=UINT_16"`
 Uint_32          int32 `parquet:"name=uint_32, type=INT32, convertedtype=UINT_32"`
 Uint_64          int64 `parquet:"name=uint_64, type=INT64, convertedtype=UINT_64"`
 Date             int32  `parquet:"name=date, type=INT32, convertedtype=DATE"`
 Date2            int32  `parquet:"name=date2, type=INT32, convertedtype=DATE, logicaltype=DATE"`
 TimeMillis       int32  `parquet:"name=timemillis, type=INT32, convertedtype=TIME_MILLIS"`
 TimeMillis2      int32  `parquet:"name=timemillis2, type=INT32, logicaltype=TIME, logicaltype.isadjustedtoutc=true, logicaltype.unit=MILLIS"`
 TimeMicros       int64  `parquet:"name=timemicros, type=INT64, convertedtype=TIME_MICROS"`
 TimeMicros2      int64  `parquet:"name=timemicros2, type=INT64, logicaltype=TIME, logicaltype.isadjustedtoutc=false, logicaltype.unit=MICROS"`
 TimestampMillis  int64  `parquet:"name=timestampmillis, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
 TimestampMillis2 int64  `parquet:"name=timestampmillis2, type=INT64, logicaltype=TIMESTAMP, logicaltype.isadjustedtoutc=true, logicaltype.unit=MILLIS"`
 TimestampMicros  int64  `parquet:"name=timestampmicros, type=INT64, convertedtype=TIMESTAMP_MICROS"`
 TimestampMicros2 int64  `parquet:"name=timestampmicros2, type=INT64, logicaltype=TIMESTAMP, logicaltype.isadjustedtoutc=false, logicaltype.unit=MICROS"`
 Interval         string `parquet:"name=interval, type=BYTE_ARRAY, convertedtype=INTERVAL"`

 Decimal1 int32  `parquet:"name=decimal1, type=INT32, convertedtype=DECIMAL, scale=2, precision=9"`
 Decimal2 int64  `parquet:"name=decimal2, type=INT64, convertedtype=DECIMAL, scale=2, precision=18"`
 Decimal3 string `parquet:"name=decimal3, type=FIXED_LEN_BYTE_ARRAY, convertedtype=DECIMAL, scale=2, precision=10, length=12"`
 Decimal4 string `parquet:"name=decimal4, type=BYTE_ARRAY, convertedtype=DECIMAL, scale=2, precision=20"`

 Decimal5 int32 `parquet:"name=decimal5, type=INT32, logicaltype=DECIMAL, logicaltype.precision=10, logicaltype.scale=2"`

 Map      map[string]int32 `parquet:"name=map, type=MAP, convertedtype=MAP, keytype=BYTE_ARRAY, keyconvertedtype=UTF8, valuetype=INT32"`
 List     []string         `parquet:"name=list, type=MAP, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8"`
 Repeated []int32          `parquet:"name=repeated, type=INT32, repetitiontype=REPEATED"`
```

## 压缩格式支持

|压缩算法|是否支持|
|-|-|
| CompressionCodec_UNCOMPRESSED | 是|
|CompressionCodec_SNAPPY|是|
|CompressionCodec_GZIP|是|
|CompressionCodec_LZO|否|
|CompressionCodec_BROTLI|否|
|CompressionCodec_LZ4 |是|
|CompressionCodec_ZSTD|是|

## ParquetFile

读写 parquet 文件需要实现 ParquetFile 接口

```go
type ParquetFile interface {
 io.Seeker
 io.Reader
 io.Writer
 io.Closer
 Open(name string) (ParquetFile, error)
 Create(name string) (ParquetFile, error)
}
```

使用该接口，go-parquet 可以在不同数据源上读写 parquet 文件。目前支持数据源都位于 [source]（./source/）。现在，它支持（本地/hdfs/s3/gcs/内存）。

## 写入

支持四种写入器： ParquetWriter、JSONWriter、CSVWriter 和 ArrowWriter。

- ParquetWriter 用于编写预定义的 Golang 结构体。
[Example of ParquetWriter](https://github.com/jimyag/go-parquet/blob/main/example/local_flat.go)

- JSONWriter 用于写入 JSON 字符串
[Example of JSONWriter](https://github.com/jimyag/go-parquet/blob/main/example/json_write.go)

- CSVWriter 用于写入与 CSV 格式类似的数据（非嵌套）。
[Example of CSVWriter](https://github.com/jimyag/go-parquet/blob/main/example/csv_write.go)

- ArrowWriter 用于使用 Arrow 模式写入镶嵌文件
[Example of ArrowWriter](https://github.com/jimyag/go-parquet/blob/main/example/arrow_to_parquet.go)

## 读取

支持两种读取器: ParquetReader, ColumnReader

- ParquetReader 用于读取预定义的 Golang 结构体
[Example of ParquetReader](https://github.com/jimyag/go-parquet/blob/main/example/local_nested.go)

- ColumnReader 用于读取原始列数据。读取函数返回记录的 3 个片段（[value], [RepetitionLevel], [DefinitionLevel]）。
[Example of ColumnReader](https://github.com/jimyag/go-parquet/blob/main/example/column_read.go)

### 注意

- 如果镶嵌文件非常大（即使镶嵌文件很小，未压缩的大小也可能很大），请不要一次读取所有行，否则可能会导致 OOM。您可以像读取面向流的文件一样，一次读取一小部分数据。

- `RowGroupSize` and `PageSize` may influence the final parquet file size. You can find the details from [here](https://github.com/apache/parquet-format). You can reset them in ParquetWriter

- `RowGroupSize`和 `PageSize`可能会影响最终拼版文件的大小。详细信息请参阅 [此处](https://github.com/apache/parquet-format)。你可以在 ParquetWriter 中重新设置它们

```go
 pw.RowGroupSize = 128 * 1024 * 1024 // default 128M
 pw.PageSize = 8 * 1024 // default 8K
```

## Schema

有四种方法可以定义模式：结构体的标签、Json、CSV、Arrow metadata。只有Schema中的项目才会被写入，其他项目将被忽略。

### Tag

```golang
type Student struct {
 Name    string  `parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
 Age     int32   `parquet:"name=age, type=INT32, encoding=PLAIN"`
 Id      int64   `parquet:"name=id, type=INT64"`
 Weight  float32 `parquet:"name=weight, type=FLOAT"`
 Sex     bool    `parquet:"name=sex, type=BOOLEAN"`
 Day     int32   `parquet:"name=day, type=INT32, convertedtype=DATE"`
 Ignored int32   //without parquet tag and won't write
}
```

[Example of tags](https://github.com/jimyag/go-parquet/blob/main/example/local_flat.go)

### JSON

JSON 模式可用于定义一些复杂的模式，这些模式无法通过标签来定义。

```golang
type Student struct {
 NameIn    string
 Age     int32
 Id      int64
 Weight  float32
 Sex     bool
 Classes []string
 Scores  map[string][]float32
 Ignored string

 Friends []struct {
  Name string
  Id   int64
 }
 Teachers []struct {
  Name string
  Id   int64
 }
}

var jsonSchema string = `
{
  "Tag": "name=parquet_go_root, repetitiontype=REQUIRED",
  "Fields": [
    {"Tag": "name=name, inname=NameIn, type=BYTE_ARRAY, convertedtype=UTF8, repetitiontype=REQUIRED"},
    {"Tag": "name=age, inname=Age, type=INT32, repetitiontype=REQUIRED"},
    {"Tag": "name=id, inname=Id, type=INT64, repetitiontype=REQUIRED"},
    {"Tag": "name=weight, inname=Weight, type=FLOAT, repetitiontype=REQUIRED"},
    {"Tag": "name=sex, inname=Sex, type=BOOLEAN, repetitiontype=REQUIRED"},

    {"Tag": "name=classes, inname=Classes, type=LIST, repetitiontype=REQUIRED",
     "Fields": [{"Tag": "name=element, type=BYTE_ARRAY, convertedtype=UTF8, repetitiontype=REQUIRED"}]
    },

    {
      "Tag": "name=scores, inname=Scores, type=MAP, repetitiontype=REQUIRED",
      "Fields": [
        {"Tag": "name=key, type=BYTE_ARRAY, convertedtype=UTF8, repetitiontype=REQUIRED"},
        {"Tag": "name=value, type=LIST, repetitiontype=REQUIRED",
         "Fields": [{"Tag": "name=element, type=FLOAT, repetitiontype=REQUIRED"}]
        }
      ]
    },

    {
      "Tag": "name=friends, inname=Friends, type=LIST, repetitiontype=REQUIRED",
      "Fields": [
       {"Tag": "name=element, repetitiontype=REQUIRED",
        "Fields": [
         {"Tag": "name=name, inname=Name, type=BYTE_ARRAY, convertedtype=UTF8, repetitiontype=REQUIRED"},
         {"Tag": "name=id, inname=Id, type=INT64, repetitiontype=REQUIRED"}
        ]}
      ]
    },

    {
      "Tag": "name=teachers, inname=Teachers, repetitiontype=REPEATED",
      "Fields": [
        {"Tag": "name=name, inname=Name, type=BYTE_ARRAY, convertedtype=UTF8, repetitiontype=REQUIRED"},
        {"Tag": "name=id, inname=Id, type=INT64, repetitiontype=REQUIRED"}
      ]
    }
  ]
}
`
```

[Example of JSON schema](https://github.com/jimyag/go-parquet/blob/main/example/json_schema.go)

### CSV metadata

```golang
 md := []string{
  "name=Name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY",
  "name=Age, type=INT32",
  "name=Id, type=INT64",
  "name=Weight, type=FLOAT",
  "name=Sex, type=BOOLEAN",
 }
```

[Example of CSV metadata](https://github.com/jimyag/go-parquet/blob/main/example/csv_write.go)

### Arrow metadata

```golang
 schema := arrow.NewSchema(
  []arrow.Field{
   {Name: "int64", Type: arrow.PrimitiveTypes.Int64},
   {Name: "float64", Type: arrow.PrimitiveTypes.Float64},
   {Name: "str", Type: arrow.BinaryTypes.String},
  },
  nil,
 )
```

[Example of Arrow metadata](https://github.com/jimyag/go-parquet/blob/main/example/arrow_to_parquet.go)

### 注意

- 在 Golang 中，go-parquet 以对象的形式读取数据，每个字段都必须是以大写字母开头的公共字段。这个字段名称为 `InName`。在 parquet 文件中的字段名称为 `ExName`。函数 `common.HeadToUpper` 将 `ExName` 转换为 `InName`。有一些限制：

1. 如果两个字段名仅在第一个字母大小写上不同，则不允许使用。如 `name` 和 `Name`。
2. `PARGO_PREFIX_` 是一个保留字符串，最好不要将其用作名称前缀。([#294](https://github.com/xitongsys/parquet-go/issues/294))
3. 使用 `\x01` 作为字段的分隔符，以支持某些字段名中的 `.`。([dot_in_name.go](https://github.com/jimyag/go-parquet/blob/main/example/dot_in_name.go), [#349](https://github.com/xitongsys/parquet-go/issues/349))

## 并发

在写入/读取过程中，Marshal/Unmarshal 是最耗时的过程。为了提高性能，go-parquet 可以使用多个 goroutines 来 Marshal/unmarshal 对象。您可以在读取/写入初始函数中设置并发数参数 `np`。

```golang
func NewParquetReader(pFile ParquetFile.ParquetFile, obj interface{}, np int64) (*ParquetReader, error)
func NewParquetWriter(pFile ParquetFile.ParquetFile, obj interface{}, np int64) (*ParquetWriter, error)
func NewJSONWriter(jsonSchema string, pfile ParquetFile.ParquetFile, np int64) (*JSONWriter, error)
func NewCSVWriter(md []string, pfile ParquetFile.ParquetFile, np int64) (*CSVWriter, error)
func NewArrowWriter(arrowSchema *arrow.Schema, pfile source.ParquetFile, np int64) (*ArrowWriter error)
```

## Examples

|Example file|Descriptions|
|-|-|
|[local_flat.go](https://github.com/jimyag/go-parquet/blob/main/example/local_flat.go)|读写无嵌套结构的parquet文件|
|[local_nested.go](https://github.com/jimyag/go-parquet/blob/main/example/local_nested.go)|读写嵌套结构的parquet文件|
|[read_partial.go](https://github.com/jimyag/go-parquet/blob/main/example/read_partial.go)|从 parquet 文件中读取部分字段|
|[read_partial2.go](https://github.com/jimyag/go-parquet/blob/main/example/read_partial2.go)|从 parquet 文件中读取部分字段|
|[read_without_schema_predefined.go](https://github.com/jimyag/go-parquet/blob/main/example/read_without_schema_predefined.go)|读取parquet，无需预定义结构/模式|
|[read_partial_without_schema_predefined.go](https://github.com/jimyag/go-parquet/blob/main/example/read_partial_without_schema_predefined.go)|从 parquet 文件读取子字段，无需预定义结构/模式|
|[json_schema.go](https://github.com/jimyag/go-parquet/blob/main/example/json_schema.go)|定义 schema 使用 json |
|[json_write.go](https://github.com/jimyag/go-parquet/blob/main/example/json_write.go)|将json转为parquet|
|[convert_to_json.go](https://github.com/jimyag/go-parquet/blob/main/example/convert_to_json.go)|将parquet转为json|
|[csv_write.go](https://github.com/jimyag/go-parquet/blob/main/example/csv_write.go)|csv 写入|
|[column_read.go](https://github.com/jimyag/go-parquet/blob/main/example/column_read.go)|读取原始列数据并返回值、重复级别、定义级别|
|[type.go](https://github.com/jimyag/go-parquet/blob/main/example/type.go)|类型示例|
|[type_alias.go](https://github.com/jimyag/go-parquet/blob/main/example/type_alias.go)|类型别名示例|
|[writer.go](https://github.com/jimyag/go-parquet/blob/main/example/writer.go)|从 io.Writer 写入|
|[keyvalue_metadata.go](https://github.com/jimyag/go-parquet/blob/main/example/keyvalue_metadata.go)|写入/读取 key/value metadata|
|[dot_in_name.go](https://github.com/jimyag/go-parquet/blob/main/example/dot_in_name.go)|字段中包含`.`|
|[arrow_to_parquet.go](https://github.com/jimyag/go-parquet/blob/main/example/arrow_to_parquet.go)|使用 arrow 定义读写parquet文件|

## 工具

- [parquet-tools](https://github.com/jimyag/parquet-tools): 帮助检查 Parquet 文件的命令行工具
