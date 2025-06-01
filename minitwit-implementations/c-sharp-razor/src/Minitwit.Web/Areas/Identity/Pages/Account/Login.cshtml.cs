// Licensed to the .NET Foundation under one or more agreements.
// The .NET Foundation licenses this file to you under the MIT license.
#nullable disable

using System.ComponentModel.DataAnnotations;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Minitwit.Core.Entities;

namespace Minitwit.Web.Areas.Identity.Pages.Account
{   
    [IgnoreAntiforgeryToken]
    public class LoginModel : PageModel
    {
        private readonly SignInManager<Author> _signInManager;

        public LoginModel(SignInManager<Author> signInManager)
        {
            _signInManager = signInManager;
        }
        
        /// <summary>
        ///     This API supports the ASP.NET Core Identity default UI infrastructure and is not intended to be used
        ///     directly from your code. This API may change or be removed in future releases.
        /// </summary>
        public IList<AuthenticationScheme> ExternalLogins { get; set; }

        /// <summary>
        ///     This API supports the ASP.NET Core Identity default UI infrastructure and is not intended to be used
        ///     directly from your code. This API may change or be removed in future releases.
        /// </summary>
        public string ReturnUrl { get; set; }

        /// <summary>
        ///     This API supports the ASP.NET Core Identity default UI infrastructure and is not intended to be used
        ///     directly from your code. This API may change or be removed in future releases.
        /// </summary>
        [TempData]
        public string ErrorMessage { get; set; }

        /// <summary>
        ///     This API supports the ASP.NET Core Identity default UI infrastructure and is not intended to be used
        ///     directly from your code. This API may change or be removed in future releases.
        /// </summary>
        
        /// <summary>
        ///     This API supports the ASP.NET Core Identity default UI infrastructure and is not intended to be used
        ///     directly from your code. This API may change or be removed in future releases.
        /// </summary>
        [BindProperty]
        [Required]
        [DataType(DataType.Text)]
        public string Username { get; set; }
            
        /// <summary>
        ///     This API supports the ASP.NET Core Identity default UI infrastructure and is not intended to be used
        ///     directly from your code. This API may change or be removed in future releases.
        /// </summary>
        [BindProperty]
        [Required]
        [DataType(DataType.Text)]
        public string Password { get; set; }

        /// <summary>
        ///     This API supports the ASP.NET Core Identity default UI infrastructure and is not intended to be used
        ///     directly from your code. This API may change or be removed in future releases.
        /// </summary>
        [Display(Name = "Remember me?")]
        public bool RememberMe { get; set; }
        
            
        public async Task<IActionResult> OnGetAsync()
        {
            if (!string.IsNullOrEmpty(ErrorMessage))
            {
                ModelState.AddModelError(string.Empty, ErrorMessage);
            }

            if (_signInManager.IsSignedIn(User))
            {
                return Redirect("/public");
            }

            // Clear the existing external cookie to ensure a clean login process
            await HttpContext.SignOutAsync(IdentityConstants.ExternalScheme);

            ExternalLogins = (
                await _signInManager.GetExternalAuthenticationSchemesAsync()
            ).ToList();

            return Page();
        }

        public async Task<IActionResult> OnPostAsync(string returnUrl = null)
        {
            ExternalLogins = (
                await _signInManager.GetExternalAuthenticationSchemesAsync()
            ).ToList();

            if (ModelState.IsValid)
            {
                var result = await _signInManager.PasswordSignInAsync(
                    Username,
                    Password,
                    RememberMe,
                    lockoutOnFailure: false
                );
                if (!result.Succeeded)
                {
                    ModelState.AddModelError(string.Empty, "Invalid username");
                    ModelState.AddModelError(string.Empty, "Invalid password");
                    return Page();
                }
                TempData["FlashMessage"] = "You were logged in";
                return RedirectToPage("/MyTimeline"); 
            }
            
            return Page();
        }
    }
}
