using System.ComponentModel.DataAnnotations;
using FluentValidation;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Minitwit.Core.Entities;
using Minitwit.Core.Repository;
using Minitwit.Web.Models;
using ValidationException = FluentValidation.ValidationException;

namespace Minitwit.Web.Pages;

public class PublicModel : PageModel
{
    private readonly IMessageService _service;
    private readonly IMessageRepository _MessageRepository;
    private readonly IAuthorRepository _authorRepository;
    private readonly IFollowRepository _followRepository;
    private readonly IValidator<CreateMessage> _validator;
    public required Author user { get; set; }
    private readonly UserManager<Author> _userManager;
    public required ICollection<MessageViewModel> Messages { get; set; }
    public required int totalPages { get; set; }
    public required int currentPage { get; set; }

    public PublicModel(
        IMessageService service,
        IMessageRepository MessageRepository,
        IAuthorRepository authorRepository,
        IFollowRepository followRepository,
        IValidator<CreateMessage> validator,
        UserManager<Author> userManager
    )
    {
        _service = service;
        _MessageRepository = MessageRepository;
        _authorRepository = authorRepository;
        _followRepository = followRepository;
        _validator = validator;
        _userManager = userManager;
    }

    public async Task<ActionResult> OnGet()
    {
        await InitializeVariables();
        return Page();
    }

    [BindProperty]
    public NewMessage? NewMessage { get; set; }

    public async Task<IActionResult> OnPostCreateMessage()
    {
        if (!ModelState.IsValid)
        {
            return Page();
        }

        var author = await _userManager.GetUserAsync(User);
        var Message = new CreateMessage(author!.Id, NewMessage!.Text!);

        await CreateMessage(Message);

        return RedirectToPage("/UserTimeline", new { author = User.Identity?.Name });
    }

    public async Task CreateMessage(CreateMessage newMessage)
    {
        var validationResult = await _validator.ValidateAsync(newMessage);

        if (!validationResult.IsValid)
        {
            Console.WriteLine(validationResult);
            //Fatal exception
            throw new ValidationException(
                "The Message must be between 5 and 160 characters.(CreateMessage)"
            );
        }

        await _MessageRepository.AddCreateMessageAsync(newMessage);
    }

    public async Task<IActionResult> OnPostReaction(
        int MessageId,
        int currentPage
    )
    {
        Author? author = await _userManager.GetUserAsync(User);
        await InitializeVariables(currentPage);
        return Page();
    }

    public async Task<IActionResult> OnPostRemoveReaction(
        int MessageId,
        int currentPage
    )
    {
        Author? author = await _userManager.GetUserAsync(User);
        await InitializeVariables(currentPage);
        return Page();
    }

    public async Task<IActionResult> OnPostFollow(int currentPage, int Author2Follow)
    {
        Author? author = await _authorRepository.GetAuthorByIdAsync(
            _userManager.GetUserAsync(User).Result!.Id
        );
        Author? authorToFollow = await _authorRepository.GetAuthorByIdAsync(Author2Follow);
        await InitializeVariables(currentPage);

        if (author == null)
            return Page();

        if (authorToFollow != null)
            await _authorRepository.AddFollowAsync(author.Id, authorToFollow.Id);
        return Page();
    }

    public async Task<IActionResult> OnPostUnfollow(int currentPage, int Author2Unfollow)
    {
        Author? author = await _authorRepository.GetAuthorByIdAsync(
            _userManager.GetUserAsync(User).Result!.Id
        );
        Author? authorToUnfollow = await _authorRepository.GetAuthorByIdAsync(Author2Unfollow);

        await InitializeVariables(currentPage);

        if (authorToUnfollow == null || author == null)
            return Page();

        await _authorRepository.RemoveFollowAsync(author.Id, authorToUnfollow.Id);
        return Page();
    }

    public async Task InitializeVariables()
    {
        int page;
        if (Request.Query.ContainsKey("page"))
        {
            page = int.Parse(Request.Query["page"]!);
        }
        else
        {
            page = 1;
        }
        await InitializeVariables(page);
    }

    public async Task InitializeVariables(int page)
    {
        Messages = await _service.GetMessagesAsync(page);

        user = _userManager.GetUserAsync(User).Result!;
        totalPages = await _MessageRepository.GetPageCountAsync();
        currentPage = page;
    }
}

public class NewMessage
{
    [Required]
    [StringLength(
        160,
        MinimumLength = 5,
        ErrorMessage = "The Message must be between 5 and 160 characters(NewMessage)."
    )]
    public string? Text { get; set; }
}
