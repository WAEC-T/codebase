using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Minitwit.Core.Entities;
using Minitwit.Core.Repository;
using Minitwit.Web;
using Minitwit.Web.Models;

namespace Minitwit.Razor.Pages;

[IgnoreAntiforgeryToken]
public class UserTimelineModel : PageModel
{
    private readonly IMessageService _service;
    private readonly UserManager<Author> _userManager;
    private readonly IAuthorRepository _authorRepository;
    private readonly SignInManager<Author> _signInManager;
    private readonly IFollowRepository _followRepository;
    
    public ICollection<MessageViewModel>? Messages { get; set; }
    
    public bool IsFollowing { get; set; }

    public required Author? user { get; set; }

    public Author TimelineAuthor { get; set; }
    
    public required int currentPage { get; set; }
    public required int totalPages { get; set; }

    public UserTimelineModel(
        IMessageService service,
        SignInManager<Author> signInManager,
        UserManager<Author> userManager,
        IAuthorRepository authorRepository,
        IFollowRepository followRepository
    )
    {
        _service = service;
        _userManager = userManager;
        _authorRepository = authorRepository;
        _signInManager = signInManager;
        _followRepository = followRepository;
        
    }
    
    public async Task<ActionResult> OnGet(string author)
    {
        user = await _userManager.GetUserAsync(User);
        if (user == null)
        {
            return Unauthorized();
        }
        await InitializeVariables(user!, author);

        return Page();
    }

    public async Task InitializeVariables(Author user, string author)
    {   
        if (Request.Query.ContainsKey("page"))
        {
            currentPage = int.Parse(Request.Query["page"]!);
        }
        else
        {
            currentPage = 1;
        }
        
        TimelineAuthor = await _authorRepository.GetAuthorByNameAsync(author); 
        if (TimelineAuthor == null)
        {
            throw new KeyNotFoundException($"Author with username '{author}' was not found.");
        }
        IsFollowing = await _followRepository.IsFollowingAsync(user.Id, TimelineAuthor.Id); 
        await LoadMessages(user, TimelineAuthor, currentPage);
    }

    private async Task LoadMessages(Author signedInAuthor, Author timelineAuthor, int page)
    {
        try
        {
            if (
                _signInManager.IsSignedIn(User)
                && signedInAuthor.UserName == timelineAuthor.UserName
            )
            {
                Messages = await _service.GetMessagesFromAuthorAndFollowingAsync(
                    signedInAuthor.Id,
                    page
                );
                totalPages = await _authorRepository.GetPageCountByAuthorAndFollowing(
                    signedInAuthor.Id
                );
            }
            else
            {
                Messages = await _service.GetMessagesFromAuthorAsync(timelineAuthor.Id, page);
                totalPages = await _authorRepository.GetPageCountByAuthor(timelineAuthor.Id);
            }
        }
        catch (Exception)
        {
            Messages = new List<MessageViewModel>();
        }
    }
    
    public async Task<IActionResult> OnGetFollow(string author)
    {
        var currentUser = await _userManager.GetUserAsync(User); // 1
        if (currentUser == null)
        {
            return Unauthorized();
        }
        Console.WriteLine("In Follow Page");
        var authorToFollow = await _authorRepository.GetAuthorByNameAsync(author); // 2
        if (authorToFollow == null)
        {
            return NotFound();
        }

        await _authorRepository.AddFollowAsync(currentUser.Id, authorToFollow.Id); // 3
        TempData["FlashMessage"] = $"You are now following {authorToFollow.UserName}";

        Response.Redirect($"/{authorToFollow.UserName}");
        return new EmptyResult();
    }
    
    public async Task<IActionResult> OnGetUnfollow(string author)
    {   
        Author? currentUser = await _userManager.GetUserAsync(User);
        Author? authorToUnfollow = await _authorRepository.GetAuthorByNameAsync(author);

        if (currentUser == null)
            return NotFound();
        
        await _authorRepository.RemoveFollowAsync(currentUser.Id, authorToUnfollow.Id);
        TempData["FlashMessage"] = "You are no longer following " + authorToUnfollow.UserName;
        Response.Redirect($"/{authorToUnfollow.UserName}");
        return new EmptyResult();
    }

}

