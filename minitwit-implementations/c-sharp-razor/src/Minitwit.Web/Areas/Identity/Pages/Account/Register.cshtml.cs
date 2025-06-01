// Licensed to the .NET Foundation under one or more agreements.
// The .NET Foundation licenses this file to you under the MIT license.
#nullable disable

using System.ComponentModel.DataAnnotations;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.UI.Services;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Minitwit.Core.Entities;

namespace Minitwit.Web.Areas.Identity.Pages.Account
{   
    [IgnoreAntiforgeryToken]
    public class RegisterModel : PageModel
    {
        private readonly SignInManager<Author> _signInManager;
        private readonly UserManager<Author> _userManager;
        
        public RegisterModel(
            UserManager<Author> userManager,
            SignInManager<Author> signInManager)
        {
            _userManager = userManager;
            _signInManager = signInManager;
        }

        /// <summary>
        ///     This API supports the ASP.NET Core Identity default UI infrastructure and is not intended to be used
        ///     directly from your code. This API may change or be removed in future releases.
        /// </summary>
        public string ReturnUrl { get; set; }

        /// <summary>
        ///     This API supports the ASP.NET Core Identity default UI infrastructure and is not intended to be used
        ///     directly from your code. This API may change or be removed in future releases.
        /// </summary>
        public IList<AuthenticationScheme> ExternalLogins { get; set; }

        public object ErrorMessage { get; set; }
        
        [BindProperty]
        [Required(ErrorMessage = "You have to enter a username")]
        [DataType(DataType.Text)]
        [Display(Name = "Username")]
        public string Username { get; set; }

        [BindProperty]
        [Required(ErrorMessage = "You have to enter a valid email address")]
        [EmailAddress]
        [Display(Name = "Email")]
        public string Email { get; set; }

        [BindProperty]
        [Required(ErrorMessage = "You have to enter a password")]
        [DataType(DataType.Password)]
        [Display(Name = "Password")]
        public string Password { get; set; }
        
        [BindProperty(Name = "password2")]
        [Required(ErrorMessage = "The two passwords do not match")]
        [DataType(DataType.Password)]
        [Display(Name = "Confirm Password")]
        [Compare("Password", ErrorMessage = "The two passwords do not match")]
        public string Password2 { get; set; }
        
        
        public async Task OnGetAsync(string returnUrl = null)
        {
            ReturnUrl = returnUrl;
            ExternalLogins = (
                await _signInManager.GetExternalAuthenticationSchemesAsync()
            ).ToList();
        }
        
        public async Task<IActionResult> OnPostAsync()
        {
            if (!ModelState.IsValid)
            {
                return Page();
            }

            if (await _userManager.FindByNameAsync(Username) != null)
            {
                ModelState.AddModelError("Input.Username", "The username is already taken");
                return Page();
            }
            
            var newUser = new Author
            {
                UserName = Username,
                Email = Email,
                NormalizedEmail = Email.ToUpperInvariant(),
            };
            var result = await _userManager.CreateAsync(newUser, Password);
            
            if (result.Succeeded)
            {
                TempData["FlashMessage"] = "You were successfully registered and can login now";
                return RedirectToPage("/Account/Login", new { area = "Identity" });
            }
            
            foreach (var error in result.Errors)
            {
                ModelState.AddModelError(string.Empty, error.Description);
            }

            return Page();
        }
    }
}
