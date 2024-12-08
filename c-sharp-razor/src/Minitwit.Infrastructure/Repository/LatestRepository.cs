using Microsoft.EntityFrameworkCore;
using Minitwit.Core.Entities;
using Minitwit.Core.Repository;

namespace Minitwit.Infrastructure.Repository;

public class LatestRepository : BaseRepository, ILatestRepository
{
    public LatestRepository(MinitwitDbContext DbContext)
        : base(DbContext) { }

    public async Task<Latest> GetLatestAsync()
    {
        return await db.LatestEntries.FirstOrDefaultAsync();
    }

    public async Task AddLatestAsync(Latest latest)
    {
        await db.LatestEntries.AddAsync(latest);
        await db.SaveChangesAsync();
    }

    public async Task UpdateLatestAsync(Latest latest)
    {
        db.LatestEntries.Update(latest);
        await db.SaveChangesAsync();
    }
}


