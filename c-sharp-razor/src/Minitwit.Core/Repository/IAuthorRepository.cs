using Minitwit.Core.Entities;

namespace Minitwit.Core.Repository;

public interface IAuthorRepository
{
    public Task<ICollection<Author>> GetAuthorsByIdAsync(IEnumerable<int> authors);
    public Task<Author?> GetAuthorByIdAsync(int authorId);
    public Task<Author> GetAuthorByNameAsync(string name);
    public Task<ICollection<Message>> GetMessagesByAuthorAsync(int id);
    public Task<ICollection<Message>> GetMessagesByAuthor(int authorId, int page);
    public Task<ICollection<Message>> GetMessagesByAuthorAndFollowing(int authorId, int page);

    public Task<int> GetMessageCountByAuthorAsync(int authorId);
    public Task<int> GetMessageCountByAuthorAndFollowing(int authorId);

    public Task<int> GetPageCountByAuthor(int authorId);
    public Task<int> GetPageCountByAuthorAndFollowing(int authorId);

    public Task<ICollection<Author>> GetFollowersByIdAsync(int authorId);
    public Task<ICollection<Author>> GetFollowingByIdAsync(int authorId);

    public Task AddFollowAsync(int followingAuthorId, int followedAuthorId);
    public Task RemoveFollowAsync(int followingAuthorId, int followedAuthorId);
}
