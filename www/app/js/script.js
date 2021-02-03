console.log("hello world");

function sidebarToggle() {
    document.getElementById('side-bar').classList.toggle('active');
    document.getElementById('side-bar-container').classList.toggle('active');
}

function sidebarClickOutside() {
    document.getElementById('side-bar').classList.toggle('active');
    document.getElementById('side-bar-container').classList.toggle('active');
}

function initDrawer() {
    let clickCatcher = document.getElementById('click-catcher');
    if (clickCatcher != null) {
        clickCatcher.onclick = sidebarClickOutside;
    }
}

initDrawer();