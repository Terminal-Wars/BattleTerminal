module github.com/Terminal-Wars/BattleTerminal

go 1.17

require (
	github.com/Terminal-Wars/TermUI v1.1.1
	gopkg.in/irc.v3 v3.1.4
)

require github.com/jezek/xgb v1.0.0

replace github.com/Terminal-Wars/TermUI v0.0.0 => /tmp/TermUITest
