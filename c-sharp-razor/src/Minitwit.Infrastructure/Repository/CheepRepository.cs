using Microsoft.EntityFrameworkCore;
using Minitwit.Core.Entities;
using Minitwit.Core.Repository;

namespace Minitwit.Infrastructure.Repository;

public class MessageRepository : BaseRepository, IMessageRepository
{
    public MessageRepository(MinitwitDbContext DbContext)
        : base(DbContext) { }

    public async Task<ICollection<Message>> GetMessagesByPageAsync(int page)
    {
        //Use EF to get the specified page of Messages from the database
        ICollection<Message> Messages = await db
            .Messages.OrderByDescending(c => c.TimeStamp)
            .Skip(PageSize * (page - 1))
            .Take(PageSize)
            .ToListAsync();

        return Messages;
    }

    public async Task<ICollection<Message>> GetMessagesByCountAsync(int count)
    {
        ICollection<Message> Messages = await db
            .Messages.OrderByDescending(c => c.TimeStamp)
            .Take(count)
            .AsNoTracking()
            .ToListAsync();

        return Messages;
    }

    public async Task<ICollection<Message>> GetMessagesFromAuthorByCountAsync(int authorId, int count)
    {
        //Use EF to get the specified count of Messages from an author from the database
        ICollection<Message> Messages = await db
            .Messages.Where(c => c.AuthorId == authorId)
            .OrderByDescending(c => c.TimeStamp)
            .Take(count)
            .ToListAsync();

        return Messages;
    }

    public async Task<int> GetMessageCountAsync()
    {
        //Use EF to get the total number of Messages from the database
        return await db.Messages.CountAsync();
    }

    public async Task<int> GetPageCountAsync()
    {
        return await GetMessageCountAsync() / PageSize + 1;
    }

    public async Task DeleteMessageByIdAsync(int MessageId)
    {
        //Delete the specified Message from the database
        Message? Message = await db.Messages.FindAsync(MessageId);
        if (Message != null)
        {
            db.Messages.Remove(Message);
        }
        else
        {
            throw new Exception("Message with id " + MessageId + " not found");
        }

        await db.SaveChangesAsync();
    }

    public async Task AddMessageAsync(Message Message)
    {
        await db.Messages.AddAsync(Message);
        await db.SaveChangesAsync();
        Console.WriteLine("Message added async");
    }

    public async Task<Message> AddCreateMessageAsync(CreateMessage Message)
    {
        Message entity = new Message()
        {
            MessageId = new int(),
            Text = Message.Text,
            TimeStamp = DateTime.Now,
            AuthorId = Message.AuthorId
        };

        await AddMessageAsync(entity);

        return entity;
    }
}
