using System.Globalization;
using Minitwit.Core.Entities;
using Minitwit.Core.Repository;
using Minitwit.Web.Models;

namespace Minitwit.Web;

public interface ICheepService
{
    public Task<ICollection<CheepViewModel>> GetCheepsAsync(int page);
    public Task<ICollection<CheepViewModel>> GetCheepsFromAuthor(string authorName, int page);
    public Task<ICollection<CheepViewModel>> GetCheepsFromAuthorAsync(int authorId, int page);
    public Task<ICollection<CheepViewModel>> GetCheepsFromAuthorAndFollowingAsync(
        int authorId,
        int page
    );
}

public class MinitwitService : ICheepService
{
    private readonly IAuthorRepository _authorRepository;
    private readonly ICheepRepository _cheepRepository;

    public MinitwitService(
        ICheepRepository cheepRepository,
        IAuthorRepository authorRepository
    )
    {
        _cheepRepository = cheepRepository;
        _authorRepository = authorRepository;
    }

    public async Task<ICollection<CheepViewModel>> GetCheepsAsync(int page)
    {
        // Fetch cheeps for the given page.
        ICollection<Cheep> cheepDtos = await _cheepRepository.GetCheepsByPageAsync(page);

        // Extract unique author IDs from the cheeps.
        var authorIds = cheepDtos.Select(c => c.AuthorId).Distinct();

        // Fetch only the authors who authored the fetched cheeps.
        ICollection<Author> authors = await _authorRepository.GetAuthorsByIdAsync(authorIds);

        // Initialize a list to hold the CheepViewModels.
        List<CheepViewModel> cheeps = new List<CheepViewModel>();

        // Process each cheepDto sequentially.
        foreach (var cheepDto in cheepDtos)
        {
            // Find the author for the current cheep.
            Author? author = authors.FirstOrDefault(a => a.Id == cheepDto.AuthorId);

            // Create and add the CheepViewModel to the list.
            cheeps.Add(
                new CheepViewModel(
                    cheepDto.CheepId,
                    new UserModel(author),
                    cheepDto.Text,
                    cheepDto.TimeStamp.ToString("o")
                )
            );
        }

        return cheeps;
    }

    public async Task<ICollection<CheepViewModel>> GetCheepsFromAuthorAsync(int id, int page)
    {
        ICollection<Cheep> cheepDtos = await _authorRepository.GetCheepsByAuthor(id, page);
        ICollection<CheepViewModel> cheeps = new List<CheepViewModel>();
        Author author = await _authorRepository.GetAuthorByIdAsync(id);

        foreach (Cheep cheepDto in cheepDtos)
        {
            cheeps.Add(
                new CheepViewModel(
                    cheepDto.CheepId,
                    new UserModel(author),
                    cheepDto.Text,
                    cheepDto.TimeStamp.ToString(CultureInfo.InvariantCulture)
                )
            );
        }

        return cheeps;
    }

    public async Task<ICollection<CheepViewModel>> GetCheepsFromAuthorAndFollowingAsync(
        int authorId,
        int page
    )
    {
        ICollection<Cheep> cheepDtos = await _authorRepository.GetCheepsByAuthorAndFollowing(
            authorId,
            page
        );
        ICollection<Author> authors = await _authorRepository.GetFollowingByIdAsync(authorId);
        authors.Add(await _authorRepository.GetAuthorByIdAsync(authorId));
        ICollection<CheepViewModel> cheeps = new List<CheepViewModel>();

        foreach (Cheep cheepDto in cheepDtos)
        {
            Author? author = authors.FirstOrDefault(a => a.Id == cheepDto.AuthorId);

            cheeps.Add(
                new CheepViewModel(
                    cheepDto.CheepId,
                    new UserModel(author!),
                    cheepDto.Text,
                    cheepDto.TimeStamp.ToString(CultureInfo.InvariantCulture)
                )
            );
        }

        return cheeps;
    }

    public async Task<ICollection<CheepViewModel>> GetCheepsFromAuthor(string authorName, int page)
    {
        Author author = await _authorRepository.GetAuthorByNameAsync(authorName);
        var cheeps = await GetCheepsFromAuthorAsync(author.Id, page);
        return cheeps;
    }
}
