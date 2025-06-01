package dk.itu.minitwit.domain;

public class SimData {

    private String content;
    private String follow;
    private String unfollow;

    public String getUnfollow() {
        return unfollow;
    }

    public void setUnfollow(String unfollow) {
        this.unfollow = unfollow;
    }

    public String getFollow() {
        return follow;
    }

    public void setFollow(String follow) {
        this.follow = follow;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    @Override
    public String toString() {
        return "SimData{" +
                "content='" + content + '\'' +
                ", follow='" + follow + '\'' +
                ", unfollow='" + unfollow + '\'' +
                '}';
    }
}
