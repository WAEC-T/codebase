using Microsoft.EntityFrameworkCore;
using Minitwit.Core.Entities;
using Minitwit.Core.Repository;

namespace Minitwit.Infrastructure.Repository;

public class AuthorRepository : BaseRepository, IAuthorRepository
{
    private readonly IFollowRepository _followRepository;

    public AuthorRepository(MinitwitDbContext minitwitDbContext, IFollowRepository followRepository)
        : base(minitwitDbContext)
    {
        _followRepository = followRepository;
    }

    // ----- Add Author Methods ----- //
    public async void AddAuthorAsync(Author authorDto)
    {
        await db.Users.AddAsync(authorDto);
        await db.SaveChangesAsync();
    }

    // ----- Get Author Methods ----- //

    public async Task<ICollection<Author>> GetAllAuthorsAsync()
    {
        return await db.Users.ToListAsync();
    }

    public async Task<ICollection<Author>> GetAuthorsByIdAsync(IEnumerable<int> authors)
    {
        return await db.Users.Where(a => authors.Contains(a.Id)).AsNoTracking().ToListAsync();
    }

    public async Task<Author?> GetAuthorByIdAsync(int authorId)
    {
        Author? author = await db.Users.FirstOrDefaultAsync(a => a.Id == authorId);
        return author!;
    }

    public async Task<Author> GetAuthorByNameAsync(string name)
    {
        Author? author = await db.Users.FirstOrDefaultAsync(a => a.UserName == name)!;
        return author!;
    }

    // ----- Get Messages By Author and Page Methods ----- //
    public async Task<ICollection<Message>> GetMessagesByAuthorAsync(int id)
    {
        return await db.Messages.Where(e => e.AuthorId == id).ToListAsync();
    }

    public async Task<ICollection<Message>> GetMessagesByAuthor(int authorId, int page)
    {
        var messages = await GetMessagesByAuthorAsync(authorId);

        //Check that author has Messages
        if (messages == null || messages.Count == 0)
            throw new InvalidOperationException("The specified author has no messages.");
        
        if (page < 1)
            page = 1;

        int pageSizeIndex = (page - 1) * PageSize;

        if (messages.Count < pageSizeIndex + PageSize)
            return messages
                .ToList()
                .GetRange(pageSizeIndex, messages.Count - pageSizeIndex)
                .OrderByDescending(c => c.TimeStamp)
                .ToList();
        if (messages.Count > PageSize)
            return messages
                .ToList()
                .GetRange(pageSizeIndex, PageSize)
                .OrderByDescending(c => c.TimeStamp)
                .ToList();
        return messages.OrderByDescending(c => c.TimeStamp).ToList();
    }

    public async Task<ICollection<Message>> GetMessagesByAuthorAndFollowingAsync(int id)
    {
        ICollection<Message> Messages = new List<Message>(await GetMessagesByAuthorAsync(id));

        foreach (Author author in await GetFollowingByIdAsync(id))
        {
            Messages = Messages.Concat(await GetMessagesByAuthorAsync(author.Id)).ToList();
        }

        return Messages;
    }

    public async Task<ICollection<Message>> GetMessagesByAuthorAndFollowing(int authorId, int page)
    {   
        //Get Messages from the author, and append Messages from followers to that list
        ICollection<Author> following = await GetFollowingByIdAsync(authorId);
        ICollection<Message> Messages = new List<Message>();

        // Add all the users Messages to the list without pagination
        foreach (var MessageDto in await GetMessagesByAuthorAsync(authorId))
            Messages.Add(MessageDto);

        foreach (Author? follower in following)
        {
            ICollection<Message> followingMessages = await GetMessagesByAuthorAsync(follower.Id);
            //If follower has no Messages, skip them
            if (followingMessages.Count == 0)
            {
                continue;
            }

            //Add each Message from the follower to the list
            foreach (var MessageDto in followingMessages)
            {
                Messages.Add(MessageDto);
            }
        }
        //Sort the Messages according to timestamp, latest first
        Messages = Messages.OrderByDescending(c => c.TimeStamp).ToList();

        int pageSizeIndex = (page - 1) * PageSize;

        if (Messages.Count < pageSizeIndex + PageSize)
            return Messages.ToList<Message>().GetRange(pageSizeIndex, Messages.Count - pageSizeIndex);
        if (Messages.Count > PageSize)
            return Messages.ToList<Message>().GetRange(pageSizeIndex, PageSize);
        return Messages;
    }

