dom.onEvent('click', 'fab', () => {
    console.log('clicked on fab');
    dom.find('files').click();
});