using Minitwit.Core.Entities;

namespace Minitwit.Core.Repository;

public interface IAuthorRepository
{
    public void AddAuthorAsync(Author authorDto);

    public Task<ICollection<Author>> GetAllAuthorsAsync();
    public Task<ICollection<Author>> GetAuthorsByIdAsync(IEnumerable<int> authors);
    public Task<Author?> GetAuthorByIdAsync(int authorId);
    public Task<Author> GetAuthorByNameAsync(string name);
    public Task<Author> GetAuthorByEmail(string email);

    public Task<ICollection<Cheep>> GetCheepsByAuthorAsync(int id);
    public Task<ICollection<Cheep>> GetCheepsByAuthor(int authorId, int page);
    public Task<ICollection<Cheep>> GetCheepsByAuthorAndFollowing(int authorId, int page);

    public Task<int> GetCheepCountByAuthorAsync(int authorId);
    public Task<int> GetCheepCountByAuthorAndFollowing(int authorId);

    public Task<int> GetPageCountByAuthor(int authorId);
    public Task<int> GetPageCountByAuthorAndFollowing(int authorId);

    public Task<ICollection<Author>> GetFollowersByIdAsync(int authorId);
    public Task<ICollection<Author>> GetFollowingByIdAsync(int authorId);

    public Task AddFollowAsync(int followingAuthorId, int followedAuthorId);
    public Task RemoveFollowAsync(int followingAuthorId, int followedAuthorId);

    public Task DeleteCheepsByAuthorIdAsync(int authorId);

    public Task RemoveAllFollowersByAuthorIdAsync(int id);

    public Task RemoveUserByIdAsync(int id);

    public Task RemoveReactionsByAuthorIdAsync(int id);
    public Task SaveContextAsync();
}
