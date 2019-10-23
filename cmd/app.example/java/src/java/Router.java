public class Router {
    Data DataObject;


    // Getter Methods

    public Data getData() {
        return DataObject;
    }

    // Setter Methods

    public void setData(Data dataObject) {
        this.DataObject = dataObject;
    }
}

class Data {
    private String uri;
    private String token;


    // Getter Methods

    public String getUri() {
        return uri;
    }

    public String getToken() {
        return token;
    }

    // Setter Methods

    public void setUri(String uri) {
        this.uri = uri;
    }

    public void setToken(String token) {
        this.token = token;
    }
}