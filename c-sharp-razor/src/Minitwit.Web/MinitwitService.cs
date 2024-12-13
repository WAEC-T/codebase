using System.Globalization;
using Minitwit.Core.Entities;
using Minitwit.Core.Repository;
using Minitwit.Web.Models;

namespace Minitwit.Web;

public interface IMessageService
{
    public Task<ICollection<MessageViewModel>> GetMessagesAsync(int page);
    public Task<ICollection<MessageViewModel>> GetMessagesFromAuthor(string authorName, int page);
    public Task<ICollection<MessageViewModel>> GetMessagesFromAuthorAsync(int authorId, int page);
    public Task<ICollection<MessageViewModel>> GetMessagesFromAuthorAndFollowingAsync(
        int authorId,
        int page
    );
}

public class MinitwitService : IMessageService
{
    private readonly IAuthorRepository _authorRepository;
    private readonly IMessageRepository _MessageRepository;

    public MinitwitService(
        IMessageRepository MessageRepository,
        IAuthorRepository authorRepository
    )
    
    {
        _MessageRepository = MessageRepository;
        _authorRepository = authorRepository;
    }

    public async Task<ICollection<MessageViewModel>> GetMessagesAsync(int page)
    {
        ICollection<Message> MessageDtos = await _MessageRepository.GetMessagesByPageAsync(page);
        var authorIds = MessageDtos.Select(c => c.AuthorId).Distinct();
        
        ICollection<Author> authors = await _authorRepository.GetAuthorsByIdAsync(authorIds);

        List<MessageViewModel> Messages = new List<MessageViewModel>();

        foreach (var MessageDto in MessageDtos)
        {
            Author? author = authors.FirstOrDefault(a => a.Id == MessageDto.AuthorId);

            Messages.Add(
                new MessageViewModel(
                    MessageDto.MessageId,
                    new UserModel(author),
                    MessageDto.Text,
                    MessageDto.TimeStamp,
                    MessageDto.Flagged
                )
            );
        }

        return Messages;
    }

    public async Task<ICollection<MessageViewModel>> GetMessagesFromAuthorAsync(int authorId, int page)
    {
        ICollection<Message> MessageDtos = await _authorRepository.GetMessagesByAuthor(authorId, page);
        ICollection<MessageViewModel> Messages = new List<MessageViewModel>();
        Author author = await _authorRepository.GetAuthorByIdAsync(authorId);
        
        foreach (Message MessageDto in MessageDtos)
        {
            Messages.Add(
                new MessageViewModel(
                    MessageDto.MessageId,
                    new UserModel(author),
                    MessageDto.Text,
                    MessageDto.TimeStamp,
                    MessageDto.Flagged
                )
            );
        }

        return Messages;
    }

    public async Task<ICollection<MessageViewModel>> GetMessagesFromAuthorAndFollowingAsync(
        int authorId,
        int page
    )
    {
        var messagesWithAuthors =await _MessageRepository.GetMessagesFromAuthor(authorId, page);
        
        var Messages = messagesWithAuthors.Select(ma => new MessageViewModel(
                ma.Message.MessageId,
                new UserModel(ma.Author),
                ma.Message.Text,
                ma.Message.TimeStamp,
                ma.Message.Flagged
            )).ToList();

            return Messages;
    }

    public async Task<ICollection<MessageViewModel>> GetMessagesFromAuthor(string authorName, int page)
    {
        Author author = await _authorRepository.GetAuthorByNameAsync(authorName);
        var Messages = await GetMessagesFromAuthorAsync(author.Id, page);
        return Messages;
    }
}
