package dk.itu.minitwit.domain;

public class SimMessage {
    private String content;
    private int pub_date;
    private String user;

    public SimMessage(String content, int pub_date, String user) {
        this.content = content;
        this.pub_date = pub_date;
        this.user = user;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public int getPub_date() {
        return pub_date;
    }

    public void setPub_date(int pub_date) {
        this.pub_date = pub_date;
    }

    public String getUser() {
        return user;
    }

    public void setUser(String user) {
        this.user = user;
    }
}
