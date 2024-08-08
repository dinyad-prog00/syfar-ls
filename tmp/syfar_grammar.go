package tmp

import "github.com/alecthomas/participle/v2/lexer"

// Package parser implements a parser for Syfar Test Language syntax.

type SyfarFile struct {
	Entries []*Entry `parser:"@@*"`
}

type Entry struct {
	Stepper       *Stepper       `parser:"@@"`
	Action        *Action        `parser:"|@@"`
	TestSet       *TestSet       `parser:"|@@"`
	Test          *Test          `parser:"|@@"`
	Variable      *Variable      `parser:"|@@"`
	MultiVariable *MultiVariable `parser:"|@@"`
	VarSet        *VarSet        `parser:"|@@"`
	SecretSet     *SecretSet     `parser:"|@@"`
	Print         *Print         `parser:"|@@"`
	Import        *Import        `parser:"|@@"`
	Pos           lexer.Position
	EndPos        lexer.Position
}

type Action struct {
	Type       string `parser:"'action' @(Ident|String)"`
	Id         string `parser:" @(Ident|String)"`
	Prefix     *string
	StartBlock StartBlock         `parser:"@@"`
	Attributes []*ActionAttribute `parser:"@@*"`
	EndBlock   EndBlock           `parser:"@@"`
	Pos        lexer.Position
	EndPos     lexer.Position
}

type ActionAttribute struct {
	Parameter *Assignment `parser:"@@"`
	Test      *Test       `parser:"|@@"`
	TestSet   *TestSet    `parser:"|@@"`
	Out       *Out        `parser:"|@@"`
	Pos       lexer.Position
	EndPos    lexer.Position
}

type Stepper struct {
	Id         string     `parser:"'steps' @(Ident|String)"`
	StartBlock StartBlock `parser:"@@"`
	Steps      []*Steps   `parser:"@@*"`
	EndBlock   EndBlock   `parser:"@@"`
	Pos        lexer.Position
	EndPos     lexer.Position
}

type Steps struct {
	Action *Action `parser:"@@"`
}

type TestSet struct {
	Description string     `parser:"'tests' @String"`
	StartBlock  StartBlock `parser:"@@"`
	Tests       []*Test    `parser:" @@*"`
	EndBlock    EndBlock   `parser:"@@"`
	Pos         lexer.Position
	EndPos      lexer.Position
}

type Test struct {
	Skipped      bool           `parser:"@'~'?"`
	Description  string         `parser:"'test' @String"`
	StartBlock   StartBlock     `parser:"@@"`
	Expectations []*Expectation `parser:"@@*"`
	EndBlock     EndBlock       `parser:"@@"`
	Pos          lexer.Position
	EndPos       lexer.Position
}

type Expectation struct {
	Pos        lexer.Position
	EndPos     lexer.Position
	Key        string             `parser:"'expect'"`
	StartBlock StartBlock         `parser:"@@"`
	Items      []*ExpectationItem `parser:"@@*"`
	EndBlock   EndBlock           `parser:"@@"`
}

type ExpectationItem struct {
	Pos      lexer.Position
	EndPos   lexer.Position
	Symbolic *SymbolicCheck `parser:"@@"`
	Chain    *ChainCheck    `parser:"|@@"`
}

type SymbolicCheck struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Key    string `parser:"@Ident (@'[' @Int @']')? ( @'.' @Ident (@'[' @Int @']')? )*"`
	Opp    string `parser:"@('=='|'<='|'>='|'<'|'>'|'!='|'eq'|'gt'|'lt'|'le'|'ge'|'ne')"`
	Value  *Value `parser:"@@"`
}

type ChainCheck struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Key    string       `parser:"@Ident ( @'.' @Ident )* ':' "`
	Chain  []*ChainItem `parser:"@@*"`
}

type ChainItem struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Start  string   `parser:"@'.'? @'to'"`
	Negate bool     `parser:"(@'.' @'not')?"`
	Deep   bool     `parser:"(@'.' @'deep')?"`
	Method string   `parser:" @'.' @('be')  @'.' @('eq'|'gt'|'lt'|'le'|'ge'|'ne')"`
	Args   []*Value `parser:"('(' @@* ')')?"`
}

