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
	AutoReadyFight: true,
	ShowInputPackets: false,
	ShowOutputPackets: true,
	DispatchMoves: false,
	FocusWindowOnCharacterTurn: true,
}