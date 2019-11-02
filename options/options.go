package options

type OptionsStruct struct {
	AutoJoinFight bool
	AutoReadyFight bool
	ShowInputPackets bool
	ShowOutputPackets bool
	DispatchMoves bool
	FocusWindowOnCharacterTurn bool
}

var Options = OptionsStruct{
	AutoJoinFight: true,
	AutoReadyFight: true,
	ShowInputPackets: true,
	ShowOutputPackets: true,
	DispatchMoves: false,
	FocusWindowOnCharacterTurn: true,
}

var ConfigSqlLitePath = "/Users/stephane/Desktop/dodo.sqlite"