package main

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
	ShowInputPackets: false,
	ShowOutputPackets: false,
	DispatchMoves: false,
	FocusWindowOnCharacterTurn: true,
}