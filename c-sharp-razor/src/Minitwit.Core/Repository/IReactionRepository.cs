using Minitwit.Core.Entities;

namespace Minitwit.Core.Repository;

public interface IReactionRepository
{
    public Task AddReaction(ReactionType reaction, int cheepId, int authorId);
    public Task RemoveReaction(ReactionType reaction, int cheepId, int authorId);

    public Task<ICollection<Reaction>> GetReactionsFromCheepIdAsync(int id);
    public Task<int> GetReactionCount(int cheepId, ReactionType reactionType);
    public Task<bool> HasUserReactedAsync(int cheepId, int authorId);
}
