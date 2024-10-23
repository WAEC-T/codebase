using Minitwit.Core.Entities;

namespace Minitwit.Core.Repository;

public interface IFollowRepository
{
    public Task<Follow> CreateFollowAsync(int followingAuthorId, int followedAuthorId);
    public Task DeleteFollowAsync(Follow follow);
    public Task<bool> IsFollowingAsync(int followingUserId, int followedAuthorId);
}
