const dom = {};

dom.onEvent = (event, id, listener) => {
    const element = document.getElementById(id);
    if (element != null) {
        element.addEventListener(event, (e) => {
            if (listener != null) {
                if (e != null && e.preventDefault != null) {
                    e.preventDefault();
                }
                listener(event);
            }
        });
    }
};

dom.node = (tagName) => {
    return document.createElement(tagName);
}

dom.text = (text) => {
    return document.createTextNode(text);
};

dom.find = (id) => {
    return document.getElementById(id);
};

dom.first = (className) => {
    const items = document.getElementsByClassName(className);
    if (items != null && items.length > 0) {
        return items[0];
    }
    return null;
};