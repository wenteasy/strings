package strings_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/wenteasy/strings"
)

func TestTable(t *testing.T) {
	tb := strings.NewTable()

	tb.Columns(
		tb.Column("", 8, false),
		tb.Column("", 5, true),
		tb.Column("", 10, false),
	)

	tb.Add("AAA", "BBB", "CCCCCCCCCCCCCC")
	tb.Add("DDDDD", "E", "FFFFF")
	tb.Add("GGGGGGGG", "HHHHHHH", "I")
	tb.Add("マルチバイト文字にも", "一応", "対応しています。")

	_, err := tb.Generate()
	if err != nil {
		t.Errorf("Table.Generate() can not be an error: %+v", err)
		return
	}
}

func TestTableWriterRealTime(t *testing.T) {
	tb := strings.NewTable()
	tb.Columns(
		tb.Column("Col1", 8, false),
		tb.Column("Col2", 6, true),
		tb.Column("Col3", 10, false),
	)

	w, err := strings.NewTableWriter(io.Discard, tb)
	if err != nil {
		t.Errorf("NewTableWriter() can not be an error: %+v", err)
	}

	w.WriteString("AAA")
	w.WriteString("BBB")
	w.WriteString("CCCCCCCCCCCCCCCC")
	w.WriteString("DDDDD")
	w.WriteString("E")
	w.WriteString("FFFFF")

	err = w.Close()
	if err != nil {
		t.Errorf("NewTableWriter() can not be an error: %+v", err)
	}
}

func TestTableWriterLazy(t *testing.T) {

	tb := strings.NewTable()
	tb.Columns(
		tb.Column("Col1", 8, false),
		tb.Column("Col2", 6, true),
		tb.Column("Col3", strings.TableExpandWidth, false),
	)

	w, err := strings.NewTableWriter(io.Discard, tb)
	if err != nil {
		t.Errorf("NewTableWriter() can not be an error: %+v", err)
	}

	w.WriteString("AAA")
	w.WriteString("BBB")
	w.WriteString("CCCCCCCCCCCCCCCC")
	w.WriteString("DDDDD")
	w.WriteString("E")
	w.WriteString("FFFFF")

	err = w.Close()
	if err != nil {
		t.Errorf("NewTableWriter() can not be an error: %+v", err)
	}
}

func Example() {

	tb := strings.NewTable()

	tb.Columns(
		tb.Column("Col1", 8, false),
		tb.Column("Col2", 6, true),
		tb.Column("Col3", 10, false),
	)

	tb.Add("AAA", "BBB", "CCCCCCCCCCCCCC")
	tb.Add("DDDDD", "E", "FFFFF")
	tb.Add("GGGGGGGG", "HHHHHHH", "I")
	tb.Add("マルチバイト文字にも", "一応", "対応しています。")

	buf, err := tb.Generate()
	if err != nil {
		return
	}

	fmt.Fprint(os.Stdout, buf)

	// Output:
	// |AAA     |   BBB|CCCCCCC...|
	// |DDDDD   |     E|FFFFF     |
	// |GGGGGGGG|HHH...|I         |
	// |マル... |  一応|対応し... |
	//

}

func ExampleTableOption_DisplayHeader() {
	tb := strings.NewTable()

	tb.Columns(
		tb.Column("Col1", 8, false),
		tb.Column("Col2", 6, true),
		tb.Column("Col3", 10, false),
	)

	tb.Options(
		strings.DisplayHeader(),
	)

	tb.Add("AAA", "BBB", "CCCCCCCCCCCCCC")
	tb.Add("DDDDD", "E", "FFFFF")

	buf, err := tb.Generate()
	if err != nil {
		return
	}

	fmt.Fprint(os.Stdout, buf)

	// Output:
	// |Col1    |  Col2|Col3      |
	// |--------|------|----------|
	// |AAA     |   BBB|CCCCCCC...|
	// |DDDDD   |     E|FFFFF     |
	//

}

func ExampleTableOption_DisplayLine() {
	tb := strings.NewTable()

	tb.Columns(
		tb.Column("Col1", 8, false),
		tb.Column("Col2", 6, true),
		tb.Column("Col3", 10, false),
	)

	tb.Options(
		strings.DisplayLine(),
	)

	tb.Add("AAA", "BBB", "CCCCCCCCCCCCCC")
	tb.Add("DDDDD", "E", "FFFFF")

	buf, err := tb.Generate()
	if err != nil {
		return
	}

	fmt.Fprint(os.Stdout, buf)

	// Output:
	// ----------------------------
	// |AAA     |   BBB|CCCCCCC...|
	// |DDDDD   |     E|FFFFF     |
	// ----------------------------
	//

}

func ExampleTableColumn_ExpandWidth() {

	tb := strings.NewTable()
	tb.Columns(
		tb.Column("Col1", 8, false),
		tb.Column("Col2", 6, true),
		tb.Column("Col3", strings.TableExpandWidth, false),
	)

	tb.Add("AAA", "BBB", "CCCCCCCCCCCCCC")
	tb.Add("DDDDD", "E", "FFFFF")

	buf, err := tb.Generate()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprint(os.Stdout, buf)

	// Output:
	// |AAA     |   BBB|CCCCCCCCCCCCCC|
	// |DDDDD   |     E|FFFFF         |
	//
}

func ExampleWriter() {

	tb := strings.NewTable()
	tb.Columns(
		tb.Column("Col1", 8, false),
		tb.Column("Col2", 6, true),
		tb.Column("Col3", 10, false),
	)

	tb.Options(
		strings.DisplayHeader(),
		strings.DisplayLine(),
	)

	w, _ := strings.NewTableWriter(os.Stdout, tb)

	defer w.Close()

	w.WriteString("AAA")
	w.WriteString("BBB")
	w.WriteString("CCCCCCCCCCCCCCCC")
	w.WriteString("DDDDD")
	w.WriteString("E")
	w.WriteString("FFFFF")

	// Output:
	// ----------------------------
	// |Col1    |  Col2|Col3      |
	// |--------|------|----------|
	// |AAA     |   BBB|CCCCCCC...|
	// |DDDDD   |     E|FFFFF     |
	// ----------------------------
	//
}

func ExampleWriter_LazyMode() {

	tb := strings.NewTable()
	tb.Columns(
		tb.Column("Col1", 8, false),
		tb.Column("Col2", 6, true),
		tb.Column("Col3", strings.TableExpandWidth, false),
	)

	tb.Options(
		strings.DisplayHeader(),
		strings.DisplayLine(),
	)

	w, _ := strings.NewTableWriter(os.Stdout, tb)

	defer w.Close()

	w.WriteString("AAA")
	w.WriteString("BBB")
	w.WriteString("CCCCCCCCCCCCCCCC")
	w.WriteString("DDDDD")
	w.WriteString("E")
	w.WriteString("FFFFF")

	// Output:
	// ----------------------------------
	// |Col1    |  Col2|Col3            |
	// |--------|------|----------------|
	// |AAA     |   BBB|CCCCCCCCCCCCCCCC|
	// |DDDDD   |     E|FFFFF           |
	// ----------------------------------
	//
}
