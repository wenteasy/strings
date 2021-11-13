package strings

import (
	"fmt"
	"io"
	"strings"

	"github.com/mattn/go-runewidth"
	"golang.org/x/xerrors"
)

type Table struct {
	headers []string
	columns []*TableColumn
	values  [][]string

	//options
	header        bool   //ヘッダ表示ありか？
	spaceChar     string //隙間の文字
	edgesChar     string //行の左右に入れる文字
	line          bool   //開始、終了にラインを入れるか？
	lineChar      string //ヘッダなどに挟まる文字
	linefeed      string //改行コード
	separatorChar string //間に入れる文字
	omitChar      string //長さが足りない時に足す文字列
}

type TableOption func(*Table)

const (
	DefaultTableLineCharacter      = "-"
	DefaultTableEdgesCharacter     = "|"
	DefaultTableSpaceCharacter     = " "
	DefaultTableSeparatorCharacter = "|"
	DefaultTableNewLineCode        = "\n"
	DefaultTableOmitCharacter      = "..."

	TableExpandWidth = 0
	TableMargeColumn = "____This_Column_Will_Merge____"
)

func NewTable() *Table {
	var t Table

	//options
	t.header = false
	t.headers = nil

	t.line = false
	t.lineChar = DefaultTableLineCharacter

	t.linefeed = DefaultTableNewLineCode
	t.edgesChar = DefaultTableEdgesCharacter
	t.spaceChar = DefaultTableSpaceCharacter
	t.separatorChar = DefaultTableSeparatorCharacter
	t.omitChar = DefaultTableOmitCharacter

	return &t
}

func (t *Table) check() error {
	//TODO 各オプション設定のチェックを行う
	return nil
}

func (t *Table) Options(opts ...TableOption) {
	for _, opt := range opts {
		opt(t)
	}
}

func Edges(b string) TableOption {
	return func(t *Table) {
		t.edgesChar = b
	}
}

func DisplayHeader() TableOption {
	return func(t *Table) {
		t.header = true
	}
}
func DisplayLine() TableOption {
	return func(t *Table) {
		t.line = true
	}
}

func (t *Table) Columns(c ...*TableColumn) {
	t.columns = append(t.columns, c...)
	for _, col := range t.columns {
		t.headers = append(t.headers, col.name)
	}
}

func (t *Table) Add(data ...string) {
	t.values = append(t.values, data)
}

func (t *Table) Generate() (string, error) {

	var b strings.Builder

	header, hl, err := t.createHeader()

	//TODO call Grow()
	if header != "" {
		b.WriteString(header)
	}

	err = t.writeBody(&b)
	if err != nil {
		return "", xerrors.Errorf("writeBody() error: %w", err)
	}

	if t.line {
		b.WriteString(hl)
		b.WriteString(t.linefeed)
	}

	return b.String(), nil
}

func (t *Table) writeBody(b *strings.Builder) error {

	for _, line := range t.values {
		data, err := t.createLine(line)
		if err != nil {
			return xerrors.Errorf("createLine() error: %w", err)
		}
		b.WriteString(data)
		b.WriteString(t.linefeed)
	}

	return nil
}

func (t *Table) createHeader() (string, string, error) {

	seps := make([]string, len(t.columns))
	width := len(t.edgesChar) * 2

	for idx, tc := range t.columns {
		if tc.digit <= TableExpandWidth {
			err := t.decideColumnDigit()
			if err != nil {
				return "", "", fmt.Errorf("Table decideColumnDigit() error: %w", err)
			}
		}
		seps[idx] = strings.Repeat(t.lineChar, tc.digit)
		width += tc.digit
	}

	var b strings.Builder

	width += len(t.columns)*len(t.separatorChar) - 1
	hl := strings.Repeat(t.lineChar, width)

	if t.line {
		b.WriteString(hl)
		b.WriteString(t.linefeed)
	} else {
		hl = ""
	}

	if t.header {

		line, err := t.createLine(t.headers)
		if err != nil {
			return "", "", xerrors.Errorf("createLine() error: %w", err)
		}
		b.WriteString(line)
		b.WriteString(t.linefeed)

		line, err = t.createLine(seps)
		if err != nil {
			return "", "", xerrors.Errorf("createLine() error: %w", err)
		}
		b.WriteString(line)
		b.WriteString(t.linefeed)
	}

	return b.String(), hl, nil
}

func (t *Table) createLine(data []string) (string, error) {

	var b strings.Builder

	if t.edgesChar != "" {
		b.WriteString(t.edgesChar)
	}

	for idx, tc := range t.columns {
		d, err := tc.Fill(data[idx])
		if err != nil {
			return "", xerrors.Errorf("TableColumn[%d] Fill() error: %w", idx, err)
		}

		b.WriteString(d)
		if idx != len(data)-1 {
			b.WriteString(t.separatorChar)
		}
	}
	if t.edgesChar != "" {
		b.WriteString(t.edgesChar)
	}
	return b.String(), nil
}

