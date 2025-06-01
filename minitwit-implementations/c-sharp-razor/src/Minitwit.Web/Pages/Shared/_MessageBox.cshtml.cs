using Microsoft.AspNetCore.Mvc.RazorPages;

namespace Minitwit.Web.Pages.Shared
{
    public class MessageBoxModel : PageModel
    {
        public string? MessageText { get; set; }
    }
}
