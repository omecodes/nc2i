class Message {

    constructor() {
        this.indeterminate = true;
        this.title = '';
        this.label = '';
        this.state = '';
        this.showing = false;

        this.element = document.getElementById('progress-panel');
        if (this.element == null) {
            this.element = dom.node('DIV');
            this.element.id = 'progress-panel';
            this.element.overflow = 'hidden';
            this.element.style.position = 'absolute';
            this.element.style.borderBox = 'box-sizing';
            this.element.style.bottom = '32px';
            this.element.style.left = '-332px';
            this.element.style.borderRadius = '4px';
            this.element.style.padding = '16px';
            this.element.style.width = '300px';
            this.element.style.height = '80px';
            this.element.style.transition = 'left 0.3s';
            this.element.style.boxShadow = '0.5px 0.5px 2px #444444aa';

            // title
            const title = dom.node('P');
            title.id = 'progress-panel-title';
            title.style.padding = '0';
            title.style.margin = '0';
            title.style.textAlign = 'start';
            title.style.height = '24px';
            title.style.color = 'black';
            title.style.fontWeight = '600';
            title.style.fontSize = '14px';

            // Info
            const info = dom.node('DIV');
            title.style.height = '24px';
            info.classList = ['flex', 'flex-jc-sb'];

            const span = dom.node('SPAN');
            span.id = 'progress-panel-label';

            const percentage = dom.node('SPAN');
            percentage.id = 'progress-panel-percentage';
            percentage.style.width = '50px';

            info.append(span);
            info.append(percentage);

            this.element.append(title);
            this.element.append(info);
        }
        document.body.appendChild(this.element);
    }

    setIndeterminate(indeterminate) {
        this.indeterminate = indeterminate;
    }

    setState(state) {
        this.indeterminate = false;
        this.state = state;

        const percentageLabel = document.getElementById('progress-panel-percentage');
        percentageLabel.innerHTML = '';
        percentageLabel.appendChild( dom.text(state))
    }

    setTitle(title) {
        this.label = title;
        const titleElement = document.getElementById('progress-panel-title');
        titleElement.innerHTML = '';
        titleElement.appendChild(dom.text(title));
    }

    setLabel(label) {
        this.label = label;
        const labelElement = document.getElementById('progress-panel-label');
        labelElement.innerHTML = '';
        labelElement.appendChild(dom.text(label));
    }

    show() {
        if (this.showing) {
            return;
        }
        this.showing = true;
        this.element.style.left = '32px';
    }

    display(duration) {
        this.show();
        const msg = this;
        setTimeout(function () {
            msg.hide();
        }, duration);
    }

    hide() {
        this.element.style.left = '-332px';
        this.showing = false;
        this.setTitle('');
        this.setLabel('');
    }
}