func (t *Table) decideColumnDigit() error {

	for idx, tc := range t.columns {

		if tc.digit > TableExpandWidth {
			continue
		}

		width := 0
		for _, row := range t.values {
			sw := runewidth.StringWidth(row[idx])
			if sw > width {
				width = sw
			}
		}

		if width <= TableExpandWidth {
			return fmt.Errorf("could not decide column[%d] width", idx)
		}
		tc.digit = width
	}

	return nil
}

func (t *Table) Allocate(rows int) {
	t.values = make([][]string, 0, rows)
}

func (t *Table) Column(n string, d int, r bool) *TableColumn {
	var tc TableColumn
	tc.owner = t
	tc.name = n
	tc.digit = d
	tc.right = r
	return &tc
}

type TableColumn struct {
	owner *Table
	name  string
	digit int
	right bool
}

func (ac *TableColumn) Fill(datum string) (string, error) {

	//まだ桁が決定してない
	if ac.digit <= 0 {
		return "", fmt.Errorf("the digit has not been decided yet.")
	}

	rtn := runewidth.Truncate(datum, ac.digit, ac.owner.omitChar)
	sw := runewidth.StringWidth(rtn)

	remain := ac.digit - sw
	if remain > 0 {
		f := strings.Repeat(ac.owner.spaceChar, remain)
		if ac.right {
			rtn = f + rtn
		} else {
			rtn = rtn + f
		}
	}

	return rtn, nil
}

type TableWriter struct {
	body  io.Writer
	table *Table

	mode    tableWriterMode
	workRow []string
	line    string
}

type tableWriterMode int

const (
	tableWriterModeRealTime tableWriterMode = iota
	tableWriterModeLazy
)

func NewTableWriter(w io.Writer, t *Table) (*TableWriter, error) {

	var tw TableWriter
	tw.mode = tableWriterModeRealTime

	err := t.decideColumnDigit()
	if err != nil {
		//TODO other error
		tw.mode = tableWriterModeLazy
	}

	tw.body = w
	tw.table = t
	tw.workRow = make([]string, 0, len(t.columns))
	tw.line = ""

	if tw.isRealTime() {
		err = tw.writeHeader()
		if err != nil {
			return nil, xerrors.Errorf("writeHeader() error: %w", err)
		}
	}

	return &tw, nil
}

func (w *TableWriter) writeHeader() error {
	header, line, err := w.table.createHeader()
	if err != nil {
		return xerrors.Errorf("Table createHeader() error: %w", err)
	}

	if header != "" {
		_, err := w.body.Write([]byte(header))
		if err != nil {
			return xerrors.Errorf("header Write() error: %w", err)
		}
	}
	w.line = line
	return nil
}

func (w *TableWriter) isRealTime() bool {
	return w.mode == tableWriterModeRealTime
}

func (w *TableWriter) Write(p []byte) (int, error) {
	return w.write(p)
}

func (w *TableWriter) WriteString(s string) (int, error) {
	//TODO unsafe?
	// https://qiita.com/mattn/items/176459728ff4f854b165
	return w.write([]byte(s))
}

func (w *TableWriter) write(p []byte) (int, error) {

	w.workRow = append(w.workRow, string(p))
	//行に達した場合
	if len(w.workRow) == len(w.table.columns) {

		if w.isRealTime() {

			line, err := w.table.createLine(w.workRow)
			if err != nil {
				return 0, xerrors.Errorf("createLine() error: %w", err)
			}

			_, err = w.body.Write([]byte(line + w.table.linefeed))
			if err != nil {
				return 0, xerrors.Errorf("body Write() line error: %w", err)
			}
		} else {
			w.table.values = append(w.table.values, w.workRow)
		}

		w.workRow = make([]string, 0, len(w.table.columns))
	}
	return len(p), nil
}

func (w *TableWriter) Close() error {

	err := w.flush()
	if err != nil {
		return xerrors.Errorf("flush() error: %w", err)
	}

	if c, ok := w.body.(io.Closer); ok {
		err := c.Close()
		if err != nil {
			return xerrors.Errorf("body Close() error: %w", err)
		}
	}

	return nil
}

func (w *TableWriter) flush() error {

	var err error
	//TODO 最終書き出しかモードをチェック
	if !w.isRealTime() {

		err = w.table.decideColumnDigit()
		if err != nil {
			return xerrors.Errorf("Table decideColumnDigit() error: %w", err)
		}

		err = w.writeHeader()
		if err != nil {
			return xerrors.Errorf("writeHaeder() error: %w", err)
		}

		var b strings.Builder
		err = w.table.writeBody(&b)
		if err != nil {
			return xerrors.Errorf("Table writeBody() error: %w", err)
		}
		w.body.Write([]byte(b.String()))
	}

	if len(w.workRow) != 0 {
		err = fmt.Errorf("the data of the last line is not available.")
	}

	if w.line != "" {
		_, err = w.body.Write([]byte(w.line))
	}

	return err
}
