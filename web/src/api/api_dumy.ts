class DummyAPI {

    public async ListCert(parent: number = 0) {
        await new Promise(resolve => setTimeout(resolve, 500)); // simulate loading time 
        return {
            "parent": {"id": 11, "name": "Hello"},
            "certs": [
                {"id": 1, "name": "Hello"},
                {"id": 2, "name": "Govv"},
            ]
        }
    }

    /* 
        Create a new certificate under given parent
    */
    public async Create(data: Object, parent: number) {
        await new Promise(resolve => setTimeout(resolve, 500)); // simulate loading time 
        throw "" // for failed condition
        // return {} // for success condition
    }

    /* 
        GetCert returns a detailed certificate given an ID
    */
    public async GetCert(id: number) {
        await new Promise(resolve => setTimeout(resolve, 500)); // simulate loading time 
        return {
            "cert": {
                "name": "Hello",
            }
        }
    }
}

export default DummyAPI;