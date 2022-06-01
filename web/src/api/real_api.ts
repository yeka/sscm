class API {
    baseURL: string

    public constructor(baseURL: string) {
        this.baseURL = baseURL
    }

    public async ListCert(parent: number = 0) {
        return fetch(this.baseURL + "/api/certs?root=" + parent)
        .then(response => response.json())
    }

    /* 
        Create a new certificate under given parent
    */
    public async Create(data: Object, parent: number) {
        return fetch(this.baseURL + "/api/cert?root=" + parent, {
            method: "POST",
            body: JSON.stringify(data),
        })
    }

    /* 
        GetCert returns a detailed certificate given an ID
    */
    public async GetCert(id: number) {
        return fetch(this.baseURL + "/api/cert/" + id)
    }

    public async DownloadCert(id: number) {
        let name = "default.cert";
        return fetch(this.baseURL + "/api/download/" + id)
        .then(response => response.blob())
        .then(blob => {
            var url = window.URL.createObjectURL(blob);
            var a = document.createElement('a');
            a.href = url;
            a.download = name;
            document.body.appendChild(a); // we need to append the element to the dom -> otherwise it will not work in firefox
            a.click();    
            a.remove();  //afterwards we remove the element again         
        });
    }
}

export default API;