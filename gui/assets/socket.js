var ws;

function print(s) {
    console.log('socket', s);
}

ws = new WebSocket("ws://localhost:8001/ws");
ws.onopen = function(evt) {
    print("OPEN");
}
ws.onclose = function(evt) {
    print("CLOSE");
    ws = null;
}
ws.onmessage = function(evt) {
    console.log(evt.data);
    const parts = evt.data.split("|");
    const type = parts[0];
    const content = parts[1] || '';

    switch (type) {
        case 'CHARACTERS':
            let characters = JSON.parse(content);
            characters = characters.filter(v => v.Name !== '');
            console.log(characters);
            app.characters = characters;

            app.characters.forEach(char => {
                app.cells = app.cells.map(c => ({...c, character: c.id === char.CellId ? 'yes' : null }));
            });
            break;
        case 'OPTIONS':
            let options = JSON.parse(content);
            console.log('options', options);
            app.options = options;
            break;
        case 'FIGHT':
            let fight = JSON.parse(content);
            console.log('fight', fight);
            break;
    }
}
ws.onerror = function(evt) {
    print("ERROR: " + evt.data);
}

