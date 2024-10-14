namespace Minitwit.Core.Entities;

public record CreateCheep(int AuthorId, string Text)
{
    public readonly int AuthorId = AuthorId;
    public readonly string Text = Text;
}
