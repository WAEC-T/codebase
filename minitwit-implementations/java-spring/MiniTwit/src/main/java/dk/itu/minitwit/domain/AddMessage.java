package dk.itu.minitwit.domain;

public class AddMessage {
    String text;

    public AddMessage(String text) {
        this.text = text;
    }

    public String getText() {
        return text;
    }

    public void setText(String text) {
        this.text = text;
    }

    @Override
    public String toString() {
        return "AddMessage{" +
                "text='" + text + '\'' +
                '}';
    }
}
