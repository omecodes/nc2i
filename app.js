const toggleMenu = () => {
    dom.find('dropdown-menu').classList.toggle('active');
    dom.find('click-catcher').classList.toggle('active');
};

dom.onEvent('click', 'hamburger', toggleMenu);
dom.onEvent('click', 'click-catcher', toggleMenu);