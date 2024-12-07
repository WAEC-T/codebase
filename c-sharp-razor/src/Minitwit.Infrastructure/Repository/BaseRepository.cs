namespace Minitwit.Infrastructure.Repository;

public abstract class BaseRepository
{
    protected MinitwitDbContext db;
    protected int PageSize { get; set; }

    protected BaseRepository(MinitwitDbContext minitwitDbContext)
    {
        db = minitwitDbContext;
        PageSize = 30;
    }
}
