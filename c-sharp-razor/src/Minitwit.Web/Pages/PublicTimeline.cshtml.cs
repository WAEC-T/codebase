using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Minitwit.Web.Models;

namespace Minitwit.Web.Pages;

public class PublicTimelineModel : PageModel
{
    private readonly IMessageService _service;
    public required ICollection<MessageViewModel> Messages { get; set; }
    public required int CurrentPage { get; set; }

    public PublicTimelineModel(IMessageService service)
    {
        _service = service;
    }
    
    public async Task<ActionResult> OnGet()
    {   
        await InitializeVariables();
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
        Messages = await _service.GetMessagesAsync(page) ?? new List<MessageViewModel>();
        CurrentPage = page;
    }
}