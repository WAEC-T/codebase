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
        // Fetch Messages for the given page.
        ICollection<Message> MessageDtos = await _MessageRepository.GetMessagesByPageAsync(page);

        // Extract unique author IDs from the Messages.
        var authorIds = MessageDtos.Select(c => c.AuthorId).Distinct();

        // Fetch only the authors who authored the fetched Messages.
        ICollection<Author> authors = await _authorRepository.GetAuthorsByIdAsync(authorIds);

        // Initialize a list to hold the MessageViewModels.
        List<MessageViewModel> Messages = new List<MessageViewModel>();

        // Process each MessageDto sequentially.
        foreach (var MessageDto in MessageDtos)
        {
            // Find the author for the current Message.
            Author? author = authors.FirstOrDefault(a => a.Id == MessageDto.AuthorId);

            // Create and add the MessageViewModel to the list.
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
        ICollection<Message> MessageDtos = await _authorRepository.GetMessagesByAuthorAndFollowing(
            authorId,
            page
        );
        ICollection<Author> authors = await _authorRepository.GetFollowingByIdAsync(authorId);
        authors.Add(await _authorRepository.GetAuthorByIdAsync(authorId));
        ICollection<MessageViewModel> Messages = new List<MessageViewModel>();

        foreach (Message MessageDto in MessageDtos)
        {
            Author? author = authors.FirstOrDefault(a => a.Id == MessageDto.AuthorId);

            Messages.Add(
                new MessageViewModel(
                    MessageDto.MessageId,
                    new UserModel(author!),
                    MessageDto.Text,
                    MessageDto.TimeStamp,
                    MessageDto.Flagged
                )
            );
        }

        return Messages;
    }

    public async Task<ICollection<MessageViewModel>> GetMessagesFromAuthor(string authorName, int page)
    {
        Author author = await _authorRepository.GetAuthorByNameAsync(authorName);
        var Messages = await GetMessagesFromAuthorAsync(author.Id, page);
        return Messages;
    }
}
