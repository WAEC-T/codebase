using Minitwit.Core.Entities;

namespace Minitwit.Core.Repository;

public interface ILatestRepository
{
    Task<Latest> GetLatestAsync();
    
    Task AddLatestAsync(Latest latest);
    Task UpdateLatestAsync(Latest latest);
}