using FluentValidation;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Minitwit.Core.Entities;
using Minitwit.Core.Repository;
using Minitwit.Web.Models;

namespace Minitwit.Web.Pages;

[IgnoreAntiforgeryToken]
public class MyTimelineModel : PageModel
{
    private readonly IMessageService _service;
    private readonly IMessageRepository _MessageRepository;
    private readonly SignInManager<Author> _signInManager;
    private readonly UserManager<Author> _userManager;
    public required Author user { get; set; }
    public required ICollection<MessageViewModel> Messages { get; set; }
    public required int currentPage { get; set; }

    public MyTimelineModel(
        IMessageService service,
        IMessageRepository  messageRepository,
        IAuthorRepository authorRepository,
        IFollowRepository followRepository,
        IValidator<CreateMessage> validator,
        UserManager<Author> userManager,
        SignInManager<Author> signInManager
    )
    {   
        _service = service;
        _MessageRepository = messageRepository;
        _userManager = userManager;
        _signInManager = signInManager;
    }
    
    public async Task<ActionResult> OnGetAsync()
    {   
        if (!_signInManager.IsSignedIn(User))
        {
            return Redirect("/public");
        }
        
        await InitializeVariables();
        return Page();
    }

    [BindProperty(Name = "text", SupportsGet = false)]
    public string Text { get; set; }
    
    public async Task<IActionResult> OnPostCreateMessage()
    {   
        if (!ModelState.IsValid)
        {   
            return Page();
        }
        if (string.IsNullOrWhiteSpace(Text))
        {
            return Page();
        }
        var author = await _userManager.GetUserAsync(User);
        if (author == null)
        {   
            return RedirectToPage("/Login");
        }
        
        user = author;
        var Message = new CreateMessage(author.Id, Text);
        await CreateMessage(Message);
        TempData["FlashMessage"] = "Your message was recorded";
        return RedirectToPage("/MyTimeline");
    }

    public async Task CreateMessage(CreateMessage newMessage)
    {
        await _MessageRepository.AddCreateMessageAsync(newMessage);
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
        user = _userManager.GetUserAsync(User).Result!;
        await LoadMessages(user, page);
        currentPage = page;
    }
    
    private async Task LoadMessages(Author signedInAuthor, int page)
    {
        try
        {
            if (
                _signInManager.IsSignedIn(User)
            )
            {
                Messages = await _service.GetMessagesFromAuthorAndFollowingAsync(
                    signedInAuthor.Id,
                    page
                );
            }
        }
        catch (Exception)
        {
            Messages = new List<MessageViewModel>();
        }
    }
}
