name=$1

osascript -e '
tell application "Dofus" to activate

tell application "System Events" to tell application process "Dofus"
	get properties of windows
	repeat with c in windows
		if name of c is "'$1' - Dofus 1.30.0" then
			perform action "AXRaise" of c
		end if
	end repeat
end tell
'
