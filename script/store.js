class Store {

    static ErrorServerInternal = 1;
    static ErrorNotAllowed = 2;
    static ErrorForbidden = 3;
    static ErrorBadRequest = 4;
    static ErrorServerNotReachable = 5;
    static ErrorRequestFailed = 6;

    static FilesPath = "/api/files"

    constructor(serverURL) {
        this.serverURL = serverURL;
        this.cache = {};
        this.apiKey = 'nc2i-website';
        this.apiSecret = 'secure-access';
    }

    _errorFromStatus(status) {
        if (status === 500) {
            return Store.ErrorServerInternal;
        } else if (status === 403) {
            return Store.ErrorForbidden;
        } else if (status === 401) {
            return Store.ErrorNotAllowed;
        } else if (status === 503) {
            return Store.ErrorServerNotReachable;
        }
        return Store.ErrorRequestFailed;
    }

    initClientApp() {
        const xhr = new XMLHttpRequest();
        xhr.open('Post', '/api/auth/sessions/client-app');
        xhr.send(JSON.stringify({
            client : {
                key: this.apiKey,
                secret: this.apiSecret
            }
        }));
    }

    getFileDataURL(path) {
        return  this.serverURL + Store.FilesPath + '/data' + path;
    }

    createClientAppSession(key, secret, onSuccess, onError = null)  {
        const data = {
            "access": {
                "key": key,
                "secret": secret
            }
        }

        const xhr = new XMLHttpRequest();
        xhr.open('Post', '/api/auth/sessions');
        xhr.withCredentials = true;
        xhr.onreadystatechange = () => {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                if (xhr.status === 200) {
                    onSuccess();
                    return;
                }

                if (onError != null) {
                    onError(this._errorFromStatus(xhr.status));
                }
            }
        };
        xhr.send(JSON.stringify(data));
    }

    getCollections(onSuccess, onError = null) {
        const xhr = new XMLHttpRequest();
        xhr.open('Get', '/api/objects/collections');
        xhr.onreadystatechange = () => {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                if (xhr.status === 200) {
                    let collections = JSON.parse(xhr.responseText);
                    onSuccess(collections);
                    return;
                }

                if (onError != null) {
                    onError(this._errorFromStatus(xhr.status));
                }
            }
        };
        xhr.send();
    }

    listObjects(collection, onSuccess, onError = null) {
        const xhr = new XMLHttpRequest();
        xhr.open('Get', '/api/objects/data/' + collection);
        xhr.onreadystatechange = () => {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                if (xhr.status === 200) {
                    let objects = JSON.parse(xhr.responseText);
                    console.log(objects);
                    onSuccess(objects);
                    return;
                }
                if (onError != null) {
                    onError(this._errorFromStatus(xhr.status));
                }
            }
        };
        xhr.send();
    }

    getObject(collection, id, onSuccess, onError = null) {
        const xhr = new XMLHttpRequest();
        xhr.open('Get', '/api/objects/' + collection + "/" + id);
        xhr.onreadystatechange = () => {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                if (xhr.status === 200) {
                    let object = JSON.parse(xhr.responseText);
                    onSuccess(object);
                    return;
                }

                if (onError != null) {
                    onError(this._errorFromStatus(xhr.status));
                }
            }
        };
        xhr.send();
    }

    saveObject(collection, object, onSuccess, onError = null) {
        const xhr = new XMLHttpRequest();
        xhr.open('Put', '/api/objects/data/' + collection);
        xhr.setRequestHeader("Content-Type", "application/json")
        xhr.onreadystatechange = () => {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                if (xhr.status === 200) {
                    let object = JSON.parse(xhr.responseText);
                    onSuccess(object);
                    return;
                }

                if (onError != null) {
                    onError(this._errorFromStatus(xhr.status));
                }
            }
        };
        xhr.send(JSON.stringify({
            "collection": collection,
            "object": {
                data: JSON.stringify(object)
            }
        }));
    }

    uploadFile(path, file, onSuccess, onProgress, onError = null) {
        const reader = new FileReader();
        const store =  this;
        reader.addEventListener('load', function (e) {
            const xhr = new XMLHttpRequest();
            xhr.open('Put', '/api/files/data' + path);
            xhr.upload.onprogress = function(e) {
                if (e.lengthComputable) {
                    let progress = (e.loaded * 100) / e.total;
                    onProgress(progress);
                }
            };
            xhr.onreadystatechange = () => {
                if (xhr.readyState === XMLHttpRequest.DONE) {
                    if (xhr.status === 200) {
                        onSuccess();
                        return;
                    }

                    if (onError != null) {
                        onError(store._errorFromStatus(xhr.status), xhr);
                    }
                }
            };
            console.log(e.target);
            xhr.send(e.target.result);
        });
        reader.readAsArrayBuffer(file);
    }
}

const store = new Store('http://localhost:8080');
store.initClientApp();