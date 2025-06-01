using System.Security.Claims;
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
        if (Request.Query.ContainsKey("page"))
        {
            currentPage = int.Parse(Request.Query["page"]!);
        }
        else
        {
            currentPage = 1;
        }
        TimelineAuthor = await _authorRepository.GetAuthorByNameAsync(author);
        if (user != null)
        {
            IsFollowing = await _followRepository.IsFollowingAsync(user.Id, TimelineAuthor.Id);
            await LoadMessages(user, TimelineAuthor, currentPage);
        }
        else
        {
            await LoadMessages(null, TimelineAuthor, currentPage);
        }
        return Page();
    }
    
    private async Task LoadMessages(Author? signedInAuthor, Author timelineAuthor, int page)
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
        var currentUserId = User.FindFirst(ClaimTypes.NameIdentifier)?.Value;
        if (currentUserId == null)
        {
            return Unauthorized();
        }
        var authorToFollow = await _authorRepository.GetAuthorByNameAsync(author); 
        if (authorToFollow == null)
        {
            return NotFound();
        }
        
        await _authorRepository.AddFollowAsync(int.Parse(currentUserId), authorToFollow.Id); 
        TempData["FlashMessage"] = $"You are now following {authorToFollow.UserName}";

        Response.Redirect($"/user/{authorToFollow.UserName}");
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
        Response.Redirect($"/user/{authorToUnfollow.UserName}");
        return new EmptyResult();
    }

}


