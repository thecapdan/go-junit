[![License][license-badge]][license-link]
[![Godoc][godoc-badge]][godoc-link]

# Go JUnit

This library is based on this [library from joshdk](https://github.com/joshdk/go-junit)
It has been changed so that junit.Ingest now takes a io.Reader interface. This means that the whole []slice of xml doesn't need
to be allocated and passed around and can be buffered by the xml unmarshaller.

🐜 Go library for ingesting JUnit XML reports

## Installing

You can fetch this library by running the following

```bash
go get -u github.com/goddenrich/go-junit
```

## Usage

### Data Ingestion

This library has a number of ingestion methods for convenience.

The simplest of which parses an io.Reader with junit XML data.

```go
xml := []byte(`
    <?xml version="1.0" encoding="UTF-8"?>
    <testsuites>
        <testsuite name="JUnitXmlReporter" errors="0" tests="0" failures="0" time="0" timestamp="2013-05-24T10:23:58" />
        <testsuite name="JUnitXmlReporter.constructor" errors="0" skipped="1" tests="3" failures="1" time="0.006" timestamp="2013-05-24T10:23:58">
            <properties>
                <property name="java.vendor" value="Sun Microsystems Inc." />
                <property name="compiler.debug" value="on" />
                <property name="project.jdk.classpath" value="jdk.classpath.1.6" />
            </properties>
            <testcase classname="JUnitXmlReporter.constructor" name="should default path to an empty string" time="0.006">
                <failure message="test failure">Assertion failed</failure>
            </testcase>
            <testcase classname="JUnitXmlReporter.constructor" name="should default consolidate to true" time="0">
                <skipped />
            </testcase>
            <testcase classname="JUnitXmlReporter.constructor" name="should default useDotNotation to true" time="0" />
        </testsuite>
    </testsuites>
`)

suites, err := junit.Ingest(bytes.NewReader(xml))
if err != nil {
    log.Fatalf("failed to ingest JUnit xml %v", err)
}
```

You can then inspect the contents of the ingested suites.

```go
for _, suite := range suites {
    fmt.Println(suite.Name)
    for _, test := range suite.Tests {
        fmt.Printf("  %s\n", test.Name)
        if test.Error != nil {
            fmt.Printf("    %s: %s\n", test.Status, test.Error.Error())
        } else {
            fmt.Printf("    %s\n", test.Status)
        }
    }
}
```

And observe some output like this.

```
JUnitXmlReporter
JUnitXmlReporter.constructor
  should default path to an empty string
    failed: Assertion failed
  should default consolidate to true
    skipped
  should default useDotNotation to true
    passed
```

### More Examples

As a file implements the io.Reader interface it can be used

```go
fileName = "test-reports/report.xml"
f, err := os.Open(fileName)
if err != nil {
    log.Fatalf("failed to open file %s", fileName)
}
suites, err := junit.Ingest("test-reports/report.xml")
if err != nil {
    log.Fatalf("failed to ingest JUnit xml %v", err)
}
```

You can use IngestFiles to ingest multiple files

```go
suites, err := junit.IngestFiles([]string{
    "test-reports/report-1.xml",
    "test-reports/report-2.xml",
})
if err != nil {
    log.Fatalf("failed to ingest JUnit xml %v", err)
}
```

Or any `.xml` files inside of a directory.

```go
suites, err := junit.IngestDir("test-reports/")
if err != nil {
    log.Fatalf("failed to ingest JUnit xml %v", err)
}
```

### Data Formats

Contrary to the library this was based on the format of the xml must have a parent node and can not have a single or multiple top level `testsuite` tag

A single top level `testsuites` tag, containing multiple `testsuite` instances.

```xml
<testsuites>
    <testsuite>
        <testcase name="Test case 1" />
        <testcase name="Test case 2" />
    </testsuite>
</testsuites>
```
In all cases, omitting (or even duplicated) the XML declaration tag is allowed.

```xml
<?xml version="1.0" encoding="UTF-8"?>
```

These formats are invalid and will throw an io.EOF error

A single top level `testsuite` tag, containing multiple `testcase` instances.

```xml
<testsuite>
    <testcase name="Test case 1" />
    <testcase name="Test case 2" />
</testsuite>
```

Multiple top level `testsuite` tags, containing multiple `testcase` instances.

```xml
<testsuite>
    <testcase name="Test case 1" />
    <testcase name="Test case 2" />
</testsuite>
<testsuite>
    <testcase name="Test case 3" />
    <testcase name="Test case 4" />
</testsuite>
```


## License

This code is distributed under the [MIT License][license-link], see [LICENSE.txt][license-file] for more information.

[circleci-badge]:   https://circleci.com/gh/goddenrich/go-junit.svg?&style=shield
[circleci-link]:    https://circleci.com/gh/goddenrich/go-junit/tree/master
[go-report-badge]:  https://goreportcard.com/badge/github.com/goddenrich/go-junit
[go-report-link]:   https://goreportcard.com/report/github.com/goddenrich/go-junit
[godoc-badge]:      https://godoc.org/github.com/goddenrich/go-junit?status.svg
[godoc-link]:       https://godoc.org/github.com/goddenrich/go-junit
[license-badge]:    https://img.shields.io/badge/license-MIT-green.svg
[license-file]:     https://github.com/goddenrich/go-junit/blob/master/LICENSE.txt
[license-link]:     https://opensource.org/licenses/MIT