    // ----- Get Messages By Author Methods ----- //
    public async Task<int> GetMessageCountByAuthorAsync(int authorId)
    {
        ICollection<Message> Messages = await GetMessagesByAuthorAsync(authorId);
        //Check that author has Messages
        if (Messages.Count == 0 || Messages == null)
            return 0;

        return Messages.Count;
    }

    public async Task<int> GetMessageCountByAuthorAndFollowing(int authorId)
    {
        ICollection<Message> Messages = await GetMessagesByAuthorAndFollowingAsync(authorId);
        return Messages.Count;
    }

    // ----- Get Page Count Methods ----- //
    public async Task<int> GetPageCountByAuthor(int authorId)
    {
        return await GetMessageCountByAuthorAsync(authorId) / PageSize + 1;
    }

    public async Task<int> GetPageCountByAuthorAndFollowing(int authorId)
    {
        return await GetMessageCountByAuthorAndFollowing(authorId) / PageSize + 1;
    }

    // ----- Get Followers and Following Methods ----- //
    public async Task<ICollection<Author>> GetFollowersByIdAsync(int authorId)
    {
        // Query to retrieve the IDs of authors followed by the specified author
        List<int> followedAuthorIds = await db
            .Follows.Where(f => f.FollowedAuthorId == authorId)
            .Select(f => f.FollowingAuthorId)
            .ToListAsync();

        // Query to retrieve the author entities based on the followed author IDs
        ICollection<Author> followedAuthors = await db
            .Users.Where(a => followedAuthorIds.Contains(a.Id))
            .ToListAsync();

        return followedAuthors;
    }

    public async Task<ICollection<Author>> GetFollowingByIdAsync(int authorId)
    {   
        // Query to retrieve the IDs of authors followed by the specified author
        var followingAuthorIds = await db
            .Follows.Where(f => f.FollowingAuthorId == authorId)
            .Select(f => f.FollowedAuthorId)
            .ToListAsync();

        // Query to retrieve the author entities based on the followed author IDs
        ICollection<Author> followingAuthors = await db
            .Users.Where(a => followingAuthorIds.Contains(a.Id))
            .ToListAsync();

        return followingAuthors;
    }

    public async Task<bool> AuthorExists(int id)
    {
        Author author = await GetAuthorByIdAsync(id);
        return author.Id != null;
    }

    // ----- Add/Remove Follow Methods ----- //
    public async Task AddFollowAsync(int followingAuthorId, int followedAuthorId)
    {
        await _followRepository.CreateFollowAsync(followingAuthorId, followedAuthorId);
    }

    public async Task RemoveFollowAsync(int followingAuthorId, int followedAuthorId)
    {
        Follow follow = await db.Follows.FirstOrDefaultAsync(e =>
            e.FollowedAuthorId == followedAuthorId && e.FollowingAuthorId == followingAuthorId
        )!;
        await _followRepository.DeleteFollowAsync(follow);
    }

    // ----- Delete Author Data Methods ----- //
    public async Task DeleteMessagesByAuthorIdAsync(int authorId)
    {
        var Messages = await GetMessagesByAuthorAsync(authorId);

        foreach (var Message in Messages)
        {
            // Delete the Message itself
            db.Messages.Remove(Message);
        }

        await db.SaveChangesAsync();
    }

    public async Task RemoveAllFollowersByAuthorIdAsync(int id)
    {
        Author? user = await GetAuthorByIdAsync(id);
        if (user == null)
        {
            throw new KeyNotFoundException($"Author with ID {id} was not found.");

        }

        var follows = await db
            .Follows.Where(f => f.FollowedAuthorId == id || f.FollowingAuthorId == id)
            .ToListAsync();
        db.Follows.RemoveRange(follows);
    }

    public async Task RemoveUserByIdAsync(int id)
    {
        Author? user = await GetAuthorByIdAsync(id);
        if (user == null)
            throw new KeyNotFoundException($"Author with ID {id} was not found.");

        db.Users.Remove(user);
        await db.SaveChangesAsync();
    }

    public Task RemoveReactionsByAuthorIdAsync(int id)
    {
        throw new NotImplementedException();
    }

    // ----- Save Context Method ----- //
    public async Task SaveContextAsync()
    {
        await db.SaveChangesAsync();
    }
}
