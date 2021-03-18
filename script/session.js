store.getCollections(() => {}, (code) => {
    console.log(code);
    if (code === Store.ErrorForbidden) {
        dom.find('login-layout').classList.toggle("active");
    }
});