class DummyAPI {
    public async GetRoot() {
        await new Promise(resolve => setTimeout(resolve, 1000)); // simulate 1000ms loading time 
        return {
            "roots": [
                {"name": "Hello"}
            ]
        }
    }
}
export default DummyAPI;