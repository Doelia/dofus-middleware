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
    // console.log(evt.data);
    const parts = evt.data.split("|");
    const type = parts[0];
    const content = parts[1] || '';

    switch (type) {
        case 'CHARACTERS':
            let characters = JSON.parse(content);
            characters = characters.filter(v => v.Name !== '');
            console.log('characters', characters);
            app.characters = characters;

            if (characters.length > 0 && characters[0].Fight) {
                let fight = characters[0].Fight;
                console.log('fight', fight);
                app.fight = fight;
                // app.cells = app.cells.map(c => ({...c, fighter: null }));
                // fight.Fighters.forEach(fighter => {
                //     app.cells = app.cells.map(c => ({...c, fighter: c.id === fighter.CellId ? fighter : c.fighter }));
                // });
            }

            break;
        case 'OPTIONS':
            let options = JSON.parse(content);
            console.log('options', options);
            app.options = options;
            break;
        case 'MAP':
            let map = JSON.parse(content);
            console.log('map', map);
            app.map = map;
            map.Cells.forEach(cell => {
                app.cells = app.cells.map(c => {
                    if (c.id === cell.CellId) {
                        return {...c, ...cell}
                    } else {
                        return c
                    }
                });
            })
            console.log('map', app.cells);
            break;
        case 'PATH':
            let ids = JSON.parse(content);
            app.cells = app.cells.map(c => ({...c, path: false }));
            ids.forEach(id => {
                app.cells = app.cells.map(c => {
                    if (c.id === id) {
                        return {...c, path: true}
                    } else {
                        return c
                    }
                });
            })

            break;
    }
}
ws.onerror = function(evt) {
    print("ERROR: " + evt.data);
}

