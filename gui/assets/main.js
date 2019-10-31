
var Data = {
    characters: [],
    message: 'Hello'
};

function refreshVue() {
    var app = new Vue({
        el: '#app',
        data: Data,
        methods: {
            focusCharacter: function(name) {
                console.log('focus', name);
                ws.send('FOCUS|' + name)
            }
        }
    })
}