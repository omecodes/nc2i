const toggleMenu = () => {
    dom.find('dropdown-menu').classList.toggle('active');
    dom.find('click-catcher').classList.toggle('active');
}
dom.onEvent('click', 'hamburger', toggleMenu);
dom.onEvent('click', 'click-catcher', toggleMenu);

function loadObjects() {
    let grid = dom.find('grid');
    const onObjectsLoaded = (objects) => {
        const keys = Object.keys(objects);
        for (let i = 0; i < keys.length; i++) {
            let object = objects[keys[i]];
            let cardData = new CardData(object.label, object.description, object.image)
            let card = new Card(cardData);
            grid.appendChild(card.dom());
        }
    };
    store.listObjects("realisations", onObjectsLoaded, (code) => {});
}

loadObjects();