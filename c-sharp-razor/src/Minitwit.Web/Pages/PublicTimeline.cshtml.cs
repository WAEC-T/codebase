using FluentValidation;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Minitwit.Core.Entities;
using Minitwit.Core.Repository;
using Minitwit.Web.Models;

namespace Minitwit.Web.Pages;

public class PublicTimelineModel : PageModel
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

    public PublicTimelineModel(
        IMessageService service,
        IMessageRepository MessageRepository,
        IAuthorRepository authorRepository,
        IFollowRepository followRepository,
        IValidator<CreateMessage> validator,
        UserManager<Author> userManager
    )
    {   
        Console.WriteLine("In PublicModel");
        _service = service;
        _MessageRepository = MessageRepository;
        _authorRepository = authorRepository;
        _followRepository = followRepository;
        _validator = validator;
        _userManager = userManager;
    }
    
    [Route("/public")]
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
        if (author == null)
        {   
            return RedirectToPage("/Login");
        }
        var Message = new CreateMessage(author.Id, NewMessage!.Text!);
        await CreateMessage(Message);
        TempData["FlashMessage"] = "Your message was recorded";
        string userTimelineUrl = $"/{User.Identity.Name}";
        return Redirect(userTimelineUrl);
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
        Messages = await _service.GetMessagesAsync(page) ?? new List<MessageViewModel>();
        user = _userManager.GetUserAsync(User).Result!;
        totalPages = await _MessageRepository.GetPageCountAsync();
        currentPage = page;
    }
}

public class NewMessage
{
    public string? Text { get; set; }
}