type ChainMethod struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Name   string `parser:"@('be' |'eq'|'gt'|'lt'|'le'|'ge'|'ne'| 'been' | 'is' | 'that' | 'which' | 'and' | 'has' | 'have' | 'with' | 'at' | 'of' | 'same' | 'but' | 'does' | 'still' | 'also' | 'not' | 'deep' )"`
}

type ChainArg struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Value  Value `parser:"@Ident"`
}

type Variable struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Name   string `parser:"'var' @Ident"`
	Type   string `parser:"(':' @('number'|'string'|'array'|'bool'|'object'))?"`
	Value  *Value `parser:"'=' @@"`
}

type MultiVariable struct {
	Pos       lexer.Position
	EndPos    lexer.Position
	Variables []*Assignment `parser:"'var' '(' @@* ')'"`
}

type VarSet struct {
	Pos        lexer.Position
	EndPos     lexer.Position
	Id         string        `parser:"'vars' @(Ident|String)"`
	StartBlock StartBlock    `parser:"@@"`
	Variables  []*Assignment `parser:"@@*"`
	EndBlock   EndBlock      `parser:"@@"`
}

type Assignment struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Name   string `parser:"@Ident"`
	Value  *Value `parser:"'=' @@"`
}

type SecretSet struct {
	Pos        lexer.Position
	EndPos     lexer.Position
	Id         string        `parser:"'secrets' @(Ident|String)"`
	StartBlock StartBlock    `parser:"@@"`
	Variables  []*Assignment `parser:"@@*"`
	EndBlock   EndBlock      `parser:"@@"`
}

type Bool bool

func (b *Bool) Capture(v []string) error { *b = v[0] == "true"; return nil }

type Value struct {
	Pos        lexer.Position
	EndPos     lexer.Position
	Boolean    *Bool    `parser:" @('true'|'false')"`
	Identifier *string  `parser:"| @Ident (@'[' @Int @']')? ( @'.' @Ident (@'[' @Int @']')? )*"`
	String     *string  `parser:"| @(String|Char|RawString)"`
	Number     *float64 `parser:"| @(Float|Int)"`
	Array      []*Value `parser:"| '[' ( @@ ','? )* ']'"`
	Json       *JSON    `parser:"|@@"`
	Map        map[string]interface{}
	Any        interface{}
}

type JSON struct {
	Pos        lexer.Position
	EndPos     lexer.Position
	Attributes []*JSONAttribute `parser:"'{' @@* '}'"`
}

type JSONAttribute struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Name   string `parser:"@(Ident|String)"`
	Value  *Value `parser:"':' @@ (',')?"`
}

type HeaderValue struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Name   string `parser:"@(Ident|String)"`
	Value  *Value `parser:"':' @@ (',')?"`
}

type VariableType struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Name   string `parser:"'number' | 'string'"`
}

type Print struct {
	Pos        lexer.Position
	EndPos     lexer.Position
	Id         string     `parser:"'print'"`
	StartBlock StartBlock `parser:"@@"`
	Variables  []*Value   `parser:"(@@ ','?)*"`
	EndBlock   EndBlock   `parser:"@@"`
}

type Out struct {
	Pos        lexer.Position
	EndPos     lexer.Position
	Key        string           `parser:"'out'"`
	StartBlock StartBlock       `parser:"@@"`
	Variables  []*OutAssignment `parser:"(@@ ','?)*"`
	EndBlock   EndBlock         `parser:"@@"`
}

type OutAssignment struct {
	Pos        lexer.Position
	EndPos     lexer.Position
	Name       string `parser:"@Ident"`
	Identifier string `parser:"'=' @Ident (@'[' @Int @']')? ( @'.' @Ident (@'[' @Int @']')? )*"`
}

type Import struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Files  []string `parser:"'import' '(' (@String ','?)* ')'"`
}

type Argument struct {
	Pos     lexer.Position
	EndPos  lexer.Position
	Name    string `parser:"@Ident"`
	Type    string `parser:"':' @('Id'|'Value')"`
	Default *Value `parser:"( '=' @@ )?"`
}

type StartBlock struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Start  string `parser:"'{'"`
}

type EndBlock struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Start  string `parser:"'}'"`
}
