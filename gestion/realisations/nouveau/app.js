const clearForm = () => {
    const labelInput = dom.find('label');
    if (labelInput != null) {
        labelInput.value = "";
    }

    const yearInput = dom.find('year');
    if (yearInput != null) {
        yearInput.value = "";
    }

    const descriptionInput = dom.find('description');
    if (labelInput != null) {
        descriptionInput.value = "";
    }

    const fileInput = dom.find('file');
    if (labelInput != null) {
        fileInput.type = "text";
        fileInput.value = "";
        fileInput.type = "file";
    }

    const visibleFileInput = dom.find('visible-file');
    if (labelInput != null) {
        visibleFileInput.value = "";
    }
};

dom.onEvent('click', 'fab', () => {
    dom.find('file').click();
});

dom.onEvent('change', 'file', () => {
    const fileInput = dom.find('file');
    let selectedFile = fileInput.files[0];
    dom.find('visible-file').value = selectedFile.name;
});

dom.onEvent('click', 'publish-button', () => {
    const isNumber = (n) => {
        return !isNaN(parseFloat(n)) && isFinite(n);
    };

    const labelInput = dom.find('label');
    if (labelInput == null || labelInput.value === '') {
        return;
    }

    const yearInput = dom.find('year');
    if (yearInput == null || yearInput.value === '' || !isNumber(yearInput.value)) {
        return;
    }

    const date = new Date();

    const year = parseInt(yearInput.value);
    if (year > date.getFullYear()) {
        return;
    }

    const descriptionInput = dom.find('description');
    if (descriptionInput == null || descriptionInput.value === '') {
        return;
    }

    const fileInput = dom.find('file');
    if (fileInput == null || fileInput.files.length === 0) {
        return;
    }

    const imageRemotePath = "/realisations/" + fileInput.files[0].name;

    const messageWidget = new Message();
    messageWidget.setTitle('Envoi');
    messageWidget.setLabel(labelInput.value);

    const onImageUploaded = () => {
        const fileURL = store.getFileDataURL(imageRemotePath);
        const object = {
            "label": labelInput.value,
            "description": descriptionInput.value,
            "year": year,
            "image": [fileURL]
        };
        store.saveObject('realisations', object, () => {
            messageWidget.setTitle("Publication");
            messageWidget.setLabel("Réalisation publiée sur le site");
            messageWidget.display(3000);

            clearForm();
        }, (code) => {
            alert("Une erreur est survenue côté serveur. CODE=" + code);
        });
    };

    const onProgress = (progress) => {
        if (progress === 100) {
            messageWidget.hide();
            return;
        }
        messageWidget.show();
        messageWidget.setState(progress);
    };
    const onImageUploadError = (code, action) => {
    };
    store.uploadFile(imageRemotePath, fileInput.files[0], onImageUploaded, onProgress, onImageUploadError);
});