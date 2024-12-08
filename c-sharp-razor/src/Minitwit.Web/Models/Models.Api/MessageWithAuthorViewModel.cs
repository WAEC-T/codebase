namespace Minitwit.Web.Models.Models.Api;

public class MessageWithAuthorViewModel
{
    public string Content { get; set; }
    public DateTime PubDate { get; set; }
    public string UserName { get; set; }
}